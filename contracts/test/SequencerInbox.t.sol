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
import {RLPEncodedTransactionsUtil} from "./utils/RLPEncodedTransactions.sol";

contract SequencerBaseSetup is Test, RLPEncodedTransactionsUtil {
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
    ////////////////////////////////////
    // Utility Function
    ////////////////////////////////////

    function createTxBatch() public pure returns (bytes memory) {
        bytes memory tx1 =
            "0xf87980830928028807db92dd3a67aec0943843c3ef2e5aeb50b18f70c05489cdf4ee02f265911d063fadb8051a25d767a4113cda640000802ea0940b36911fcb70d982f5ba1c919600acb86412d90a01186185ed960006864ba3a02bc5b5c686ca8593a69dc27961aecac49b4e9ea9cc374fdc981799bb09887a82";
        bytes memory tx2 =
            "0xf8798083092801880c5b572ca09a1ab094a43a08212ba8b14e3b1cf9ccaf854a40973511be9117ffc93897f5d9a601a3142ade1c640000802ea0c7f40eff33312a972b7e5eaa277345ae9c9be1afa70e991db798da24dfa03782a05ef4400508c80bf56e089d9458c6d49ba4fa21a97e0f2f164387994e57472ff0";
        bytes memory tx3 =
            "0xf87980830928018809431dd5a59585dc94f45d59c6e6fabb0e927ee496975e2ccdb47ed754911b3288a788f83591e788770843e3640000802ea0c6ae9f3f6c1de0de9c6b34016c75f052ffb740f214ca979f217b6b27c66cddaea03a8b4e8e270969d0b11afee29bd5a25c65b725da478604ba0d54b9f3b767c7df";
        bytes memory tx4 =
            "0xf879808309280188064606a4f40a2c8894f0f66416890155860c7763312c61c04ba0662cbc91085966cd6130788c6cfa6300af1fe40000802da06ce8941605db844d69d6435aaa17e8a318f7ffeaa1ddb1875a4dc43f5e5c5fc6a032a20d32ccf119686cc6b69ca5f61f92864a2d57e409cac48aa411ae55483e7a";
        bytes memory tx5 =
            "0xf8798083092801880b3bdd5e4cc2f2e49499b7d78a719bb1a7c42fd137029704b6cae4fa449104a412c41700122abf98c1f054d9e40000802da0a8ebcdc1d48d1b425d26f75c56095868fa79600cb28c5f57b3aef4090225576ea0472d2a86c00dfabda26fa6cdd8a4c4df394e762b7d8117346fff44371c8584bd";
        bytes memory tx6 =
            "0xf87980830928018809a25f73849892b894db16cfe9acc3b7da9e5aaace521483c50c43cfec910755f22d8ff7ee357f900c7f2bfee40000802ea0f08397648419c67baf8d8df16d13f86cc806ca429664e4fb0aba6425907a825fa06074c6910dc61bd6f73d7db1e4f0d1a3cf5bcecb29389e6b2260b4327487c643";

        bytes[6] memory txnArray = [tx1, tx2, tx3, tx4, tx5, tx6];

        return abi.encode(txnArray);
    }

    function findTxLength() public pure returns (uint256) {
        bytes memory tx1 = "0xf87980830928028807db92dd3a67aec0943843c3ef2e5aeb50b18f70c05489cdf4ee02f265911d063fadb8051a25d767a4113cda640000802ea0940b36911fcb70d982f5ba1c919600acb86412d90a01186185ed960006864ba3a02bc5b5c686ca8593a69dc27961aecac49b4e9ea9cc374fdc981799bb09887a82";
        return tx1.length;
    }

    function createTxBatchDiff() public view returns (bytes memory) {
        return abi.encode(
            RLPEncodedTransactionsUtil.rlpEncodedTransactions[0],
            RLPEncodedTransactionsUtil.rlpEncodedTransactions[1],
            RLPEncodedTransactionsUtil.rlpEncodedTransactions[2],
            RLPEncodedTransactionsUtil.rlpEncodedTransactions[3],
            RLPEncodedTransactionsUtil.rlpEncodedTransactions[4],
            RLPEncodedTransactionsUtil.rlpEncodedTransactions[5]
        );
    }

    function findTxLengthDiff() public view returns(uint256[] memory) {
        uint256[] memory transactionLengthArray = new uint256[](6);
        for (uint i; i < 6; ) {
            transactionLengthArray[i] = RLPEncodedTransactionsUtil.rlpEncodedTransactions[i].length;
            unchecked {
                ++i;
            }
        }
        return transactionLengthArray;
    }

    /////////////////////////////////
    // SequencerInbox Setup
    /////////////////////////////////

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

    /////////////////////////////////
    // SequencerInbox Tests
    /////////////////////////////////

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
    function test_appendTxBatch_positiveCase_1(uint256 numTxnsPerBlock) public {
        // Since we are assuming that we will have two transaction blocks and we have a total of 300 sample transactions right now.
        numTxnsPerBlock = bound(numTxnsPerBlock, 1, 150); 
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
        (bytes memory txBatch, uint256[] memory txLengths) = _helper_sequencerInbox_appendTx(numTxns);

        // Pranking as the sequencer and calling appendTxBatch
        vm.prank(sequencer);
        seqIn.appendTxBatch(contexts, txLengths, txBatch);

        uint256 inboxSizeFinal = seqIn.getInboxSize();

        assertGt(inboxSizeFinal, inboxSizeInitial);

        uint256 expectedInboxSize = numTxns;
        assertEq(inboxSizeFinal, expectedInboxSize);
    }

    function test_appendTxBatch_positiveCase_2(uint256 numTxnsPerBlock) public {
        // Since we are assuming that we will have two transaction blocks and we have a total of 300 sample transactions right now.
        numTxnsPerBlock = 3;
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
        bytes memory txBatch = createTxBatchDiff();
        uint256[] memory txLengths = findTxLengthDiff();

        // Pranking as the sequencer and calling appendTxBatch
        vm.prank(sequencer);
        seqIn.appendTxBatch(contexts, txLengths, txBatch);

        uint256 inboxSizeFinal = seqIn.getInboxSize();

        assertGt(inboxSizeFinal, inboxSizeInitial);

        uint256 expectedInboxSize = numTxns;
        assertEq(inboxSizeFinal, expectedInboxSize);
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
        vm.prank(sequencer);
        seqIn.appendTxBatch(contexts, txLengths, txBatch);

        uint256 inboxSizeFinal = seqIn.getInboxSize();

        assertGt(inboxSizeFinal, inboxSizeInitial);
        assertEq(inboxSizeFinal, numTxnsPerBlock); // Since the timestamp and block.number were not included for the 2nd block, only 1st block's 3 txns are included.
    }
}
