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

pragma solidity ^0.8.0;

// import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "./ISequencerInbox.sol";
import "./libraries/DeserializationLib.sol";
import "./libraries/Errors.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";

contract SequencerInbox is ISequencerInbox, Initializable {
    string private constant EMPTY_BATCH = "EMPTY_BATCH";

    // Total number of transactions
    uint256 private inboxSize;
    // accumulators[i] is an accumulator of transactions in txBatch i.
    bytes32[] public accumulators;

    address public sequencerAddress;

    function initialize(address _sequencerAddress) public initializer {
        if (_sequencerAddress == address(0)) {
            revert ZeroAddress();
        }
        sequencerAddress = _sequencerAddress;
    }

    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }

    function getInboxSize() external view override returns (uint256) {
        return inboxSize;
    }

    function appendTxBatch(uint256[] calldata contexts, uint256[] calldata txLengths, bytes calldata txBatch)
        external
        override
    {
        if (msg.sender != sequencerAddress) {
            revert NotSequencer(msg.sender, sequencerAddress);
        }

        uint256 start = inboxSize;
        uint256 numTxs = inboxSize;
        uint256 numProcessedTxs = 0;
        bytes32 runningAccumulator;
        if (accumulators.length > 0) {
            runningAccumulator = accumulators[accumulators.length - 1];
        }

        uint256 dataOffset;
        assembly {
            dataOffset := txBatch.offset
        }

        for (uint256 i = 0; i + 3 <= contexts.length; i += 3) {
            // TODO: consider adding L1 context.
            uint256 l2BlockNumber = contexts[i + 1];
            uint256 l2Timestamp = contexts[i + 2];
            bytes32 prefixHash = keccak256(abi.encodePacked(msg.sender, l2BlockNumber, l2Timestamp));

            uint256 numCtxTxs = contexts[i];
            for (uint256 j = 0; j < numCtxTxs; j++) {
                uint256 txLength = txLengths[numProcessedTxs];
                bytes32 txDataHash;
                assembly {
                    txDataHash := keccak256(dataOffset, txLength)
                }
                runningAccumulator = keccak256(abi.encodePacked(runningAccumulator, numTxs, prefixHash, txDataHash));
                dataOffset += txLength;
                numTxs++;
            }
            numProcessedTxs += numCtxTxs;
        }

        if (numTxs <= inboxSize) revert EmptyBatch();
        inboxSize = numTxs;
        accumulators.push(runningAccumulator);

        emit TxBatchAppended(accumulators.length - 1, start, inboxSize);
    } // -> hashchain of transactions

    // proof -> what comes before the txn, the txn and then what comes after the txn
    function verifyTxInclusion(bytes memory proof) external view override {
        uint256 offset = 0;

        // Deserialize tx and tx info.
        address sender;
        uint256 l2BlockNumber;
        uint256 l2Timestamp;
        uint256 txDataLength;
        bytes32 txDataHash;
        (offset, sender) = DeserializationLib.deserializeAddress(proof, offset);
        (offset, l2BlockNumber) = DeserializationLib.deserializeUint256(proof, offset);
        (offset, l2Timestamp) = DeserializationLib.deserializeUint256(proof, offset);
        (offset, txDataLength) = DeserializationLib.deserializeUint256(proof, offset);
        assembly {
            // TODO: check if off-by-32.
            txDataHash := keccak256(add(proof, offset), txDataLength)
        }
        offset += txDataLength;

        // Deserialize inbox info.
        uint256 batchNum;
        uint256 numTxs;
        uint256 numTxsAfterInBatch;
        bytes32 acc;
        (offset, batchNum) = DeserializationLib.deserializeUint256(proof, offset);
        (offset, numTxs) = DeserializationLib.deserializeUint256(proof, offset);
        (offset, numTxsAfterInBatch) = DeserializationLib.deserializeUint256(proof, offset);
        (offset, acc) = DeserializationLib.deserializeBytes32(proof, offset);

        // Start accumulator at the tx.
        bytes32 prefixHash = keccak256(abi.encodePacked(sender, l2BlockNumber, l2Timestamp));
        acc = keccak256(abi.encodePacked(acc, numTxs, prefixHash, txDataHash));
        numTxs++;

        // Compute final accumulator value.
        for (uint256 i = 0; i < numTxsAfterInBatch; i++) {
            (offset, prefixHash) = DeserializationLib.deserializeBytes32(proof, offset);
            (offset, txDataHash) = DeserializationLib.deserializeBytes32(proof, offset);
            acc = keccak256(abi.encodePacked(acc, numTxs, prefixHash, txDataHash));
            numTxs++;
        }

        if (acc != accumulators[batchNum]) {
            revert IncorrectAccOrBatch();
        }
    }

    /**
     *
     *     EXTREMELY DANGEROUS FUNCTION(S). DO NOT DEPLOY BEFORE DELETING THESE FUNCTION(S) FROM CONTRACT     *********
     *
     */
    function dangerousIncreaseSequencerInboxSize(uint256 newInboxSize) external returns (uint256) {
        if (msg.sender != sequencerAddress) {
            revert NotSequencer(msg.sender, sequencerAddress);
        }

        inboxSize = newInboxSize;

        uint256 changedInboxSize = inboxSize;
        return changedInboxSize;
    }
}
/*
context -> Metadata about which blocks the txn it belongs to
txnBatch -> blob of transactions 
each txnBatch is of 5 bytes and then the txnBatch is a repeating strings of 5 batches!!*/
