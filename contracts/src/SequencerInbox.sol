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

import "./ISequencerInbox.sol";
import "./libraries/DeserializationLib.sol";
import "./libraries/Errors.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";

contract SequencerInbox is ISequencerInbox, Initializable {
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

        for (uint256 i = 0; i + 3 <= contexts.length; i += 3) {
            // TODO: consider adding L1 context.
            uint256 l2BlockNumber = contexts[i + 1];
            uint256 l2Timestamp = contexts[i + 2];
            bytes32 prefixHash = keccak256(abi.encodePacked(msg.sender, l2BlockNumber, l2Timestamp));

            uint256 numCtxTxs = contexts[i];
            for (uint256 j = 0; j < numCtxTxs; j++) {
                uint256 txLength = txLengths[numTxs - inboxSize];
                bytes32 txDataHash;
                assembly {
                    txDataHash := keccak256(dataOffset, txLength)
                }
                runningAccumulator = keccak256(abi.encodePacked(runningAccumulator, numTxs, prefixHash, txDataHash));
                dataOffset += txLength;
                if (dataOffset - initialDataOffset > txBatch.length) {
                    revert TxBatchDataOverflow();
                }
                numTxs++;
            }
        }

        if (numTxs <= inboxSize) revert EmptyBatch();
        uint256 start = inboxSize;
        inboxSize = numTxs;
        accumulators.push(runningAccumulator);

        emit TxBatchAppended(accumulators.length - 1, start, inboxSize);
    }

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
}
