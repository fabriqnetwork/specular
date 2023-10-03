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

import "forge-std/console.sol";
import "forge-std/Test.sol";
import "../src/ISequencerInbox.sol";
import "../src/libraries/Errors.sol";
import "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";
import {PausableUpgradeable} from "@openzeppelin/contracts-upgradeable/security/PausableUpgradeable.sol";
import {SequencerInbox} from "../src/SequencerInbox.sol";
import {Utils} from "./utils/Utils.sol";
import {RLPEncodedTransactionsUtil} from "./utils/RLPEncodedTransactions.sol";

contract SequencerBaseSetup is Test, RLPEncodedTransactionsUtil {
    address internal alice;
    address internal bob;
    address internal sequencerOwner;
    address internal sequencerAddress;

    function setUp() public virtual {
        sequencerOwner = makeAddr("sequencerOwner");
        sequencerAddress = makeAddr("sequencerAddress");
    }
}

contract SequencerInboxTest is SequencerBaseSetup {
    /////////////////////////////////
    // SequencerInbox Setup
    /////////////////////////////////

    SequencerInbox public seqIn;
    SequencerInbox public implementationSequencer;

    function setUp() public virtual override {
        SequencerBaseSetup.setUp();

        bytes memory seqInInitData = abi.encodeWithSignature("initialize(address)", sequencerAddress);

        vm.startPrank(sequencerOwner);

        implementationSequencer = new SequencerInbox();
        seqIn = SequencerInbox(address(new ERC1967Proxy(address(implementationSequencer), seqInInitData)));
        vm.stopPrank();
    }

    /////////////////////////////////
    // SequencerInbox Tests
    /////////////////////////////////
    function test_initializeSequencerInbox_withSequencerAddressZero_reverts() external {
        bytes memory seqInInitData = abi.encodeWithSignature("initialize(address)", address(0));
        vm.startPrank(sequencerOwner);
        vm.expectRevert(ZeroAddress.selector);
        seqIn = SequencerInbox(address(new ERC1967Proxy(address(implementationSequencer), seqInInitData)));
    }

    function test_sequencerAddress_works() public {
        assertEq(seqIn.sequencerAddress(), sequencerAddress, "Sequencer Address is not as expected");
    }

    function test_intializeSequencerInbox_cannotBeCalledTwice_reverts() external {
        vm.expectRevert("Initializable: contract is already initialized");

        // Since seqIn has already been initialised once(in the setUp function), if we try to initialize it again, it should fail
        seqIn.initialize(makeAddr("random address"));
    }

    function test_verifyTxInclusion_withEmptyProof_fails() external {
        vm.expectRevert("too short");
        bytes memory emptyProof = bytes("");
        seqIn.verifyTxInclusion(bytes(""), emptyProof);
    }

    // Only sequencer can append transaction batches to the sequencerInbox
    function test_appendTxBatch_invalidSequencer_reverts() public {
        vm.expectRevert(abi.encodeWithSelector(NotSequencer.selector, alice, sequencerAddress));
        vm.prank(alice);
        uint256[] memory contexts = new uint256[](1);
        uint256[] memory txLengths = new uint256[](1);
        uint256 txBatchVersion = _helper_sequencerInbox_appendTx_Version();
        seqIn.appendTxBatch(contexts, txLengths, 1, txBatchVersion, "0x");
    }

    function test_appendTxBatch_emptyBatch_reverts() public {
        vm.expectRevert(ISequencerInbox.EmptyBatch.selector);
        vm.prank(sequencerAddress);
        uint256[] memory contexts = new uint256[](1);
        uint256[] memory txLengths = new uint256[](1);
        uint256 txBatchVersion = _helper_sequencerInbox_appendTx_Version();
        seqIn.appendTxBatch(contexts, txLengths, 1, txBatchVersion, "0x");
    }

    //////////////////////////////
    // appendTxBatch
    //////////////////////////////
    function test_appendTxBatch_case1_succeeds(uint256 numTxnsPerBlock, uint256 txnBlocks) public {
        // We will operate at a limit of transactionsPerBlock = 30 and number of transactionBlocks = 10.
        numTxnsPerBlock = bound(numTxnsPerBlock, 1, 30);
        txnBlocks = bound(txnBlocks, 1, 10);

        uint256 inboxSizeInitial = seqIn.getInboxSize();

        // Each context corresponds to a single "L2 block"
        uint256 numTxns = numTxnsPerBlock * txnBlocks;
        // Each `context` is represented with uint256 2-tuple: (numTxs, l2Timestamp)
        uint256 numContextsArrEntries = 2 * txnBlocks;

        // Making sure that the block.timestamp is a reasonable value (> txnBlocks)
        vm.warp(block.timestamp + (4 * txnBlocks));
        uint256 txnBlockTimestamp = block.timestamp - (2 * txnBlocks); // Subtracing just `txnBlocks` would have sufficed. However we are subtracting 2 times txnBlocks for some margin of error.
        // The objective for this subtraction is that while building the `contexts` array, no timestamp should go higher than the current block.timestamp

        uint256 firstL2BlockNumber = block.timestamp / 20;

        // Let's create an array of contexts
        uint256[] memory contexts = new uint256[](numContextsArrEntries);
        for (uint256 i = 0; i < numContextsArrEntries; i += 2) {
            // The first entry for `contexts` for each txnBlock is `numTxns` which we are keeping as constant for all blocks for this test
            contexts[i] = numTxnsPerBlock;

            // Formula used for blockTimestamp: (current block.timestamp) / 5x
            contexts[i + 1] = txnBlockTimestamp;

            // The only requirement for timestamps for the transaction blocks is that, these timestamps are monotonically increasing.
            // So, let's increase the value of txnBlock's timestamp monotonically, in a way that is does not exceed current block.timestamp
            ++txnBlockTimestamp;
        }

        // txLengths is defined as: Array of lengths of each encoded tx in txBatch
        // txBatch is defined as: Batch of RLP-encoded transactions
        uint256 txBatchVersion = _helper_sequencerInbox_appendTx_Version();
        (bytes memory txBatch, uint256[] memory txLengths) = _helper_sequencerInbox_appendTx(numTxns);

        // Pranking as the sequencer and calling appendTxBatch
        vm.prank(sequencerAddress);
        seqIn.appendTxBatch(contexts, txLengths, firstL2BlockNumber, txBatchVersion, txBatch);

        uint256 inboxSizeFinal = seqIn.getInboxSize();
        assertGt(inboxSizeFinal, inboxSizeInitial);

        uint256 expectedInboxSize = numTxns;
        assertEq(inboxSizeFinal, expectedInboxSize);
    }

    function test_appendTxBatch_case2Hardcoded_succeeds() public {
        //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
        // Here, we are assuming we have 2 transaction blocks with 3 transactions each (initial lower load hardcoded test)
        /////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

        uint256 numTxnsPerBlock = 3;
        uint256 inboxSizeInitial = seqIn.getInboxSize();

        uint256 firstL2BlockNumber = block.timestamp / 20;

        // Each context corresponds to a single "L2 block"
        // `contexts` is represented with uint256 2-tuple: (numTxs, l2Timestamp)
        // Let's create an array of contexts
        uint256 numTxns = numTxnsPerBlock * 2;
        uint256 timeStamp1 = block.timestamp / 10;
        uint256 timeStamp2 = block.timestamp / 5;

        uint256[] memory contexts = new uint256[](6);

        // Let's assume that we had 2 blocks and each had 3 transactions
        contexts[0] = (numTxnsPerBlock);
        contexts[1] = (timeStamp1);
        contexts[2] = (numTxnsPerBlock);
        contexts[3] = (timeStamp2);

        // txLengths is defined as: Array of lengths of each encoded tx in txBatch
        // txBatch is defined as: Batch of RLP-encoded transactions
        uint256 txBatchVersion = _helper_sequencerInbox_appendTx_Version();
        bytes memory txBatch = _helper_createTxBatch_hardcoded();
        uint256[] memory txLengths = _helper_findTxLength_hardcoded();

        // Pranking as the sequencer and calling appendTxBatch
        vm.prank(sequencerAddress);
        seqIn.appendTxBatch(contexts, txLengths, firstL2BlockNumber, txBatchVersion, txBatch);

        uint256 inboxSizeFinal = seqIn.getInboxSize();

        assertGt(inboxSizeFinal, inboxSizeInitial);

        uint256 expectedInboxSize = numTxns;
        assertEq(inboxSizeFinal, expectedInboxSize);
    }

    function test_appendTxBatch_txBatchDataOverflow_reverts(uint256 numTxnsPerBlock, uint256 txnBlocks) public {
        // We will operate at a limit of transactionsPerBlock = 30 and number of transactionBlocks = 10.
        numTxnsPerBlock = bound(numTxnsPerBlock, 1, 30);
        txnBlocks = bound(txnBlocks, 1, 10);

        // Each context corresponds to a single "L2 block"
        uint256 numTxns = numTxnsPerBlock * txnBlocks;
        uint256 numContextsArrEntries = 2 * txnBlocks; // Since each `context` is represented with uint256 3-tuple: (numTxs, l2BlockNumber, l2Timestamp)

        // Making sure that the block.timestamp is a reasonable value (> txnBlocks)
        vm.warp(block.timestamp + (4 * txnBlocks));
        uint256 txnBlockTimestamp = block.timestamp - (2 * txnBlocks); // Subtracing just `txnBlocks` would have sufficed. However we are subtracting 2 times txnBlocks for some margin of error.
        // The objective for this subtraction is that while building the `contexts` array, no timestamp should go higher than the current block.timestamp

        uint256 firstL2BlockNumber = block.timestamp / 20;

        // Let's create an array of contexts
        uint256[] memory contexts = new uint256[](numContextsArrEntries);
        for (uint256 i = 0; i < numContextsArrEntries; i += 2) {
            // The first entry for `contexts` for each txnBlock is `numTxns` which we are keeping as constant for all blocks for this test
            contexts[i] = numTxnsPerBlock;

            // Formula used for blockTimestamp: (current block.timestamp) / 5x
            contexts[i + 1] = txnBlockTimestamp;

            // The only requirement for timestamps for the transaction blocks is that, these timestamps are monotonically increasing.
            // So, let's increase the value of txnBlock's timestamp monotonically, in a way that is does not exceed current block.timestamp
            ++txnBlockTimestamp;
        }

        // txLengths is defined as: Array of lengths of each encoded tx in txBatch
        // txBatch is defined as: Batch of RLP-encoded transactions
        (bytes memory txBatch, uint256[] memory txLengths) = _helper_sequencerInbox_appendTx(numTxns);
        uint256 txBatchVersion = _helper_sequencerInbox_appendTx_Version();

        // Now, we want to trigger the `txnBatchDataOverflow`, so we want to disturn the values receieved in the txLengths array.
        for (uint256 i; i < numTxns; i++) {
            txLengths[i] = txLengths[i] + 1;
        }

        // Pranking as the sequencer and calling appendTxBatch (should throw the TxBatchDataOverflow error)
        vm.expectRevert(ISequencerInbox.TxBatchDataOverflow.selector);
        vm.prank(sequencerAddress);
        seqIn.appendTxBatch(contexts, txLengths, firstL2BlockNumber, txBatchVersion, txBatch);
    }

    function test_appendTxBatch_paused_reverts(uint256 numTxnsPerBlock, uint256 txnBlocks) public {
        // We will operate at a limit of transactionsPerBlock = 30 and number of transactionBlocks = 10.
        numTxnsPerBlock = bound(numTxnsPerBlock, 1, 30);
        txnBlocks = bound(txnBlocks, 1, 10);

        vm.prank(sequencerOwner);
        seqIn.pause();

        // Each context corresponds to a single "L2 block"
        uint256 numTxns = numTxnsPerBlock * txnBlocks;
        // Each `context` is represented with uint256 2-tuple: (numTxs, l2Timestamp)
        uint256 numContextsArrEntries = 2 * txnBlocks;

        // Making sure that the block.timestamp is a reasonable value (> txnBlocks)
        vm.warp(block.timestamp + (4 * txnBlocks));
        uint256 txnBlockTimestamp = block.timestamp - (2 * txnBlocks); // Subtracing just `txnBlocks` would have sufficed. However we are subtracting 2 times txnBlocks for some margin of error.
        // The objective for this subtraction is that while building the `contexts` array, no timestamp should go higher than the current block.timestamp

        uint256 firstL2BlockNumber = block.timestamp / 20;

        // Let's create an array of contexts
        uint256[] memory contexts = new uint256[](numContextsArrEntries);
        for (uint256 i = 0; i < numContextsArrEntries; i += 2) {
            // The first entry for `contexts` for each txnBlock is `numTxns` which we are keeping as constant for all blocks for this test
            contexts[i] = numTxnsPerBlock;

            // Formula used for blockTimestamp: (current block.timestamp) / 5x
            contexts[i + 1] = txnBlockTimestamp;

            // The only requirement for timestamps for the transaction blocks is that, these timestamps are monotonically increasing.
            // So, let's increase the value of txnBlock's timestamp monotonically, in a way that is does not exceed current block.timestamp
            ++txnBlockTimestamp;
        }

        // txLengths is defined as: Array of lengths of each encoded tx in txBatch
        // txBatch is defined as: Batch of RLP-encoded transactions
        (bytes memory txBatch, uint256[] memory txLengths) = _helper_sequencerInbox_appendTx(numTxns);
        uint256 txBatchVersion = _helper_sequencerInbox_appendTx_Version();

        // Pranking as the sequencer and calling appendTxBatch
        vm.prank(sequencerAddress);
        vm.expectRevert("Pausable: paused");
        seqIn.appendTxBatch(contexts, txLengths, firstL2BlockNumber, txBatchVersion, txBatch);
    }

    /////////////////////////////////////////////////////////////////////////////////////////
    // TENTATIVE CODE CHANGE. TEST SUBJECT TO CHANGE BASED ON CHANGE IN CODE
    // Probable change: `appendTxBatch` passes even with malformed `contexts` array
    /////////////////////////////////////////////////////////////////////////////////////////
    function test_appendTxBatch_incompleteDataInContextsArray_succeeds(uint256 numTxnsPerBlock) public {
        // Since we are assuming that we will have two transaction blocks and we have a total of 300 sample transactions right now.
        numTxnsPerBlock = bound(numTxnsPerBlock, 1, 150);
        uint256 inboxSizeInitial = seqIn.getInboxSize();

        uint256 firstL2BlockNumber = block.timestamp / 20;

        // Each context corresponds to a single "L2 block"
        // `contexts` is represented with uint256 3-tuple: (numTxs, l2BlockNumber, l2Timestamp)
        // Let's create an array of contexts
        uint256 numTxns = numTxnsPerBlock * 2;
        uint256 timeStamp1 = block.timestamp / 10;

        uint256[] memory contexts = new uint256[](3);

        // Let's assume that we had 2 blocks and each had 3 transactions, but we fail to pass the block.timestamp and block.number of the 2nd transaction block.
        contexts[0] = (numTxnsPerBlock);
        contexts[1] = (timeStamp1);
        contexts[2] = (numTxnsPerBlock);

        // txLengths is defined as: Array of lengths of each encoded tx in txBatch
        // txBatch is defined as: Batch of RLP-encoded transactions
        (bytes memory txBatch, uint256[] memory txLengths) = _helper_sequencerInbox_appendTx(numTxns);
        uint256 txBatchVersion = _helper_sequencerInbox_appendTx_Version();

        // Pranking as the sequencer and calling appendTxBatch
        vm.prank(sequencerAddress);
        seqIn.appendTxBatch(contexts, txLengths, firstL2BlockNumber, txBatchVersion, txBatch);

        uint256 inboxSizeFinal = seqIn.getInboxSize();

        assertGt(inboxSizeFinal, inboxSizeInitial);
        assertEq(inboxSizeFinal, numTxnsPerBlock); // Since the timestamp and block.number were not included for the 2nd block, only 1st block's 3 txns are included.
    }

    //////////////////////////////
    // verifyTxInclusion
    // commit a batch and verify the Nth transaction
    //////////////////////////////
    function test_verifyTxInclusion_succeeds(uint256 numTxPerBlock, uint256 numBlocks, uint256 txToVerify) public {
        numTxPerBlock = bound(numTxPerBlock, 1, 30);
        numBlocks = bound(numBlocks, 1, 5);
        uint256 numTx = numTxPerBlock * numBlocks;

        // append a batch of transactions to the sequencer
        assertEq(seqIn.getInboxSize(), 0);
        vm.warp(block.timestamp + (4 * numBlocks));

        uint256[] memory contexts = generateContexts(numBlocks, numTxPerBlock);

        (bytes memory txBatch, uint256[] memory txLengths) = _helper_sequencerInbox_appendTx(numTx);
        uint256 txBatchVersion = _helper_sequencerInbox_appendTx_Version();

        vm.prank(sequencerAddress);
        seqIn.appendTxBatch(contexts, txLengths, 0, txBatchVersion, txBatch);
        assertEq(seqIn.getInboxSize(), numTx);

        // randomly choose a transaction to verify and prepare the proof
        txToVerify = bound(txToVerify, 0, numTx - 1);

        bytes memory txAfterData = generateTxAfterData(txToVerify, numTx, numTxPerBlock, contexts);

        // prepare the encoded transaction we want to verify and its context hash
        bytes memory encodedTx = rlpEncodedTransactions[txToVerify % 10];

        bytes32 proofContextHash = generateProofContextHash(txToVerify, numTxPerBlock, contexts);

        // prepare the accumulator hash of the preceding transactions in the batch
        bytes32 accBefore = generateAccumulator(txToVerify, numTxPerBlock, contexts);

        {
            uint256 batchNum = 0;
            uint256 numTxBefore = txToVerify;
            uint256 numTxAfter = numTx - txToVerify - 1;

            bytes memory batchInfo = abi.encodePacked(batchNum, numTxBefore, numTxAfter, accBefore);

            bytes memory proof = abi.encodePacked(proofContextHash, batchInfo, txAfterData);

            seqIn.verifyTxInclusion(encodedTx, proof);
        }
    }

    /////////////////////////
    // Auxillary Functions
    /////////////////////////

    function generateContexts(uint256 numBlocks, uint256 numTxPerBlock) public view returns (uint256[] memory) {
        uint256 txBlockTimestamp = block.timestamp - (2 * numBlocks);
        uint256 numContextsArrEntries = numBlocks * 2;

        uint256[] memory contexts = new uint256[](numContextsArrEntries);
        for (uint256 i = 0; i < numContextsArrEntries; i += 2) {
            contexts[i] = numTxPerBlock;
            contexts[i + 1] = txBlockTimestamp;
            txBlockTimestamp++;
        }

        return contexts;
    }

    function generateTxAfterData(uint256 txToVerify, uint256 numTx, uint256 numTxPerBlock, uint256[] memory contexts)
        public
        view
        returns (bytes memory)
    {
        bytes memory txAfterData;
        bytes32 txContextHash;

        for (uint256 i = txToVerify + 1; i < numTx; i++) {
            uint256 currentBlock = i / numTxPerBlock;

            txContextHash = keccak256(
                abi.encodePacked(
                    sequencerAddress,
                    currentBlock,
                    contexts[(2 * currentBlock) + 1] // L2 Timestamp
                )
            );

            bytes memory txData = rlpEncodedTransactions[i % 10];
            bytes32 txDataHash = keccak256(txData);

            txAfterData = abi.encodePacked(txAfterData, abi.encodePacked(abi.encodePacked(txContextHash, txDataHash)));
        }

        return txAfterData;
    }

    function generateProofContextHash(uint256 txToVerify, uint256 numTxPerBlock, uint256[] memory contexts)
        public
        view
        returns (bytes32)
    {
        uint256 blockToVerify = txToVerify / numTxPerBlock;
        uint256 txL2Timestamp = contexts[(2 * blockToVerify) + 1];

        return keccak256(abi.encodePacked(sequencerAddress, blockToVerify, txL2Timestamp));
    }

    function generateAccumulator(uint256 txToVerify, uint256 numTxPerBlock, uint256[] memory contexts)
        public
        view
        returns (bytes32)
    {
        bytes32 txContextHash;
        bytes32 accBefore;

        for (uint256 i = 0; i < txToVerify; i++) {
            uint256 currentBlock = i / numTxPerBlock;

            txContextHash = keccak256(
                abi.encodePacked(
                    sequencerAddress,
                    currentBlock, // L2 Block Number
                    contexts[(2 * currentBlock) + 1] // L2 Timestamp
                )
            );

            bytes memory txData = rlpEncodedTransactions[i % 10];
            bytes32 txDataHash = keccak256(txData);
            accBefore = keccak256(abi.encodePacked(accBefore, i, txContextHash, txDataHash));
        }

        return accBefore;
    }
}
