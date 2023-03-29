// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "../libraries/Types.sol";

interface IL1Portal {
    /**
     * @notice Emitted any time a deposit is initiated.
     *
     * @param nonce    Unique value corresponding to each deposit.
     * @param sender   The L2 account address which initiated the deposit.
     * @param target   The L1 account address the call will be send to.
     * @param value    The ETH value submitted for deposit, to be forwarded to the target.
     * @param gasLimit The minimum amount of gas that must be provided when depositing on L1.
     * @param data     The data to be forwarded to the target on L1.
     * @param depositHash     The hash of the deposit.
     */
    event DepositInitiated(
        uint256 indexed nonce,
        address indexed sender,
        address indexed target,
        uint256 value,
        uint256 gasLimit,
        bytes data,
        bytes32 depositHash
    );

    /**
     * @notice Emitted when a withdrawal transaction is finalized.
     *
     * @param withdrawalHash Hash of the withdrawal transaction.
     * @param success        Whether the withdrawal transaction was successful.
     */
    event WithdrawalFinalized(bytes32 indexed withdrawalHash, bool success);

    /**
     * @notice Sends a message from L1 to L2.
     *
     * @param target   Address to call on L2 execution.
     * @param gasLimit Minimum gas limit for executing the message on L2.
     * @param data     Data to forward to L2 target.
     */
    function initiateDeposit(address target, uint256 gasLimit, bytes memory data) external payable;

    /**
     * @notice Finalizes a withdrawal transaction.
     *
     * @param withdrawalTx           Withdrawal transaction to finalize.
     * @param assertionID            ID of the assertion that can be used to finalize the withdrawal.
     * @param withdrawalAccountProof Inclusion proof of the L2Portal contract's storage root.
     * @param withdrawalProof        Inclusion proof of the withdrawal in L2Portal contract.
     */
    function finalizeWithdrawalTransaction(
        Types.CrossDomainMessage memory withdrawalTx,
        uint256 assertionID,
        // bytes calldata encodedBlockHeader,
        bytes[] calldata withdrawalAccountProof,
        bytes[] calldata withdrawalProof
    ) external;

    /**
     * @notice Determine if a given assertion is finalized and confirmed.
     *
     * @param assertionID The ID of the assertion.
     */
    function isAssertionConfirmed(uint256 assertionID) external view returns (bool);
}
