// SPDX-License-Identifier: MIT
pragma solidity ^0.8.4;

import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";

import {StandardBridge} from "./StandardBridge.sol";
import {L2Portal} from "./L2Portal.sol";
import {IMintableERC20} from "./mintable/IMintableERC20.sol";
import {Predeploys} from "../libraries/Predeploys.sol";
import {AddressAliasHelper} from "../vendor/AddressAliasHelper.sol";

contract L2StandardBridge is StandardBridge {
    using SafeERC20 for IERC20;

    L2Portal internal constant l2Portal = L2Portal(payable(Predeploys.L2_PORTAL));

    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }

    /// @inheritdoc StandardBridge
    modifier onlyOtherBridge() override {
        address origSender = AddressAliasHelper.undoL1ToL2Alias(l2Portal.l1Sender());
        require(
            msg.sender == address(PORTAL_ADDRESS) && origSender == address(OTHER_BRIDGE),
            "StandardBridge: function can only be called from the other bridge"
        );

        _;
    }

    /// @notice Initializer;
    function initialize(address _owner, address payable _otherBridge) public initializer {
        __StandardBridge_init(payable(Predeploys.L2_PORTAL), _otherBridge);

        _transferOwnership(_owner);
    }

    function _authorizeUpgrade(address) internal override onlyOwner whenPaused {}

    /// @inheritdoc StandardBridge
    receive() external payable override whenNotPaused {
        _initiateBridgeETH(msg.sender, msg.sender, msg.value, RECEIVE_DEFAULT_GAS_LIMIT, bytes(""));
    }

    /// @inheritdoc StandardBridge
    function _initiateBridgeETH(
        address _from,
        address _to,
        uint256 _amount,
        uint32 _minGasLimit,
        bytes memory _extraData
    ) internal override {
        emit ETHBridgeInitiated(_from, _to, _amount, _extraData);

        l2Portal.initiateWithdrawal{value: _amount}(
            address(OTHER_BRIDGE),
            _minGasLimit,
            abi.encodeWithSelector(this.finalizeBridgeETH.selector, _from, _to, _amount, _extraData)
        );
    }

    /// @inheritdoc StandardBridge
    function _initiateBridgeERC20(
        address _localToken,
        address _remoteToken,
        address _from,
        address _to,
        uint256 _amount,
        uint32 _minGasLimit,
        bytes memory _extraData
    ) internal override {
        if (_isNonNativeTokenPair(_localToken, _remoteToken)) {
            IMintableERC20(_localToken).burn(_from, _amount);
        } else {
            IERC20(_localToken).safeTransferFrom(_from, address(this), _amount);
            deposits[_localToken][_remoteToken] = deposits[_localToken][_remoteToken] + _amount;
        }

        emit ERC20BridgeInitiated(_localToken, _remoteToken, _from, _to, _amount, _extraData);

        l2Portal.initiateWithdrawal(
            address(OTHER_BRIDGE),
            _minGasLimit,
            abi.encodeWithSelector(
                this.finalizeBridgeERC20.selector,
                // Because this call will be executed on the remote chain,
                // we reverse the order of
                // the remote and local token addresses relative to their order in the
                // finalizeBridgeERC20 function.
                _remoteToken,
                _localToken,
                _from,
                _to,
                _amount,
                _extraData
            )
        );
    }
}
