// SPDX-License-Identifier: Apache-2.0

/*
 * Modifications Copyright 2022, Specular contributors
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

pragma solidity ^0.8.13;

import "forge-std/Test.sol";
import "@openzeppelin/contracts/proxy/transparent/TransparentUpgradeableProxy.sol";
import "../src/ISequencerInbox.sol";
import "../src/libraries/Errors.sol";
import {SequencerInbox} from "../src/SequencerInbox.sol";
import {Utils} from "./utils/Utils.sol";

contract SequencerBaseSetup is Test {
    Utils internal utils;
    address payable[] internal users;

    address internal sequencer;
    address internal alice;
    address internal bob;

    function setUp() public virtual {
        sequencer = makeAddr("sequencer");

        alice = makeAddr("alice");

        bob = makeAddr("bob");
    }
}

contract SequencerInboxTest is SequencerBaseSetup {
    SequencerInbox private seqIn;
    SequencerInbox private seqIn2;

    function setUp() public virtual override {
        SequencerBaseSetup.setUp();

        SequencerInbox _impl = new SequencerInbox();
        bytes memory data = abi.encodeWithSelector(SequencerInbox.initialize.selector, sequencer);
        address admin = address(47); // Random admin
        TransparentUpgradeableProxy proxy = new TransparentUpgradeableProxy(address(_impl), admin, data);

        seqIn = SequencerInbox(address(proxy));
    }

    function test_initializeSequencerInbox_withSequencerAddressZero() external {
        SequencerInbox _impl2 = new SequencerInbox();
        bytes memory data = abi.encodeWithSelector(SequencerInbox.initialize.selector, address(0));
        address proxyAdmin = makeAddr("Proxy Admin");

        vm.expectRevert(ZeroAddress.selector);
        TransparentUpgradeableProxy proxy2 = new TransparentUpgradeableProxy(address(_impl2), proxyAdmin, data);

        seqIn2 = SequencerInbox(address(proxy2));
    }

    function test_SequencerAddress() public {
        assertEq(seqIn.sequencerAddress(), sequencer, "Sequencer Address is not as expected");
    }

    function test_intializeSequencerInbox_cannotBeCalledTwice() external {
        vm.expectRevert("Initializable: contract is already initialized");

        // Since seqIn has already been initialised once(in the setUp function), if we try to initialize it again, it should fail
        seqIn.initialize(makeAddr("random address"));
    }

    function testFail_verifyTxInclusion_withEmptyProof() external {
        bytes memory emptyProof = bytes("");
        seqIn.verifyTxInclusion(emptyProof);
    }

    // Only sequencer can append transaction batches to the sequencerInbox
    function test_RevertWhen_InvalidSequencer() public {
        vm.expectRevert(abi.encodeWithSelector(NotSequencer.selector, alice, sequencer));
        vm.prank(alice);
        uint256[] memory contexts = new uint256[](1);
        uint256[] memory txLengths = new uint256[](1);
        seqIn.appendTxBatch(contexts, txLengths, "0x");
    }

    function test_RevertWhen_EmptyBatch() public {
        vm.expectRevert(ISequencerInbox.EmptyBatch.selector);
        vm.prank(sequencer);
        uint256[] memory contexts = new uint256[](1);
        uint256[] memory txLengths = new uint256[](1);
        seqIn.appendTxBatch(contexts, txLengths, "0x");
    }
}
