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

import "./IDAProvider.sol";

/**
 * @notice On-chain DA provider.
 */
interface ISequencerInbox is IDAProvider {
    event TxBatchAppended(uint256 batchNumber);

    /// @dev Thrown when the given tx inclusion proof couldn't be verified.
    error ProofVerificationFailed();

    /// @dev Thrown when sequencer tries to append an empty batch
    error EmptyBatch();

    /// @dev Thrown when underflow occurs reading txBatch
    error TxBatchDataUnderflow();

    /// @dev Thrown when overflow occurs reading txBatch
    error TxBatchDataOverflow();

    /// @dev Thrown when a transaction batch has an incorrect version
    error TxBatchVersionIncorrect();

    /**
     * @notice Appends a batch of transactions (stored in calldata) and emits a TxBatchAppended event.
     * @param txBatchData Batch of RLP-encoded transactions, encoded as:
     * txBatchData format:
     *   txBatchData = version || batchData (|| is concatenation)
     *   where:
     *   - version: uint8
     *   - data: bytes
     * batchData format:
     *   batchData = RLP([firstL2BlockNum, batchList])
     *   where:
     *   - firstL2BlockNum: uint256
     *   - batchList: List[blockData]
     *   blockData = [timestamp, txList]
     *   where:
     *   - timestamp: uint256
     *   - txList: bytes
     */
    function appendTxBatch(bytes calldata txBatchData) external;
}
