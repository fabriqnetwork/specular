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

    // 1. external/public functions (user data)
    // 2. functions that call an external contract
    // 3. calculateINterest, calculateWithdrawabeCollateral
    // 4. test pause and upgrade mechanism of the contracts.
    // 5. internal/private functions: Anchor contracts -> public 

    function test_initializeRollup_vaultAddressZero() external {
        // Preparing the initializing data for the (proxied)Rollup contract
        bytes memory initializingData = abi.encodeWithSelector(
            Rollup.initialize.selector,
            address(0), // vault
            address(seqIn),
            address(verifier),
            0, //confirmationPeriod
            0, //challengPeriod
            0, // minimumAssertionPeriod
            type(uint256).max, // maxGasPerAssertion
            0, //baseStakeAmount,
            0, // initialAssertionID
            0, // initialInboxSize
            bytes32(""),
            0 // initialGasUsed
        );

        vm.startPrank(rollupOwner);

        Rollup implementationRollup = new Rollup(); // implementation contract

        vm.expectRevert(ZeroAddress.selector);
        specularProxy = new SpecularProxy(address(implementationRollup), initializingData);

        rollup = Rollup(address(specularProxy));

        vm.stopPrank();
    }

    function test_initializeRollup_verifierAddressZero() external {
        bytes memory initializingData = abi.encodeWithSelector(
            Rollup.initialize.selector,
            vault, // vault
            address(seqIn),
            address(0),
            0, //confirmationPeriod
            0, //challengPeriod
            0, // minimumAssertionPeriod
            type(uint256).max, // maxGasPerAssertion
            0, //baseStakeAmount,
            0, // initialAssertionID
            0, // initialInboxSize
            bytes32(""),
            0 // initialGasUsed
        );

        vm.startPrank(rollupOwner);

        Rollup implementationRollup = new Rollup(); // implementation contract

        vm.expectRevert(ZeroAddress.selector);
        specularProxy = new SpecularProxy(address(implementationRollup), initializingData);

        rollup = Rollup(address(specularProxy));

        vm.stopPrank();
    }

    function test_initializeRollup_sequencerInboxAddressZero() external {
        bytes memory initializingData = abi.encodeWithSelector(
            Rollup.initialize.selector,
            vault, // vault
            address(0),
            address(verifier),
            0, //confirmationPeriod
            0, //challengPeriod
            0, // minimumAssertionPeriod
            type(uint256).max, // maxGasPerAssertion
            0, //baseStakeAmount,
            0, // initialAssertionID
            0, // initialInboxSize
            bytes32(""),
            0 // initialGasUsed
        );

        vm.startPrank(rollupOwner);

        Rollup implementationRollup = new Rollup(); // implementation contract

        vm.expectRevert(ZeroAddress.selector);
        specularProxy = new SpecularProxy(address(implementationRollup), initializingData);

        rollup = Rollup(address(specularProxy));

        vm.stopPrank();
    }

    function test_initializeRollup_cannotBeCalledTwice() external {
        bytes memory initializingData = abi.encodeWithSelector(
            Rollup.initialize.selector,
            vault, // vault
            address(seqIn),
            address(verifier),
            0, //confirmationPeriod
            0, //challengPeriod
            0, // minimumAssertionPeriod
            type(uint256).max, // maxGasPerAssertion
            0, //baseStakeAmount,
            0, // initialAssertionID
            0, // initialInboxSize
            bytes32(""),
            0 // initialGasUsed
        );

        vm.startPrank(rollupOwner);

        Rollup implementationRollup = new Rollup(); // implementation contract
        specularProxy = new SpecularProxy(address(implementationRollup), initializingData);
        rollup = Rollup(address(specularProxy)); // The rollup contract (proxy, not implementation should have been initialized by now)

        vm.stopPrank();

        // Trying to call initialize for the second time
        vm.expectRevert("Initializable: contract is already initialized");
        vm.prank(alice); // Someone random, who is not the proxy owner.

        rollup.initialize(
            vault, // vault
            address(seqIn),
            address(verifier),
            0, //confirmationPeriod
            0, //challengPeriod
            0, // minimumAssertionPeriod
            type(uint256).max, // maxGasPerAssertion
            0, //baseStakeAmount,
            0, // initialAssertionID
            0, // initialInboxSize
            bytes32(""),
            0 // initialGasUsed
        );
    }

    /*
    function test_initializeRollup_valuesAfterInit(
        address _vault,
        uint256 confirmationPeriod,
        uint256 challengePeriod,
        uint256 minimumAssertionPeriod,
        uint256 maxGasPerAssertion,
        uint256 baseStakeAmount,
        uint256 initialInboxSize,
        uint256 initialL2GasUsed,
        uint256 initialAssertionID
    ) external {
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

        vm.startPrank(rollupOwner);

        Rollup implementationRollup = new Rollup(); // implementation contract
        specularProxy = new SpecularProxy(address(implementationRollup), initializingData);
        rollup = Rollup(address(specularProxy)); // The rollup contract (proxy, not implementation should have been initialized by now)

        vm.stopPrank();


        // Putting in different scope to do away with the stack too deep error.
        {
            // Check if the value of the address owner was set correctly
            address _rollupOwner = rollup.owner();
            assertEq(_rollupOwner, rollupOwner, "Rollup.initialize failed to update owner correctly");

            // Check if the value of SequencerInbox was set correctly
            address rollupSeqIn = address(rollup.sequencerInbox());
            assertEq(rollupSeqIn, address(seqIn), "Rollup.initialize failed to update Sequencer Inbox correctly");

            // Check if the value of the verifier was set correctly
            address rollupVerifier = address(rollup.verifier());
            assertEq(rollupVerifier, address(verifier), "Rollup.initialize failed to update verifier value correctly");
        }

        {
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
                "Rollup.initialize failed to update challengePeriod value correctly"
            );
            assertEq(
                rollupMinimumAssertionPeriod,
                minimumAssertionPeriod,
                "Rollup.initialize failed to update minimumAssertionPeriod value correctly"
            );
            assertEq(
                rollupMaxGasPerAssertion,
                maxGasPerAssertion,
                "Rollup.initialize failed to update maxGasPerAssertion value correctly"
            );
            assertEq(
                rollupBaseStakeAmount,
                baseStakeAmount,
                "Rollup.initialize failed to update baseStakeAmount value correctly"
            );
        }
    }
    */

    ////////////////
    // Staking
    ///////////////

    function test_stake_isStaked(
        uint256 confirmationPeriod,
        uint256 challengePeriod,
        uint256 minimumAssertionPeriod,
        uint256 maxGasPerAssertion,
        uint256 baseStakeAmount,
        uint256 initialAssertionID,
        uint256 initialInboxSize,
        uint256 initialL2GasUsed
    ) external {
        address _vault = makeAddr("vault");
        _initializeRollup(
            _vault,
            confirmationPeriod, challengePeriod, minimumAssertionPeriod, maxGasPerAssertion, baseStakeAmount,
            initialAssertionID,
            initialInboxSize,
            initialL2GasUsed
        );

        // Alice has not staked yet and therefore, this function should return `false`
        (bool isAliceStaked,,,) = rollup.stakers(alice);
        assertTrue(!isAliceStaked);
    }

    function test_stake_insufficentAmountStaking(
        uint256 confirmationPeriod,
        uint256 challengePeriod,
        uint256 minimumAssertionPeriod,
        uint256 maxGasPerAssertion,
        uint256 initialAssertionID,
        uint256 initialInboxSize,
        uint256 initialL2GasUsed
    ) external {
        address _vault = makeAddr("vault");
        _initializeRollup(
            _vault,
            confirmationPeriod, challengePeriod, minimumAssertionPeriod, maxGasPerAssertion, type(uint256).max,
            initialAssertionID,
            initialInboxSize,
            initialL2GasUsed
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

        (bool isAliceStaked,,,) = rollup.stakers(alice);
        assertTrue(!isAliceStaked);
    }

    function test_stake_sufficientAmountStakingAndNumStakersIncrement(
        uint256 confirmationPeriod,
        uint256 challengePeriod,
        uint256 minimumAssertionPeriod,
        uint256 maxGasPerAssertion,
        uint256 initialAssertionID,
        uint256 initialInboxSize,
        uint256 initialL2GasUsed
    ) external {
        address _vault = makeAddr("vault");
        _initializeRollup(
            _vault,
            confirmationPeriod, challengePeriod, minimumAssertionPeriod, maxGasPerAssertion, 1000,
            initialAssertionID,
            initialInboxSize,
            initialL2GasUsed
        );

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
        (bool isAliceStaked,,,) = rollup.stakers(alice);
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
        uint256 initialAssertionID,
        uint256 initialInboxSize,
        uint256 initialL2GasUsed
    ) external {
        address _vault = makeAddr("vault");
        _initializeRollup(
            _vault,
            confirmationPeriod, challengePeriod, minimumAssertionPeriod, maxGasPerAssertion, 1000,
            initialAssertionID,
            initialInboxSize,
            initialL2GasUsed
        );

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
        (bool isAliceStaked,,,) = rollup.stakers(alice);
        assertTrue(isAliceStaked);

        // stakers mapping gets updated
        (isAliceStaked, amountStaked, assertionID, challengeAddress) = rollup.stakers(alice);

        uint256 aliceBalanceFinal = alice.balance;

        assertEq(alice.balance, 0, "Alice should not have any balance left");
        assertGt(bob.balance, 0, "Bob should have a non-zero native currency balance");

        vm.prank(bob);
        (bool sent, ) = alice.call{value: bob.balance}("");
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
        uint256 baseStakeAmount,
        uint256 initialAssertionID,
        uint256 initialInboxSize,
        uint256 initialL2GasUsed
    ) external {
        address _vault = makeAddr("vault");
        _initializeRollup(
            _vault,
            confirmationPeriod, challengePeriod, minimumAssertionPeriod, maxGasPerAssertion, baseStakeAmount,
            initialAssertionID,
            initialInboxSize,
            initialL2GasUsed
        );

        // Alice has not staked yet and therefore, this function should return `false`
        (bool isAliceStaked,,,) = rollup.stakers(alice);
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
        uint256 amountToWithdraw,
        uint256 initialAssertionID,
        uint256 initialInboxSize,
        uint256 initialL2GasUsed
    ) external {
        address _vault = makeAddr("vault");
        _initializeRollup(
            _vault,
            confirmationPeriod, challengePeriod, minimumAssertionPeriod, maxGasPerAssertion, 100000,
            initialAssertionID,
            initialInboxSize,
            initialL2GasUsed
        );

        // Alice has not staked yet and therefore, this function should return `false`
        (bool isAliceStaked,,,) = rollup.stakers(alice);
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
        (isAliceStaked,,,) = rollup.stakers(alice);
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
        uint256 amountToWithdraw,
        uint256 initialAssertionID,
        uint256 initialInboxSize,
        uint256 initialL2GasUsed
    ) external {
        address _vault = makeAddr("vault");
        _initializeRollup(
            _vault,
            confirmationPeriod, challengePeriod, minimumAssertionPeriod, maxGasPerAssertion, 100000,
            initialAssertionID,
            initialInboxSize,
            initialL2GasUsed
        );

        // Alice has not staked yet and therefore, this function should return `false`
        (bool isAliceStaked,,,) = rollup.stakers(alice);
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
        (isAliceStaked,,,) = rollup.stakers(alice);
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
        uint256 confirmationPeriod,
        uint256 challengePeriod
    ) external {
        // Bounding it otherwise, function `newAssertionDeadline()` overflows
        address _vault = makeAddr("vault");
        confirmationPeriod = bound(confirmationPeriod, 1, type(uint128).max);
        _initializeRollup(
            _vault,
            confirmationPeriod, 
            challengePeriod, 
            1 days, 
            500, 
            1 ether,
            0,
            5,
            0
        );

        // Alice has not staked yet and therefore, this function should return `false`
        (bool isAliceStaked,,uint256 assertionID1,) = rollup.stakers(alice);
        assertEq(assertionID1, 0);
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

        // stakers mapping gets updated
        (bool isAliceStaked2, ,uint256 assertionID2, ) = rollup.stakers(alice);
        emit log_named_uint("assertion id 2", assertionID2);
        assertEq(assertionID2, 0);
        assertTrue(isAliceStaked2);

        // Checking previous Sequencer Inbox Size
        uint256 seqInboxSizeInitial = seqIn.getInboxSize();
        emit log_named_uint("Sequencer Inbox Size", seqInboxSizeInitial);

        _increaseSequencerInboxSize();

        uint256 seqInboxSizeFinal = seqIn.getInboxSize();
        emit log_named_uint("Sequencer Inbox Size", seqInboxSizeFinal);

        assertGt(seqInboxSizeFinal, seqInboxSizeInitial);

        bytes32 mockVmHash = bytes32("");
        uint256 mockInboxSize = 6;
        uint256 mockL2GasUsed = 342;
        bytes32 mockPrevVMHash = bytes32("");
        uint256 mockPrevL2GasUsed = 0;

        // To avoid the MinimumAssertionPeriodNotPassed error, increase block.number
        vm.warp(block.timestamp + 50 days);
        vm.roll(block.number + (50 * 86400) / 20);

        // assertEq(rollup.lastCreatedAssertionID(), 2, "The lastCreatedAssertionID should be 1 (genesis)");
        (,, uint256 assertionIDInitial,) = rollup.stakers(address(alice));
        assertEq(assertionIDInitial, 0, "Alice has not yet created an assertionID, so it should be the same as genesis(0)");

        vm.prank(alice);
        rollup.createAssertion(mockVmHash, mockInboxSize, mockL2GasUsed, mockPrevVMHash, mockPrevL2GasUsed);

        // The assertionID of alice should change after she called `createAssertion`
        (, uint256 stakedAmount, uint256 assertionIDFinal,) = rollup.stakers(address(alice));

        assertEq(assertionIDFinal, 1); // Alice is now staked on assertionID = 1 instead of assertionID = 0.

        assertGt(assertionIDFinal, assertionIDInitial, "AssertionID should increase");

        // Alice tries to unstake
        vm.prank(alice);
        vm.expectRevert(IRollup.StakedOnUnconfirmedAssertion.selector);
        rollup.unstake(stakedAmount);
    }

    //////////////////////
    // Remove Stake
    /////////////////////

    function test_removeStake_forNonStaker(
        uint256 confirmationPeriod,
        uint256 challengePeriod,
        uint256 minimumAssertionPeriod,
        uint256 maxGasPerAssertion,
        uint256 baseStakeAmount,
        uint256 initialAssertionID,
        uint256 initialInboxSize,
        uint256 initialL2GasUsed
    ) external {
        address _vault = makeAddr("vault");
        _initializeRollup(
            _vault,
            confirmationPeriod, challengePeriod, minimumAssertionPeriod, maxGasPerAssertion, baseStakeAmount,
            initialAssertionID,
            initialInboxSize,
            initialL2GasUsed
        );

        // Alice has not staked yet and therefore, this function should return `false`
        (bool isAliceStaked,,,) = rollup.stakers(alice);
        assertTrue(!isAliceStaked);

        // Since Alice is not staked, function unstake should also revert
        vm.expectRevert(IRollup.NotStaked.selector);
        vm.prank(alice);

        rollup.removeStake(address(alice));
    }

    function test_removeStake_forNonStaker_thirdPartyCall(
        uint256 confirmationPeriod,
        uint256 challengePeriod,
        uint256 minimumAssertionPeriod,
        uint256 maxGasPerAssertion,
        uint256 baseStakeAmount,
        uint256 initialAssertionID,
        uint256 initialInboxSize,
        uint256 initialL2GasUsed
    ) external {
        address _vault = makeAddr("vault");
        _initializeRollup(
            _vault,
            confirmationPeriod, challengePeriod, minimumAssertionPeriod, maxGasPerAssertion, baseStakeAmount,
            initialAssertionID,
            initialInboxSize,
            initialL2GasUsed
        );

        // Alice has not staked yet and therefore, this function should return `false`
        (bool isAliceStaked,,,) = rollup.stakers(alice);
        assertTrue(!isAliceStaked);

        // Since Alice is not staked, function unstake should also revert
        vm.expectRevert(IRollup.NotStaked.selector);
        vm.prank(bob);

        rollup.removeStake(address(alice));
    }

    function test_removeStake_positiveCase(
        uint256 confirmationPeriod,
        uint256 challengePeriod,
        uint256 minimumAssertionPeriod,
        uint256 maxGasPerAssertion,
        uint256 initialAssertionID,
        uint256 initialInboxSize,
        uint256 initialL2GasUsed
    ) external {
        address _vault = makeAddr("vault");
        
        _initializeRollup(
            _vault,
            confirmationPeriod, challengePeriod, minimumAssertionPeriod, maxGasPerAssertion, 1 ether,
            initialAssertionID,
            initialInboxSize,
            initialL2GasUsed
        );

        // Alice has not staked yet and therefore, this function should return `false`
        (bool isAliceStaked,,,) = rollup.stakers(alice);
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
        (isAliceStaked,,,) = rollup.stakers(alice);
        assertTrue(isAliceStaked);

        uint256 aliceBalanceBeforeRemoveStake = alice.balance;

        vm.prank(alice);
        rollup.removeStake(address(alice));

        (bool isStakedAfterRemoveStake, , , ) =
            rollup.stakers(address(alice));

        uint256 aliceBalanceAfterRemoveStake = alice.balance;

        assertGt(aliceBalanceAfterRemoveStake, aliceBalanceBeforeRemoveStake);
        assertEq((aliceBalanceAfterRemoveStake - aliceBalanceBeforeRemoveStake), aliceAmountToStake);

        assertTrue(!isStakedAfterRemoveStake);
    }

    function test_removeStake_fromUnconfirmedAssertionID(uint256 confirmationPeriod, uint256 challengePeriod,
        uint256 initialAssertionID,
        uint256 initialInboxSize,
        uint256 initialL2GasUsed)
        external
    {
        // Bounding it otherwise, function `newAssertionDeadline()` overflows
        address _vault = makeAddr("vault");
        confirmationPeriod = bound(confirmationPeriod, 1, type(uint128).max);
        _initializeRollup(
            _vault,
            confirmationPeriod, 
            challengePeriod, 
            1 days, 
            500, 
            1 ether,
            0,
            5,
            0
        );

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
        assertEq(stakerAssertionID, 0);

        // Checking previous Sequencer Inbox Size
        uint256 seqInboxSize = seqIn.getInboxSize();
        emit log_named_uint("Sequencer Inbox Size", seqInboxSize);

        _increaseSequencerInboxSize();

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
        uint256 assertionID,
        uint256 initialAssertionID,
        uint256 initialInboxSize,
        uint256 initialL2GasUsed
    ) external {
        address _vault = makeAddr("vault");
        
        _initializeRollup(
            _vault,
            confirmationPeriod, challengePeriod, minimumAssertionPeriod, maxGasPerAssertion, baseStakeAmount,
            initialAssertionID,
            initialInboxSize,
            initialL2GasUsed
        );

        // Alice has not staked yet and therefore, this function should return `false`
        (bool isAliceStaked,,,) = rollup.stakers(alice);
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
        uint256 assertionID,
        uint256 initialAssertionID,
        uint256 initialInboxSize,
        uint256 initialL2GasUsed
    ) external {
        address _vault = makeAddr("vault");
        _initializeRollup(
            _vault,
            confirmationPeriod, challengePeriod, minimumAssertionPeriod, maxGasPerAssertion, 1 ether,
            initialAssertionID,
            initialInboxSize,
            initialL2GasUsed
        );

        // Alice has not staked yet and therefore, this function should return `false`
        (bool isAliceStaked,,,) = rollup.stakers(alice);
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
        (isAliceStaked,,,) = rollup.stakers(alice);
        assertTrue(isAliceStaked);

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
        uint256 assertionID,
        uint256 initialAssertionID,
        uint256 initialInboxSize,
        uint256 initialL2GasUsed
    ) external {
        address _vault = makeAddr("vault");
        _initializeRollup(
            _vault,
            confirmationPeriod, challengePeriod, minimumAssertionPeriod, maxGasPerAssertion, 1 ether,
            initialAssertionID,
            initialInboxSize,
            initialL2GasUsed
        );

        // Alice has not staked yet and therefore, this function should return `false`
        (bool isAliceStaked,,,) = rollup.stakers(alice);
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
        (isAliceStaked,,,) = rollup.stakers(alice);
        assertTrue(isAliceStaked);

        vm.expectRevert(IRollup.AssertionOutOfRange.selector);
        vm.prank(alice);
        rollup.advanceStake(assertionID);
    }

    function test_advanceStake_positiveCase(uint256 confirmationPeriod, uint256 challengePeriod,
        uint256 initialAssertionID,
        uint256 initialInboxSize,
        uint256 initialL2GasUsed)
        external
    {
        // Bounding it otherwise, function `newAssertionDeadline()` overflows
        confirmationPeriod = bound(confirmationPeriod, 1, type(uint128).max);
        address _vault = makeAddr("vault");

        _initializeRollup(
            _vault,
            confirmationPeriod, challengePeriod, 1 days, 500, 1 ether,
            0,
            5,
            0
        );

        // Alice has not staked yet and therefore, this function should return `false`
        (bool isAliceStaked,,,) = rollup.stakers(alice);
        assertTrue(!isAliceStaked);

        uint256 minimumAmount = rollup.baseStakeAmount();
        uint256 aliceBalance = alice.balance;

        emit log_named_uint("AB", aliceBalance);

        // Bob also wants to stake on this assertion
        (bool isBobStaked,,,) = rollup.stakers(bob);
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
        uint256 mockInboxSize = 6; // Which is smaller than the previously set sequencerInboxSize with the function _increaseSequencerInboxSize

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
            confirmationPeriod, challengePeriod, minimumAssertionPeriod, maxGasPerAssertion, type(uint256).max,
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
            confirmationPeriod, challengePeriod, minimumAssertionPeriod, maxGasPerAssertion, type(uint256).max,
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

    function test_challengeAssertion_alreadyResolvedAssertionID(uint256 confirmationPeriod, uint256 challengePeriod,
        uint256 initialAssertionID,
        uint256 initialInboxSize,
        uint256 initialL2GasUsed)
        public
    {
        // Initializing the rollup
        confirmationPeriod = bound(confirmationPeriod, 1, type(uint128).max);
        address _vault = makeAddr("vault");
        _initializeRollup(
            _vault,
            confirmationPeriod, challengePeriod, 1 days, 500, 1 ether,
            0,
            5,
            0
        );

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
        if(_vault == address(0)) {
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
