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

import "../src/libraries/Errors.sol";
import {SequencerInbox} from "../src/SequencerInbox.sol";
import {Utils} from "./utils/Utils.sol";

contract BaseSetup is Test {
    Utils internal utils;
    address payable[] internal users;

    address internal sequencer;
    address internal alice;
    address internal bob;

    function setUp() public virtual {
        utils = new Utils();
        users = utils.createUsers(3);

        sequencer = users[0];
        vm.label(sequencer, "Sequencer");

        alice = users[1];
        vm.label(alice, "Alice");

        bob = users[2];
        vm.label(bob, "Bob");
    }
}

contract SequencerInboxTest is BaseSetup {
    SequencerInbox private seqIn;

    function setUp() public virtual override {
        BaseSetup.setUp();

        SequencerInbox _impl = new SequencerInbox();
        bytes memory data = abi.encodeWithSelector(SequencerInbox.initialize.selector, sequencer);
        address admin = address(47); // Random admin
        TransparentUpgradeableProxy proxy = new TransparentUpgradeableProxy(address(_impl), admin, data);

        seqIn = SequencerInbox(address(proxy));
    }

    function test_SequencerAddress() public {
        assertEq(seqIn.sequencerAddress(), sequencer, "Sequencer Address is not as expected");
    }

    function test_RevertWhen_InvalidSequencer() public {
        vm.expectRevert(abi.encodeWithSelector(NotSequencer.selector, alice, sequencer));
        vm.prank(alice);
        uint256[] memory contexts = new uint256[](1);
        uint256[] memory txLengths = new uint256[](1);
        seqIn.appendTxBatch(contexts, txLengths, "0x");
    }

    function test_RevertWhen_EmptyBatch() public {
        vm.expectRevert(EmptyBatch.selector);
        vm.prank(sequencer);
        uint256[] memory contexts = new uint256[](1);
        uint256[] memory txLengths = new uint256[](1);
        seqIn.appendTxBatch(contexts, txLengths, "0x");
    }
}

/**
    To-Dos:
    1. Write more tests (optimize for well thought-out tests rather than code coverage)
    2. Work on Test Harnesses (Harness contracts inherit from the Contracts under Test and expose the internal functions as external ones.)
        This is useful in testing internal functions.
    3. Workaround functions which expose functionality or information otherwise unavailable in the original smart contract.
    4. Make sure we have a healthy mix of Integration tests, Fork Tests, positive & negative unit tests for each code path
    5. Ensure that no tainted data (malicious user input) ever reaches a sink (contract logic where something important is happening)
    6. Write tests for deployment scripts (assuming they will be written in Solidity)
    7. Inspect potential cases for front running (not much relevant in our case)
 */
