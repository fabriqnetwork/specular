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
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "../src/ISequencerInbox.sol";
import "../src/libraries/Errors.sol";
import {SpecularProxy} from "test/SpecularProxy.sol";
import {MockToken} from "./utils/MockToken.sol";
import {Utils} from "./utils/Utils.sol";
import {IRollup} from "../src/IRollup.sol";
import {Verifier} from "../src/challenge/verifier/Verifier.sol";
import {Rollup} from "../src/Rollup.sol";
import {SequencerInbox} from "../src/SequencerInbox.sol";
import {RLPEncodedTransactionsUtil} from "./utils/RLPEncodedTransactions.sol";

contract RollupBaseSetup is Test, RLPEncodedTransactionsUtil {
    Utils internal utils;
    address payable[] internal users;

    address internal sequencer;
    address internal alice;
    address internal bob;
    address internal vault;
    address internal sequencerOwner;
    address internal sequencerAddress;
    address internal rollupOwner;

    Verifier verifier = new Verifier();

    function setUp() public virtual {
        utils = new Utils();
        users = utils.createUsers(3);

        sequencer = users[0];
        vm.label(sequencer, "Sequencer");

        alice = users[1];
        vm.label(alice, "Alice");

        bob = users[2];
        vm.label(bob, "Bob");

        vault = makeAddr("vault");
        sequencerOwner = makeAddr("sequencerOwner");
        sequencerAddress = makeAddr("sequencerAddress");
        rollupOwner = makeAddr("rollupOwner");
    }
}

contract RollupTest is RollupBaseSetup {
    Rollup public rollup;
    uint256 randomNonce;
    SpecularProxy public specularProxy;
    SequencerInbox public seqIn;
    SequencerInbox public implementationSequencer;

    address defender = makeAddr("defender");
    address challenger = makeAddr("challenger");

    function setUp() public virtual override {
        // Parent contract setup
        RollupBaseSetup.setUp();

        // Deploying the SequencerInbox with sequencerOwner
        bytes memory seqInInitData = abi.encodeWithSignature("initialize(address)", sequencerAddress);
        vm.startPrank(sequencerOwner);
        implementationSequencer = new SequencerInbox();
        specularProxy = new SpecularProxy(address(implementationSequencer), seqInInitData);
        seqIn = SequencerInbox(address(specularProxy));
        vm.stopPrank();

        // Making sure that the proxy returns the correct proxy owner and sequencerAddress
        address sequencerInboxOwner = seqIn.owner();
        assertEq(sequencerInboxOwner, sequencerOwner);

        address fetchedSequencerAddress = seqIn.sequencerAddress();
        assertEq(fetchedSequencerAddress, sequencerAddress);
    }

    /////////////////////////
    // Challenge Assertion
    /////////////////////////

    function test_challengeAssertion_wrongOrderOfAssertionIDs(
        uint256 confirmationPeriod,
        uint256 challengePeriod,
        uint256 minimumAssertionPeriod,
        uint256 maxGasPerAssertion,
        uint256 defenderAssertionID,
        uint256 challengerAssertionID,
        uint256 initialAssertionID,
        uint256 initialInboxSize,
        uint256 initialL2GasUsed
    ) public {
        // Initializing the rollup
        address _vault = makeAddr("vault");
        _initializeRollup(
            _vault,
            confirmationPeriod,
            challengePeriod,
            minimumAssertionPeriod,
            maxGasPerAssertion,
            type(uint256).max,
            initialAssertionID,
            initialInboxSize,
            initialL2GasUsed
        );

        defenderAssertionID = bound(defenderAssertionID, challengerAssertionID, type(uint256).max);

        address[2] memory players;
        uint256[2] memory assertionIDs;

        players[0] = defender;
        players[1] = challenger;

        assertionIDs[0] = defenderAssertionID;
        assertionIDs[1] = challengerAssertionID;

        vm.expectRevert(IRollup.WrongOrder.selector);
        rollup.challengeAssertion(players, assertionIDs);
    }

    function test_challengeAssertion_unproposedAssertionID(
        uint256 confirmationPeriod,
        uint256 challengePeriod,
        uint256 minimumAssertionPeriod,
        uint256 maxGasPerAssertion,
        uint256 initialAssertionID,
        uint256 initialInboxSize,
        uint256 initialL2GasUsed
    ) public {
        // Initializing the rollup
        address _vault = makeAddr("vault");
        initialAssertionID = bound(initialAssertionID, 0, (type(uint256).max - 10));
        _initializeRollup(
            _vault,
            confirmationPeriod,
            challengePeriod,
            minimumAssertionPeriod,
            maxGasPerAssertion,
            type(uint256).max,
            initialAssertionID,
            initialInboxSize,
            initialL2GasUsed
        );

        uint256 lastCreatedAssertionID = rollup.lastCreatedAssertionID();
        uint256 challengerAssertionID;
        uint256 defenderAssertionID;

        challengerAssertionID = bound(challengerAssertionID, lastCreatedAssertionID + 1, type(uint256).max);
        defenderAssertionID = bound(defenderAssertionID, 0, challengerAssertionID - 1);

        address[2] memory players;
        uint256[2] memory assertionIDs;

        players[0] = defender;
        players[1] = challenger;

        assertionIDs[0] = defenderAssertionID;
        assertionIDs[1] = challengerAssertionID;

        vm.expectRevert(IRollup.UnproposedAssertion.selector);
        rollup.challengeAssertion(players, assertionIDs);
    }

    function test_challengeAssertion_alreadyResolvedAssertionID(
        uint256 confirmationPeriod,
        uint256 challengePeriod,
        uint256 initialAssertionID,
        uint256 initialInboxSize,
        uint256 initialL2GasUsed
    ) public {
        // Initializing the rollup
        confirmationPeriod = bound(confirmationPeriod, 1, type(uint128).max);
        address _vault = makeAddr("vault");
        _initializeRollup(_vault, confirmationPeriod, challengePeriod, 1 days, 500, 1 ether, 0, 5, 0);

        uint256 lastConfirmedAssertionID = rollup.lastConfirmedAssertionID();

        // Let's increase the lastCreatedAssertionID
        {
            // Alice has not staked yet and therefore, this function should return `false`
            (bool isAliceStaked,,,) = rollup.stakers(alice);
            assertTrue(!isAliceStaked);

            uint256 minimumAmount = rollup.baseStakeAmount();
            uint256 aliceBalance = alice.balance;

            // Let's stake something on behalf of Alice
            uint256 aliceAmountToStake = minimumAmount * 10;

            vm.prank(alice);
            require(aliceBalance >= aliceAmountToStake, "Increase balance of Alice to proceed");

            // Calling the staking function as Alice
            //slither-disable-next-line arbitrary-send-eth
            rollup.stake{value: aliceAmountToStake}();

            // Now Alice should be staked
            uint256 stakerAssertionID;

            // stakers mapping gets updated
            (isAliceStaked,, stakerAssertionID,) = rollup.stakers(alice);
            assertTrue(isAliceStaked);

            // Checking previous Sequencer Inbox Size
            uint256 seqInboxSize = seqIn.getInboxSize();
            emit log_named_uint("Sequencer Inbox Size", seqInboxSize);

            _increaseSequencerInboxSize();

            seqInboxSize = seqIn.getInboxSize();
            emit log_named_uint("Sequencer Inbox Size Final", seqInboxSize);

            bytes32 mockVmHash = bytes32("");
            uint256 mockInboxSize = 6;
            uint256 mockL2GasUsed = 342;
            bytes32 mockPrevVMHash = bytes32("");
            uint256 mockPrevL2GasUsed = 0;

            // To avoid the MinimumAssertionPeriodNotPassed error, increase block.number
            vm.warp(block.timestamp + 50 days);
            vm.roll(block.number + (50 * 86400) / 20);

            assertEq(rollup.lastCreatedAssertionID(), 0, "The lastCreatedAssertionID should be 0 (genesis)");
            (,, uint256 assertionIDInitial,) = rollup.stakers(address(alice));

            assertEq(assertionIDInitial, 0);

            vm.prank(alice);
            rollup.createAssertion(mockVmHash, mockInboxSize, mockL2GasUsed, mockPrevVMHash, mockPrevL2GasUsed);

            // The assertionID of alice should change after she called `createAssertion`
            (,, uint256 assertionIDFinal,) = rollup.stakers(address(alice));

            assertEq(assertionIDFinal, 1); // Alice is now staked on assertionID = 1 instead of assertionID = 0.
        }

        uint256 defenderAssertionID = lastConfirmedAssertionID; //would be 0 in this case. cannot assign anything lower
        uint256 challengerAssertionID = lastConfirmedAssertionID + 1; // that would mean 1

        address[2] memory players;
        uint256[2] memory assertionIDs;

        players[0] = defender;
        players[1] = challenger;

        assertionIDs[0] = defenderAssertionID;
        assertionIDs[1] = challengerAssertionID;

        vm.expectRevert(IRollup.AssertionAlreadyResolved.selector);
        rollup.challengeAssertion(players, assertionIDs);

        assertTrue(true);
    }

    /////////////////////////
    // Auxillary Functions
    /////////////////////////

    function checkRange(uint256 _lower, uint256 _upper, uint256 _random) external {
        uint256 test = _generateRandomUintInRange(_lower, _upper, _random);

        require(test >= _lower && test <= _upper, "Cheat didn't work as expected");
        assertEq(uint256(2), uint256(2));
    }

    function _generateRandomUintInRange(uint256 _lower, uint256 _upper, uint256 randomUint)
        internal
        view
        returns (uint256)
    {
        uint256 boundedUint = bound(randomUint, _lower, _upper);
        return boundedUint;
    }

    // This function increases the inbox size by 6
    function _increaseSequencerInboxSize() internal {
        uint256 numTxnsPerBlock = 3;

        // Each context corresponds to a single "L2 block"
        // `contexts` is represented with uint256 3-tuple: (numTxs, l2BlockNumber, l2Timestamp)
        // Let's create an array of contexts
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
    }

    function _initializeRollup(
        address _vault,
        uint256 confirmationPeriod,
        uint256 challengePeriod,
        uint256 minimumAssertionPeriod,
        uint256 maxGasPerAssertion,
        uint256 baseStakeAmount,
        uint256 initialAssertionID,
        uint256 initialInboxSize,
        uint256 initialL2GasUsed
    ) internal {
        if (_vault == address(0)) {
            _vault = address(uint160(32456));
        }
        bytes memory initializingData = abi.encodeWithSelector(
            Rollup.initialize.selector,
            _vault,
            address(seqIn), // sequencerInbox
            address(verifier),
            confirmationPeriod, //confirmationPeriod
            challengePeriod, //challengePeriod
            minimumAssertionPeriod, // minimumAssertionPeriod
            maxGasPerAssertion, // maxGasPerAssertion
            baseStakeAmount, //baseStakeAmount
            initialAssertionID,
            initialInboxSize,
            bytes32(""), //initialVMHash
            initialL2GasUsed
        );

        // Deploying the rollup contract as the rollup owner/deployer
        vm.startPrank(rollupOwner);
        Rollup implementationRollup = new Rollup();
        specularProxy = new SpecularProxy(address(implementationRollup), initializingData);
        rollup = Rollup(address(specularProxy));
        vm.stopPrank();
    }
}
