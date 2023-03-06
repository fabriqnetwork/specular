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
import "../src/ISequencerInbox.sol";
import "../src/libraries/Errors.sol";
import {SpecularProxy} from "test/SpecularProxy.sol";
import {SequencerInbox} from "../src/SequencerInbox.sol";
import {Utils} from "./utils/Utils.sol";
import {RLPEncodedTransactionsUtil} from "./utils/RLPEncodedTransactions.sol";

contract SequencerBaseSetup is Test, RLPEncodedTransactionsUtil {
    address internal alice;
    address internal bob;
    address internal sequencerOwner;
    address internal sequencerAddress;

    function setUp() public virtual {
        alice = makeAddr("alice");
        bob = makeAddr("bob");
        sequencerOwner = makeAddr("sequencerOwner");
        sequencerAddress = makeAddr("sequencerAddress");
    }
}

contract SequencerInboxTest is SequencerBaseSetup {
    /////////////////////////////////
    // SequencerInbox Setup
    /////////////////////////////////

    SequencerInbox public seqIn;
    SpecularProxy public specularProxy;
    SequencerInbox public implementationSequencer;

    function setUp() public virtual override {
        SequencerBaseSetup.setUp();

        bytes memory seqInInitData = abi.encodeWithSignature("initialize(address)", sequencerAddress);

        vm.startPrank(sequencerOwner);

        implementationSequencer = new SequencerInbox();
        specularProxy = new SpecularProxy(address(implementationSequencer), seqInInitData);
        seqIn = SequencerInbox(address(specularProxy));

        vm.stopPrank();
    }

    /////////////////////////////////
    // SequencerInbox Tests
    /////////////////////////////////
    function test_initializeSequencerInbox_withSequencerAddressZero() external {
        bytes memory seqInInitData = abi.encodeWithSignature("initialize(address)", address(0));
        vm.startPrank(sequencerOwner);
        SequencerInbox implementationSequencer2 = new SequencerInbox();
        vm.expectRevert(ZeroAddress.selector);
        specularProxy = new SpecularProxy(address(implementationSequencer), seqInInitData);
        seqIn = SequencerInbox(address(specularProxy));
        vm.stopPrank();
    }

    function test_SequencerAddress() public {
        assertEq(seqIn.sequencerAddress(), sequencerAddress, "Sequencer Address is not as expected");
    }

    function test_intializeSequencerInbox_cannotBeCalledTwice() external {
        vm.expectRevert("Initializable: contract is already initialized");

        // Since seqIn has already been initialised once(in the setUp function), if we try to initialize it again, it should fail
        seqIn.initialize(makeAddr("random address"));
    }

    function testFail_verifyTxInclusion_withEmptyProof() external view {
        bytes memory emptyProof = bytes("");
        seqIn.verifyTxInclusion(emptyProof);
    }

    // Only sequencer can append transaction batches to the sequencerInbox
    function test_RevertWhen_InvalidSequencer() public {
        vm.expectRevert(abi.encodeWithSelector(NotSequencer.selector, alice, sequencerAddress));
        vm.prank(alice);
        uint256[] memory contexts = new uint256[](1);
        uint256[] memory txLengths = new uint256[](1);
        seqIn.appendTxBatch(contexts, txLengths, "0x");
    }

    function test_RevertWhen_EmptyBatch() public {
        vm.expectRevert(ISequencerInbox.EmptyBatch.selector);
        vm.prank(sequencerAddress);
        uint256[] memory contexts = new uint256[](1);
        uint256[] memory txLengths = new uint256[](1);
        seqIn.appendTxBatch(contexts, txLengths, "0x");
    }

    /////////////////////////////
    // appendTxBatch
    /////////////////////////////
    function test_appendTxBatch_positiveCase_1(uint256 numTxnsPerBlock, uint256 txnBlocks) public {
        // We will operate at a limit of transactionsPerBlock = 30 and number of transactionBlocks = 10.
        numTxnsPerBlock = bound(numTxnsPerBlock, 1, 30);
        txnBlocks = bound(txnBlocks, 1, 10);

        uint256 inboxSizeInitial = seqIn.getInboxSize();

        // Each context corresponds to a single "L2 block"
        uint256 numTxns = numTxnsPerBlock * txnBlocks;
        uint256 numContextsArrEntries = 3 * txnBlocks; // Since each `context` is represented with uint256 3-tuple: (numTxs, l2BlockNumber, l2Timestamp)

        // Making sure that the block.timestamp is a reasonable value (> txnBlocks)
        vm.warp(block.timestamp + (4 * txnBlocks));
        uint256 txnBlockTimestamp = block.timestamp - (2 * txnBlocks); // Subtracing just `txnBlocks` would have sufficed. However we are subtracting 2 times txnBlocks for some margin of error.
            // The objective for this subtraction is that while building the `contexts` array, no timestamp should go higher than the current block.timestamp

        // Let's create an array of contexts
        uint256[] memory contexts = new uint256[](numContextsArrEntries);
        for (uint256 i; i < numContextsArrEntries; i += 3) {
            // The first entry for `contexts` for each txnBlock is `numTxns` which we are keeping as constant for all blocks for this test
            contexts[i] = numTxnsPerBlock;

            // Formual Used for blockNumber: (txnBlock's block.timestamp) / 20;
            contexts[i + 1] = txnBlockTimestamp / 20;

            // Formula used for blockTimestamp: (current block.timestamp) / 5x
            contexts[i + 2] = txnBlockTimestamp;

            // The only requirement for timestamps for the transaction blocks is that, these timestamps are monotonically increasing.
            // So, let's increase the value of txnBlock's timestamp monotonically, in a way that is does not exceed current block.timestamp
            ++txnBlockTimestamp;
        }

        // txLengths is defined as: Array of lengths of each encoded tx in txBatch
        // txBatch is defined as: Batch of RLP-encoded transactions
        (bytes memory txBatch, uint256[] memory txLengths) = _helper_sequencerInbox_appendTx(numTxns);

        // Pranking as the sequencer and calling appendTxBatch
        vm.prank(sequencerAddress);
        seqIn.appendTxBatch(contexts, txLengths, txBatch);

        uint256 inboxSizeFinal = seqIn.getInboxSize();
        assertGt(inboxSizeFinal, inboxSizeInitial);

        uint256 expectedInboxSize = numTxns;
        assertEq(inboxSizeFinal, expectedInboxSize);
    }

    function test_appendTxBatch_positiveCase_2_hardcoded() public {
        //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
        // Here, we are assuming we have 2 transaction blocks with 3 transactions each (initial lower load hardcoded test)
        /////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

        uint256 numTxnsPerBlock = 3;
        uint256 inboxSizeInitial = seqIn.getInboxSize();

        // Each context corresponds to a single "L2 block"
        // `contexts` is represented with uint256 3-tuple: (numTxs, l2BlockNumber, l2Timestamp)
        // Let's create an array of contexts
        uint256 numTxns = numTxnsPerBlock * 2;
        uint256 timeStamp1 = block.timestamp / 10;
        uint256 timeStamp2 = block.timestamp / 5;
        uint256 blockNumber1 = timeStamp1 / 20;
        uint256 blockNumber2 = timeStamp2 / 20;

        uint256[] memory contexts = new uint256[](6);

        // Let's assume that we had 2 blocks and each had 3 transactions
        contexts[0] = (numTxnsPerBlock);
        contexts[1] = (blockNumber1);
        contexts[2] = (timeStamp1);
        contexts[3] = (numTxnsPerBlock);
        contexts[4] = (blockNumber2);
        contexts[5] = (timeStamp2);

        // txLengths is defined as: Array of lengths of each encoded tx in txBatch
        // txBatch is defined as: Batch of RLP-encoded transactions
        bytes memory txBatch = _helper_createTxBatch_hardcoded();
        uint256[] memory txLengths = _helper_findTxLength_hardcoded();

        // Pranking as the sequencer and calling appendTxBatch
        vm.prank(sequencerAddress);
        seqIn.appendTxBatch(contexts, txLengths, txBatch);

        uint256 inboxSizeFinal = seqIn.getInboxSize();

        assertGt(inboxSizeFinal, inboxSizeInitial);

        uint256 expectedInboxSize = numTxns;
        assertEq(inboxSizeFinal, expectedInboxSize);
    }

    function test_appendTxBatch_revert_txBatchDataOverflow(uint256 numTxnsPerBlock, uint256 txnBlocks) public {
        // We will operate at a limit of transactionsPerBlock = 30 and number of transactionBlocks = 10.
        numTxnsPerBlock = bound(numTxnsPerBlock, 1, 30);
        txnBlocks = bound(txnBlocks, 1, 10);

        // Each context corresponds to a single "L2 block"
        uint256 numTxns = numTxnsPerBlock * txnBlocks;
        uint256 numContextsArrEntries = 3 * txnBlocks; // Since each `context` is represented with uint256 3-tuple: (numTxs, l2BlockNumber, l2Timestamp)

        // Making sure that the block.timestamp is a reasonable value (> txnBlocks)
        vm.warp(block.timestamp + (4 * txnBlocks));
        uint256 txnBlockTimestamp = block.timestamp - (2 * txnBlocks); // Subtracing just `txnBlocks` would have sufficed. However we are subtracting 2 times txnBlocks for some margin of error.
            // The objective for this subtraction is that while building the `contexts` array, no timestamp should go higher than the current block.timestamp

        // Let's create an array of contexts
        uint256[] memory contexts = new uint256[](numContextsArrEntries);
        for (uint256 i; i < numContextsArrEntries; i += 3) {
            // The first entry for `contexts` for each txnBlock is `numTxns` which we are keeping as constant for all blocks for this test
            contexts[i] = numTxnsPerBlock;

            // Formual Used for blockNumber: (txnBlock's block.timestamp) / 20;
            contexts[i + 1] = txnBlockTimestamp / 20;

            // Formula used for blockTimestamp: (current block.timestamp) / 5x
            contexts[i + 2] = txnBlockTimestamp;

            // The only requirement for timestamps for the transaction blocks is that, these timestamps are monotonically increasing.
            // So, let's increase the value of txnBlock's timestamp monotonically, in a way that is does not exceed current block.timestamp
            ++txnBlockTimestamp;
        }

        // txLengths is defined as: Array of lengths of each encoded tx in txBatch
        // txBatch is defined as: Batch of RLP-encoded transactions
        (bytes memory txBatch, uint256[] memory txLengths) = _helper_sequencerInbox_appendTx(numTxns);

        // Now, we want to trigger the `txnBatchDataOverflow`, so we want to disturn the values receieved in the txLengths array.
        for (uint256 i; i < numTxns; i++) {
            txLengths[i] = txLengths[i] + 1;
        }

        // Pranking as the sequencer and calling appendTxBatch (should throw the TxBatchDataOverflow error)
        vm.expectRevert(ISequencerInbox.TxBatchDataOverflow.selector);
        vm.prank(sequencerAddress);
        seqIn.appendTxBatch(contexts, txLengths, txBatch);
    }

    /////////////////////////////////////////////////////////////////////////////////////////
    // TENTATIVE CODE CHANGE. TEST SUBJECT TO CHANGE BASED ON CHANGE IN CODE
    // Probable change: `appendTxBatch` passes even with malformed `contexts` array
    /////////////////////////////////////////////////////////////////////////////////////////
    function test_appendTxBatch_incompleteDataInContextsArray(uint256 numTxnsPerBlock) public {
        // Since we are assuming that we will have two transaction blocks and we have a total of 300 sample transactions right now.
        numTxnsPerBlock = bound(numTxnsPerBlock, 1, 150);
        uint256 inboxSizeInitial = seqIn.getInboxSize();

        // Each context corresponds to a single "L2 block"
        // `contexts` is represented with uint256 3-tuple: (numTxs, l2BlockNumber, l2Timestamp)
        // Let's create an array of contexts
        uint256 numTxns = numTxnsPerBlock * 2;
        uint256 timeStamp1 = block.timestamp / 10;
        uint256 blockNumber1 = timeStamp1 / 20;

        uint256[] memory contexts = new uint256[](4);

        // Let's assume that we had 2 blocks and each had 3 transactions, but we fail to pass the block.timestamp and block.number of the 2nd transaction block.
        contexts[0] = (numTxnsPerBlock);
        contexts[1] = (blockNumber1);
        contexts[2] = (timeStamp1);
        contexts[3] = (numTxnsPerBlock);

        // txLengths is defined as: Array of lengths of each encoded tx in txBatch
        // txBatch is defined as: Batch of RLP-encoded transactions
        (bytes memory txBatch, uint256[] memory txLengths) = _helper_sequencerInbox_appendTx(numTxns);

        // Pranking as the sequencer and calling appendTxBatch
        vm.prank(sequencerAddress);
        seqIn.appendTxBatch(contexts, txLengths, txBatch);

        uint256 inboxSizeFinal = seqIn.getInboxSize();

        assertGt(inboxSizeFinal, inboxSizeInitial);
        assertEq(inboxSizeFinal, numTxnsPerBlock); // Since the timestamp and block.number were not included for the 2nd block, only 1st block's 3 txns are included.
    }
}
