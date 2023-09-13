// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

// Testing utilities
import {Test, stdStorage, StdStorage} from "forge-std/Test.sol";
import {ERC20} from "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import {StandardBridge_Initializer} from "./Common.t.sol";

// Target contract dependencies
import {Types} from "../../src/libraries/Types.sol";
import {StandardBridge} from "../../src/bridge/StandardBridge.sol";
import {L1StandardBridge} from "../../src/bridge/L1StandardBridge.sol";
import {L2StandardBridge} from "../../src/bridge/L2StandardBridge.sol";
import {L1Portal} from "../../src/bridge/L1Portal.sol";
import {AddressAliasHelper} from "../../src/vendor/AddressAliasHelper.sol";

contract L1StandardBridge_Getter_Test is StandardBridge_Initializer {
    /// @dev Test that the accessors return the correct initialized values.
    function test_getters_succeeds() external view {
        assert(l1StandardBridge.L1_PORTAL() == l1Portal);
        assert(l1StandardBridge.PORTAL_ADDRESS() == l1PortalAddress);
        assert(l1StandardBridge.OTHER_BRIDGE() == l2StandardBridge);
    }
}

contract L1StandardBridge_Receive_Test is StandardBridge_Initializer {
    /// @dev Tests receive bridges ETH successfully.
    function test_receive_succeeds() external {
        assertEq(l1PortalAddress.balance, 0);

        vm.expectEmit(true, true, true, true, l1StandardBridgeAddress);
        emit ETHBridgeInitiated(alice, alice, 100, hex"");

        vm.expectCall(
            l1PortalAddress,
            abi.encodeWithSelector(
                l1Portal.initiateDeposit.selector,
                l2StandardBridgeAddress,
                200_000,
                abi.encodeWithSelector(StandardBridge.finalizeBridgeETH.selector, alice, alice, 100, "")
            )
        );

        vm.prank(alice, alice);
        (bool success,) = l1StandardBridgeAddress.call{value: 100}(hex"");
        assertEq(success, true);
        assertEq(l1PortalAddress.balance, 100);
    }
}

contract PreBridgeETH is StandardBridge_Initializer {
    /// @dev Asserts the expected calls and events for bridging ETH
    function _preBridgeETH() internal {
        assertEq(l1StandardBridgeAddress.balance, 0);

        bytes memory message =
            abi.encodeWithSelector(StandardBridge.finalizeBridgeETH.selector, alice, alice, 500, hex"dead");

        vm.expectCall(
            l1StandardBridgeAddress, 500, abi.encodeWithSelector(StandardBridge.bridgeETH.selector, 50000, hex"dead")
        );

        vm.expectCall(
            l1PortalAddress,
            500,
            abi.encodeWithSelector(L1Portal.initiateDeposit.selector, l2StandardBridgeAddress, 50000, message)
        );

        vm.expectEmit(true, true, true, true, l1StandardBridgeAddress);
        emit ETHBridgeInitiated(alice, alice, 500, hex"dead");

        uint256 nonce = l1Portal.nonce();
        emit DepositInitiated(nonce, alice, alice, 500, 5000, hex"dead", "");

        vm.prank(alice, alice);
    }
}

contract L1StandardBridge_DepositETH_Test is PreBridgeETH {
    /// @dev Tests that depositing ETH succeeds.
    ///      Emits ETHDepositInitiated and ETHBridgeInitiated events.
    function test_depositETH_succeeds() external {
        _preBridgeETH();
        l1StandardBridge.bridgeETH{value: 500}(50000, hex"dead");
        assertEq(l1PortalAddress.balance, 500);
    }
}

contract PreBridgeETHTo is StandardBridge_Initializer {
    /// @dev Asserts the expected calls and events for bridging ETH to a different address
    function _preBridgeETHTo() internal {
        assertEq(l1StandardBridgeAddress.balance, 0);

        bytes memory message =
            abi.encodeWithSelector(StandardBridge.finalizeBridgeETH.selector, alice, bob, 600, hex"dead");

        vm.expectCall(
            l1StandardBridgeAddress,
            600,
            abi.encodeWithSelector(StandardBridge.bridgeETHTo.selector, bob, 60000, hex"dead")
        );

        vm.expectCall(
            l1PortalAddress,
            600,
            abi.encodeWithSelector(L1Portal.initiateDeposit.selector, l2StandardBridgeAddress, 60000, message)
        );

        vm.expectEmit(true, true, true, true, l1StandardBridgeAddress);
        emit ETHBridgeInitiated(alice, bob, 600, hex"dead");

        uint256 nonce = l1Portal.nonce();
        vm.expectEmit(true, true, true, true, l1PortalAddress);
        emit DepositInitiated(
            nonce,
            AddressAliasHelper.applyL1ToL2Alias(l1StandardBridgeAddress),
            l2StandardBridgeAddress,
            600,
            60000,
            message,
            hex"efdfee2f5058687ebc6d0c1f711885c5cf1f24aa32933d67b97e8c49bccabd00"
        );

        vm.prank(alice, alice);
    }
}

contract L1StandardBridge_DepositETHTo_Test is PreBridgeETHTo {
    /// @dev Tests that depositing ETH to a different address succeeds.
    ///      Emits ETHDepositInitiated event.
    function test_depositETHTo_succeeds() external {
        _preBridgeETHTo();
        l1StandardBridge.bridgeETHTo{value: 600}(bob, 60000, hex"dead");
        assertEq(address(l1PortalAddress).balance, 600);
    }
}

contract L1StandardBridge_DepositERC20_Test is StandardBridge_Initializer {
    using stdStorage for StdStorage;

    /// @dev Tests that depositing ERC20 to the bridge succeeds.
    ///      Bridge deposits are updated.
    ///      Emits ERC20DepositInitiated event.
    function test_depositERC20_succeeds() external {
        uint256 nonce = l1Portal.nonce();

        // Deal Alice's ERC20 State
        deal(address(l1Token), alice, 100000, true);
        vm.prank(alice);
        l1Token.approve(l1StandardBridgeAddress, type(uint256).max);

        // The L1Bridge should transfer alice's tokens to itself
        vm.expectCall(
            address(l1Token), abi.encodeWithSelector(ERC20.transferFrom.selector, alice, l1StandardBridgeAddress, 100)
        );

        bytes memory message = abi.encodeWithSelector(
            StandardBridge.finalizeBridgeERC20.selector, address(l2Token), address(l1Token), alice, alice, 100, hex""
        );

        // the L1 bridge should call L1Portal.initiateDeposit
        vm.expectCall(
            l1PortalAddress,
            abi.encodeWithSelector(L1Portal.initiateDeposit.selector, l2StandardBridgeAddress, 10000, message)
        );

        vm.expectEmit(true, true, true, true, l1StandardBridgeAddress);
        emit ERC20BridgeInitiated(address(l1Token), address(l2Token), alice, alice, 100, hex"");
        vm.expectEmit(true, true, true, true, l1PortalAddress);
        emit DepositInitiated(
            nonce,
            AddressAliasHelper.applyL1ToL2Alias(l1StandardBridgeAddress),
            l2StandardBridgeAddress,
            0,
            10000,
            message,
            hex"eff2698675bc67cd2fec9b070280d12442b0156ae58f133f242b550c81b0c946"
        );

        vm.prank(alice);
        l1StandardBridge.bridgeERC20(address(l1Token), address(l2Token), 100, 10000, hex"");
        assertEq(l1StandardBridge.deposits(address(l1Token), address(l2Token)), 100);
    }
}

contract L1StandardBridge_DepositERC20To_Test is StandardBridge_Initializer {
    /// @dev Tests that depositing ERC20 to the bridge succeeds when
    ///      sent to a different address.
    ///      Bridge deposits are updated.
    ///      Emits ERC20DepositInitiated event.
    function test_depositERC20To_succeeds() external {
        uint256 nonce = l1Portal.nonce();

        // Deal Alice's ERC20 State
        deal(address(l1Token), alice, 100000, true);
        vm.prank(alice);
        l1Token.approve(l1StandardBridgeAddress, type(uint256).max);

        // The L1Bridge should transfer alice's tokens to itself
        vm.expectCall(
            address(l1Token), abi.encodeWithSelector(ERC20.transferFrom.selector, alice, l1StandardBridgeAddress, 100)
        );

        bytes memory message = abi.encodeWithSelector(
            StandardBridge.finalizeBridgeERC20.selector, address(l2Token), address(l1Token), alice, bob, 100, hex""
        );

        // the L1 bridge should call L1Portal.initiateDeposit
        vm.expectCall(
            l1PortalAddress,
            abi.encodeWithSelector(L1Portal.initiateDeposit.selector, l2StandardBridgeAddress, 10000, message)
        );

        vm.expectEmit(true, true, true, true, l1StandardBridgeAddress);
        emit ERC20BridgeInitiated(address(l1Token), address(l2Token), alice, bob, 100, hex"");
        vm.expectEmit(true, true, true, true, l1PortalAddress);
        emit DepositInitiated(
            nonce,
            AddressAliasHelper.applyL1ToL2Alias(l1StandardBridgeAddress),
            l2StandardBridgeAddress,
            0,
            10000,
            message,
            hex"e4b2e5e4417057999367136a09bb17ce400dd06984e78bb926776419fd15cfb3"
        );

        vm.prank(alice);
        l1StandardBridge.bridgeERC20To(address(l1Token), address(l2Token), bob, 100, 10000, hex"");
        assertEq(l1StandardBridge.deposits(address(l1Token), address(l2Token)), 100);
    }
}

contract L1StandardBridge_FinalizeETHWithdrawal_Test is StandardBridge_Initializer {
    /// @dev Tests that finalizing an ETH withdrawal succeeds.
    ///      Emits ETHWithdrawalFinalized event.
    ///      Only callable by the L2 bridge.
    function test_finalizeETHWithdrawal_succeeds() external {
        uint256 aliceBalance = alice.balance;

        vm.expectEmit(true, true, true, true, l1StandardBridgeAddress);
        emit ETHBridgeFinalized(alice, alice, 100, hex"");

        vm.expectCall(alice, hex"");

        //address l2StandardBridgeAlias = AddressAliasHelper.applyL1ToL2Alias(l2StandardBridgeAddress);
        vm.mockCall(
            l1PortalAddress, abi.encodeWithSelector(l1Portal.l2Sender.selector), abi.encode(l2StandardBridgeAddress)
        );

        vm.deal(l1PortalAddress, 100);
        vm.prank(l1PortalAddress);
        l1StandardBridge.finalizeBridgeETH{value: 100}(alice, alice, 100, hex"");

        assertEq(l1PortalAddress.balance, 0);
        assertEq(aliceBalance + 100, alice.balance);
    }
}

contract L1StandardBridge_FinalizeERC20Withdrawal_Test is StandardBridge_Initializer {
    using stdStorage for StdStorage;

    /// @dev Tests that finalizing an ERC20 withdrawal succeeds.
    ///      Bridge deposits are updated.
    ///      Emits ERC20WithdrawalFinalized event.
    ///      Only callable by the L2 bridge.
    function test_finalizeERC20Withdrawal_succeeds() external {
        deal(address(l1Token), l1StandardBridgeAddress, 100, true);

        uint256 slot = stdstore.target(l1StandardBridgeAddress).sig("deposits(address,address)").with_key(
            address(l1Token)
        ).with_key(address(l2Token)).find();

        // Give the L1 bridge some ERC20 tokens
        vm.store(l1StandardBridgeAddress, bytes32(slot), bytes32(uint256(100)));
        assertEq(l1StandardBridge.deposits(address(l1Token), address(l2Token)), 100);

        vm.expectEmit(true, true, true, true, l1StandardBridgeAddress);
        emit ERC20BridgeFinalized(address(l1Token), address(l2Token), alice, alice, 100, hex"");

        vm.expectCall(address(l1Token), abi.encodeWithSelector(ERC20.transfer.selector, alice, 100));

        vm.mockCall(
            l1PortalAddress, abi.encodeWithSelector(l1Portal.l2Sender.selector), abi.encode(l2StandardBridgeAddress)
        );

        vm.prank(l1PortalAddress);
        l1StandardBridge.finalizeBridgeERC20(address(l1Token), address(l2Token), alice, alice, 100, hex"");

        assertEq(l1Token.balanceOf(l1StandardBridgeAddress), 0);
        assertEq(l1Token.balanceOf(address(alice)), 100);
    }
}

contract L1StandardBridge_FinalizeERC20Withdrawal_TestFail is StandardBridge_Initializer {
    /// @dev Tests that finalizing an ERC20 withdrawal reverts if the caller is not the L2 bridge.
    function test_finalizeERC20Withdrawal_notMessenger_reverts() external {
        vm.mockCall(
            l1PortalAddress, abi.encodeWithSelector(l1Portal.l2Sender.selector), abi.encode(l2StandardBridgeAddress)
        );
        vm.prank(address(28));
        vm.expectRevert("StandardBridge: function can only be called from the other bridge");
        l1StandardBridge.finalizeBridgeERC20(address(l1Token), address(l2Token), alice, alice, 100, hex"");
    }

    /// @dev Tests that finalizing an ERC20 withdrawal reverts if the caller is not the L2 bridge.
    function test_finalizeERC20Withdrawal_notOtherBridge_reverts() external {
        vm.mockCall(l1PortalAddress, abi.encodeWithSelector(l1Portal.l2Sender.selector), abi.encode(address(28)));
        vm.prank(l1PortalAddress);

        vm.expectRevert("StandardBridge: function can only be called from the other bridge");
        l1StandardBridge.finalizeBridgeERC20(address(l1Token), address(l2Token), alice, alice, 100, hex"");
    }

    /// @dev Tests that finalizing bridged ETH reverts if the amount is incorrect.
    function test_finalizeBridgeETH_incorrectValue_reverts() external {
        vm.mockCall(
            l1PortalAddress, abi.encodeWithSelector(l1Portal.l2Sender.selector), abi.encode(l2StandardBridgeAddress)
        );

        vm.deal(l1PortalAddress, 100);
        vm.prank(l1PortalAddress);

        vm.expectRevert("StandardBridge: amount sent does not match amount required");
        l1StandardBridge.finalizeBridgeETH{value: 50}(alice, alice, 100, hex"");
    }

    /// @dev Tests that finalizing bridged ETH reverts if the destination is the L1 bridge.
    function test_finalizeBridgeETH_sendToSelf_reverts() external {
        vm.mockCall(
            l1PortalAddress, abi.encodeWithSelector(l1Portal.l2Sender.selector), abi.encode(l2StandardBridgeAddress)
        );

        vm.deal(l1PortalAddress, 100);
        vm.prank(l1PortalAddress);

        vm.expectRevert("StandardBridge: cannot send to self");
        l1StandardBridge.finalizeBridgeETH{value: 100}(alice, l1StandardBridgeAddress, 100, hex"");
    }

    /// @dev Tests that finalizing bridged ETH reverts if the destination is the portal.
    function test_finalizeBridgeETH_sendToPortal_reverts() external {
        vm.mockCall(
            l1PortalAddress, abi.encodeWithSelector(l1Portal.l2Sender.selector), abi.encode(l2StandardBridgeAddress)
        );

        vm.deal(l1PortalAddress, 100);
        vm.prank(l1PortalAddress);

        vm.expectRevert("StandardBridge: cannot send to portal");
        l1StandardBridge.finalizeBridgeETH{value: 100}(alice, l1PortalAddress, 100, hex"");
    }
}
