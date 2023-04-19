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
import "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";
import {SequencerInbox} from "../src/SequencerInbox.sol";
import {Utils} from "./utils/Utils.sol";
import {RLPEncodedTransactionsUtil} from "./utils/RLPEncodedTransactions.sol";

contract SequencerBaseSetup is Test, RLPEncodedTransactionsUtil {
    address internal alice;
    address internal bob;
    address internal sequencerOwner;
    address internal sequencerAddress;

    function setUp() public virtual {
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
    function test_initializeSequencerInbox_withSequencerAddressZero() external {
        bytes memory seqInInitData = abi.encodeWithSignature("initialize(address)", address(0));
        vm.startPrank(sequencerOwner);
        vm.expectRevert(ZeroAddress.selector);
        seqIn = SequencerInbox(address(new ERC1967Proxy(address(implementationSequencer), seqInInitData)));
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
        seqIn.verifyTxInclusion(bytes(""), emptyProof);
    }

    // Only sequencer can append transaction batches to the sequencerInbox
    // function test_RevertWhen_InvalidSequencer() public {
    //     vm.expectRevert(abi.encodeWithSelector(NotSequencer.selector, alice, sequencerAddress));
    //     vm.prank(alice);
    //     uint256[] memory contexts = new uint256[](1);
    //     uint256[] memory txLengths = new uint256[](1);
    //     seqIn.appendTxBatch(contexts, txLengths, "0x");
    // }

    // function test_RevertWhen_EmptyBatch() public {
    //     vm.expectRevert(ISequencerInbox.EmptyBatch.selector);
    //     vm.prank(sequencerAddress);
    //     uint256[] memory contexts = new uint256[](1);
    //     uint256[] memory txLengths = new uint256[](1);
    //     seqIn.appendTxBatch(contexts, txLengths, "0x");
    // }

    //////////////////////////////
    // forceInclusion
    //////////////////////////////

    function test_force_inclusion() external {
        vm.prank(bob);
        vm.deal(bob, 1 ether);

        uint256 baseFee;
        uint256 gasLimit = 726097;
        uint256 maxFeePerGas = 505;
        uint256 nonce = 0;
        address to = address(alice);
        uint256 value = 0.1 ether;
        bytes memory data = "0xcb90549900000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000003e000000000000000000000000000000000000000000000000000000000000007400000000000000000000000000000000000000000000000000000000000000a0000000000000000000000000000000000000000000000000000000000000000e00b75a2bb8967f067ba368dfac9723cec3a52f2372f9beea4cdca4a98f897c97f1feb7b29f7257707e3c4d267ccf0d8115218523b6caed6bcec87385f892f52637a6df42950a6e2747d2e81355686a3ece8678aff1014ea1c1985d5cbe2a74ce7000000000000000000000000000000000000000000000000000000000000138dd89594d7e4253f933e49f0c62b78c1bd0e32c0b67f17eef1830459c409b3900500000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000001c66b2eccfc15d71513b10a762eae615e5f09aa4ca33afe058575337ed5aa32b80638b6237ae9e79efe6bcfc0505fefed44cf6a01dfbc0e0bc7a62bbcb7d4d5b000000000000000000000000000000000000000000000000000000000000013b200000000000000000000000000000000000000000000000000000000642538e4000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000007a12000000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001260000008ff88d821794830f4240829bcb94530000000000000000000000000000000000000280a4bede39b50000000000000000000000000000000000000000000000000000003b75b8eb8485027cffa9a0a0c5e9c681dce4160020154cc7c30b6ca1f902cd63e08b69ea0e6f740cd9540fd0a04b6ea19a518ec5014e09cf2d9122d9301f251d5312dc283f443dac2baa3d54b00000008ff88d821795830f4240829bcb94530000000000000000000000000000000000000280a4bede39b50000000000000000000000000000000000000000000000000000004aa5c0d00885027cffa99fa0a01129f2ed7b68bf78be970a45e7a6180f0fe83fa2e9ef81197554ad0447c530a00340d398109253530423892e183bfaf9db6b88e8b09c18ac1411a7acf77fbcf8000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000e01feb7b29f7257707e3c4d267ccf0d8115218523b6caed6bcec87385f892f52632e07f8d6cb3cbb9bee8f99742b2bb829a8615d54d2222a84138ac6568e511fce7a6df42950a6e2747d2e81355686a3ece8678aff1014ea1c1985d5cbe2a74ce7000000000000000000000000000000000000000000000000000000000000138ea03fd453531143c17c2cf68911397e851343cfdae9102a71de8ba7007d286b53000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000000017f5bce0304a3a0b0c48e6a356bcb6b7e13009a709a2d543f34573136bbae2b4dc66b2eccfc15d71513b10a762eae615e5f09aa4ca33afe058575337ed5aa32b800000000000000000000000000000000000000000000000000000000000013b30000000000000000000000000000000000000000000000000000000064253a64000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000007a12000000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001260000008ff88d821796830f4240829bcb94530000000000000000000000000000000000000280a4bede39b50000000000000000000000000000000000000000000000000000005da3fd29e985027cffa99fa0ef6f3d6b17bad2a82704ceb3bca3025457348a45ca5455915ab45f454885c18fa06f6e84960c4ce75b4974b731a4bd2d01e72195b90fab1d46fee41be408543cef0000008ff88d821797830f4240829bcb94530000000000000000000000000000000000000280a4bede39b50000000000000000000000000000000000000000000000000000004f6360414f85027cffa9a0a0485c5fc9eb33a6fe5ea8ed8c9b5fe7f800b7ef1cccc75d6eca5407e0715e8d40a0263016339ccbd519e1b727c001a56d5d3858b274e1a37d9ced457c83ecdc8877000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000e02e07f8d6cb3cbb9bee8f99742b2bb829a8615d54d2222a84138ac6568e511fce2fd92079ed8efe2f52018da39aa9c7a948759d320464ee5c2c6cabd6bc2bf3ce7a6df42950a6e2747d2e81355686a3ece8678aff1014ea1c1985d5cbe2a74ce7000000000000000000000000000000000000000000000000000000000000138fb53d10383607d6c7778c5d5e4e8523c00173cdb7352d187aca83e157ac01ddcb00000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000001c7b42481e3447080bb30bfa96ff0f4714d4f26286dfe40d4c42ec1e86b3aaf8e7f5bce0304a3a0b0c48e6a356bcb6b7e13009a709a2d543f34573136bbae2b4d00000000000000000000000000000000000000000000000000000000000013b40000000000000000000000000000000000000000000000000000000064253a67000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000007a12000000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000930000008ff88d821798830f4240829bcb94530000000000000000000000000000000000000280a4bede39b50000000000000000000000000000000000000000000000000000004a213816ac85027cffa99fa0bcf84d148d8fb2c64092d9b187ba7c138bef07c7bbf6b9849a5d454a5272f2d7a04cf5feca9b1132f0e2955e1063ba5243ab1fd79148d62594f76de9826a0258ff0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000e02fd92079ed8efe2f52018da39aa9c7a948759d320464ee5c2c6cabd6bc2bf3ce29de806427c3abf320edf2b8b1ce4af3676c23083623a37ac70d2ac34d6ed9ca7a6df42950a6e2747d2e81355686a3ece8678aff1014ea1c1985d5cbe2a74ce700000000000000000000000000000000000000000000000000000000000013900ea5bd6f10bff321d91cf1e0cd8c1926fff598afb12456837e7b824aa86022db0000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000000158f053faec41c7a6837b855f9d130c065ef911476a774bb5a97b480bda9ba7c9c7b42481e3447080bb30bfa96ff0f4714d4f26286dfe40d4c42ec1e86b3aaf8e00000000000000000000000000000000000000000000000000000000000013b50000000000000000000000000000000000000000000000000000000064253be4000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000007a12000000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000930000008ff88d821799830f424082ce2f94530000000000000000000000000000000000000280a4bede39b50000000000000000000000000000000000000000000000000000005f7a1cebce85027cffa9a0a09bc1b7d6e013cb9504a34f6fc55ad57f0363173693b07e0fb68bb2539b5c8532a058b4def30586b44b2a3212a045d8d902c481d11c93903659c4b4d43ce16d897700000000000000000000000000";

        bytes32 dataHashCalc = keccak256(
            abi.encodePacked(
                gasLimit,
                maxFeePerGas,
                nonce,
                to,
                value,
                data
            )
        );

        console.log("dataHashCalc");
        console.logBytes32(dataHashCalc);

        seqIn.sendUnsignedTx(
            gasLimit,
            maxFeePerGas,
            nonce,
            to,
            value,
            data
        );

        bytes32 dataHash = seqIn.messageDataHashes(0);
        uint count = seqIn.delayedMessageCounter();

        uint64[2] memory l1BlockAndTime;

        l1BlockAndTime[0] = seqIn.l1BlockAndTime(0);
        l1BlockAndTime[1] = seqIn.l1BlockAndTime(1);
        baseFee = seqIn.baseFee();

        assertEq(count,1);
        assertEq(dataHash, dataHashCalc);

        vm.warp(block.timestamp + 86401);
        vm.roll(block.number + 5761);

        console.log(block.timestamp);


        bytes32 messageHashCalc = keccak256(
            abi.encodePacked(
                address(bob),
                l1BlockAndTime[0],
                l1BlockAndTime[1],
                uint(0),
                baseFee,
                dataHashCalc
            )
        );

        console.log("printing calculated message hash");
        console.logBytes32(messageHashCalc);

        seqIn.forceInclusion(
            0,
            address(bob),
            baseFee,
            l1BlockAndTime,
            dataHash
        );

        assertEq(seqIn.delayedMessagesRead(), 1);

    }




    //////////////////////////////
    // appendTxBatch
    //////////////////////////////
    function test_appendTxBatch_positiveCase_1(uint256 numTxnsPerBlock, uint256 txnBlocks) public {
        // We will operate at a limit of transactionsPerBlock = 30 and number of transactionBlocks = 10.
        numTxnsPerBlock = bound(numTxnsPerBlock, 1, 30);
        txnBlocks = bound(txnBlocks, 1, 10);


        console.log("initial", seqIn.getInboxSize());
        // Each context corresponds to a single "L2 block"
        uint256 numTxns = numTxnsPerBlock * txnBlocks;
        uint256 numContextsArrEntries = 3 * txnBlocks; // Since each `context` is represented with uint256 triplet: (numTxs, l2BlockNumber, l2Timestamp)

        // Making sure that the block.timestamp is a reasonable value (> txnBlocks)
        vm.warp(block.timestamp + (4 * txnBlocks));
        uint256 txnBlockTimestamp = block.timestamp - (2 * txnBlocks); // Subtracing just `txnBlocks` would have sufficed. However we are subtracting 2 times txnBlocks for some margin of error.
            // The objective for this subtraction is that while building the `contexts` array, no timestamp should go higher than the current block.timestamp

        // Let's create an array of contexts
        uint256[] memory contexts = new uint256[](numContextsArrEntries);
        for (uint256 i; i < numContextsArrEntries; i += 3) {
            // The first entry for `contexts` for each txnBlock is `numTxns` which we are keeping as constant for all blocks for this test
            contexts[i] = numTxnsPerBlock;

            // Formula Used for blockNumber: (txnBlock's block.timestamp) / 20;
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



        vm.prank(bob);
        vm.deal(bob, 1 ether);

        uint256 gasLimit = 726097;
        uint256 maxFeePerGas = 505;
        uint256 nonce = 0;
        address to = address(alice);
        uint256 value = 0.1 ether;
        bytes memory data = "0xcb90549900000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000003e000000000000000000000000000000000000000000000000000000000000007400000000000000000000000000000000000000000000000000000000000000a0000000000000000000000000000000000000000000000000000000000000000e00b75a2bb8967f067ba368dfac9723cec3a52f2372f9beea4cdca4a98f897c97f1feb7b29f7257707e3c4d267ccf0d8115218523b6caed6bcec87385f892f52637a6df42950a6e2747d2e81355686a3ece8678aff1014ea1c1985d5cbe2a74ce7000000000000000000000000000000000000000000000000000000000000138dd89594d7e4253f933e49f0c62b78c1bd0e32c0b67f17eef1830459c409b3900500000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000001c66b2eccfc15d71513b10a762eae615e5f09aa4ca33afe058575337ed5aa32b80638b6237ae9e79efe6bcfc0505fefed44cf6a01dfbc0e0bc7a62bbcb7d4d5b000000000000000000000000000000000000000000000000000000000000013b200000000000000000000000000000000000000000000000000000000642538e4000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000007a12000000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001260000008ff88d821794830f4240829bcb94530000000000000000000000000000000000000280a4bede39b50000000000000000000000000000000000000000000000000000003b75b8eb8485027cffa9a0a0c5e9c681dce4160020154cc7c30b6ca1f902cd63e08b69ea0e6f740cd9540fd0a04b6ea19a518ec5014e09cf2d9122d9301f251d5312dc283f443dac2baa3d54b00000008ff88d821795830f4240829bcb94530000000000000000000000000000000000000280a4bede39b50000000000000000000000000000000000000000000000000000004aa5c0d00885027cffa99fa0a01129f2ed7b68bf78be970a45e7a6180f0fe83fa2e9ef81197554ad0447c530a00340d398109253530423892e183bfaf9db6b88e8b09c18ac1411a7acf77fbcf8000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000e01feb7b29f7257707e3c4d267ccf0d8115218523b6caed6bcec87385f892f52632e07f8d6cb3cbb9bee8f99742b2bb829a8615d54d2222a84138ac6568e511fce7a6df42950a6e2747d2e81355686a3ece8678aff1014ea1c1985d5cbe2a74ce7000000000000000000000000000000000000000000000000000000000000138ea03fd453531143c17c2cf68911397e851343cfdae9102a71de8ba7007d286b53000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000000017f5bce0304a3a0b0c48e6a356bcb6b7e13009a709a2d543f34573136bbae2b4dc66b2eccfc15d71513b10a762eae615e5f09aa4ca33afe058575337ed5aa32b800000000000000000000000000000000000000000000000000000000000013b30000000000000000000000000000000000000000000000000000000064253a64000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000007a12000000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001260000008ff88d821796830f4240829bcb94530000000000000000000000000000000000000280a4bede39b50000000000000000000000000000000000000000000000000000005da3fd29e985027cffa99fa0ef6f3d6b17bad2a82704ceb3bca3025457348a45ca5455915ab45f454885c18fa06f6e84960c4ce75b4974b731a4bd2d01e72195b90fab1d46fee41be408543cef0000008ff88d821797830f4240829bcb94530000000000000000000000000000000000000280a4bede39b50000000000000000000000000000000000000000000000000000004f6360414f85027cffa9a0a0485c5fc9eb33a6fe5ea8ed8c9b5fe7f800b7ef1cccc75d6eca5407e0715e8d40a0263016339ccbd519e1b727c001a56d5d3858b274e1a37d9ced457c83ecdc8877000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000e02e07f8d6cb3cbb9bee8f99742b2bb829a8615d54d2222a84138ac6568e511fce2fd92079ed8efe2f52018da39aa9c7a948759d320464ee5c2c6cabd6bc2bf3ce7a6df42950a6e2747d2e81355686a3ece8678aff1014ea1c1985d5cbe2a74ce7000000000000000000000000000000000000000000000000000000000000138fb53d10383607d6c7778c5d5e4e8523c00173cdb7352d187aca83e157ac01ddcb00000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000001c7b42481e3447080bb30bfa96ff0f4714d4f26286dfe40d4c42ec1e86b3aaf8e7f5bce0304a3a0b0c48e6a356bcb6b7e13009a709a2d543f34573136bbae2b4d00000000000000000000000000000000000000000000000000000000000013b40000000000000000000000000000000000000000000000000000000064253a67000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000007a12000000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000930000008ff88d821798830f4240829bcb94530000000000000000000000000000000000000280a4bede39b50000000000000000000000000000000000000000000000000000004a213816ac85027cffa99fa0bcf84d148d8fb2c64092d9b187ba7c138bef07c7bbf6b9849a5d454a5272f2d7a04cf5feca9b1132f0e2955e1063ba5243ab1fd79148d62594f76de9826a0258ff0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000e02fd92079ed8efe2f52018da39aa9c7a948759d320464ee5c2c6cabd6bc2bf3ce29de806427c3abf320edf2b8b1ce4af3676c23083623a37ac70d2ac34d6ed9ca7a6df42950a6e2747d2e81355686a3ece8678aff1014ea1c1985d5cbe2a74ce700000000000000000000000000000000000000000000000000000000000013900ea5bd6f10bff321d91cf1e0cd8c1926fff598afb12456837e7b824aa86022db0000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000000158f053faec41c7a6837b855f9d130c065ef911476a774bb5a97b480bda9ba7c9c7b42481e3447080bb30bfa96ff0f4714d4f26286dfe40d4c42ec1e86b3aaf8e00000000000000000000000000000000000000000000000000000000000013b50000000000000000000000000000000000000000000000000000000064253be4000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000007a12000000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000930000008ff88d821799830f424082ce2f94530000000000000000000000000000000000000280a4bede39b50000000000000000000000000000000000000000000000000000005f7a1cebce85027cffa9a0a09bc1b7d6e013cb9504a34f6fc55ad57f0363173693b07e0fb68bb2539b5c8532a058b4def30586b44b2a3212a045d8d902c481d11c93903659c4b4d43ce16d897700000000000000000000000000";

        bytes32 dataHashCalc = keccak256(
            abi.encodePacked(
                gasLimit,
                maxFeePerGas,
                nonce,
                to,
                value,
                data
            )
        );

        console.log("dataHashCalc");
        console.logBytes32(dataHashCalc);

        seqIn.sendUnsignedTx(
            gasLimit,
            maxFeePerGas,
            nonce,
            to,
            value,
            data
        );



        // Pranking as the sequencer and calling appendTxBatch
        vm.prank(sequencerAddress);
        seqIn.appendTxBatch(contexts, txLengths, txBatch, 0);

        assertEq(seqIn.getInboxSize(), numTxns);
    }


    // testing the appendTxBatch with the delayedMessages


    uint256 txnL2BlockNumber;
    uint256 txnL2Timestamp;

    uint256 batchNum;
    uint256 numTxsBefore;
    uint256 numTxsAfterInBatch;
    bytes32 accBefore;


    function test_appendTxBatch_delayed(uint256 numTxnsPerBlock, uint256 txnBlocks) public {
        // We will operate at a limit of transactionsPerBlock = 30 and number of transactionBlocks = 10.
        numTxnsPerBlock = bound(numTxnsPerBlock, 1, 30);
        txnBlocks = bound(txnBlocks, 1, 10);

        console.log("initial", seqIn.getInboxSize());

        // Each context corresponds to a single "L2 block"
        uint256 numTxns = numTxnsPerBlock * txnBlocks;
        uint256 numContextsArrEntries = 3 * txnBlocks; // Since each `context` is represented with uint256 triplet: (numTxs, l2BlockNumber, l2Timestamp)

        // Making sure that the block.timestamp is a reasonable value (> txnBlocks)
        vm.warp(block.timestamp + (4 * txnBlocks));
        uint256 txnBlockTimestamp = block.timestamp - (2 * txnBlocks); // Subtracing just `txnBlocks` would have sufficed. However we are subtracting 2 times txnBlocks for some margin of error.
            // The objective for this subtraction is that while building the `contexts` array, no timestamp should go higher than the current block.timestamp

        // Let's create an array of contexts
        uint256[] memory contexts = new uint256[](numContextsArrEntries);
        for (uint256 i; i < numContextsArrEntries; i += 3) {
            // The first entry for `contexts` for each txnBlock is `numTxns` which we are keeping as constant for all blocks for this test
            contexts[i] = numTxnsPerBlock;

            // Formula Used for blockNumber: (txnBlock's block.timestamp) / 20;
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


        vm.prank(bob);
        vm.deal(bob, 1 ether);

        uint256 gasLimit = 726097;
        uint256 maxFeePerGas = 505;
        uint256 nonce = 0;
        address to = address(alice);
        uint256 value = 0.1 ether;
        bytes memory data = "0xcb90549900000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000003e000000000000000000000000000000000000000000000000000000000000007400000000000000000000000000000000000000000000000000000000000000a0000000000000000000000000000000000000000000000000000000000000000e00b75a2bb8967f067ba368dfac9723cec3a52f2372f9beea4cdca4a98f897c97f1feb7b29f7257707e3c4d267ccf0d8115218523b6caed6bcec87385f892f52637a6df42950a6e2747d2e81355686a3ece8678aff1014ea1c1985d5cbe2a74ce7000000000000000000000000000000000000000000000000000000000000138dd89594d7e4253f933e49f0c62b78c1bd0e32c0b67f17eef1830459c409b3900500000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000001c66b2eccfc15d71513b10a762eae615e5f09aa4ca33afe058575337ed5aa32b80638b6237ae9e79efe6bcfc0505fefed44cf6a01dfbc0e0bc7a62bbcb7d4d5b000000000000000000000000000000000000000000000000000000000000013b200000000000000000000000000000000000000000000000000000000642538e4000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000007a12000000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001260000008ff88d821794830f4240829bcb94530000000000000000000000000000000000000280a4bede39b50000000000000000000000000000000000000000000000000000003b75b8eb8485027cffa9a0a0c5e9c681dce4160020154cc7c30b6ca1f902cd63e08b69ea0e6f740cd9540fd0a04b6ea19a518ec5014e09cf2d9122d9301f251d5312dc283f443dac2baa3d54b00000008ff88d821795830f4240829bcb94530000000000000000000000000000000000000280a4bede39b50000000000000000000000000000000000000000000000000000004aa5c0d00885027cffa99fa0a01129f2ed7b68bf78be970a45e7a6180f0fe83fa2e9ef81197554ad0447c530a00340d398109253530423892e183bfaf9db6b88e8b09c18ac1411a7acf77fbcf8000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000e01feb7b29f7257707e3c4d267ccf0d8115218523b6caed6bcec87385f892f52632e07f8d6cb3cbb9bee8f99742b2bb829a8615d54d2222a84138ac6568e511fce7a6df42950a6e2747d2e81355686a3ece8678aff1014ea1c1985d5cbe2a74ce7000000000000000000000000000000000000000000000000000000000000138ea03fd453531143c17c2cf68911397e851343cfdae9102a71de8ba7007d286b53000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000000017f5bce0304a3a0b0c48e6a356bcb6b7e13009a709a2d543f34573136bbae2b4dc66b2eccfc15d71513b10a762eae615e5f09aa4ca33afe058575337ed5aa32b800000000000000000000000000000000000000000000000000000000000013b30000000000000000000000000000000000000000000000000000000064253a64000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000007a12000000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001260000008ff88d821796830f4240829bcb94530000000000000000000000000000000000000280a4bede39b50000000000000000000000000000000000000000000000000000005da3fd29e985027cffa99fa0ef6f3d6b17bad2a82704ceb3bca3025457348a45ca5455915ab45f454885c18fa06f6e84960c4ce75b4974b731a4bd2d01e72195b90fab1d46fee41be408543cef0000008ff88d821797830f4240829bcb94530000000000000000000000000000000000000280a4bede39b50000000000000000000000000000000000000000000000000000004f6360414f85027cffa9a0a0485c5fc9eb33a6fe5ea8ed8c9b5fe7f800b7ef1cccc75d6eca5407e0715e8d40a0263016339ccbd519e1b727c001a56d5d3858b274e1a37d9ced457c83ecdc8877000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000e02e07f8d6cb3cbb9bee8f99742b2bb829a8615d54d2222a84138ac6568e511fce2fd92079ed8efe2f52018da39aa9c7a948759d320464ee5c2c6cabd6bc2bf3ce7a6df42950a6e2747d2e81355686a3ece8678aff1014ea1c1985d5cbe2a74ce7000000000000000000000000000000000000000000000000000000000000138fb53d10383607d6c7778c5d5e4e8523c00173cdb7352d187aca83e157ac01ddcb00000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000001c7b42481e3447080bb30bfa96ff0f4714d4f26286dfe40d4c42ec1e86b3aaf8e7f5bce0304a3a0b0c48e6a356bcb6b7e13009a709a2d543f34573136bbae2b4d00000000000000000000000000000000000000000000000000000000000013b40000000000000000000000000000000000000000000000000000000064253a67000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000007a12000000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000930000008ff88d821798830f4240829bcb94530000000000000000000000000000000000000280a4bede39b50000000000000000000000000000000000000000000000000000004a213816ac85027cffa99fa0bcf84d148d8fb2c64092d9b187ba7c138bef07c7bbf6b9849a5d454a5272f2d7a04cf5feca9b1132f0e2955e1063ba5243ab1fd79148d62594f76de9826a0258ff0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000e02fd92079ed8efe2f52018da39aa9c7a948759d320464ee5c2c6cabd6bc2bf3ce29de806427c3abf320edf2b8b1ce4af3676c23083623a37ac70d2ac34d6ed9ca7a6df42950a6e2747d2e81355686a3ece8678aff1014ea1c1985d5cbe2a74ce700000000000000000000000000000000000000000000000000000000000013900ea5bd6f10bff321d91cf1e0cd8c1926fff598afb12456837e7b824aa86022db0000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000000158f053faec41c7a6837b855f9d130c065ef911476a774bb5a97b480bda9ba7c9c7b42481e3447080bb30bfa96ff0f4714d4f26286dfe40d4c42ec1e86b3aaf8e00000000000000000000000000000000000000000000000000000000000013b50000000000000000000000000000000000000000000000000000000064253be4000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000007a12000000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000930000008ff88d821799830f424082ce2f94530000000000000000000000000000000000000280a4bede39b50000000000000000000000000000000000000000000000000000005f7a1cebce85027cffa9a0a09bc1b7d6e013cb9504a34f6fc55ad57f0363173693b07e0fb68bb2539b5c8532a058b4def30586b44b2a3212a045d8d902c481d11c93903659c4b4d43ce16d897700000000000000000000000000";

        bytes32 dataHashCalc = keccak256(
            abi.encodePacked(
                gasLimit,
                maxFeePerGas,
                nonce,
                to,
                value,
                data
            )
        );

        console.log("dataHashCalc");
        console.logBytes32(dataHashCalc);

        seqIn.sendUnsignedTx(
            gasLimit,
            maxFeePerGas,
            nonce,
            to,
            value,
            data
        );



        // Pranking as the sequencer and calling appendTxBatch
        vm.prank(sequencerAddress);
        seqIn.appendTxBatch(contexts, txLengths, txBatch, 1);


        console.log(seqIn.getInboxSize());


        console.log(numTxns);
        assertEq(seqIn.getInboxSize(), numTxns + 1);

        assertEq(seqIn.delayedMessagesRead(), 1);
        // console.log("delayed messages read: ",seqIn.delayedMessagesRead());




        // transaction verification

        {
            bytes memory batchTransactionsMetadata;

            txnL2BlockNumber = contexts[numContextsArrEntries-1];
            txnL2Timestamp = contexts[numContextsArrEntries-2];

            bytes32 txContextHash = keccak256(
                abi.encodePacked(
                    sequencerAddress,
                    txnL2BlockNumber,
                    txnL2Timestamp
                )
            );

            batchNum = 0; //seqIn.accumulators(0) returns a non-zero value while seqIn.accumulators(1) reverts.
            numTxsBefore = numTxns-1; // since we are verifying the inclusion of the 1st transaction of the first block.
            numTxsAfterInBatch = 0; // since we are verifying the inclusion of the 1st transaction of the first block and the total number of transactions are numTxnsPerBlock * txnBlocks
            accBefore = 0x168ce4eb8a565b51ae25d2bb8ad768afedec13bd9ffc481090de96852f4bb3fd; // again, since this is our first inclusion, acc before that would have been empty bytes too.




            bytes memory batchInfo = abi.encodePacked(
                batchNum,
                numTxsBefore,
                accBefore
            );


            uint numDelayedTxsBefore = 0;
            uint numDelayedTxsAfter = 0;
            bytes memory delayedAccBefore = bytes("");

            bytes memory delayedInfo = abi.encodePacked(
                numDelayedTxsBefore,
                numDelayedTxsAfter,
                delayedAccBefore
            );

            bytes memory proof = abi.encodePacked(
                txContextHash,
                batchInfo,
                delayedInfo,
                batchTransactionsMetadata
            );

            console.logBytes(proof);
            vm.prank(alice); // Alice wants to verify if the 1st transaction is indeed included or not
            seqIn.verifyTxInclusion(
                "0x1d96f2f6bef1202e4ce1ff6dad0c2cb002861d3e00000000000000010000000000000001000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000e07b72d05562e093c16c044479393196ea4a7a3b37b0228af92431ca6638fb7",
                proof
            );


        }






    }

    /////////////////////////////////
    // verifyTxInclusion
    /////////////////////////////////

    // List of variables used inside of the test: test_verifyTxInclusion_positiveCase_1. Declared globally so as to avoid the *stack too deep* error while compilation.
    // address txnSender;
    // uint256 txnL2BlockNumber;
    // uint256 txnL2Timestamp;
    // uint256 txnDataLength;
    // bytes txnDataBytes;

    // uint256 batchNum;
    // uint256 numTxsBefore;
    // uint256 numTxsAfterInBatch;
    // bytes accBefore;

    // bytes encodedTxn; // This is the transaction whose inclusion we want to check
    // uint currentBlock;

    // uint256 inboxSizeInitial;
    // uint256 numTxns;
    // uint256 numContextsArrEntries;
    // uint256 inboxSizeFinal;
    // uint256 expectedInboxSize;

    // function test_verifyTxInclusion_firstTransactionIncluded_positiveCase_1(uint256 numTxnsPerBlock, uint256 txnBlocks) public {
    //     // We have to test this function's correctness by passing a correct bytes `proof` to ensure
    //     // that the function works correctly.

    //     // For the particular transaction whose inclusion needs to be verified, among others, we need the following data
    //     /**
    //         proof := txContextHash || batchInfo || {foreach tx in batch: (txContextHash || KEC(txData)), ...} where,
    //         * {foreach tx in batch: (txContextHash || KEC(txData)), ...} = txContextHash1 || KEC(txData1) || ... || txContextHash10 || KEC(txData10) if you have 10 txs ✅
    //         * batchInfo := (batchNum || numTxsBefore || numTxsAfterInBatch || accBefore) ✅
    //         * txContextHash := KEC(sequencerAddress || l2BlockNumber || l2Timestamp) ✅
    //     */

    //     // Let us first create a batch, append that txn batch and then check for the inclusion of a particular txn

    //     //////////////////////////////
    //     // Appending a txn batch
    //     //////////////////////////////

    //     // We will operate at a limit of transactionsPerBlock = 30 and number of transactionBlocks = 10.

    //     // the ^ comment is incorrect, as per the code the txsPerBlock = 5 and the blocks = 2
    //     numTxnsPerBlock = bound(numTxnsPerBlock, 1, 5);
    //     txnBlocks = bound(txnBlocks, 1, 2);

    //     inboxSizeInitial = seqIn.getInboxSize();

    //     // Each context corresponds to a single "L2 block"
    //     numTxns = numTxnsPerBlock * txnBlocks;
    //     numContextsArrEntries = 3 * txnBlocks; // Since each `context` is represented with uint256 triplet: (numTxs, l2BlockNumber, l2Timestamp)

    //     // Making sure that the block.timestamp is a reasonable value (> txnBlocks)
    //     vm.warp(block.timestamp + (4 * txnBlocks));
    //     uint256 txnBlockTimestamp = block.timestamp - (2 * txnBlocks); // Subtracing just `txnBlocks` would have sufficed. However we are subtracting 2 times txnBlocks for some margin of error.
    //         // The objective for this subtraction is that while building the `contexts` array, no timestamp should go higher than the current block.timestamp

    //     // Let's create an array of contexts
    //     uint256[] memory contexts = new uint256[](numContextsArrEntries);
    //     for (uint256 i; i < numContextsArrEntries; i += 3) {
    //         // The first entry for `contexts` for each txnBlock is `numTxns` which we are keeping as constant for all blocks for this test
    //         contexts[i] = numTxnsPerBlock;

    //         // Formual Used for blockNumber: (txnBlock's block.timestamp) / 20;
    //         contexts[i + 1] = txnBlockTimestamp / 20;

    //         // Formula used for blockTimestamp: (current block.timestamp) / 5x
    //         contexts[i + 2] = txnBlockTimestamp;

    //         // The only requirement for timestamps for the transaction blocks is that, these timestamps are monotonically increasing.
    //         // So, let's increase the value of txnBlock's timestamp monotonically, in a way that is does not exceed current block.timestamp
    //         ++txnBlockTimestamp;
    //     }

    //     bytes memory batchTransactionsMetadata;

    //     {
    //         // Let's construct the blockTransactionsHashArray and the txnDataHashArray
    //         bytes32[] memory blockTransactionsHashArray = new bytes32[](numTxns);
    //         bytes32[] memory txnDataHashArray = new bytes32[](numTxns);

    //         // There will be numTxns number of entries, ie, one entry for one transaction and we know: numTxns = numTxnsPerBlock * txnBlocks
    //         // And every txnBlock will have the same block.timestamp and block.number for each of it's own transactions
    //         // And that is how we will populate this array, so that it contains the block.number and block.timestamp for each transaction that will be used.
    //         for(uint i; i < numTxns; i++) {
    //             ( , txnDataBytes, , ) = _helper_getTxnInfo_fromTxnID(i % 10); // Since we have only 10 sample RLPEncodedTransactions (we could increase them in the future)
    //             txnDataHashArray[i] = keccak256(txnDataBytes);

    //             currentBlock = i / numTxnsPerBlock;

    //             blockTransactionsHashArray[i] = keccak256 (
    //                 abi.encodePacked(
    //                     sequencerAddress,
    //                     contexts[currentBlock + 1],
    //                     contexts[currentBlock + 2]
    //                 )
    //             );

    //             // Since we are trying to verify the presence of the first transaction and the proof
    //             // should include the txns after the tx you’re proving inclusion for
    //             if(i != 0) {
    //                 batchTransactionsMetadata = bytes.concat(
    //                     batchTransactionsMetadata,
    //                     abi.encodePacked(
    //                         blockTransactionsHashArray[i],
    //                         txnDataHashArray[i]
    //                     )
    //                 );
    //             }
    //         }

    //         assertEq(blockTransactionsHashArray.length, numTxns, "blockTransactionsHashArray should have the txContextHash for every transaction of the block");
    //     }

    //     // txLengths is defined as: Array of lengths of each encoded tx in txBatch
    //     // txBatch is defined as: Batch of RLP-encoded transactions
    //     (bytes memory txBatch, uint256[] memory txLengths) = _helper_sequencerInbox_appendTx(numTxns);

    //     // Pranking as the sequencer and calling appendTxBatch
    //     vm.prank(sequencerAddress);
    //     seqIn.appendTxBatch(contexts, txLengths, txBatch, 0);

    //     inboxSizeFinal = seqIn.getInboxSize();
    //     assertGt(inboxSizeFinal, inboxSizeInitial);

    //     expectedInboxSize = numTxns;
    //     assertEq(inboxSizeFinal, expectedInboxSize);

    //     ////////////////////////////////////////
    //     // Verify Transaction Inclusion
    //     ///////////////////////////////////////

    //     // We will be verifying the inclusion of the very first transaction
    //     (encodedTxn , txnDataBytes, txnSender, txnDataLength) = _helper_getTxnInfo_fromTxnID(0);

    //     // Since, we will be verifying the inclusion of the very first transaction. blockNumber and timestamp can be fetched from 1 and 2 indices of the contexts array.
    //     txnL2BlockNumber = contexts[1];
    //     txnL2Timestamp = contexts[2];

    //     bytes32 txContextHash = keccak256(
    //         abi.encodePacked(
    //             sequencerAddress,
    //             txnL2BlockNumber,
    //             txnL2Timestamp
    //         )
    //     );

    //     batchNum = 0; //seqIn.accumulators(0) returns a non-zero value while seqIn.accumulators(1) reverts.
    //     numTxsBefore = 0; // since we are verifying the inclusion of the 1st transaction of the first block.
    //     numTxsAfterInBatch = numTxns - 1; // since we are verifying the inclusion of the 1st transaction of the first block and the total number of transactions are numTxnsPerBlock * txnBlocks
    //     accBefore = bytes(""); // again, since this is our first inclusion, acc before that would have been empty bytes too.

    //     bytes memory batchInfo = abi.encodePacked(
    //         batchNum,
    //         numTxsBefore,
    //         numTxsAfterInBatch,
    //         accBefore
    //     );

    //     bytes memory proof = abi.encodePacked(
    //         txContextHash,
    //         batchInfo,
    //         batchTransactionsMetadata
    //     );

    //     vm.prank(alice); // Alice wants to verify if the 1st transaction is indeed included or not
    //     seqIn.verifyTxInclusion(
    //         encodedTxn,
    //         proof
    //     );
    // }



    // function test_appendTxBatch_positiveCase_2_hardcoded() public {
    //     //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    //     // Here, we are assuming we have 2 transaction blocks with 3 transactions each (initial lower load hardcoded test)
    //     /////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    //     uint256 numTxnsPerBlock = 3;
    //     uint256 inboxSizeInitial = seqIn.getInboxSize();

    //     // Each context corresponds to a single "L2 block"
    //     // `contexts` is represented with uint256 3-tuple: (numTxs, l2BlockNumber, l2Timestamp)
    //     // Let's create an array of contexts
    //     uint256 numTxns = numTxnsPerBlock * 2;
    //     uint256 timeStamp1 = block.timestamp / 10;
    //     uint256 timeStamp2 = block.timestamp / 5;
    //     uint256 blockNumber1 = timeStamp1 / 20;
    //     uint256 blockNumber2 = timeStamp2 / 20;

    //     uint256[] memory contexts = new uint256[](6);

    //     // Let's assume that we had 2 blocks and each had 3 transactions
    //     contexts[0] = (numTxnsPerBlock);
    //     contexts[1] = (blockNumber1);
    //     contexts[2] = (timeStamp1);
    //     contexts[3] = (numTxnsPerBlock);
    //     contexts[4] = (blockNumber2);
    //     contexts[5] = (timeStamp2);

    //     // txLengths is defined as: Array of lengths of each encoded tx in txBatch
    //     // txBatch is defined as: Batch of RLP-encoded transactions
    //     bytes memory txBatch = _helper_createTxBatch_hardcoded();
    //     uint256[] memory txLengths = _helper_findTxLength_hardcoded();

    //     // Pranking as the sequencer and calling appendTxBatch
    //     vm.prank(sequencerAddress);
    //     seqIn.appendTxBatch(contexts, txLengths, txBatch);

    //     uint256 inboxSizeFinal = seqIn.getInboxSize();

    //     assertGt(inboxSizeFinal, inboxSizeInitial);

    //     uint256 expectedInboxSize = numTxns;
    //     assertEq(inboxSizeFinal, expectedInboxSize);
    // }

    // function test_appendTxBatch_revert_txBatchDataOverflow(uint256 numTxnsPerBlock, uint256 txnBlocks) public {
    //     // We will operate at a limit of transactionsPerBlock = 30 and number of transactionBlocks = 10.
    //     numTxnsPerBlock = bound(numTxnsPerBlock, 1, 30);
    //     txnBlocks = bound(txnBlocks, 1, 10);

    //     // Each context corresponds to a single "L2 block"
    //     uint256 numTxns = numTxnsPerBlock * txnBlocks;
    //     uint256 numContextsArrEntries = 3 * txnBlocks; // Since each `context` is represented with uint256 3-tuple: (numTxs, l2BlockNumber, l2Timestamp)

    //     // Making sure that the block.timestamp is a reasonable value (> txnBlocks)
    //     vm.warp(block.timestamp + (4 * txnBlocks));
    //     uint256 txnBlockTimestamp = block.timestamp - (2 * txnBlocks); // Subtracing just `txnBlocks` would have sufficed. However we are subtracting 2 times txnBlocks for some margin of error.
    //         // The objective for this subtraction is that while building the `contexts` array, no timestamp should go higher than the current block.timestamp

    //     // Let's create an array of contexts
    //     uint256[] memory contexts = new uint256[](numContextsArrEntries);
    //     for (uint256 i; i < numContextsArrEntries; i += 3) {
    //         // The first entry for `contexts` for each txnBlock is `numTxns` which we are keeping as constant for all blocks for this test
    //         contexts[i] = numTxnsPerBlock;

    //         // Formual Used for blockNumber: (txnBlock's block.timestamp) / 20;
    //         contexts[i + 1] = txnBlockTimestamp / 20;

    //         // Formula used for blockTimestamp: (current block.timestamp) / 5x
    //         contexts[i + 2] = txnBlockTimestamp;

    //         // The only requirement for timestamps for the transaction blocks is that, these timestamps are monotonically increasing.
    //         // So, let's increase the value of txnBlock's timestamp monotonically, in a way that is does not exceed current block.timestamp
    //         ++txnBlockTimestamp;
    //     }

    //     // txLengths is defined as: Array of lengths of each encoded tx in txBatch
    //     // txBatch is defined as: Batch of RLP-encoded transactions
    //     (bytes memory txBatch, uint256[] memory txLengths) = _helper_sequencerInbox_appendTx(numTxns);

    //     // Now, we want to trigger the `txnBatchDataOverflow`, so we want to disturn the values receieved in the txLengths array.
    //     for (uint256 i; i < numTxns; i++) {
    //         txLengths[i] = txLengths[i] + 1;
    //     }

    //     // Pranking as the sequencer and calling appendTxBatch (should throw the TxBatchDataOverflow error)
    //     vm.expectRevert(ISequencerInbox.TxBatchDataOverflow.selector);
    //     vm.prank(sequencerAddress);
    //     seqIn.appendTxBatch(contexts, txLengths, txBatch);
    // }

    // /////////////////////////////////////////////////////////////////////////////////////////
    // // TENTATIVE CODE CHANGE. TEST SUBJECT TO CHANGE BASED ON CHANGE IN CODE
    // // Probable change: `appendTxBatch` passes even with malformed `contexts` array
    // /////////////////////////////////////////////////////////////////////////////////////////
    // function test_appendTxBatch_incompleteDataInContextsArray(uint256 numTxnsPerBlock) public {
    //     // Since we are assuming that we will have two transaction blocks and we have a total of 300 sample transactions right now.
    //     numTxnsPerBlock = bound(numTxnsPerBlock, 1, 150);
    //     uint256 inboxSizeInitial = seqIn.getInboxSize();

    //     // Each context corresponds to a single "L2 block"
    //     // `contexts` is represented with uint256 3-tuple: (numTxs, l2BlockNumber, l2Timestamp)
    //     // Let's create an array of contexts
    //     uint256 numTxns = numTxnsPerBlock * 2;
    //     uint256 timeStamp1 = block.timestamp / 10;
    //     uint256 blockNumber1 = timeStamp1 / 20;

    //     uint256[] memory contexts = new uint256[](4);

    //     // Let's assume that we had 2 blocks and each had 3 transactions, but we fail to pass the block.timestamp and block.number of the 2nd transaction block.
    //     contexts[0] = (numTxnsPerBlock);
    //     contexts[1] = (blockNumber1);
    //     contexts[2] = (timeStamp1);
    //     contexts[3] = (numTxnsPerBlock);

    //     // txLengths is defined as: Array of lengths of each encoded tx in txBatch
    //     // txBatch is defined as: Batch of RLP-encoded transactions
    //     (bytes memory txBatch, uint256[] memory txLengths) = _helper_sequencerInbox_appendTx(numTxns);

    //     // Pranking as the sequencer and calling appendTxBatch
    //     vm.prank(sequencerAddress);
    //     seqIn.appendTxBatch(contexts, txLengths, txBatch);

    //     uint256 inboxSizeFinal = seqIn.getInboxSize();

    //     assertGt(inboxSizeFinal, inboxSizeInitial);
    //     assertEq(inboxSizeFinal, numTxnsPerBlock); // Since the timestamp and block.number were not included for the 2nd block, only 1st block's 3 txns are included.
    // }


    // /////////////////////////////////
    // // verifyTxInclusion
    // /////////////////////////////////

    // // List of variables used inside of the test: test_verifyTxInclusion_positiveCase_1. Declared globally so as to avoid the *stack too deep* error while compilation.
    // address txnSender;
    // uint256 txnL2BlockNumber;
    // uint256 txnL2Timestamp;
    // uint256 txnDataLength;
    // bytes txnDataBytes;

    // uint256 batchNum;
    // uint256 numTxsBefore;
    // uint256 numTxsAfterInBatch;
    // bytes accBefore;

    // bytes encodedTxn; // This is the transaction whose inclusion we want to check
    // uint currentBlock;

    // uint256 inboxSizeInitial;
    // uint256 numTxns;
    // uint256 numContextsArrEntries;
    // uint256 inboxSizeFinal;
    // uint256 expectedInboxSize;

    // function test_verifyTxInclusion_firstTransactionIncluded_positiveCase_1(uint256 numTxnsPerBlock, uint256 txnBlocks) public {
    //     // We have to test this function's correctness by passing a correct bytes `proof` to ensure
    //     // that the function works correctly.

    //     // For the particular transaction whose inclusion needs to be verified, among others, we need the following data
    //     /**
    //         proof := txContextHash || batchInfo || {foreach tx in batch: (txContextHash || KEC(txData)), ...} where,
    //         * {foreach tx in batch: (txContextHash || KEC(txData)), ...} = txContextHash1 || KEC(txData1) || ... || txContextHash10 || KEC(txData10) if you have 10 txs ✅
    //         * batchInfo := (batchNum || numTxsBefore || numTxsAfterInBatch || accBefore) ✅
    //         * txContextHash := KEC(sequencerAddress || l2BlockNumber || l2Timestamp) ✅
    //     */

    //     // Let us first create a batch, append that txn batch and then check for the inclusion of a particular txn

    //     //////////////////////////////
    //     // Appending a txn batch
    //     //////////////////////////////

    //     // We will operate at a limit of transactionsPerBlock = 30 and number of transactionBlocks = 10.
    //     numTxnsPerBlock = bound(numTxnsPerBlock, 1, 5);
    //     txnBlocks = bound(txnBlocks, 1, 2);

    //     inboxSizeInitial = seqIn.getInboxSize();

    //     // Each context corresponds to a single "L2 block"
    //     numTxns = numTxnsPerBlock * txnBlocks;
    //     numContextsArrEntries = 3 * txnBlocks; // Since each `context` is represented with uint256 3-tuple: (numTxs, l2BlockNumber, l2Timestamp)

    //     // Making sure that the block.timestamp is a reasonable value (> txnBlocks)
    //     vm.warp(block.timestamp + (4 * txnBlocks));
    //     uint256 txnBlockTimestamp = block.timestamp - (2 * txnBlocks); // Subtracing just `txnBlocks` would have sufficed. However we are subtracting 2 times txnBlocks for some margin of error.
    //         // The objective for this subtraction is that while building the `contexts` array, no timestamp should go higher than the current block.timestamp

    //     // Let's create an array of contexts
    //     uint256[] memory contexts = new uint256[](numContextsArrEntries);
    //     for (uint256 i; i < numContextsArrEntries; i += 3) {
    //         // The first entry for `contexts` for each txnBlock is `numTxns` which we are keeping as constant for all blocks for this test
    //         contexts[i] = numTxnsPerBlock;

    //         // Formual Used for blockNumber: (txnBlock's block.timestamp) / 20;
    //         contexts[i + 1] = txnBlockTimestamp / 20;

    //         // Formula used for blockTimestamp: (current block.timestamp) / 5x
    //         contexts[i + 2] = txnBlockTimestamp;

    //         // The only requirement for timestamps for the transaction blocks is that, these timestamps are monotonically increasing.
    //         // So, let's increase the value of txnBlock's timestamp monotonically, in a way that is does not exceed current block.timestamp
    //         ++txnBlockTimestamp;
    //     }

    //     bytes memory batchTransactionsMetadata;

    //     {
    //         // Let's construct the blockTransactionsHashArray and the txnDataHashArray
    //         bytes32[] memory blockTransactionsHashArray = new bytes32[](numTxns);
    //         bytes32[] memory txnDataHashArray = new bytes32[](numTxns);

    //         // There will be numTxns number of entries, ie, one entry for one transaction and we know: numTxns = numTxnsPerBlock * txnBlocks
    //         // And every txnBlock will have the same block.timestamp and block.number for each of it's own transactions
    //         // And that is how we will populate this array, so that it contains the block.number and block.timestamp for each transaction that will be used.
    //         for(uint i; i < numTxns; i++) {
    //             ( , txnDataBytes, , ) = _helper_getTxnInfo_fromTxnID(i % 10); // Since we have only 10 sample RLPEncodedTransactions (we could increase them in the future)
    //             txnDataHashArray[i] = keccak256(txnDataBytes);

    //             currentBlock = i / numTxnsPerBlock;

    //             blockTransactionsHashArray[i] = keccak256 (
    //                 abi.encodePacked(
    //                     sequencerAddress,
    //                     contexts[currentBlock + 1],
    //                     contexts[currentBlock + 2]
    //                 )
    //             );

    //             // Since we are trying to verify the presence of the first transaction and the proof
    //             // should include the txns after the tx you’re proving inclusion for
    //             if(i != 0) {
    //                 batchTransactionsMetadata = bytes.concat(
    //                     batchTransactionsMetadata,
    //                     abi.encodePacked(
    //                         blockTransactionsHashArray[i],
    //                         txnDataHashArray[i]
    //                     )
    //                 );
    //             }
    //         }

    //         assertEq(blockTransactionsHashArray.length, numTxns, "blockTransactionsHashArray should have the txContextHash for every transaction of the block");
    //     }

    //     // txLengths is defined as: Array of lengths of each encoded tx in txBatch
    //     // txBatch is defined as: Batch of RLP-encoded transactions
    //     (bytes memory txBatch, uint256[] memory txLengths) = _helper_sequencerInbox_appendTx(numTxns);

    //     // Pranking as the sequencer and calling appendTxBatch
    //     vm.prank(sequencerAddress);
    //     seqIn.appendTxBatch(contexts, txLengths, txBatch, 0);

    //     inboxSizeFinal = seqIn.getInboxSize();
    //     assertGt(inboxSizeFinal, inboxSizeInitial);

    //     expectedInboxSize = numTxns;
    //     assertEq(inboxSizeFinal, expectedInboxSize);

    //     ////////////////////////////////////////
    //     // Verify Transaction Inclusion
    //     ///////////////////////////////////////

    //     // We will be verifying the inclusion of the very first transaction
    //     (encodedTxn , txnDataBytes, txnSender, txnDataLength) = _helper_getTxnInfo_fromTxnID(0);

    //     // Since, we will be verifying the inclusion of the very first transaction. blockNumber and timestamp can be fetched from 1 and 2 indices of the contexts array.
    //     txnL2BlockNumber = contexts[1];
    //     txnL2Timestamp = contexts[2];

    //     bytes32 txContextHash = keccak256(
    //         abi.encodePacked(
    //             sequencerAddress,
    //             txnL2BlockNumber,
    //             txnL2Timestamp
    //         )
    //     );

    //     batchNum = 0; //seqIn.accumulators(0) returns a non-zero value while seqIn.accumulators(1) reverts.
    //     numTxsBefore = 0; // since we are verifying the inclusion of the 1st transaction of the first block.
    //     numTxsAfterInBatch = numTxns - 1; // since we are verifying the inclusion of the 1st transaction of the first block and the total number of transactions are numTxnsPerBlock * txnBlocks
    //     accBefore = bytes(""); // again, since this is our first inclusion, acc before that would have been empty bytes too.

    //     bytes memory batchInfo = abi.encodePacked(
    //         batchNum,
    //         numTxsBefore,
    //         numTxsAfterInBatch,
    //         accBefore
    //     );

    //     bytes memory proof = abi.encodePacked(
    //         txContextHash,
    //         batchInfo,
    //         batchTransactionsMetadata
    //     );

    //     vm.prank(alice); // Alice wants to verify if the 1st transaction is indeed included or not
    //     seqIn.verifyTxInclusion(
    //         encodedTxn,
    //         proof
    //     );
    // }
}
