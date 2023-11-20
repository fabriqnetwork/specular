// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

// Testing utilities
import {Test, stdStorage, StdStorage} from "forge-std/Test.sol";
import {ERC20} from "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import {StandardBridge_Initializer} from "./Common.t.sol";

// Target contract dependencies
import {Types} from "../../src/libraries/Types.sol";
import {Hashing} from "../../src/libraries/Hashing.sol";
import {StandardBridge} from "../../src/bridge/StandardBridge.sol";
import {L1StandardBridge} from "../../src/bridge/L1StandardBridge.sol";
import {L2StandardBridge} from "../../src/bridge/L2StandardBridge.sol";
import {L2Portal} from "../../src/bridge/L2Portal.sol";
import {AddressAliasHelper} from "../../src/vendor/AddressAliasHelper.sol";
import {MintableERC20} from "../../src/bridge/mintable/MintableERC20.sol";

contract L2StandardBridge_Getter_Test is StandardBridge_Initializer {
    /// @dev Test that the accessors return the correct initialized values.
    function test_getters_succeeds() external view {
        assert(l2StandardBridge.PORTAL_ADDRESS() == l2PortalAddress);
        assert(l2StandardBridge.OTHER_BRIDGE() == l1StandardBridge);
    }
}

contract L2StandardBridge_Receive_Test is StandardBridge_Initializer {
    /// @dev Tests receive bridges ETH successfully.
    function test_receive_succeeds() external {
        assertEq(l2PortalAddress.balance, 0);

        vm.expectEmit(true, true, true, true, l2StandardBridgeAddress);
        emit ETHBridgeInitiated(alice, alice, 100, hex"");

        vm.expectCall(
            l2PortalAddress,
            abi.encodeWithSelector(
                l2Portal.initiateWithdrawal.selector,
                l1StandardBridgeAddress,
                200_000,
                abi.encodeWithSelector(StandardBridge.finalizeBridgeETH.selector, alice, alice, 100, "")
            )
        );

        vm.prank(alice, alice);
        (bool success,) = l2StandardBridgeAddress.call{value: 100}(hex"");
        assertEq(success, true);
        assertEq(l2PortalAddress.balance, 100);
    }
}

contract PreBridgeETH is StandardBridge_Initializer {
    /// @dev Asserts the expected calls and events for bridging ETH
    function _preBridgeETH() internal {
        assertEq(l2StandardBridgeAddress.balance, 0);

        uint256 nonce = l2Portal.nonce();
        uint256 value = 600;
        uint256 gasLimit = 60000;
        bytes memory data = hex"dead";

        bytes memory message =
            abi.encodeWithSelector(StandardBridge.finalizeBridgeETH.selector, alice, alice, value, data);

        bytes32 depositHash = Hashing.hashCrossDomainMessage(
            Types.CrossDomainMessage({
                version: 0,
                nonce: nonce,
                sender: l2StandardBridgeAddress,
                target: l1StandardBridgeAddress,
                value: value,
                gasLimit: gasLimit,
                data: message
            })
        );

        vm.expectCall(
            l2StandardBridgeAddress, value, abi.encodeWithSelector(StandardBridge.bridgeETH.selector, gasLimit, data)
        );

        vm.expectCall(
            l2PortalAddress,
            value,
            abi.encodeWithSelector(L2Portal.initiateWithdrawal.selector, l1StandardBridgeAddress, gasLimit, message)
        );

        vm.expectEmit(true, true, true, true, l2StandardBridgeAddress);
        emit ETHBridgeInitiated(alice, alice, value, data);

        vm.expectEmit(true, true, true, true, l2PortalAddress);
        emit WithdrawalInitiated(
            nonce, l2StandardBridgeAddress, l1StandardBridgeAddress, value, gasLimit, message, depositHash
        );

        vm.prank(alice, alice);
    }
}

contract L2StandardBridge_DepositETH_Test is PreBridgeETH {
    /// @dev Tests that depositing ETH succeeds.
    ///      Emits ETHDepositInitiated and ETHBridgeInitiated events.
    function test_depositETH_succeeds() external {
        _preBridgeETH();
        l2StandardBridge.bridgeETH{value: 600}(60000, hex"dead");
        assertEq(l2PortalAddress.balance, 600);
    }
}

contract PreBridgeETHTo is StandardBridge_Initializer {
    /// @dev Asserts the expected calls and events for bridging ETH to a different address
    function _preBridgeETHTo() internal {
        assertEq(l2StandardBridgeAddress.balance, 0);

        uint256 nonce = l2Portal.nonce();
        uint256 value = 600;
        uint256 gasLimit = 60000;
        bytes memory data = hex"dead";

        bytes memory message =
            abi.encodeWithSelector(StandardBridge.finalizeBridgeETH.selector, alice, bob, value, data);

        bytes32 depositHash = Hashing.hashCrossDomainMessage(
            Types.CrossDomainMessage({
                version: 0,
                nonce: nonce,
                sender: l2StandardBridgeAddress,
                target: l1StandardBridgeAddress,
                value: value,
                gasLimit: gasLimit,
                data: message
            })
        );

        vm.expectCall(
            l2StandardBridgeAddress,
            value,
            abi.encodeWithSelector(StandardBridge.bridgeETHTo.selector, bob, gasLimit, data)
        );

        vm.expectCall(
            l2PortalAddress,
            value,
            abi.encodeWithSelector(L2Portal.initiateWithdrawal.selector, l1StandardBridgeAddress, gasLimit, message)
        );

        vm.expectEmit(true, true, true, true, l2StandardBridgeAddress);
        emit ETHBridgeInitiated(alice, bob, value, data);

        vm.expectEmit(true, true, true, true, l2PortalAddress);
        emit WithdrawalInitiated(
            nonce, l2StandardBridgeAddress, l1StandardBridgeAddress, value, gasLimit, message, depositHash
        );

        vm.prank(alice, alice);
    }
}

contract L2StandardBridge_DepositETHTo_Test is PreBridgeETHTo {
    /// @dev Tests that depositing ETH to a different address succeeds.
    ///      Emits ETHDepositInitiated event.
    function test_depositETHTo_succeeds() external {
        _preBridgeETHTo();
        l2StandardBridge.bridgeETHTo{value: 600}(bob, 60000, hex"dead");
        assertEq(address(l2PortalAddress).balance, 600);
    }
}

contract L2StandardBridge_DepositERC20_Test is StandardBridge_Initializer {
    using stdStorage for StdStorage;

    /// @dev Tests that depositing ERC20 to the bridge succeeds.
    ///      Bridge deposits are updated.
    ///      Emits ERC20DepositInitiated event.
    function test_depositERC20_succeeds() external {
        uint256 nonce = l2Portal.nonce();
        uint256 value = 100;
        uint256 gasLimit = 10000;

        // Deal Alice's ERC20 State
        deal(address(l2Token), alice, gasLimit, true);
        vm.prank(alice);
        l2Token.approve(l2StandardBridgeAddress, type(uint256).max);

        uint256 slot = stdstore.target(l2StandardBridgeAddress).sig("deposits(address,address)").with_key(
            address(l2Token)
        ).with_key(address(l1Token)).find();

        // Give the L1 bridge some ERC20 tokens
        vm.store(l2StandardBridgeAddress, bytes32(slot), bytes32(uint256(value)));
        assertEq(l2StandardBridge.deposits(address(l2Token), address(l1Token)), value);

        bytes memory message = abi.encodeWithSelector(
            StandardBridge.finalizeBridgeERC20.selector, address(l1Token), address(l2Token), alice, alice, value, hex""
        );

        bytes32 depositHash = Hashing.hashCrossDomainMessage(
            Types.CrossDomainMessage({
                version: 0,
                nonce: nonce,
                sender: l2StandardBridgeAddress,
                target: l1StandardBridgeAddress,
                value: 0,
                gasLimit: gasLimit,
                data: message
            })
        );

        vm.expectCall(
            l2PortalAddress,
            abi.encodeWithSelector(L2Portal.initiateWithdrawal.selector, l1StandardBridgeAddress, gasLimit, message)
        );

        vm.expectEmit(true, true, true, true, l2StandardBridgeAddress);
        emit ERC20BridgeInitiated(address(l2Token), address(l1Token), alice, alice, value, hex"");

        vm.expectEmit(true, true, true, true, l2PortalAddress);
        emit WithdrawalInitiated(
            nonce, l2StandardBridgeAddress, l1StandardBridgeAddress, 0, gasLimit, message, depositHash
        );

        vm.prank(alice);
        l2StandardBridge.bridgeERC20(address(l2Token), address(l1Token), value, uint32(gasLimit), hex"");
        assertEq(l2StandardBridge.deposits(address(l2Token), address(l1Token)), value);
    }
}

contract L2StandardBridge_DepositERC20To_Test is StandardBridge_Initializer {
    using stdStorage for StdStorage;

    /// @dev Tests that depositing ERC20 to the bridge succeeds when
    ///      sent to a different address.
    ///      Bridge deposits are updated.
    ///      Emits ERC20DepositInitiated event.
    function test_depositERC20_succeeds() external {
        uint256 nonce = l2Portal.nonce();
        uint256 value = 100;
        uint256 gasLimit = 10000;

        // Deal Alice's ERC20 State
        deal(address(l2Token), alice, gasLimit, true);
        vm.prank(alice);
        l2Token.approve(l2StandardBridgeAddress, type(uint256).max);

        uint256 slot = stdstore.target(l2StandardBridgeAddress).sig("deposits(address,address)").with_key(
            address(l2Token)
        ).with_key(address(l1Token)).find();

        // Give the L1 bridge some ERC20 tokens
        vm.store(l2StandardBridgeAddress, bytes32(slot), bytes32(uint256(value)));
        assertEq(l2StandardBridge.deposits(address(l2Token), address(l1Token)), value);

        bytes memory message = abi.encodeWithSelector(
            StandardBridge.finalizeBridgeERC20.selector, address(l1Token), address(l2Token), alice, bob, value, hex""
        );

        bytes32 depositHash = Hashing.hashCrossDomainMessage(
            Types.CrossDomainMessage({
                version: 0,
                nonce: nonce,
                sender: l2StandardBridgeAddress,
                target: l1StandardBridgeAddress,
                value: 0,
                gasLimit: gasLimit,
                data: message
            })
        );

        vm.expectCall(
            l2PortalAddress,
            abi.encodeWithSelector(L2Portal.initiateWithdrawal.selector, l1StandardBridgeAddress, gasLimit, message)
        );

        vm.expectEmit(true, true, true, true, l2StandardBridgeAddress);
        emit ERC20BridgeInitiated(address(l2Token), address(l1Token), alice, bob, value, hex"");

        vm.expectEmit(true, true, true, true, l2PortalAddress);
        emit WithdrawalInitiated(
            nonce, l2StandardBridgeAddress, l1StandardBridgeAddress, 0, gasLimit, message, depositHash
        );

        vm.prank(alice);
        l2StandardBridge.bridgeERC20To(address(l2Token), address(l1Token), bob, value, uint32(gasLimit), hex"");
        assertEq(l2StandardBridge.deposits(address(l2Token), address(l1Token)), value);
    }
}

contract L2StandardBridge_FinalizeETHWithdrawal_Test is StandardBridge_Initializer {
    /// @dev Tests that finalizing an ETH withdrawal succeeds.
    ///      Emits ETHWithdrawalFinalized event.
    ///      Only callable by the L2 bridge.
    function test_finalizeETHWithdrawal_succeeds() external {
        uint256 aliceBalance = alice.balance;

        vm.expectEmit(true, true, true, true, l2StandardBridgeAddress);
        emit ETHBridgeFinalized(alice, alice, 100, hex"");

        vm.expectCall(alice, hex"");

        address l1StandardBridgeAlias = AddressAliasHelper.applyL1ToL2Alias(l1StandardBridgeAddress);

        vm.mockCall(
            l2PortalAddress, abi.encodeWithSelector(l2Portal.l1Sender.selector), abi.encode(l1StandardBridgeAlias)
        );

        vm.deal(l2PortalAddress, 100);
        vm.prank(l2PortalAddress);
        l2StandardBridge.finalizeBridgeETH{value: 100}(alice, alice, 100, hex"");

        assertEq(l2PortalAddress.balance, 0);
        assertEq(aliceBalance + 100, alice.balance);
    }
}

contract L2StandardBridge_FinalizeERC20Withdrawal_Test is StandardBridge_Initializer {
    using stdStorage for StdStorage;

    /// @dev Tests that finalizing an ERC20 withdrawal succeeds.
    ///      Bridge deposits are updated.
    ///      Emits ERC20WithdrawalFinalized event.
    ///      Only callable by the L2 bridge.
    function test_finalizeERC20Withdrawal_succeeds() external {
        assertEq(l2Token.balanceOf(l2StandardBridgeAddress), 0);

        uint256 slot = stdstore.target(l2StandardBridgeAddress).sig("deposits(address,address)").with_key(
            address(l2Token)
        ).with_key(address(l1Token)).find();

        // Give the L2 bridge some ERC20 tokens
        vm.store(l2StandardBridgeAddress, bytes32(slot), bytes32(uint256(100)));
        assertEq(l2StandardBridge.deposits(address(l2Token), address(l1Token)), 100);

        vm.expectEmit(true, true, true, true, l2StandardBridgeAddress);
        emit ERC20BridgeFinalized(address(l2Token), address(l1Token), alice, alice, 100, hex"");

        vm.expectCall(address(l2Token), abi.encodeWithSelector(MintableERC20.mint.selector, alice, 100));

        address l1StandardBridgeAlias = AddressAliasHelper.applyL1ToL2Alias(l1StandardBridgeAddress);

        vm.mockCall(
            l2PortalAddress, abi.encodeWithSelector(l2Portal.l1Sender.selector), abi.encode(l1StandardBridgeAlias)
        );

        vm.prank(l2PortalAddress);
        l2StandardBridge.finalizeBridgeERC20(address(l2Token), address(l1Token), alice, alice, 100, hex"");

        assertEq(l2Token.balanceOf(l2StandardBridgeAddress), 0);
        assertEq(l2Token.balanceOf(address(alice)), 100);
    }
}

contract L2StandardBridge_FinalizeERC20Withdrawal_TestFail is StandardBridge_Initializer {
    /// @dev Tests that finalizing an ERC20 withdrawal reverts if the caller is not the L2 bridge.
    function test_finalizeERC20Withdrawal_notMessenger_reverts() external {
        address l1StandardBridgeAlias = AddressAliasHelper.applyL1ToL2Alias(l1StandardBridgeAddress);

        vm.mockCall(
            l2PortalAddress, abi.encodeWithSelector(l2Portal.l1Sender.selector), abi.encode(l1StandardBridgeAlias)
        );

        vm.prank(address(28));
        vm.expectRevert("StandardBridge: function can only be called from the other bridge");
        l2StandardBridge.finalizeBridgeERC20(address(l2Token), address(l1Token), alice, alice, 100, hex"");
    }

    /// @dev Tests that finalizing an ERC20 withdrawal reverts if the caller is not the L2 bridge.
    function test_finalizeERC20Withdrawal_notOtherBridge_reverts() external {
        vm.mockCall(l2PortalAddress, abi.encodeWithSelector(l2Portal.l1Sender.selector), abi.encode(address(28)));
        vm.prank(l1PortalAddress);

        vm.expectRevert("StandardBridge: function can only be called from the other bridge");
        l2StandardBridge.finalizeBridgeERC20(address(l2Token), address(l1Token), alice, alice, 100, hex"");
    }

    /// @dev Tests that finalizing bridged ETH reverts if the amount is incorrect.
    function test_finalizeBridgeETH_incorrectValue_reverts() external {
        address l1StandardBridgeAlias = AddressAliasHelper.applyL1ToL2Alias(l1StandardBridgeAddress);

        vm.mockCall(
            l2PortalAddress, abi.encodeWithSelector(l2Portal.l1Sender.selector), abi.encode(l1StandardBridgeAlias)
        );

        vm.deal(l2PortalAddress, 100);
        vm.prank(l2PortalAddress);

        vm.expectRevert("StandardBridge: amount sent does not match amount required");
        l2StandardBridge.finalizeBridgeETH{value: 50}(alice, alice, 100, hex"");
    }

    /// @dev Tests that finalizing bridged ETH reverts if the destination is the L1 bridge.
    function test_finalizeBridgeETH_sendToSelf_reverts() external {
        address l1StandardBridgeAlias = AddressAliasHelper.applyL1ToL2Alias(l1StandardBridgeAddress);

        vm.mockCall(
            l2PortalAddress, abi.encodeWithSelector(l2Portal.l1Sender.selector), abi.encode(l1StandardBridgeAlias)
        );

        vm.deal(l2PortalAddress, 100);
        vm.prank(l2PortalAddress);

        vm.expectRevert("StandardBridge: cannot send to self");
        l2StandardBridge.finalizeBridgeETH{value: 100}(alice, l2StandardBridgeAddress, 100, hex"");
    }

    /// @dev Tests that finalizing bridged ETH reverts if the destination is the portal.
    function test_finalizeBridgeETH_sendToPortal_reverts() external {
        address l1StandardBridgeAlias = AddressAliasHelper.applyL1ToL2Alias(l1StandardBridgeAddress);

        vm.mockCall(
            l2PortalAddress, abi.encodeWithSelector(l2Portal.l1Sender.selector), abi.encode(l1StandardBridgeAlias)
        );

        vm.deal(l2PortalAddress, 100);
        vm.prank(l2PortalAddress);

        vm.expectRevert("StandardBridge: cannot send to portal");
        l2StandardBridge.finalizeBridgeETH{value: 100}(alice, l2PortalAddress, 100, hex"");
    }
}
