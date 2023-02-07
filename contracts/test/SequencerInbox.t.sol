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

    /////////////////////////////
    // appendTxBatch
    /////////////////////////////
    function test_appendTxBatch_positiveCase() public {
        // function appendTxBatch(uint256[] calldata contexts, uint256[] calldata txLengths, bytes calldata txBatch)
        // To be able to call this function, we first need an array of `contexts`
        // So, what is `contexts`? => Each context corresponds to a single "L2 block"
        // `contexts` is represented with uint256 3-tuple: (numTxs, l2BlockNumber, l2Timestamp)
        // Ok, so it means number of transactions in block, the l2 block.number and the l2 block.timestamp of that block

        // Let's create an array of contexts
        uint256 numTxns = 3;
        uint256 timeStamp1 = block.timestamp / 10;
        uint256 timeStamp2 = block.timestamp / 5;
        uint256 blockNumber1 = timestamp1 / 20;
        uint256 blockNumber2 = timestamp2 / 20;

        // Let's assume that we had 2 blocks and each had 3 transactions
        uint256[6] memory contexts = [numTxns, blockNumber1, timeStamp1, numTxns, blockNumber2, timeStamp2];

        assertTrue(false);
    }
}

// These are two examples of RLP encoded transactions and they both have the same length... I think all RLP transactions have the same length as a property.
// So, what does `txLengths` array comprise of?

// Also, given different RLP encoded transactions in a batch, how do we calculate the txBatch thing?
// Is it something like, keccak256(abi.encodePacked(tx1, tx2, tx3)) {for all blocks in one}?

// 0xf876801888073bb609d9bc81b894830a495171034ebbd89c236e7a90e24a49d38a819108683d8cf929fd690df519ee3488640000802da071c72753fa2081067f338ca0799e952dc23ac341fe62ddfe3e598279d34931d4a01a9b35b3d886b0621bad6b1c922b997e2496691cb2b55130ec95f1b6a8ad6230
// 0xf8768018880dae933732f365e894a4f21f9296a37bba080b93c2110487a230041ef79105f8dbf082345488ac3526185ce3640000802ea03a56c5c72c0859760e62ae6a1e089f53d0c47a595aa47afd0362b98b300d84b9a062dda585b52188a14b9cd31cae95fd0844cef83711240e4d633045c4ff70aa61
