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

import {Initializable} from "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import {UUPSUpgradeable} from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import {OwnableUpgradeable} from "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import {PausableUpgradeable} from "@openzeppelin/contracts-upgradeable/security/PausableUpgradeable.sol";

import "./ISequencerInbox.sol";
import "./libraries/DeserializationLib.sol";
import "./libraries/Errors.sol";

contract SequencerInbox is ISequencerInbox, Initializable, UUPSUpgradeable, OwnableUpgradeable, PausableUpgradeable {
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
        __Ownable_init();
        __Pausable_init();
        __UUPSUpgradeable_init();
    }

    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }

    function pause() public onlyOwner {
        _pause();
    }

    function unpause() public onlyOwner {
        _unpause();
    }

    function _authorizeUpgrade(address) internal override onlyOwner whenPaused {}

    /// @inheritdoc IDAProvider
    function getInboxSize() external view override returns (uint256) {
        return inboxSize;
    }

    /// @inheritdoc ISequencerInbox
    function appendTxBatch(
        uint256[] calldata contexts,
        uint256[] calldata txLengths,
        uint256 firstL2BlockNumber,
        bytes calldata txBatch
    ) external override whenNotPaused {
        if (msg.sender != sequencerAddress) {
            revert NotSequencer(msg.sender, sequencerAddress);
        }

        uint256 numTxs = inboxSize;
        bytes32 runningAccumulator;
        if (accumulators.length > 0) {
            runningAccumulator = accumulators[accumulators.length - 1];
        }

        uint256 dataOffset = 0;
        uint256 l2BlockNumber = firstL2BlockNumber;

        for (uint256 i = 0; i + 2 <= contexts.length; i += 2) {
            // TODO: consider adding L1 context.
            uint256 l2Timestamp = contexts[i + 1];
            bytes32 txContextHash = keccak256(abi.encodePacked(sequencerAddress, l2BlockNumber, l2Timestamp));

            uint256 numCtxTxs = contexts[i];

            for (uint256 j = 0; j < numCtxTxs; j++) {
                uint256 txLength = txLengths[numTxs - inboxSize];
                if (dataOffset + txLength > txBatch.length) {
                    revert TxBatchDataOverflow();
                }
                bytes32 txDataHash = keccak256(txBatch[dataOffset:dataOffset + txLength]);

                runningAccumulator = keccak256(abi.encodePacked(runningAccumulator, numTxs, txContextHash, txDataHash));

                dataOffset += txLength;
                numTxs++;
            }

            // block numbers get incremented by one
            // we can reconstruct all block numbers of the batch since we know the first one
            l2BlockNumber++;
        }

        if (numTxs <= inboxSize) revert EmptyBatch();
        uint256 start = inboxSize;
        inboxSize = numTxs;
        accumulators.push(runningAccumulator);

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
