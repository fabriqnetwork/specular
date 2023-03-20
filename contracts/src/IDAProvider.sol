// SPDX-License-Identifier: Apache-2.0

/*
 * Modifications Copyright 2022, Specular contributors
 */

pragma solidity ^0.8.0;

/**
 * @notice Data availability layer interface to rollup contracts.
 */
interface IDAProvider {
    /**
     * @notice Gets inbox size (total number of messages stored).
     */
    function getInboxSize() external view returns (uint256);
    /**
     * Verifies proof of inclusion of a transaction by the data availability provider.
     * If verification fails, the function reverts.
     * @param encodedTx RLP-encoded transaction.
     * @param proof DA-specific membership proof.
     */
    function verifyTxInclusion(bytes memory encodedTx, bytes calldata proof) external view;
}