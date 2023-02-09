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
import {Utils} from "./utils/Utils.sol";
import {MockToken} from "./utils/MockToken.sol";
import {IRollup} from "../src/IRollup.sol";
import {Verifier} from "../src/challenge/verifier/Verifier.sol";
import {Rollup} from "../src/Rollup.sol";
import {AssertionMap} from "../src/AssertionMap.sol";
import {SequencerInbox} from "../src/SequencerInbox.sol";

contract RollupBaseSetup is Test {
    Utils internal utils;
    address payable[] internal users;

    address internal sequencer;
    address internal alice;
    address internal bob;

    address owner = makeAddr("Owner");
    Verifier verifier = new Verifier();

    IERC20 public stakeToken;

    function setUp() public virtual {
        utils = new Utils();
        users = utils.createUsers(3);

        sequencer = users[0];
        vm.label(sequencer, "Sequencer");

        alice = users[1];
        vm.label(alice, "Alice");

        bob = users[2];
        vm.label(bob, "Bob");

        stakeToken = new MockToken (
                            "Stake Token",
                            "SPEC",
                            1e40,
                            address(owner)
                        );
    }
}

contract RollupTest is RollupBaseSetup {
    SequencerInbox private seqIn;
    Rollup private rollup;
    uint256 randomNonce;
    AssertionMap rollupAssertion;

    function setUp() public virtual override {
        RollupBaseSetup.setUp();

        SequencerInbox _impl = new SequencerInbox();
        bytes memory data = abi.encodeWithSelector(SequencerInbox.initialize.selector, sequencer);
        address admin = address(47); // Random admin
        TransparentUpgradeableProxy proxy = new TransparentUpgradeableProxy(address(_impl), admin, data);

        seqIn = SequencerInbox(address(proxy));
    }

    function test_initializeRollup_ownerAddressZero() external {
        Rollup _tempRollup = new Rollup();
        bytes memory initializingData = abi.encodeWithSelector(
            Rollup.initialize.selector,
            address(0), // owner
            address(seqIn),
            address(verifier),
            address(stakeToken),
            0, //confirmationPeriod
            0, //challengPeriod
            0, // minimumAssertionPeriod
            type(uint256).max, // maxGasPerAssertion
            0, //baseStakeAmount
            bytes32("")
        );

        vm.expectRevert(ZeroAddress.selector);

        address proxyAdmin = makeAddr("Proxy Admin");
        TransparentUpgradeableProxy proxy = new TransparentUpgradeableProxy(
            address(_tempRollup), 
            proxyAdmin, 
            initializingData
        );

        rollup = Rollup(address(proxy));
    }

    function test_initializeRollup_verifierAddressZero() external {
        emit log_named_address("Stake Token", address(stakeToken));

        Rollup _tempRollup = new Rollup();
        bytes memory initializingData = abi.encodeWithSelector(
            Rollup.initialize.selector,
            owner, // owner
            address(seqIn),
            address(0),
            address(stakeToken),
            0, //confirmationPeriod
            0, //challengPeriod
            0, // minimumAssertionPeriod
            type(uint256).max, // maxGasPerAssertion
            0, //baseStakeAmount
            bytes32("")
        );

        vm.expectRevert(ZeroAddress.selector);

        address proxyAdmin = makeAddr("Proxy Admin");
        TransparentUpgradeableProxy proxy = new TransparentUpgradeableProxy(
            address(_tempRollup), 
            proxyAdmin, 
            initializingData
        );

        rollup = Rollup(address(proxy));
    }

    function test_initializeRollup_sequencerInboxAddressZero() external {
        Rollup _tempRollup = new Rollup();
        bytes memory initializingData = abi.encodeWithSelector(
            Rollup.initialize.selector,
            owner, // owner
            address(0),
            address(verifier),
            address(stakeToken),
            0, //confirmationPeriod
            0, //challengPeriod
            0, // minimumAssertionPeriod
            type(uint256).max, // maxGasPerAssertion
            0, //baseStakeAmount
            bytes32("")
        );

        vm.expectRevert(ZeroAddress.selector);

        address proxyAdmin = makeAddr("Proxy Admin");
        TransparentUpgradeableProxy proxy = new TransparentUpgradeableProxy(
            address(_tempRollup), 
            proxyAdmin, 
            initializingData
        );

        rollup = Rollup(address(proxy));
    }

    function test_initializeRollup_cannotBeCalledTwice() external {
        Rollup _tempRollup = new Rollup();

        bytes memory initializingData = abi.encodeWithSelector(
            Rollup.initialize.selector,
            owner, // owner
            address(seqIn), // sequencerInbox
            address(verifier),
            address(stakeToken),
            0, //confirmationPeriod
            0, //challengPeriod
            0, // minimumAssertionPeriod
            type(uint256).max, // maxGasPerAssertion
            0, //baseStakeAmount
            bytes32("")
        );

        address proxyAdmin = makeAddr("Proxy Admin");

        TransparentUpgradeableProxy proxy = new TransparentUpgradeableProxy(
            address(_tempRollup), 
            proxyAdmin, 
            initializingData
        );

        // Initialize is called here for the first time.
        rollup = Rollup(address(proxy));

        // Trying to call initialize for the second time
        vm.expectRevert("Initializable: contract is already initialized");
        vm.prank(alice);

        rollup.initialize(
            owner, // owner
            address(seqIn), // sequencerInbox
            address(verifier),
            address(stakeToken),
            0, //confirmationPeriod
            0, //challengPeriod
            0, // minimumAssertionPeriod
            type(uint256).max, // maxGasPerAssertion
            0, //baseStakeAmount
            bytes32("")
        );
    }

    function test_initializeRollup_valuesAfterInit(
        uint256 confirmationPeriod,
        uint256 challengePeriod,
        uint256 minimumAssertionPeriod,
        uint256 maxGasPerAssertion,
        uint256 baseStakeAmount
    ) external {
        Rollup _tempRollup = new Rollup();

        /*
            emit log_named_uint("confirmationPeriod", confirmationPeriod);
            emit log_named_uint("CP", challengePeriod);
            emit log_named_uint("BSA", baseStakeAmount);
        */

        bytes memory initializingData = abi.encodeWithSelector(
            Rollup.initialize.selector,
            owner, // owner
            address(seqIn), // sequencerInbox
            address(verifier),
            address(stakeToken),
            confirmationPeriod, //confirmationPeriod
            challengePeriod, //challengePeriod
            minimumAssertionPeriod, // minimumAssertionPeriod
            maxGasPerAssertion, // maxGasPerAssertion
            baseStakeAmount, //baseStakeAmount
            bytes32("")
        );

        address proxyAdmin = makeAddr("Proxy Admin");

        TransparentUpgradeableProxy proxy = new TransparentUpgradeableProxy(
            address(_tempRollup), 
            proxyAdmin, 
            initializingData
        );

        // Initialize is called here for the first time.
        rollup = Rollup(address(proxy));

        // Putting in different scope to do away with the stack too deep error.
        {
            // Check if the value of the address owner was set correctly
            address rollupOwner = rollup.owner();
            assertEq(rollupOwner, owner, "Rollup.initialize failed to update owner correctly");

            // Check if the value of SequencerInbox was set correctly
            address rollupSeqIn = address(rollup.sequencerInbox());
            assertEq(rollupSeqIn, address(seqIn), "Rollup.initialize failed to update Sequencer Inbox correctly");

            // Check if the value of the stakeToken was set correctly
            address rollupToken = address(rollup.stakeToken());
            assertEq(rollupToken, address(stakeToken), "Rollup.initialize failed to update StakeToken value correctly");

            // Check if the value of the verifier was set correctly
            address rollupVerifier = address(rollup.verifier());
            assertEq(rollupVerifier, address(verifier), "Rollup.initialize failed to update verifier value correctly");
        }

        // Check if the various durations and uint values were set correctly
        uint256 rollupConfirmationPeriod = rollup.confirmationPeriod();
        uint256 rollupChallengePeriod = rollup.challengePeriod();
        uint256 rollupMinimumAssertionPeriod = rollup.minimumAssertionPeriod();
        uint256 rollupMaxGasPerAssertion = rollup.maxGasPerAssertion();
        uint256 rollupBaseStakeAmount = rollup.baseStakeAmount();

        assertEq(
            rollupConfirmationPeriod,
            confirmationPeriod,
            "Rollup.initialize failed to update confirmationPeriod value correctly"
        );
        assertEq(
            rollupChallengePeriod,
            challengePeriod,
            "Rollup.initialize failed to update confirmationPeriod value correctly"
        );
        assertEq(
            rollupMinimumAssertionPeriod,
            minimumAssertionPeriod,
            "Rollup.initialize failed to update confirmationPeriod value correctly"
        );
        assertEq(
            rollupMaxGasPerAssertion,
            maxGasPerAssertion,
            "Rollup.initialize failed to update confirmationPeriod value correctly"
        );
        assertEq(
            rollupBaseStakeAmount,
            baseStakeAmount,
            "Rollup.initialize failed to update confirmationPeriod value correctly"
        );

        // Make sure an assertion was created
        rollupAssertion = rollup.assertions();

        uint256 rollupAssertionParentID = rollupAssertion.getParentID(0);
        assertEq(rollupAssertionParentID, 0);

        // AssertionMap was created by the correct rollup address
        assertEq(rollupAssertion.rollupAddress(), address(rollup));
    }

    ////////////////
    // Staking
    ///////////////

    function test_stake_isStaked(
        uint256 confirmationPeriod,
        uint256 challengePeriod,
        uint256 minimumAssertionPeriod,
        uint256 maxGasPerAssertion,
        uint256 baseStakeAmount
    ) external {
        _initializeRollup(
            confirmationPeriod, challengePeriod, minimumAssertionPeriod, maxGasPerAssertion, baseStakeAmount
        );

        // Alice has not staked yet and therefore, this function should return `false`
        bool isAliceStaked = rollup.isStaked(alice);
        assertTrue(!isAliceStaked);
    }

    function test_stake_insufficentAmountStaking(
        uint256 confirmationPeriod,
        uint256 challengePeriod,
        uint256 minimumAssertionPeriod,
        uint256 maxGasPerAssertion,
        uint256 baseStakeAmount
    ) external {
        _initializeRollup(
            confirmationPeriod, challengePeriod, minimumAssertionPeriod, maxGasPerAssertion, type(uint256).max
        );

        uint256 minimumAmount = rollup.baseStakeAmount();
        uint256 aliceBalance = alice.balance;

        /*
            emit log_named_uint("BSA", minimumAmount);
        */

        if (aliceBalance > minimumAmount) {
            aliceBalance = minimumAmount / 10;
        }

        vm.expectRevert(IRollup.InsufficientStake.selector);

        vm.prank(alice);
        //slither-disable-next-line arbitrary-send-eth
        rollup.stake{value: aliceBalance}();

        bool isAliceStaked = rollup.isStaked(alice);
        assertTrue(!isAliceStaked);
    }

    function test_stake_sufficientAmountStakingAndNumStakersIncrement(
        uint256 confirmationPeriod,
        uint256 challengePeriod,
        uint256 minimumAssertionPeriod,
        uint256 maxGasPerAssertion,
        uint256 baseStakeAmount
    ) external {
        _initializeRollup(confirmationPeriod, challengePeriod, minimumAssertionPeriod, maxGasPerAssertion, 1000);

        uint256 initialStakers = rollup.numStakers();

        uint256 minimumAmount = rollup.baseStakeAmount();
        uint256 aliceBalance = alice.balance;

        /*
            emit log_named_uint("BSA", minimumAmount);
        */

        assertGt(aliceBalance, minimumAmount, "Alice's Balance should be greater than stake amount for this test");

        vm.prank(alice);
        //slither-disable-next-line arbitrary-send-eth
        rollup.stake{value: aliceBalance}();

        uint256 finalStakers = rollup.numStakers();

        assertEq(alice.balance, 0, "Alice should not have any balance left");
        assertEq(finalStakers, (initialStakers + 1), "Number of stakers should increase by 1");

        // isStaked should return true for Alice now
        bool isAliceStaked = rollup.isStaked(alice);
        assertTrue(isAliceStaked);

        uint256 amountStaked;
        uint256 assertionID;
        address challengeAddress;

        // stakers mapping gets updated
        (isAliceStaked, amountStaked, assertionID, challengeAddress) = rollup.stakers(alice);

        assertEq(amountStaked, aliceBalance, "amountStaked not updated properly");
        assertEq(assertionID, rollup.lastConfirmedAssertionID(), "assertionID not updated properly");
        assertEq(challengeAddress, address(0), "challengeAddress not updated properly");
    }

    function test_stake_increaseStake(
        uint256 confirmationPeriod,
        uint256 challengePeriod,
        uint256 minimumAssertionPeriod,
        uint256 maxGasPerAssertion,
        uint256 baseStakeAmount
    ) external {
        _initializeRollup(confirmationPeriod, challengePeriod, minimumAssertionPeriod, maxGasPerAssertion, 1000);

        uint256 minimumAmount = rollup.baseStakeAmount();
        uint256 aliceBalanceInitial = alice.balance;
        uint256 bobBalance = bob.balance;

        /*
            emit log_named_uint("BSA", minimumAmount);
        */

        assertGt(
            aliceBalanceInitial, minimumAmount, "Alice's Balance should be greater than stake amount for this test"
        );

        vm.prank(alice);
        //slither-disable-next-line arbitrary-send-eth
        rollup.stake{value: aliceBalanceInitial}();

        uint256 initialStakers = rollup.numStakers();

        uint256 amountStaked;
        uint256 assertionID;
        address challengeAddress;

        // isStaked should return true for Alice now
        bool isAliceStaked = rollup.isStaked(alice);
        assertTrue(isAliceStaked);

        // stakers mapping gets updated
        (isAliceStaked, amountStaked, assertionID, challengeAddress) = rollup.stakers(alice);

        uint256 aliceBalanceFinal = alice.balance;

        assertEq(alice.balance, 0, "Alice should not have any balance left");
        assertGt(bob.balance, 0, "Bob should have a non-zero native currency balance");

        vm.prank(bob);
        (bool sent, bytes memory data) = alice.call{value: bob.balance}("");
        require(sent, "Failed to send Ether");

        assertEq((aliceBalanceInitial - aliceBalanceFinal), bobBalance, "Tokens transferred successfully");

        vm.prank(alice);
        //slither-disable-next-line arbitrary-send-eth
        rollup.stake{value: alice.balance}();

        uint256 finalStakers = rollup.numStakers();

        uint256 amountStakedFinal;
        uint256 assertionIDFinal;
        address challengeAddressFinal;

        // stakers mapping gets updated (only the relevant values)
        (isAliceStaked, amountStakedFinal, assertionIDFinal, challengeAddressFinal) = rollup.stakers(alice);

        assertEq(challengeAddress, challengeAddressFinal, "Challenge Address should not change with more staking");
        assertEq(assertionID, assertionIDFinal, "Challenge Address should not change with more staking");
        assertEq(amountStakedFinal, (amountStaked + bobBalance), "Additional stake not updated correctly");
        assertEq(initialStakers, finalStakers, "Number of stakers should not increase");
    }

    ////////////////
    // Unstaking
    ///////////////

    function test_unstake_asANonStaker(
        uint256 randomAmount,
        uint256 confirmationPeriod,
        uint256 challengePeriod,
        uint256 minimumAssertionPeriod,
        uint256 maxGasPerAssertion,
        uint256 baseStakeAmount
    ) external {
        _initializeRollup(
            confirmationPeriod, challengePeriod, minimumAssertionPeriod, maxGasPerAssertion, baseStakeAmount
        );

        // Alice has not staked yet and therefore, this function should return `false`
        bool isAliceStaked = rollup.isStaked(alice);
        assertTrue(!isAliceStaked);

        // Since Alice is not staked, function unstake should also revert
        vm.expectRevert(IRollup.NotStaked.selector);
        vm.prank(alice);

        rollup.unstake(randomAmount);
    }

    function test_unstake_positiveCase(
        uint256 confirmationPeriod,
        uint256 challengePeriod,
        uint256 minimumAssertionPeriod,
        uint256 maxGasPerAssertion,
        uint256 baseStakeAmount,
        uint256 amountToWithdraw
    ) external {
        _initializeRollup(confirmationPeriod, challengePeriod, minimumAssertionPeriod, maxGasPerAssertion, 100000);

        // Alice has not staked yet and therefore, this function should return `false`
        bool isAliceStaked = rollup.isStaked(alice);
        assertTrue(!isAliceStaked);

        uint256 minimumAmount = rollup.baseStakeAmount();
        uint256 aliceBalance = alice.balance;

        emit log_named_uint("AB", aliceBalance);

        // Let's stake something on behalf of Alice
        uint256 aliceAmountToStake = minimumAmount * 10;

        vm.prank(alice);
        require(aliceBalance >= aliceAmountToStake, "Increase balance of Alice to proceed");

        // Calling the staking function as Alice
        //slither-disable-next-line arbitrary-send-eth
        rollup.stake{value: aliceAmountToStake}();

        // Now Alice should be staked
        isAliceStaked = rollup.isStaked(alice);
        assertTrue(isAliceStaked);

        uint256 aliceBalanceInitial = alice.balance;

        /*
            emit log_named_address("MSGS" , msg.sender);
            emit log_named_address("Alice", alice);
            emit log_named_address("Rollup", address(rollup));
        */

        amountToWithdraw = _generateRandomUintInRange(1, (aliceAmountToStake - minimumAmount), amountToWithdraw);

        vm.prank(alice);
        rollup.unstake(amountToWithdraw);

        uint256 aliceBalanceFinal = alice.balance;

        assertEq((aliceBalanceFinal - aliceBalanceInitial), amountToWithdraw, "Desired amount could not be withdrawn.");
    }

    function test_unstake_moreThanStakedAmount(
        uint256 confirmationPeriod,
        uint256 challengePeriod,
        uint256 minimumAssertionPeriod,
        uint256 maxGasPerAssertion,
        uint256 baseStakeAmount,
        uint256 amountToWithdraw
    ) external {
        _initializeRollup(confirmationPeriod, challengePeriod, minimumAssertionPeriod, maxGasPerAssertion, 100000);

        // Alice has not staked yet and therefore, this function should return `false`
        bool isAliceStaked = rollup.isStaked(alice);
        assertTrue(!isAliceStaked);

        uint256 minimumAmount = rollup.baseStakeAmount();
        uint256 aliceBalance = alice.balance;

        emit log_named_uint("AB", aliceBalance);

        // Let's stake something on behalf of Alice
        uint256 aliceAmountToStake = minimumAmount * 10;

        vm.prank(alice);
        require(aliceBalance >= aliceAmountToStake, "Increase balance of Alice to proceed");

        // Calling the staking function as Alice
        //slither-disable-next-line arbitrary-send-eth
        rollup.stake{value: aliceAmountToStake}();

        // Now Alice should be staked
        isAliceStaked = rollup.isStaked(alice);
        assertTrue(isAliceStaked);

        /*
            emit log_named_address("MSGS" , msg.sender);
            emit log_named_address("Alice", alice);
            emit log_named_address("Rollup", address(rollup));
        */

        amountToWithdraw =
            _generateRandomUintInRange((aliceAmountToStake - minimumAmount) + 1, type(uint256).max, amountToWithdraw);

        vm.expectRevert(IRollup.InsufficientStake.selector);
        vm.prank(alice);
        rollup.unstake(amountToWithdraw);
    }

    function test_unstake_fromUnconfirmedAssertionID(
        uint256 randomAmount,
        uint256 confirmationPeriod,
        uint256 challengePeriod
    ) external {
        // Bounding it otherwise, function `newAssertionDeadline()` overflows
        confirmationPeriod = bound(confirmationPeriod, 1, type(uint128).max);
        _initializeRollup(confirmationPeriod, challengePeriod, 1 days, 500, 1 ether);

        // Alice has not staked yet and therefore, this function should return `false`
        bool isAliceStaked = rollup.isStaked(alice);
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
        uint256 seqInboxSizeInitial = seqIn.getInboxSize();
        emit log_named_uint("Sequencer Inbox Size", seqInboxSizeInitial);

        _increaseSequencerInboxSize();

        uint256 seqInboxSizeFinal = seqIn.getInboxSize();
        emit log_named_uint("Sequencer Inbox Size", seqInboxSizeFinal);

        assertGt(seqInboxSizeFinal, seqInboxSizeInitial);

        bytes32 mockVmHash = bytes32("");
        uint256 mockInboxSize = 5;
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
        (, uint256 stakedAmount, uint256 assertionIDFinal,) = rollup.stakers(address(alice));

        assertEq(assertionIDFinal, 1); // Alice is now staked on assertionID = 1 instead of assertionID = 0.

        // Alice tries to unstake
        vm.expectRevert(IRollup.StakedOnUnconfirmedAssertion.selector);
        vm.prank(alice);
        rollup.unstake(stakedAmount);
    }

    //////////////////////
    // Remove Stake
    /////////////////////

    function test_removeStake_forNonStaker(
        uint256 randomAmount,
        uint256 confirmationPeriod,
        uint256 challengePeriod,
        uint256 minimumAssertionPeriod,
        uint256 maxGasPerAssertion,
        uint256 baseStakeAmount
    ) external {
        _initializeRollup(
            confirmationPeriod, challengePeriod, minimumAssertionPeriod, maxGasPerAssertion, baseStakeAmount
        );

        // Alice has not staked yet and therefore, this function should return `false`
        bool isAliceStaked = rollup.isStaked(alice);
        assertTrue(!isAliceStaked);

        // Since Alice is not staked, function unstake should also revert
        vm.expectRevert(IRollup.NotStaked.selector);
        vm.prank(alice);

        rollup.removeStake(address(alice));
    }

    function test_removeStake_forNonStaker_thirdPartyCall(
        uint256 randomAmount,
        uint256 confirmationPeriod,
        uint256 challengePeriod,
        uint256 minimumAssertionPeriod,
        uint256 maxGasPerAssertion,
        uint256 baseStakeAmount
    ) external {
        _initializeRollup(
            confirmationPeriod, challengePeriod, minimumAssertionPeriod, maxGasPerAssertion, baseStakeAmount
        );

        // Alice has not staked yet and therefore, this function should return `false`
        bool isAliceStaked = rollup.isStaked(alice);
        assertTrue(!isAliceStaked);

        // Since Alice is not staked, function unstake should also revert
        vm.expectRevert(IRollup.NotStaked.selector);
        vm.prank(bob);

        rollup.removeStake(address(alice));
    }

    function test_removeStake_positiveCase(
        uint256 randomAmount,
        uint256 confirmationPeriod,
        uint256 challengePeriod,
        uint256 minimumAssertionPeriod,
        uint256 maxGasPerAssertion
    ) external {
        _initializeRollup(confirmationPeriod, challengePeriod, minimumAssertionPeriod, maxGasPerAssertion, 1 ether);

        // Alice has not staked yet and therefore, this function should return `false`
        bool isAliceStaked = rollup.isStaked(alice);
        assertTrue(!isAliceStaked);

        uint256 minimumAmount = rollup.baseStakeAmount();
        uint256 aliceBalance = alice.balance;

        emit log_named_uint("AB", aliceBalance);

        // Let's stake something on behalf of Alice
        uint256 aliceAmountToStake = minimumAmount * 10;

        vm.prank(alice);
        require(aliceBalance >= aliceAmountToStake, "Increase balance of Alice to proceed");

        // Calling the staking function as Alice
        //slither-disable-next-line arbitrary-send-eth
        rollup.stake{value: aliceAmountToStake}();

        // Now Alice should be staked
        isAliceStaked = rollup.isStaked(alice);
        assertTrue(isAliceStaked);

        uint256 aliceBalanceBeforeRemoveStake = alice.balance;

        (, uint256 amountStakedInitial, uint256 assertionIDInitial, address challengeIDInitial) =
            rollup.stakers(address(alice));

        vm.prank(alice);
        rollup.removeStake(address(alice));

        (bool isStakedAfterRemoveStake, uint256 amountStakedFinal, uint256 assertionIDFinal, address challengeIDFinal) =
            rollup.stakers(address(alice));

        uint256 aliceBalanceAfterRemoveStake = alice.balance;

        assertGt(aliceBalanceAfterRemoveStake, aliceBalanceBeforeRemoveStake);
        assertEq((aliceBalanceAfterRemoveStake - aliceBalanceBeforeRemoveStake), aliceAmountToStake);

        assertTrue(!isStakedAfterRemoveStake);
    }

    function test_removeStake_fromUnconfirmedAssertionID(
        uint256 randomAmount,
        uint256 confirmationPeriod,
        uint256 challengePeriod
    ) external {
        // Bounding it otherwise, function `newAssertionDeadline()` overflows
        confirmationPeriod = bound(confirmationPeriod, 1, type(uint128).max);
        _initializeRollup(confirmationPeriod, challengePeriod, 1 days, 500, 1 ether);

        // Alice has not staked yet and therefore, this function should return `false`
        bool isAliceStaked = rollup.isStaked(alice);
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

        bytes32 mockVmHash = bytes32("");
        uint256 mockInboxSize = 5;
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

        // Try to remove Alice's stake
        vm.expectRevert(IRollup.StakedOnUnconfirmedAssertion.selector);
        rollup.removeStake(address(alice));
    }

    /////////////////////////
    // Advance Stake
    /////////////////////////

    function test_advanceStake_calledByNonStaker(
        uint256 confirmationPeriod,
        uint256 challengePeriod,
        uint256 minimumAssertionPeriod,
        uint256 maxGasPerAssertion,
        uint256 baseStakeAmount,
        uint256 assertionID
    ) external {
        _initializeRollup(
            confirmationPeriod, challengePeriod, minimumAssertionPeriod, maxGasPerAssertion, baseStakeAmount
        );

        // Alice has not staked yet and therefore, this function should return `false`
        bool isAliceStaked = rollup.isStaked(alice);
        assertTrue(!isAliceStaked);

        // Since Alice is not staked, function advanceStake should also revert
        vm.expectRevert(IRollup.NotStaked.selector);
        vm.prank(alice);

        rollup.advanceStake(assertionID);
    }

    function test_advanceStake_calledWithRandomAssertionID(
        uint256 confirmationPeriod,
        uint256 challengePeriod,
        uint256 minimumAssertionPeriod,
        uint256 maxGasPerAssertion,
        uint256 baseStakeAmount,
        uint256 assertionID
    ) external {
        _initializeRollup(confirmationPeriod, challengePeriod, minimumAssertionPeriod, maxGasPerAssertion, 1 ether);

        // Alice has not staked yet and therefore, this function should return `false`
        bool isAliceStaked = rollup.isStaked(alice);
        assertTrue(!isAliceStaked);

        uint256 minimumAmount = rollup.baseStakeAmount();
        uint256 aliceBalance = alice.balance;

        emit log_named_uint("AB", aliceBalance);

        // Let's stake something on behalf of Alice
        uint256 aliceAmountToStake = minimumAmount * 10;

        vm.prank(alice);
        require(aliceBalance >= aliceAmountToStake, "Increase balance of Alice to proceed");

        // Calling the staking function as Alice
        //slither-disable-next-line arbitrary-send-eth
        rollup.stake{value: aliceAmountToStake}();

        // Now Alice should be staked
        isAliceStaked = rollup.isStaked(alice);
        assertTrue(isAliceStaked);

        uint256 aliceBalanceBeforeRemoveStake = alice.balance;

        (,, uint256 stakerAssertionID,) = rollup.stakers(address(alice));

        uint256 lastCreatedAssertionID = rollup.lastCreatedAssertionID();

        if (assertionID > stakerAssertionID && assertionID <= lastCreatedAssertionID) {
            assertionID = lastCreatedAssertionID + 10;
        }

        vm.expectRevert(IRollup.AssertionOutOfRange.selector);
        vm.prank(alice);

        rollup.advanceStake(assertionID);
    }

    function test_advanceStake_illegalAssertionID(
        uint256 confirmationPeriod,
        uint256 challengePeriod,
        uint256 minimumAssertionPeriod,
        uint256 maxGasPerAssertion,
        uint256 baseStakeAmount,
        uint256 assertionID
    ) external {
        _initializeRollup(confirmationPeriod, challengePeriod, minimumAssertionPeriod, maxGasPerAssertion, 1 ether);

        // Alice has not staked yet and therefore, this function should return `false`
        bool isAliceStaked = rollup.isStaked(alice);
        assertTrue(!isAliceStaked);

        uint256 minimumAmount = rollup.baseStakeAmount();
        uint256 aliceBalance = alice.balance;

        emit log_named_uint("AB", aliceBalance);

        // Let's stake something on behalf of Alice
        uint256 aliceAmountToStake = minimumAmount * 10;

        vm.prank(alice);
        require(aliceBalance >= aliceAmountToStake, "Increase balance of Alice to proceed");

        // Calling the staking function as Alice
        //slither-disable-next-line arbitrary-send-eth
        rollup.stake{value: aliceAmountToStake}();

        // Now Alice should be staked
        isAliceStaked = rollup.isStaked(alice);
        assertTrue(isAliceStaked);

        vm.expectRevert(IRollup.AssertionOutOfRange.selector);
        vm.prank(alice);
        rollup.advanceStake(assertionID);
    }

    function test_advanceStake_positiveCase(uint256 randomAmount, uint256 confirmationPeriod, uint256 challengePeriod)
        external
    {
        // Bounding it otherwise, function `newAssertionDeadline()` overflows
        confirmationPeriod = bound(confirmationPeriod, 1, type(uint128).max);
        _initializeRollup(confirmationPeriod, challengePeriod, 1 days, 500, 1 ether);

        // Alice has not staked yet and therefore, this function should return `false`
        bool isAliceStaked = rollup.isStaked(alice);
        assertTrue(!isAliceStaked);

        uint256 minimumAmount = rollup.baseStakeAmount();
        uint256 aliceBalance = alice.balance;

        emit log_named_uint("AB", aliceBalance);

        // Bob also wants to stake on this assertion
        bool isBobStaked = rollup.isStaked(bob);
        assertTrue(!isBobStaked);

        uint256 bobBalance = bob.balance;

        // Let's stake something on behalf of Alice
        uint256 aliceAmountToStake = minimumAmount * 10;
        uint256 bobAmountToStake = minimumAmount * 10;

        vm.prank(alice);
        require(aliceBalance >= aliceAmountToStake, "Increase balance of Alice to proceed");

        // Calling the staking function as Alice
        //slither-disable-next-line arbitrary-send-eth
        rollup.stake{value: aliceAmountToStake}();

        vm.prank(bob);
        require(bobBalance >= bobAmountToStake, "Increase balance of Bob to proceed");

        // slither-disable-next-line arbitrary-send-eth
        rollup.stake{value: bobAmountToStake}();

        // Now Alice should be staked
        uint256 stakerAssertionID;

        // stakers mapping gets updated
        (isAliceStaked,, stakerAssertionID,) = rollup.stakers(alice);
        assertTrue(isAliceStaked);

        // Bob should be marked as staked now
        (isBobStaked,,,) = rollup.stakers(bob);
        assertTrue(isBobStaked);

        // Comparing lastConfirmedAssertionID and stakerAssertionID
        emit log_named_uint("Last Confirmed Assertion ID", rollup.lastConfirmedAssertionID());
        emit log_named_uint("Staker Assertion ID", stakerAssertionID);
        emit log_named_uint("Last Created Assertion ID", rollup.lastCreatedAssertionID());

        // Let's create a brand new assertion, so that the lastCreatedAssertionID goes up and we can successfully advance stake to the new ID after that

        // To create a brand new assertion, we will need to call Rollup.createAssertion()
        // To call Rollup.createAssertion(), we will need to pass in a param called uint256 inboxSize
        // The thing about inboxSize is that it needs to fulfill 2 require statements and to fulfill the 2nd one, the inboxSize from the SequencerInbox needs to be increased

        // Checking previous Sequencer Inbox Size
        uint256 seqInboxSize = seqIn.getInboxSize();
        emit log_named_uint("Sequencer Inbox Size", seqInboxSize);

        // Increasing the sequencerInbox inboxSize
        _increaseSequencerInboxSize();

        emit log_named_uint("Changed Sequencer Inbox Size", seqIn.getInboxSize());

        bytes32 mockVmHash = bytes32("");
        uint256 mockInboxSize = 5; // Which is smaller than the previously set sequencerInboxSize with the function dangerousIncreaseSequencerInboxSize

        // Need to figure out values of mockL2GasUsed so that the following condition is satisfied:
        // if (assertionGasUsed > maxGasPerAssertion) {
        //    revert MaxGasLimitExceeded();
        // }
        // where, uint256 assertionGasUsed = l2GasUsed - prevL2GasUsed
        uint256 mockL2GasUsed = 342;
        bytes32 mockPrevVMHash = bytes32("");
        uint256 mockPrevL2GasUsed = 0;

        emit log_named_uint("BN-1", block.number);

        /**
         * This error is popping up. Let's figure out how to tackle this:
         *         if (block.number - assertions.getProposalTime(parentID) < minimumAssertionPeriod) {
         *             revert MinimumAssertionPeriodNotPassed();
         *         }
         */
        // To avoid the MinimumAssertionPeriodNotPassed error, increase block.number
        vm.warp(block.timestamp + 50 days);
        vm.roll(block.number + (50 * 86400) / 20);

        // The method to mock startState(prevL2GasUsed, prevVmHash) to be equal to assertions.getStateHash(parentId)
        // is not known yet, so, let's assume they won't match and move forward.
        // ^ The above problem is solved because coincidentally we are on the 0th assertionID and the values of creating that can
        // be seen from the function `Rollup.initialize()`
        /**
         * assertions.createAssertion(
         *             0, // assertionID
         *             RollupLib.stateHash(RollupLib.ExecutionState(0, _initialVMhash)),
         *             0, // inboxSize (genesis)
         *             0, // parentID
         *             block.number // deadline (unchallengeable)
         *         );
         */

        /*
            rollupAssertion = rollupAssertion = rollup.assertions();
            uint256 proposalTime = rollupAssertion.getProposalTime(0);

            emit log_named_uint("Proposal Time", proposalTime);
            emit log_named_uint("BN-2", block.number);
            emit log_named_uint("MAP", 1 days);
        */

        // Now getting this error:
        /**
         * if (assertionGasUsed > maxGasPerAssertion) {
         *             revert MaxGasLimitExceeded();
         *         }
         *
         *         We've set maxGasPerAssertion as 500
         *
         *         And, assertionGasUsed = l2GasUsed - prevL2GasUsed
         */

        assertEq(rollup.lastCreatedAssertionID(), 0, "The lastCreatedAssertionID should be 0 (genesis)");
        (, uint256 amountStakedInitial, uint256 assertionIDInitial,) = rollup.stakers(address(alice));

        assertEq(assertionIDInitial, 0);

        vm.prank(alice);
        rollup.createAssertion(mockVmHash, mockInboxSize, mockL2GasUsed, mockPrevVMHash, mockPrevL2GasUsed);

        // Now assuming that the last assertion was created successfully, the lastCreatedAssertionID should have bumped to 1.
        assertEq(rollup.lastCreatedAssertionID(), 1, "LastCreatedAssertionID not updated correctly");

        // The assertionID of alice should change after she called `createAssertion`
        (, uint256 amountStakedFinal, uint256 assertionIDFinal,) = rollup.stakers(address(alice));

        assertEq(amountStakedInitial, amountStakedFinal);
        assertEq(assertionIDFinal, 1);

        // Advance stake of the staker
        // Since Alice's stake was already advanced when she called createAssertion, her call to `rollup.advanceStake` should fail
        vm.expectRevert(IRollup.AssertionOutOfRange.selector);
        vm.prank(alice);
        rollup.advanceStake(1);

        // Bob's call to `rollup.advanceStake` should succeed as he is still staked on the previous assertion
        vm.prank(bob);
        rollup.advanceStake(1);

        (,, uint256 bobAssertionID,) = rollup.stakers(address(alice));

        assertEq(bobAssertionID, 1);
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

    function _createTxBatch() internal pure returns (bytes memory) {
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

        return bytes.concat(keccak256(abi.encode(tx1, tx2, tx3, tx4, tx5, tx6)));
    }

    // This function increases the inbox size by 6
    function _increaseSequencerInboxSize() internal {
        // Each context corresponds to a single "L2 block"
        // `contexts` is represented with uint256 3-tuple: (numTxs, l2BlockNumber, l2Timestamp)
        // Let's create an array of contexts
        uint256 numTxns = 3;
        uint256 timeStamp1 = block.timestamp / 10;
        uint256 timeStamp2 = block.timestamp / 5;
        uint256 blockNumber1 = timeStamp1 / 20;
        uint256 blockNumber2 = timeStamp2 / 20;

        uint256[] memory contexts = new uint256[](6);
        uint256[] memory txLengths = new uint256[](6);

        // Let's assume that we had 2 blocks and each had 3 transactions
        contexts[0] = (numTxns);
        contexts[1] = (blockNumber1);
        contexts[2] = (timeStamp1);
        contexts[3] = (numTxns);
        contexts[4] = (blockNumber2);
        contexts[5] = (timeStamp2);

        // txLengths is defined as: Array of lengths of each encoded tx in txBatch
        for (uint256 i; i < 6; i++) {
            txLengths[i] = (uint256(268));
        }

        // txBatch is defined as: Batch of RLP-encoded transactions
        bytes memory txBatch = _createTxBatch();

        // Pranking as the sequencer and calling appendTxBatch
        vm.prank(sequencer);
        seqIn.appendTxBatch(contexts, txLengths, txBatch);
    }

    function _initializeRollup(
        uint256 confirmationPeriod,
        uint256 challengePeriod,
        uint256 minimumAssertionPeriod,
        uint256 maxGasPerAssertion,
        uint256 baseStakeAmount
    ) internal {
        Rollup _tempRollup = new Rollup();

        bytes memory initializingData = abi.encodeWithSelector(
            Rollup.initialize.selector,
            owner, // owner
            address(seqIn), // sequencerInbox
            address(verifier),
            address(stakeToken),
            confirmationPeriod, //confirmationPeriod
            challengePeriod, //challengePeriod
            minimumAssertionPeriod, // minimumAssertionPeriod
            maxGasPerAssertion, // maxGasPerAssertion
            baseStakeAmount, //baseStakeAmount
            bytes32("")
        );

        address proxyAdmin = makeAddr("Proxy Admin");

        TransparentUpgradeableProxy proxy = new TransparentUpgradeableProxy(
            address(_tempRollup), 
            proxyAdmin, 
            initializingData
        );

        // Initialize is called here for the first time.
        rollup = Rollup(address(proxy));
    }
}
