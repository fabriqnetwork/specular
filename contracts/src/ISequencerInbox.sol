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

    /// @dev Thrown when the given tx inlcusion proof couldn't be verified.
    error ProofVerificationFailed();

    /// @dev Thrown when sequencer tries to append an empty batch
    error EmptyBatch();

    /// @dev Thrown when overflow occurs reading txBatch (likely due to malformed txLengths)
    error TxBatchDataOverflow();

    /// @dev Thrown when overflow occurs reading txBatch (likely due to malformed txLengths)
    error TxBatchVersionIncorrect();

    /**
     * @notice Appends a batch of transactions (stored in calldata) and emits a TxBatchAppended event.
     * @param txBatchVersion The serialization version of the submitted tx batch
     * @param txBatch Batch of RLP-encoded transactions.
     */
    function appendTxBatch(
        uint256 txBatchVersion,
        bytes calldata txBatch
    ) external;
}
