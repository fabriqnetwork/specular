// SPDX-License-Identifier: MIT
pragma solidity ^0.8.4;

import {UUPSUpgradeable} from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import {OwnableUpgradeable} from "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import {PausableUpgradeable} from "@openzeppelin/contracts-upgradeable/security/PausableUpgradeable.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import {ERC165Checker} from "@openzeppelin/contracts/utils/introspection/ERC165Checker.sol";

import {IStandardBridge} from "./IStandardBridge.sol";
import {SafeCall} from "../libraries/SafeCall.sol";
import {IMintableERC20} from "./mintable/IMintableERC20.sol";

/// @title StandardBridge
/// @notice StandardBridge is a base contract for the L1 and L2 standard ERC20 bridges. It handles
///         the core bridging logic, including escrowing tokens that are native to the local chain
///         and minting/burning tokens that are native to the remote chain.
abstract contract StandardBridge is IStandardBridge, UUPSUpgradeable, OwnableUpgradeable, PausableUpgradeable {
    using SafeERC20 for IERC20;

    /// @notice The L2 gas limit set when eth is depoisited using the receive() function.
    uint32 internal constant RECEIVE_DEFAULT_GAS_LIMIT = 200_000;

    /// @notice Corresponding bridge on the other network.
    StandardBridge public OTHER_BRIDGE;

    /// @notice Portal contract on this network.
    address public PORTAL_ADDRESS;

    /// @notice Mapping that stores deposits for a given pair of local and remote tokens.
    mapping(address => mapping(address => uint256)) public deposits;

    /// @notice Ensures that the caller is a cross-chain message from the other bridge.
    modifier onlyOtherBridge() virtual;

    function __StandardBridge_init(address payable _portalAddress, address payable _otherBridge) internal {
        PORTAL_ADDRESS = _portalAddress;
        OTHER_BRIDGE = StandardBridge(_otherBridge);
        __Ownable_init();
        __Pausable_init();
        __UUPSUpgradeable_init();
    }
    
    function pause() public onlyOwner {
        _pause();
    }

    function unpause() public onlyOwner {
        _unpause();
    }

    /// @notice Allows EOAs to bridge ETH by sending directly to the bridge.
    ///         Must be implemented by contracts that inherit.
    receive() external payable virtual;

    /// @inheritdoc IStandardBridge
    function bridgeETH(uint32 _minGasLimit, bytes calldata _extraData) public payable whenNotPaused {
        _initiateBridgeETH(msg.sender, msg.sender, msg.value, _minGasLimit, _extraData);
    }

    /// @inheritdoc IStandardBridge
    function bridgeETHTo(address _to, uint32 _minGasLimit, bytes calldata _extraData) public payable whenNotPaused {
        _initiateBridgeETH(msg.sender, _to, msg.value, _minGasLimit, _extraData);
    }

    /// @inheritdoc IStandardBridge
    function bridgeERC20(
        address _localToken,
        address _remoteToken,
        uint256 _amount,
        uint32 _minGasLimit,
        bytes calldata _extraData
    ) public virtual whenNotPaused {
        _initiateBridgeERC20(_localToken, _remoteToken, msg.sender, msg.sender, _amount, _minGasLimit, _extraData);
    }

    /// @inheritdoc IStandardBridge
    function bridgeERC20To(
        address _localToken,
        address _remoteToken,
        address _to,
        uint256 _amount,
        uint32 _minGasLimit,
        bytes calldata _extraData
    ) public virtual whenNotPaused {
        _initiateBridgeERC20(_localToken, _remoteToken, msg.sender, _to, _amount, _minGasLimit, _extraData);
    }

    /// @inheritdoc IStandardBridge
    function finalizeBridgeETH(address _from, address _to, uint256 _amount, bytes calldata _extraData)
        external
        payable
        onlyOtherBridge
        whenNotPaused
    {
        require(msg.value == _amount, "StandardBridge: amount sent does not match amount required");
        require(_to != address(this), "StandardBridge: cannot send to self");
        require(_to != PORTAL_ADDRESS, "StandardBridge: cannot send to portal");

        bool success = SafeCall.call(_to, gasleft(), _amount, hex"");
        require(success, "StandardBridge: ETH transfer failed");

        emit ETHBridgeFinalized(_from, _to, _amount, _extraData);
    }

    /// @inheritdoc IStandardBridge
    function finalizeBridgeERC20(
        address _localToken,
        address _remoteToken,
        address _from,
        address _to,
        uint256 _amount,
        bytes calldata _extraData
    ) public onlyOtherBridge whenNotPaused {
        if (_isNonNativeTokenPair(_localToken, _remoteToken)) {
            IMintableERC20(_localToken).mint(_to, _amount);
        } else {
            deposits[_localToken][_remoteToken] = deposits[_localToken][_remoteToken] - _amount;
            IERC20(_localToken).safeTransfer(_to, _amount);
        }

        emit ERC20BridgeFinalized(_localToken, _remoteToken, _from, _to, _amount, _extraData);
    }

    /// @notice Initiates a bridge of ETH through the CrossDomainMessenger.
    /// @param _from        Address of the sender.
    /// @param _to          Address of the receiver.
    /// @param _amount      Amount of ETH being bridged.
    /// @param _minGasLimit Minimum amount of gas that the bridge can be relayed with.
    /// @param _extraData   Extra data to be sent with the transaction. Note that the recipient will
    ///                     not be triggered with this data, but it will be emitted and can be used
    ///                     to identify the transaction.
    function _initiateBridgeETH(
        address _from,
        address _to,
        uint256 _amount,
        uint32 _minGasLimit,
        bytes memory _extraData
    ) internal virtual {}

    /// @notice Sends ERC20 tokens to a receiver's address on the other chain.
    /// @param _localToken  Address of the ERC20 on this chain.
    /// @param _remoteToken Address of the corresponding token on the remote chain.
    /// @param _to          Address of the receiver.
    /// @param _amount      Amount of local tokens to deposit.
    /// @param _minGasLimit Minimum amount of gas that the bridge can be relayed with.
    /// @param _extraData   Extra data to be sent with the transaction. Note that the recipient will
    ///                     not be triggered with this data, but it will be emitted and can be used
    ///                     to identify the transaction.
    function _initiateBridgeERC20(
        address _localToken,
        address _remoteToken,
        address _from,
        address _to,
        uint256 _amount,
        uint32 _minGasLimit,
        bytes memory _extraData
    ) internal virtual {}

    /// @notice Checks if the "other token" is the correct pair token for the MintableERC20.
    /// @param _localToken  MintableERC20 to check against.
    /// @param _remoteToken Pair token to check.
    function _isNonNativeTokenPair(address _localToken, address _remoteToken) internal view returns (bool) {
        if (!ERC165Checker.supportsInterface(_localToken, type(IMintableERC20).interfaceId)) {
            return false;
        }

        if (_remoteToken != IMintableERC20(_localToken).REMOTE_TOKEN()) {
            return false;
        }

        return true;
    }
}
