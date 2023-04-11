// SPDX-License-Identifier: Apache-2.0

/*
 * Modifications Copyright 2022, Specular contributors
 *
 * This file was changed in accordance to Apache License, Version 2.0.
 *
 * Copyright 2021, Offchain Labs, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

pragma solidity ^0.8.4;

import "./ISequencerInbox.sol";
import "./libraries/DeserializationLib.sol";
import "./libraries/Errors.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import "hardhat/console.sol";

contract SequencerInbox is ISequencerInbox, Initializable, UUPSUpgradeable, OwnableUpgradeable {
    // Total number of transactions
    uint256 private inboxSize;
    // accumulators[i] is an accumulator of transactions in txBatch i.
    bytes32[] public accumulators;
    // delayedMessages successfully read from the delayedInbox
    uint256 public delayedMessagesRead;
    // delayedInbox hashes array
    bytes32[] public delayedInboxAccumulator;

    uint256 public delayedMessageCounter;

    bytes32[] public messageDataHashes;

    address public sequencerAddress;

    uint64[2] public l2BlockAndTime;

    uint64[2] public l1BlockAndTime;

    uint256 public baseFee;

    function initialize(address _sequencerAddress) public initializer {
        if (_sequencerAddress == address(0)) {
            revert ZeroAddress();
        }
        sequencerAddress = _sequencerAddress;
        __Ownable_init();
        __UUPSUpgradeable_init();
    }

    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }

    function _authorizeUpgrade(address) internal override onlyOwner {}

    /// @inheritdoc IDAProvider
    function getInboxSize() external view override returns (uint256) {
        return inboxSize;
    }

    function sendUnsignedTx(
        uint256 gasLimit,
        uint256 maxFeePerGas,
        uint256 nonce,
        address to,
        uint256 value,
        bytes calldata data
    ) external {
        console.log("something");
        if (gasLimit > type(uint256).max) revert();

        addToDelayedInbox(msg.sender, keccak256(abi.encodePacked(gasLimit, maxFeePerGas, nonce, to, value, data)));
    }

    function addToDelayedInbox(address _sender, bytes32 _messageDataHash)
        internal
        returns (uint256 delayedMessageCount)
    {

        messageDataHashes.push(_messageDataHash);
        delayedMessageCount = delayedInboxAccumulator.length;
        console.log("delayed message count %d", delayedMessageCount);
        // generating a message hash
        bytes32 messageHash = keccak256(
            abi.encodePacked(
                _sender, uint64(block.number), uint64(block.timestamp), delayedMessageCount, block.basefee, _messageDataHash
            )
        );

        console.log("Printing message Hash");
        console.logBytes32(messageHash);

        bytes32 prevAcc = 0;
        if (delayedMessageCount > 0) {
            prevAcc = delayedInboxAccumulator[delayedMessageCount - 1];
        }

        // saving the block data for forceInclusion (it is for testing purposes as of now)
        l1BlockAndTime[0] = uint64(block.number);
        l1BlockAndTime[1] = uint64(block.timestamp);
        baseFee = block.basefee;

        // adding the message to delayedInbox
        delayedInboxAccumulator.push(keccak256(abi.encodePacked(prevAcc, messageHash)));
        delayedMessageCounter++;
        return delayedMessageCount;
    }

    function forceInclusion(
        uint256 delayedMessageIndex,
        address _sender,
        uint256 baseFeeL1,
        uint64[2] calldata _l1BlockAndTime,
        bytes32 messageDataHash
    ) external {
        // check to avoid invalid index
        if (delayedMessageIndex >= delayedInboxAccumulator.length) revert();
        // check to avoid already included transactions
        if (delayedMessageIndex < delayedMessagesRead) revert ();

        // calculating messageHash with all the given message data
        bytes32 messageHash = keccak256(
            abi.encodePacked(
                _sender, _l1BlockAndTime[0], _l1BlockAndTime[1], delayedMessageIndex, baseFeeL1, messageDataHash
            )
        );
        console.logBytes32(messageHash);

        // enforcing the 1 day time limit
        if (_l1BlockAndTime[0] + 5760 >= block.number) revert ();
        if (_l1BlockAndTime[1] + 86400 >= block.timestamp) revert ();

        console.logBytes32(messageHash);
        console.logBytes32(delayedInboxAccumulator[delayedMessageIndex]);

        bytes32 prevDelayedAcc = 0;
        if (delayedInboxAccumulator.length > 1) {
            prevDelayedAcc = delayedInboxAccumulator[delayedMessageIndex - 1];
        }

        // messageHash should be identical to hash stored in delayedInbox
        if (keccak256(abi.encodePacked(prevDelayedAcc, messageHash)) != delayedInboxAccumulator[delayedMessageIndex]) revert EmptyBatch();

        bytes32 txContentHash = keccak256(abi.encodePacked(_sender, l2BlockAndTime[0], l2BlockAndTime[1]));

        bytes32 runningAccumulator;
        if (accumulators.length > 0) {
            runningAccumulator = accumulators[accumulators.length - 1];
        }
        uint256 numTxs = inboxSize;
        for (uint256 i = delayedMessagesRead; i <= delayedMessageIndex; i++) {
            bytes32 txDataHash = delayedInboxAccumulator[i];
            runningAccumulator = keccak256(abi.encodePacked(runningAccumulator, numTxs, txContentHash, txDataHash));
            numTxs++;
        }
        if (numTxs <= inboxSize) revert EmptyBatch();
        inboxSize = numTxs;

        // pushing it to main accumulator
        accumulators.push(runningAccumulator);

        delayedMessagesRead = delayedMessageIndex + 1;
    }

    /// @inheritdoc ISequencerInbox
    function appendTxBatch(
        uint256[] calldata contexts,
        uint256[] calldata txLengths,
        bytes calldata txBatch,
        uint256 _totalDelayedMessagesRead
    ) external override {
        if (msg.sender != sequencerAddress) {
            revert NotSequencer(msg.sender, sequencerAddress);
        }
        // check - if messages already read
        if (_totalDelayedMessagesRead < delayedMessagesRead) revert();
        // check - if new index is valid
        if (_totalDelayedMessagesRead > delayedInboxAccumulator.length) revert();

        uint256 numTxs = inboxSize;

        bytes32 runningAccumulator;
        if (accumulators.length > 0) {
            runningAccumulator = accumulators[accumulators.length - 1];
        }

        uint256 initialDataOffset;

        assembly {
            initialDataOffset := txBatch.offset
        }


        uint256 dataOffset = initialDataOffset;

        uint256 l2BlockNumber;
        uint256 l2Timestamp;
        bytes32 txContextHash;

        for (uint256 i = 0; i + 3 <= contexts.length; i += 3) {
            // TODO: consider adding L1 context.
            l2BlockNumber = contexts[i + 1];
            l2Timestamp = contexts[i + 2];
            txContextHash = keccak256(abi.encodePacked(sequencerAddress, l2BlockNumber, l2Timestamp));

            uint256 numCtxTxs = contexts[i];
            for (uint256 j = 0; j < numCtxTxs; j++) {
                uint256 txLength = txLengths[numTxs - inboxSize];
                bytes32 txDataHash;
                assembly {
                    txDataHash := keccak256(dataOffset, txLength)
                }
                runningAccumulator = keccak256(abi.encodePacked(runningAccumulator, numTxs, txContextHash, txDataHash));
                dataOffset += txLength;
                if (dataOffset - initialDataOffset > txBatch.length) {
                    revert TxBatchDataOverflow();
                }
                numTxs++;
            }
        }


        bytes32 prevDelayedAccumulator = 0;
        if (prevDelayedAccumulator.length > 0) {
            prevDelayedAccumulator = prevDelayedAccumulator[prevDelayedAccumulator.length - 1];
        }

        // number of delayed Tx being added.
        numTxs +=  _totalDelayedMessagesRead - delayedMessagesRead;

        runningAccumulator = keccak256(abi.encodePacked(runningAccumulator, numTxs, txContextHash, prevDelayedAccumulator));

        if (numTxs <= inboxSize) revert EmptyBatch();
        uint256 start = inboxSize;
        inboxSize = numTxs;
        accumulators.push(runningAccumulator);

        delayedMessagesRead = _totalDelayedMessagesRead;
        l2BlockAndTime[0] = uint64(l2BlockNumber);
        l2BlockAndTime[1] = uint64(l2Timestamp);

        emit TxBatchAppended(accumulators.length - 1, start, inboxSize);
    }

    // TODO post EIP-4844: KZG proof verification
    // https://eips.ethereum.org/EIPS/eip-4844#point-evaluation-precompile

    /**
     * @notice Verifies that a transaction is included in a batch, at the expected offset.
     * @param encodedTx Transaction to verify inclusion of.
     * @param proof Proof of inclusion, in the form:
     * proof := txContextHash || batchInfo || {foreach tx in batch: (txContextHash || KEC(txData)), ...} where,
     * batchInfo := (batchNum || numTxsBefore || numTxsAfterInBatch || accBefore)
     * txContextHash := KEC(sequencerAddress || l2BlockNumber || l2Timestamp)
     */
    function verifyTxInclusion(bytes calldata encodedTx, bytes calldata proof) external view override {
        uint256 offset = 0;
        // Deserialize tx context of `encodedTx`.
        bytes32 txContextHash;
        (offset, txContextHash) = DeserializationLib.deserializeBytes32(proof, offset);
        // Deserialize batch info.
        uint256 batchNum;
        uint256 numTxs;
        uint256 numTxsAfterInBatch;
        bytes32 acc;
        (offset, batchNum) = DeserializationLib.deserializeUint256(proof, offset);
        (offset, numTxs) = DeserializationLib.deserializeUint256(proof, offset);
        (offset, numTxsAfterInBatch) = DeserializationLib.deserializeUint256(proof, offset);
        (offset, acc) = DeserializationLib.deserializeBytes32(proof, offset);

        // Start accumulator at the tx.
        bytes32 txDataHash = keccak256(encodedTx);
        acc = keccak256(abi.encodePacked(acc, numTxs, txContextHash, txDataHash));
        numTxs++;

        // Compute final accumulator value.
        for (uint256 i = 0; i < numTxsAfterInBatch; i++) {
            (offset, txContextHash) = DeserializationLib.deserializeBytes32(proof, offset);
            (offset, txDataHash) = DeserializationLib.deserializeBytes32(proof, offset);
            acc = keccak256(abi.encodePacked(acc, numTxs, txContextHash, txDataHash));
            numTxs++;
        }

        if (acc != accumulators[batchNum]) {
            revert ProofVerificationFailed();
        }
    }
}
