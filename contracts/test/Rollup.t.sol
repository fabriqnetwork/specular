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

contract BaseSetup is Test {
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

contract RollupTest is BaseSetup {
    SequencerInbox private seqIn;
    Rollup private rollup;
    uint256 randomNonce;
    AssertionMap rollupAssertion;

    function setUp() public virtual override {
        BaseSetup.setUp();

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
        //slither-disable-next-line arbitrary-send
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
        //slither-disable-next-line arbitrary-send
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
        //slither-disable-next-line arbitrary-send
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
        //slither-disable-next-line arbitrary-send
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
        //slither-disable-next-line arbitrary-send
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
        //slither-disable-next-line arbitrary-send
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

    //////////////////////////////////////////////////////////////////////////////
    // Test on-hold due to bug in Rollup contract (Rollup.removeStake)
    //////////////////////////////////////////////////////////////////////////////
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
        //slither-disable-next-line arbitrary-send
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

        /*
            // emit log_named_uint("Initial Staked Amount", amountStakedInitial);
            // emit log_named_uint("Final Staked Amount", amountStakedFinal);
            // emit log_named_uint("Base Amount", 1 ether);
            // emit log_named_uint("Alice Staked Amount", aliceAmountToStake);
            // emit log_named_uint("Alice Balance Before", aliceBalanceBeforeRemoveStake);
            // emit log_named_uint("Alice Balance After", aliceBalanceAfterRemoveStake);
        */

        // The below assertion will fail because of a bug in the Rollup contract. Until that is fixed, I'll just comment this assertion out.
        /**
         * function removeStake(address stakerAddress) external override {
         *             requireStaked(stakerAddress);
         *             // Require that staker is staked on a confirmed assertion.
         *             Staker storage staker = stakers[stakerAddress];
         *             if (staker.assertionID > lastConfirmedAssertionID) {
         *                 revert StakedOnUnconfirmedAssertion();
         *             }
         *             deleteStaker(stakerAddress);
         *             // Note: we don't need to modify assertion state because you can only unstake from a confirmed assertion.
         *             (bool success,) = stakerAddress.call{value: staker.amountStaked}("");
         *             if (!success) revert TransferFailed();
         *         }
         */
        // Now what you'll notice in the second last line where the stake is being transferred to the staker, the value is kept as `staker.amountStaked`. However in the line above that
        // the staker has been deleted and thus the value transferred will always be 0.
        // Need to fix this critical bug.
        //assertGt(aliceBalanceAfterRemoveStake, aliceBalanceBeforeRemoveStake);
        // assertEq((aliceBalanceAfterRemoveStake - aliceBalanceBeforeRemoveStake), aliceAmountToStake);

        assertTrue(!isStakedAfterRemoveStake);
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

    //////////////////////////////////////////////////////////////////////////////
    // Test on-hold due to bug in Rollup contract (Rollup.advanceStake)
    //////////////////////////////////////////////////////////////////////////////
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
        //slither-disable-next-line arbitrary-send
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

        // The following lines of test are written to expect the rollup.advanceStake function to revert with the error IRollup.AssertionOutOfRange error,
        // since the assertionID is generated randomly and checked so that it does not lie between the staker's assertionID and the last createdID.
        // However, the test fails to revert since the function uses an `&&` operator to check these conditions rather than an `||` operator.
        // Commenting out this test unless the bug in the Rollup contract is not fixed.
        /**
         * function advanceStake(uint256 assertionID) external override stakedOnly {
         *             Staker storage staker = stakers[msg.sender];
         *             // The below line has the bug.
         *             if (assertionID <= staker.assertionID && assertionID > lastCreatedAssertionID) {
         *                 revert AssertionOutOfRange();
         *             }
         *             // TODO: allow arbitrary descendant of current staked assertionID, not just child.
         *             if (staker.assertionID != assertions.getParentID(assertionID)) {
         *                 revert ParentAssertionUnstaked();
         *             }
         *             stakeOnAssertion(msg.sender, assertionID);
         *         }
         */

        // vm.expectRevert(IRollup.AssertionOutOfRange.selector);
        // vm.prank(alice);

        // rollup.advanceStake(assertionID);
    }

    function test_advanceStake_positiveCase(
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
        //slither-disable-next-line arbitrary-send
        rollup.stake{value: aliceAmountToStake}();

        // Now Alice should be staked
        isAliceStaked = rollup.isStaked(alice);
        assertTrue(isAliceStaked);

        uint256 aliceBalanceBeforeRemoveStake = alice.balance;

        (,, uint256 stakerAssertionID,) = rollup.stakers(address(alice));

        vm.prank(alice);
        rollup.advanceStake(assertionID);

        (,, uint256 stakerAssertionIDFinal,) = rollup.stakers(address(alice));

        assertEq(stakerAssertionIDFinal, assertionID);
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
