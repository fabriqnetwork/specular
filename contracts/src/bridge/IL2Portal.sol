// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "../libraries/Types.sol";

interface IL2Portal {
    /**
     * @notice Emitted any time a withdrawal is initiated.
     *
     * @param nonce    Unique value corresponding to each withdrawal.
     * @param sender   The L2 account address which initiated the withdrawal.
     * @param target   The L1 account address the call will be send to.
     * @param value    The ETH value submitted for withdrawal, to be forwarded to the target.
     * @param gasLimit The minimum amount of gas that must be provided when withdrawing on L1.
     * @param data     The data to be forwarded to the target on L1.
     * @param withdrawalHash     The hash of the withdrawal.
     */
    event WithdrawalInitiated(
        uint256 indexed nonce,
        address indexed sender,
        address indexed target,
        uint256 value,
        uint256 gasLimit,
        bytes data,
        bytes32 withdrawalHash
    );

    /**
     * @notice Emitted when a deposit transaction is finalized.
     *
     * @param depositHash Hash of the deposit transaction.
     * @param success     Whether the deposit transaction was successful.
     */
    event DepositFinalized(bytes32 indexed depositHash, bool success);

    /**
     * @notice Sends a message from L2 to L1.
     *
     * @param _target   Address to call on L1 execution.
     * @param _gasLimit Minimum gas limit for executing the message on L1.
     * @param _data     Data to forward to L1 target.
     */
    function initiateWithdrawal(address _target, uint256 _gasLimit, bytes memory _data) external payable;

    /**
     * @notice Finalizes a deposit transaction.
     *
     * @param depositTx           Deposit transaction to finalize.
     * @param l1BlockNumber       L1 Block Number of the deposit.
     * @param depositAccountProof Inclusion proof of the L1Portal contract's storage root.
     * @param depositProof        Inclusion proof of the deposit in L1Portal contract.
     */
    function finalizeDepositTransaction(
        Types.CrossDomainMessage memory depositTx,
        uint256 l1BlockNumber,
        bytes[] calldata depositAccountProof,
        bytes[] calldata depositProof
    ) external;
}
