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
    // for testing purposes
    uint256 public delayedMessageCounter;

    // for testing purpose
    bytes32 public testRunAcc;

    address public sequencerAddress;

    uint64[2] public l2BlockAndTime; // proof
    // for testing purposes
    uint64[2] public l1BlockAndTime;


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

    // function sendUnsignedTx(
    //     uint256 gasLimit,
    //     uint256 maxFeePerGas,
    //     uint256 nonce,
    //     address to,
    //     uint256 value,
    //     bytes calldata data
    // ) external {
    //     console.log("something");
    //     if (gasLimit > type(uint256).max) revert();

    //     addToDelayedInbox(msg.sender, keccak256(abi.encodePacked(gasLimit, maxFeePerGas, nonce, to, value, data)));
    // }

    function sendRLPEncodedTx(
        bytes calldata messageDataHash
    ) external {

        addToDelayedInbox(keccak256(messageDataHash));
    }

    function addToDelayedInbox(
        bytes32 messageDataHash
    ) internal {

        bytes32 messageHash = keccak256(
            abi.encodePacked(
                uint64(block.number), uint64(block.timestamp), delayedInboxAccumulator.length, messageDataHash
            )
        );

        // for testing purposes
        l1BlockAndTime[0] = uint64(block.number);
        l1BlockAndTime[1] = uint64(block.timestamp);

        bytes32 prevAcc = 0;
        if (delayedInboxAccumulator.length > 0) {
            prevAcc = delayedInboxAccumulator[delayedInboxAccumulator.length - 1];
        }

        delayedInboxAccumulator.push(keccak256(abi.encodePacked(prevAcc, messageHash)));
        delayedMessageCounter++;

    }


    // function addToDelayedInbox(address _sender, bytes32 _messageDataHash)
    //     internal
    //     returns (uint256 delayedMessageCount)
    // {
    //     messageDataHashes.push(_messageDataHash);
    //     delayedMessageCount = delayedInboxAccumulator.length;
    //     // console.log("delayed message count %d", delayedMessageCount);
    //     // generating a message hash
    //     bytes memory messageHash_1 =
    //         abi.encodePacked(
    //             _sender, uint64(block.number), uint64(block.timestamp), delayedMessageCount, block.basefee, _messageDataHash
    //         );

    //     bytes32 messageHash = keccak256(messageHash_1);
    //     console.logBytes32(messageHash);
    //     bytes32 prevAcc = 0;
    //     if (delayedMessageCount > 0) {
    //         prevAcc = delayedInboxAccumulator[delayedMessageCount - 1];
    //     }

    //     // saving the block data for forceInclusion (it is for testing purposes as of now)
    //     l1BlockAndTime[0] = uint64(block.number);
    //     l1BlockAndTime[1] = uint64(block.timestamp);
    //     baseFee = block.basefee;

    //     // adding the message to delayedInbox
    //     delayedInboxAccumulator.push(keccak256(abi.encodePacked(prevAcc, messageHash)));
    //     delayedMessageCounter++;
    //     return delayedMessageCount;
    // }


    // proof => prevAcc || numTxs || txContextHash || txDataHash
    // txContextHash = (sequencerAddress || l2Block || l2time)
    function forceInclusion(
        bytes calldata blockdataProof,
        uint256 delayedMessageIndex,
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
                 _l1BlockAndTime[0], _l1BlockAndTime[1], delayedMessageIndex, messageDataHash
            )
        );
        // console.logBytes32(messageHash);

        // enforcing the 1 day time limit
        if (_l1BlockAndTime[0] + 5760 >= block.number) revert ();
        if (_l1BlockAndTime[1] + 86400 >= block.timestamp) revert ();

        // console.logBytes32(messageHash);
        // console.logBytes32(delayedInboxAccumulator[delayedMessageIndex]);

        bytes32 prevDelayedAcc = 0;
        if (delayedInboxAccumulator.length > 1) {
            prevDelayedAcc = delayedInboxAccumulator[delayedMessageIndex - 1];
        }

        // messageHash should be identical to hash stored in delayedInbox
        if (keccak256(abi.encodePacked(prevDelayedAcc, messageHash)) != delayedInboxAccumulator[delayedMessageIndex]) revert EmptyBatch();


        // checking proof correctness
        uint256 offset = 0;
        // Deserialize tx context of `encodedTx`.
        bytes32 prevAcc;
        (offset, prevAcc) = DeserializationLib.deserializeBytes32(blockdataProof, offset);

        uint256 prevNumTxs;
        (offset, prevNumTxs) = DeserializationLib.deserializeUint256(blockdataProof, offset);

        uint256 lastL2Block;
        (offset, lastL2Block) = DeserializationLib.deserializeUint256(blockdataProof, offset);

        uint256 lastL2Timestamp;
        (offset, lastL2Timestamp) = DeserializationLib.deserializeUint256(blockdataProof, offset);

        bytes32 prevTxDataHash;
        (offset, prevTxDataHash) = DeserializationLib.deserializeBytes32(blockdataProof, offset);

        bytes32 lastContextHash = keccak256(abi.encodePacked(sequencerAddress, lastL2Block, lastL2Timestamp));

        if(accumulators.length > 0) {

            if( keccak256(abi.encodePacked(prevAcc, prevNumTxs, lastContextHash, prevTxDataHash)) !=  accumulators[accumulators.length - 1]) revert();

        }


        // bytes32 txContentHash = keccak256(abi.encodePacked(_sender, l2BlockAndTime[0], l2BlockAndTime[1]));

        bytes32 runningAccumulator;
        if (accumulators.length > 0) {
            runningAccumulator = accumulators[accumulators.length - 1];
        }
        uint256 numTxs = inboxSize;
        for (uint256 i = delayedMessagesRead; i <= delayedMessageIndex; i++) {
            bytes32 txDataHash = delayedInboxAccumulator[i];
            runningAccumulator = keccak256(abi.encodePacked(runningAccumulator, numTxs, lastContextHash, txDataHash));
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
        console.log("AppendTxBatch");
        console.logBytes(txBatch);
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


        uint256 offset = 0;

        // uint256 initialDataOffset = 0;
        // // assembly {
        // //     initialDataOffset := txBatch.offset
        // // }


        // uint256 dataOffset = initialDataOffset;
        // console.log("dataOffset", dataOffset);

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
                console.log(txLength);

                bytes memory trx;
                bytes32 txDataHash;
                trx = BytesLib.slice(txBatch, offset, txLength);
                offset += txLength;

                txDataHash = keccak256(trx);
                // assembly {
                //     txDataHash := keccak256(dataOffset, txLength)
                // }
                console.log("TxData Hash:  ");
                console.logBytes32(txDataHash);

                runningAccumulator = keccak256(abi.encodePacked(runningAccumulator, numTxs, txContextHash, txDataHash));

                // dataOffset += txLength;
                // if (dataOffset - initialDataOffset > txBatch.length) {
                //     revert TxBatchDataOverflow();
                // }

                numTxs++;
            }
        }
        // console.log("APPENDTX BATCH DATA::::: ");
        // console.logBytes32(runningAccumulator);

        testRunAcc = runningAccumulator;

        if (_totalDelayedMessagesRead != delayedMessagesRead) {

            bytes32 prevDelayedAccumulator = 0;

            if (delayedInboxAccumulator.length > 0) {
                prevDelayedAccumulator = delayedInboxAccumulator[delayedInboxAccumulator.length - 1];
            }

            // number of delayed Tx being added.
            numTxs +=  _totalDelayedMessagesRead - delayedMessagesRead;

            // console.log(numTxs);
            // console.log("PrevDelayedAccumulator: " );
            // console.logBytes32(txContextHash);
            // console.logBytes32(prevDelayedAccumulator);

            runningAccumulator = keccak256(abi.encodePacked(runningAccumulator, numTxs, txContextHash, prevDelayedAccumulator));
        }

        if (numTxs <= inboxSize) revert EmptyBatch();
        uint256 start = inboxSize;
        inboxSize = numTxs;

        accumulators.push(runningAccumulator);

        console.log("ACCUMULATOR ::::::");
        console.logBytes32(runningAccumulator);

        delayedMessagesRead = _totalDelayedMessagesRead;
        // l2BlockAndTime[0] = uint64(l2BlockNumber);
        // l2BlockAndTime[1] = uint64(l2Timestamp);

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
        console.logBytes(proof);
        console.log("VERIFICATION :: ");
        uint256 offset = 0;
        // Deserialize tx context of `encodedTx`.
        bytes32 txContextHash;
        (offset, txContextHash) = DeserializationLib.deserializeBytes32(proof, offset);
        console.log("txContextHash: ");
        console.logBytes32(txContextHash);
        // Deserialize batch info.
        uint256 batchNum;
        uint256 numTxs;
        uint256 numTxsAfterInBatch;
        bytes32 acc;

        (offset, batchNum) = DeserializationLib.deserializeUint256(proof, offset);
        console.log("batch Number: ",batchNum);

        (offset, numTxs) = DeserializationLib.deserializeUint256(proof, offset);
        console.log("Num of txs: ",numTxs);

        (offset, numTxsAfterInBatch) = DeserializationLib.deserializeUint256(proof, offset);
        console.log("Txs after in batch: ",numTxsAfterInBatch);

        (offset, acc) = DeserializationLib.deserializeBytes32(proof, offset);
        console.log("Accumulator before value: ");
        console.logBytes32(acc);

        bytes32 txDataHash = keccak256(encodedTx);
        console.log("Tx Data Hash:  ");
        console.logBytes32(txDataHash);

        acc = keccak256(abi.encodePacked(acc, numTxs, txContextHash, txDataHash));
        console.log("Accumulator after value: ");
        console.logBytes32(acc);

        numTxs++;

        // Compute final accumulator value.
        for (uint256 i = 0; i < numTxsAfterInBatch; i++) {
            (offset, txContextHash) = DeserializationLib.deserializeBytes32(proof, offset);
            (offset, txDataHash) = DeserializationLib.deserializeBytes32(proof, offset);

            acc = keccak256(abi.encodePacked(acc, numTxs, txContextHash, txDataHash));
            numTxs++;
        }

        console.log("Accumulator final value: ");
        console.logBytes32(acc);

        console.log("Accumulator real value: ");
        console.logBytes32(accumulators[batchNum]);

        if (acc != accumulators[batchNum]) {
            revert ProofVerificationFailed();
        }

        console.log("ENDED");
    }

    /**
     * @notice Verifies that a transaction is included in a batch, at the expected offset.
     * @param encodedTx Transaction to verify inclusion of.
     * @param proof Proof of inclusion, in the form:
     * proof := txContextHash || batchInfo || delayedInfo || {foreach delayedTx in batch: ( messageHash ), ...} where,
     * batchInfo := (batchNum || numTxsBefore || accBefore)
     * delayedInfo := (numTxsBefore || numTxsAfterInDelayedAcc || delayedAccBefore)
     * txContextHash := KEC(sequencerAddress || l2BlockNumber || l2Timestamp)
     */
    function verifyDelayedTxInclusion(bytes calldata encodedTx, bytes calldata proof) external view override {
        uint256 offset = 0;
        // Deserialize tx context of `encodedTx`.
        bytes32 txContextHash;
        // console.log("TxContexHash: ");
        // console.logBytes32(txContextHash);
        (offset, txContextHash) = DeserializationLib.deserializeBytes32(proof, offset);
        // console.logBytes32(txContextHash);

        // Deserialize batch info.
        uint256 batchNum;
        // console.log("BatchNum", batchNum);
        uint256 numTxs; // num txs before  // might not need it for delayed verification
        // console.log("numTxsBefore", numTxs);
        bytes32 acc; // main accumulator before
        // console.log("acc before");
        // console.logBytes32(acc);
        (offset, batchNum) = DeserializationLib.deserializeUint256(proof, offset);
        // console.log("BatchNum", batchNum);

        (offset, numTxs) = DeserializationLib.deserializeUint256(proof, offset);
        // console.log("numTxsBefore", numTxs);

        (offset, acc) = DeserializationLib.deserializeBytes32(proof, offset);
        // console.log("acc before");
        // console.logBytes32(acc);

        // delayed txs related data
        uint256 numTxsBefore;
        uint256 numTxsAfter;
        bytes32 delayedAcc; // before

        (offset, numTxsBefore) = DeserializationLib.deserializeUint256(proof, offset);
        // console.log("Delayed: numTxsBefore", numTxsBefore);
        (offset, numTxsAfter) = DeserializationLib.deserializeUint256(proof, offset);
        // console.log("Delayed: numTxsAfter", numTxsAfter);
        (offset, delayedAcc) = DeserializationLib.deserializeBytes32(proof, offset);
        // console.log("Delayed: delayedAcc before");
        // console.logBytes32(delayedAcc);


        // Start accumulator at the tx.
        bytes32 txDataHash = keccak256(encodedTx);
        delayedAcc = keccak256(abi.encodePacked(delayedAcc, txDataHash));
        // console.log("VERIFICATION DATAAAAAAAA::");
        // console.log("Verification Delayed Accumulator");
        // console.logBytes32(delayedAcc);

        // Compute final delayed accumulator value.
        for (uint256 i = 0; i < numTxsAfter; i++) {
            (offset, txDataHash) = DeserializationLib.deserializeBytes32(proof, offset);
            delayedAcc = keccak256(abi.encodePacked(delayedAcc, txDataHash));
        }
        uint256 totalDelayedTxs = numTxsBefore + numTxsAfter + 1 + numTxs;

        // console.logBytes32(acc);
        // console.log(totalDelayedTxs);
        // console.logBytes32(txContextHash);
        // console.logBytes32(delayedAcc);

        acc = keccak256(abi.encodePacked(acc, totalDelayedTxs, txContextHash, delayedAcc));

        // console.log("VERFICATION ACCUMULATOR RESULT");
        // console.logBytes32(acc);

        if (acc != accumulators[batchNum]) {
            revert ProofVerificationFailed();
        }
    }

}
