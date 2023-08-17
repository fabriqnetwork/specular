// SPDX-License-Identifier: MIT
pragma solidity ^0.8.4;

interface IStandardBridge {
    /// @notice Emitted when an ETH bridge is initiated to the other chain.
    /// @param from      Address of the sender.
    /// @param to        Address of the receiver.
    /// @param amount    Amount of ETH sent.
    /// @param extraData Extra data sent with the transaction.
    event ETHBridgeInitiated(address indexed from, address indexed to, uint256 amount, bytes extraData);

    /// @notice Emitted when an ETH bridge is finalized on this chain.
    /// @param from      Address of the sender.
    /// @param to        Address of the receiver.
    /// @param amount    Amount of ETH sent.
    /// @param extraData Extra data sent with the transaction.
    event ETHBridgeFinalized(address indexed from, address indexed to, uint256 amount, bytes extraData);

    /// @notice Emitted when an ERC20 bridge is initiated to the other chain.
    /// @param localToken  Address of the ERC20 on this chain.
    /// @param remoteToken Address of the ERC20 on the remote chain.
    /// @param from        Address of the sender.
    /// @param to          Address of the receiver.
    /// @param amount      Amount of the ERC20 sent.
    /// @param extraData   Extra data sent with the transaction.
    event ERC20BridgeInitiated(
        address indexed localToken,
        address indexed remoteToken,
        address indexed from,
        address to,
        uint256 amount,
        bytes extraData
    );

    /// @notice Emitted when an ERC20 bridge is finalized on this chain.
    /// @param localToken  Address of the ERC20 on this chain.
    /// @param remoteToken Address of the ERC20 on the remote chain.
    /// @param from        Address of the sender.
    /// @param to          Address of the receiver.
    /// @param amount      Amount of the ERC20 sent.
    /// @param extraData   Extra data sent with the transaction.
    event ERC20BridgeFinalized(
        address indexed localToken,
        address indexed remoteToken,
        address indexed from,
        address to,
        uint256 amount,
        bytes extraData
    );

    /// @notice Allows EOAs to bridge ETH by sending directly to the bridge.
    ///         Must be implemented by contracts that inherit.
    receive() external payable;

    /// @notice Sends ETH to the sender's address on the other chain.
    /// @param _minGasLimit Minimum amount of gas that the bridge can be relayed with.
    /// @param _extraData   Extra data to be sent with the transaction. Note that the recipient will
    ///                     not be triggered with this data, but it will be emitted and can be used
    ///                     to identify the transaction.
    function bridgeETH(uint32 _minGasLimit, bytes calldata _extraData) external payable;

    /// @notice Sends ETH to a receiver's address on the other chain. Note that if ETH is sent to a
    ///         smart contract and the call fails, the ETH will be temporarily locked in the
    ///         StandardBridge on the other chain until the call is replayed. If the call cannot be
    ///         replayed with any amount of gas (call always reverts), then the ETH will be
    ///         permanently locked in the StandardBridge on the other chain. ETH will also
    ///         be locked if the receiver is the other bridge, because finalizeBridgeETH will revert
    ///         in that case.
    /// @param _to          Address of the receiver.
    /// @param _minGasLimit Minimum amount of gas that the bridge can be relayed with.
    /// @param _extraData   Extra data to be sent with the transaction. Note that the recipient will
    ///                     not be triggered with this data, but it will be emitted and can be used
    ///                     to identify the transaction.
    function bridgeETHTo(address _to, uint32 _minGasLimit, bytes calldata _extraData) external payable;

    /// @notice Sends ERC20 tokens to the sender's address on the other chain. Note that if the
    ///         ERC20 token on the other chain does not recognize the local token as the correct
    ///         pair token, the ERC20 bridge will fail and the tokens will be returned to sender on
    ///         this chain.
    /// @param _localToken  Address of the ERC20 on this chain.
    /// @param _remoteToken Address of the corresponding token on the remote chain.
    /// @param _amount      Amount of local tokens to deposit.
    /// @param _minGasLimit Minimum amount of gas that the bridge can be relayed with.
    /// @param _extraData   Extra data to be sent with the transaction. Note that the recipient will
    ///                     not be triggered with this data, but it will be emitted and can be used
    ///                     to identify the transaction.
    function bridgeERC20(
        address _localToken,
        address _remoteToken,
        uint256 _amount,
        uint32 _minGasLimit,
        bytes calldata _extraData
    ) external;

    /// @notice Sends ERC20 tokens to a receiver's address on the other chain. Note that if the
    ///         ERC20 token on the other chain does not recognize the local token as the correct
    ///         pair token, the ERC20 bridge will fail and the tokens will be returned to sender on
    ///         this chain.
    /// @param _localToken  Address of the ERC20 on this chain.
    /// @param _remoteToken Address of the corresponding token on the remote chain.
    /// @param _to          Address of the receiver.
    /// @param _amount      Amount of local tokens to deposit.
    /// @param _minGasLimit Minimum amount of gas that the bridge can be relayed with.
    /// @param _extraData   Extra data to be sent with the transaction. Note that the recipient will
    ///                     not be triggered with this data, but it will be emitted and can be used
    ///                     to identify the transaction.
    function bridgeERC20To(
        address _localToken,
        address _remoteToken,
        address _to,
        uint256 _amount,
        uint32 _minGasLimit,
        bytes calldata _extraData
    ) external;

    /// @notice Finalizes an ETH bridge on this chain. Can only be triggered by the other
    ///         StandardBridge contract on the remote chain.
    /// @param _from      Address of the sender.
    /// @param _to        Address of the receiver.
    /// @param _amount    Amount of ETH being bridged.
    /// @param _extraData Extra data to be sent with the transaction. Note that the recipient will
    ///                   not be triggered with this data, but it will be emitted and can be used
    ///                   to identify the transaction.
    function finalizeBridgeETH(address _from, address _to, uint256 _amount, bytes calldata _extraData)
        external
        payable;

    /// @notice Finalizes an ERC20 bridge on this chain. Can only be triggered by the other
    ///         StandardBridge contract on the remote chain.
    /// @param _localToken  Address of the ERC20 on this chain.
    /// @param _remoteToken Address of the corresponding token on the remote chain.
    /// @param _from        Address of the sender.
    /// @param _to          Address of the receiver.
    /// @param _amount      Amount of the ERC20 being bridged.
    /// @param _extraData   Extra data to be sent with the transaction. Note that the recipient will
    ///                     not be triggered with this data, but it will be emitted and can be used
    ///                     to identify the transaction.
    function finalizeBridgeERC20(
        address _localToken,
        address _remoteToken,
        address _from,
        address _to,
        uint256 _amount,
        bytes calldata _extraData
    ) external;
}
