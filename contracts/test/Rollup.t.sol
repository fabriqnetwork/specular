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
import {PausableUpgradeable} from "@openzeppelin/contracts-upgradeable/security/PausableUpgradeable.sol";
import {StringsUpgradeable} from "@openzeppelin/contracts-upgradeable/utils/StringsUpgradeable.sol";

import {Utils} from "./utils/Utils.sol";
import {IRollup, RollupData} from "../src/IRollup.sol";
import {Verifier} from "../src/challenge/verifier/Verifier.sol";
import {Rollup, RollupData} from "../src/Rollup.sol";
import {ISequencerInbox} from "../src/ISequencerInbox.sol";
import {SequencerInbox} from "../src/SequencerInbox.sol";
import {RLPEncodedTransactionsUtil} from "./utils/RLPEncodedTransactions.sol";

contract RollupBaseSetup is Test, RLPEncodedTransactionsUtil, RollupData {
    Utils internal utils;
    address payable[] internal users;

    address internal alice;
    address internal bob;
    address internal deployer;
    address internal sequencerAddress;
    address internal defender;
    address internal challenger;

    Config internal validConfig;
    InitialRollupState internal validState;

    Verifier verifier = new Verifier();

    event TxBatchAppended();

    function setUp() public virtual {
        utils = new Utils();
        users = utils.createUsers(6);

        deployer = users[0];
        sequencerAddress = users[1];
        alice = users[2];
        bob = users[3];
        defender = users[4];
        challenger = users[5];

        validState = InitialRollupState(0, 0, bytes32(""), bytes32(""));
    }
}

contract RollupTest is RollupBaseSetup {
    Rollup public rollup;
    uint256 randomNonce;
    SequencerInbox public seqIn;
    SequencerInbox public implementationSequencer;

    function setUp() public virtual override {
        // Parent contract setup
        RollupBaseSetup.setUp();

        // Deploying the SequencerInbox
        bytes memory seqInInitData = abi.encodeWithSignature("initialize(address)", sequencerAddress);
        vm.startPrank(deployer);
        implementationSequencer = new SequencerInbox();
        seqIn = SequencerInbox(address(new ERC1967Proxy(address(implementationSequencer), seqInInitData)));
        vm.stopPrank();

        // Making sure that the proxy returns the correct proxy owner and sequencerAddress
        address sequencerInboxDeployer = seqIn.owner();
        assertEq(sequencerInboxDeployer, deployer);

        address fetchedSequencerAddress = seqIn.sequencerAddress();
        assertEq(fetchedSequencerAddress, sequencerAddress);

        validConfig = Config(
            sequencerAddress, // vault
            address(seqIn),
            address(verifier),
            0,
            0,
            0,
            0,
            new address[](0)
        );
    }

    function test_constructRollup_zeroValues_reverts() external {
        vm.startPrank(deployer);

        Config[] memory cfgs = new Config[](3);
        cfgs[0] = Config(address(0), address(1), address(1), 0, 0, 0, 0, new address[](0));
        cfgs[1] = Config(address(1), address(0), address(1), 0, 0, 0, 0, new address[](0));
        cfgs[2] = Config(address(1), address(1), address(0), 0, 0, 0, 0, new address[](0));

        for (uint256 i = 0; i < cfgs.length; i++) {
            bytes memory initializingData = abi.encodeWithSelector(Rollup.initialize.selector, cfgs[i]);

            Rollup implementationRollup = new Rollup(); // implementation contract

            vm.expectRevert();
            rollup = Rollup(address(new ERC1967Proxy(address(implementationRollup), initializingData)));
        }
    }

    function test_initialize_reinitializeRollup_reverts() external {
        bytes memory initializingData = abi.encodeWithSelector(Rollup.initialize.selector, validConfig);

        vm.startPrank(deployer);

        Rollup implementationRollup = new Rollup(); // implementation contract
        rollup = Rollup(address(new ERC1967Proxy(address(implementationRollup), initializingData))); // The rollup contract (proxy, not implementation should have been initialized by now)
        rollup.initializeGenesis(validState);

        // Trying to call initialize for the second time
        vm.expectRevert("Initializable: contract is already initialized");

        rollup.initialize(validConfig);
    }

    function testFuzz_initialize_valuesAfterInit_succeeds(
        uint256 confirmationPeriod,
        uint256 challengePeriod,
        uint256 minimumAssertionPeriod,
        uint256 baseStakeAmount,
        uint256 initialInboxSize,
        uint256 initialAssertionID
    ) external {
        {
            Config memory cfg = Config(
                sequencerAddress,
                address(seqIn), // sequencerInbox
                address(verifier),
                confirmationPeriod, //confirmationPeriod
                challengePeriod, //challengePeriod
                minimumAssertionPeriod, // minimumAssertionPeriod
                baseStakeAmount, //baseStakeAmount
                new address[](0) // validators
            );
            InitialRollupState memory state =
                InitialRollupState(initialAssertionID, initialInboxSize, bytes32(""), bytes32(""));
            bytes memory initializingData = abi.encodeWithSelector(Rollup.initialize.selector, cfg);

            vm.startPrank(deployer);

            Rollup implementationRollup = new Rollup(); // implementation contract
            rollup = Rollup(address(new ERC1967Proxy(address(implementationRollup), initializingData))); // The rollup contract (proxy, not implementation should have been initialized by now)
            rollup.initializeGenesis(state);

            vm.stopPrank();
        }

        // Putting in different scope to do away with the stack too deep error.
        {
            // Check if the value of the address owner was set correctly
            address _rollupDeployer = rollup.owner();
            assertEq(_rollupDeployer, deployer, "Rollup.initialize failed to update owner correctly");

            // Check if the value of SequencerInbox was set correctly
            address rollupSeqIn = address(rollup.daProvider());
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
                rollupBaseStakeAmount,
                baseStakeAmount,
                "Rollup.initialize failed to update baseStakeAmount value correctly"
            );
        }
    }

    ////////////////
    // Admin
    ///////////////

    function test_addValidators_succeeds() external {
        // Initialize rollup with an empty validator whitelist
        _initializeRollup(
            0, // confirmationPeriod
            0, // challengePeriod
            0, // minimumAssertionPeriod
            1, // baseStakeAmount,
            0, // initialAssertionID
            0, // initialInboxSize
            new address[](0) // validator whitelist
        );

        // Adding a validator to whitelist as non-admin should fail
        vm.expectRevert(
            abi.encodePacked(
                "AccessControl: account ",
                StringsUpgradeable.toHexString(alice),
                " is missing role ",
                StringsUpgradeable.toHexString(uint256(rollup.DEFAULT_ADMIN_ROLE()), 32)
            )
        );
        vm.prank(alice);
        rollup.addValidator(alice);

        // Update validator whitelist as deployer, adding alice to whitelist
        vm.prank(deployer);
        rollup.addValidator(alice);

        // Alice should now be whitelisted
        assertTrue(rollup.hasRole(rollup.VALIDATOR_ROLE(), alice), "Expected address to be in validator whitelist");
        // Bob should not be in whitelist
        assertFalse(rollup.hasRole(rollup.VALIDATOR_ROLE(), bob), "Expected address not to be in validator whitelist");
    }

    function test_removeValidators_succeeds() external {
        // Initialize rollup with alice in the validator whitelist
        address[] memory validators = new address[](1);
        validators[0] = alice;
        _initializeRollup(
            0, // confirmationPeriod
            0, // challengePeriod
            0, // minimumAssertionPeriod
            1, // baseStakeAmount,
            0, // initialAssertionID
            0, // initialInboxSize
            validators // validator whitelist
        );

        // Removing validators from whitelist as non-admin should fail
        vm.expectRevert(
            abi.encodePacked(
                "AccessControl: account ",
                StringsUpgradeable.toHexString(alice),
                " is missing role ",
                StringsUpgradeable.toHexString(uint256(rollup.DEFAULT_ADMIN_ROLE()), 32)
            )
        );
        vm.prank(alice);
        rollup.removeValidator(alice);

        // Update validator whitelist as deployer, removing alice from whitelist
        vm.prank(deployer);
        rollup.removeValidator(alice);

        assertFalse(rollup.hasRole(rollup.VALIDATOR_ROLE(), alice), "Expected address to be in validator whitelist");
    }

    function test_removeOwnValidatorRole_succeeds() external {
        address[] memory validators = new address[](1);
        validators[0] = alice;
        // Initialize rollup with alice in validator whitelist
        _initializeRollup(
            0, // confirmationPeriod
            0, // challengePeriod
            0, // minimumAssertionPeriod
            1, // baseStakeAmount,
            0, // initialAssertionID
            0, // initialInboxSize
            validators // validator whitelist
        );

        // Alice removes themself from validator whitelist
        vm.prank(alice);
        rollup.removeOwnValidatorRole();

        assertFalse(rollup.hasRole(rollup.VALIDATOR_ROLE(), alice), "Expected address to be in validator whitelist");
    }

    function test_whitelistedFunctions_succeeds() external {
        address[] memory validators = new address[](1);
        validators[0] = alice;
        // Initialize rollup with alice in validator whitelist
        _initializeRollup(
            0, // confirmationPeriod
            0, // challengePeriod
            0, // minimumAssertionPeriod
            1, // baseStakeAmount,
            0, // initialAssertionID
            0, // initialInboxSize
            validators // validator whitelist
        );

        uint256 baseStakeAmount = rollup.baseStakeAmount();

        // Alice should be allowed to stake
        vm.prank(alice);
        //slither-disable-next-line arbitrary-send-eth
        rollup.stake{value: baseStakeAmount}();

        // Bob should not be allowed to stake
        vm.expectRevert(
            abi.encodePacked(
                "AccessControl: account ",
                StringsUpgradeable.toHexString(bob),
                " is missing role ",
                StringsUpgradeable.toHexString(uint256(rollup.VALIDATOR_ROLE()), 32)
            )
        );
        vm.prank(bob);
        rollup.stake{value: baseStakeAmount}();
    }

    function test_addValidators_roleAlreadyGranted_reverts() external {
        address[] memory validators = new address[](1);
        validators[0] = alice;
        // Initialize rollup with alice in validator whitelist
        _initializeRollup(
            0, // confirmationPeriod
            0, // challengePeriod
            0, // minimumAssertionPeriod
            1, // baseStakeAmount,
            0, // initialAssertionID
            0, // initialInboxSize
            validators // validator whitelist
        );

        // Add already whitelisted validator
        vm.expectRevert(IRollup.RoleAlreadyGranted.selector);
        vm.prank(deployer);
        rollup.addValidator(alice);
    }

    function test_removeValidators_noRoleToRevoke_reverts() external {
        // Initialize rollup with empty validator whitelist
        _initializeRollup(
            0, // confirmationPeriod
            0, // challengePeriod
            0, // minimumAssertionPeriod
            1, // baseStakeAmount,
            0, // initialAssertionID
            0, // initialInboxSize
            new address[](0) // validator whitelist
        );

        address[] memory validators = new address[](1);
        validators[0] = alice;

        // Remove validator not in the whitelist
        vm.expectRevert(IRollup.NoRoleToRevoke.selector);
        vm.prank(deployer);
        rollup.removeValidator(alice);
    }

    ////////////////
    // Staking
    ///////////////

    function testFuzz_stakers_notStaked_succeeds(
        uint256 confirmationPeriod,
        uint256 challengePeriod,
        uint256 minimumAssertionPeriod,
        uint256 baseStakeAmount,
        uint256 initialAssertionID
    ) external {
        _initializeRollup(
            confirmationPeriod, challengePeriod, minimumAssertionPeriod, baseStakeAmount, initialAssertionID
        );

        // Alice has not staked yet and therefore, this function should return `false`
        (bool isAliceStaked,,,) = rollup.stakers(alice);
        assertTrue(!isAliceStaked);
    }

    function testFuzz_stake_insufficentStake_reverts(
        uint256 confirmationPeriod,
        uint256 challengePeriod,
        uint256 minimumAssertionPeriod,
        uint256 initialAssertionID
    ) external {
        _initializeRollup(
            confirmationPeriod, challengePeriod, minimumAssertionPeriod, type(uint256).max, initialAssertionID
        );

        uint256 minimumAmount = rollup.baseStakeAmount();
        uint256 aliceBalance = alice.balance;

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

    function testFuzz_stake_succeeds(
        uint256 confirmationPeriod,
        uint256 challengePeriod,
        uint256 minimumAssertionPeriod,
        uint256 initialAssertionID
    ) external {
        _initializeRollup(confirmationPeriod, challengePeriod, minimumAssertionPeriod, 1000, initialAssertionID);

        uint256 initialStakers = rollup.numStakers();
        uint256 minimumAmount = rollup.baseStakeAmount();
        uint256 aliceBalance = alice.balance;

        assertGt(aliceBalance, minimumAmount, "Alice's Balance should be greater than stake amount for this test");

        _stake(alice, aliceBalance);

        uint256 finalStakers = rollup.numStakers();

        assertEq(alice.balance, 0, "Alice should not have any balance left");
        assertEq(finalStakers, (initialStakers + 1), "Number of stakers should increase by 1");

        uint256 amountStaked;
        uint256 assertionID;
        address challengeAddress;

        // stakers mapping gets updated
        (, amountStaked, assertionID, challengeAddress) = rollup.stakers(alice);

        assertEq(amountStaked, aliceBalance, "amountStaked not updated properly");
        assertEq(assertionID, rollup.lastConfirmedAssertionID(), "assertionID not updated properly");
        assertEq(challengeAddress, address(0), "challengeAddress not updated properly");
    }

    function test_stake_stakeWhenPaused_succeeds(
        uint256 confirmationPeriod,
        uint256 challengePeriod,
        uint256 minimumAssertionPeriod,
        uint256 initialAssertionID
    ) external {
        _initializeRollup(confirmationPeriod, challengePeriod, minimumAssertionPeriod, 1000, initialAssertionID);

        uint256 initialStakers = rollup.numStakers();
        uint256 minimumAmount = rollup.baseStakeAmount();
        uint256 aliceBalance = alice.balance;

        assertGt(aliceBalance, minimumAmount, "Alice's Balance should be greater than stake amount for this test");

        vm.prank(deployer);
        rollup.pause();

        _stake(alice, aliceBalance);

        uint256 finalStakers = rollup.numStakers();

        assertEq(alice.balance, 0, "Alice should not have any balance left");
        assertEq(finalStakers, (initialStakers + 1), "Number of stakers should increase by 1");

        uint256 amountStaked;
        uint256 assertionID;
        address challengeAddress;

        // stakers mapping gets updated
        (, amountStaked, assertionID, challengeAddress) = rollup.stakers(alice);

        assertEq(amountStaked, aliceBalance, "amountStaked not updated properly");
        assertEq(assertionID, rollup.lastConfirmedAssertionID(), "assertionID not updated properly");
        assertEq(challengeAddress, address(0), "challengeAddress not updated properly");
    }

    function testFuzz_stake_increaseStake_succeeds(
        uint256 confirmationPeriod,
        uint256 challengePeriod,
        uint256 minimumAssertionPeriod,
        uint256 initialAssertionID
    ) external {
        _initializeRollup(confirmationPeriod, challengePeriod, minimumAssertionPeriod, 1000, initialAssertionID);

        uint256 minimumAmount = rollup.baseStakeAmount();
        uint256 aliceBalanceInitial = alice.balance;
        uint256 bobBalance = bob.balance;

        assertGt(
            aliceBalanceInitial, minimumAmount, "Alice's Balance should be greater than stake amount for this test"
        );

        _stake(alice, aliceBalanceInitial);

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
        (bool sent,) = alice.call{value: bob.balance}("");
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

    //////////////////////
    // Remove Stake
    /////////////////////

    function testFuzz_removeStake_notStaked_reverts(
        uint256 confirmationPeriod,
        uint256 challengePeriod,
        uint256 minimumAssertionPeriod,
        uint256 baseStakeAmount,
        uint256 initialAssertionID
    ) external {
        _initializeRollup(
            confirmationPeriod, challengePeriod, minimumAssertionPeriod, baseStakeAmount, initialAssertionID
        );

        // Alice has not staked yet and therefore, this function should return `false`
        (bool isAliceStaked,,,) = rollup.stakers(alice);
        assertTrue(!isAliceStaked);

        // Since Alice is not staked, function unstake should also revert
        vm.expectRevert(IRollup.NotStaked.selector);
        vm.prank(alice);

        rollup.removeStake(address(alice));
    }

    function testFuzz_removeStake_notStakedThirdPartyCall_reverts(
        uint256 confirmationPeriod,
        uint256 challengePeriod,
        uint256 minimumAssertionPeriod,
        uint256 baseStakeAmount,
        uint256 initialAssertionID
    ) external {
        _initializeRollup(
            confirmationPeriod, challengePeriod, minimumAssertionPeriod, baseStakeAmount, initialAssertionID
        );

        // Alice has not staked yet and therefore, this function should return `false`
        (bool isAliceStaked,,,) = rollup.stakers(alice);
        assertTrue(!isAliceStaked);

        // Since Alice is not staked, function unstake should also revert
        vm.expectRevert(IRollup.NotStaked.selector);
        vm.prank(bob);

        rollup.removeStake(address(alice));
    }

    function testFuzz_removeStake_succeeds(
        uint256 confirmationPeriod,
        uint256 challengePeriod,
        uint256 minimumAssertionPeriod,
        uint256 initialAssertionID
    ) external {
        _initializeRollup(confirmationPeriod, challengePeriod, minimumAssertionPeriod, 1 ether, initialAssertionID);

        uint256 minimumAmount = rollup.baseStakeAmount();
        uint256 aliceBalance = alice.balance;

        // Let's stake something on behalf of Alice
        uint256 aliceAmountToStake = minimumAmount * 10;

        require(aliceBalance >= aliceAmountToStake, "Increase balance of Alice to proceed");
        _stake(alice, aliceAmountToStake);

        uint256 aliceBalanceBeforeRemoveStake = alice.balance;

        vm.prank(alice);
        rollup.removeStake(address(alice));

        (bool isStakedAfterRemoveStake,,,) = rollup.stakers(address(alice));

        uint256 aliceBalanceAfterRemoveStake = alice.balance;

        assertGt(aliceBalanceAfterRemoveStake, aliceBalanceBeforeRemoveStake);
        assertEq((aliceBalanceAfterRemoveStake - aliceBalanceBeforeRemoveStake), aliceAmountToStake);

        assertTrue(!isStakedAfterRemoveStake);
    }

    function testFuzz_removeStakePaused_succeeds(
        uint256 confirmationPeriod,
        uint256 challengePeriod,
        uint256 minimumAssertionPeriod,
        uint256 initialAssertionID
    ) external {
        _initializeRollup(confirmationPeriod, challengePeriod, minimumAssertionPeriod, 1 ether, initialAssertionID);

        uint256 minimumAmount = rollup.baseStakeAmount();
        uint256 aliceBalance = alice.balance;

        uint256 aliceAmountToStake = minimumAmount * 10;

        require(aliceBalance >= aliceAmountToStake, "Increase balance of Alice to proceed");
        _stake(alice, aliceAmountToStake);

        uint256 aliceBalanceBeforeRemoveStake = alice.balance;

        // as owner pause
        vm.prank(deployer);
        rollup.pause();

        // as alice, attempt to remove stake
        vm.prank(alice);
        rollup.removeStake(address(alice));

        (bool isStakedAfterRemoveStake,,,) = rollup.stakers(address(alice));

        uint256 aliceBalanceAfterRemoveStake = alice.balance;

        assertGt(aliceBalanceAfterRemoveStake, aliceBalanceBeforeRemoveStake);
        assertEq((aliceBalanceAfterRemoveStake - aliceBalanceBeforeRemoveStake), aliceAmountToStake);

        assertTrue(!isStakedAfterRemoveStake);
    }

    function testFuzz_removeStake_stakedOnUnconfirmedAssertion_reverts(
        uint256 confirmationPeriod,
        uint256 challengePeriod
    ) external {
        // Bounding it otherwise, function `newAssertionDeadline()` overflows
        confirmationPeriod = bound(confirmationPeriod, 1, type(uint128).max);
        _initializeRollup(confirmationPeriod, challengePeriod, 1 days, 1 ether, 0);

        uint256 minimumAmount = rollup.baseStakeAmount();
        uint256 aliceBalance = alice.balance;

        // Let's stake something on behalf of Alice
        uint256 aliceAmountToStake = minimumAmount * 10;

        require(aliceBalance >= aliceAmountToStake, "Increase balance of Alice to proceed");

        _stake(alice, aliceAmountToStake);

        // Now Alice should be staked
        uint256 stakerAssertionID;

        // stakers mapping gets updated
        (,, stakerAssertionID,) = rollup.stakers(alice);
        assertEq(stakerAssertionID, 0);

        _appendTxBatch();

        bytes32 mockStateCommitment = bytes32(hex"00");
        uint256 mockBlockNum = 1;

        // To avoid the MinimumAssertionPeriodNotPassed error, increase block.number
        vm.roll(block.number + rollup.minimumAssertionPeriod());

        assertEq(rollup.lastCreatedAssertionID(), 0, "The lastCreatedAssertionID should be 0 (genesis)");
        (,, uint256 assertionIDInitial,) = rollup.stakers(address(alice));

        assertEq(assertionIDInitial, 0);

        vm.prank(alice);
        rollup.createAssertion(mockStateCommitment, mockBlockNum, bytes32(0), 0);

        // The assertionID of alice should change after she called `createAssertion`
        (,, uint256 assertionIDFinal,) = rollup.stakers(address(alice));

        assertEq(assertionIDFinal, 1); // Alice is now staked on assertionID = 1 instead of assertionID = 0.

        // Try to remove Alice's stake
        vm.expectRevert(IRollup.StakedOnUnconfirmedAssertion.selector);
        vm.prank(bob); // validator only
        rollup.removeStake(address(alice));
    }

    ////////////////
    // Unstaking
    ///////////////

    function testFuzz_unstake_notStaked_reverts(
        uint256 randomAmount,
        uint256 confirmationPeriod,
        uint256 challengePeriod,
        uint256 minimumAssertionPeriod,
        uint256 baseStakeAmount,
        uint256 initialAssertionID
    ) external {
        _initializeRollup(
            confirmationPeriod, challengePeriod, minimumAssertionPeriod, baseStakeAmount, initialAssertionID
        );

        // Alice has not staked yet and therefore, this function should return `false`
        (bool isAliceStaked,,,) = rollup.stakers(alice);
        assertTrue(!isAliceStaked);

        // Since Alice is not staked, function unstake should also revert
        vm.expectRevert(IRollup.NotStaked.selector);
        vm.prank(alice);
        rollup.unstake(randomAmount);
    }

    function testFuzz_unstake_succeeds(
        uint256 confirmationPeriod,
        uint256 challengePeriod,
        uint256 minimumAssertionPeriod,
        uint256 amountToWithdraw,
        uint256 initialAssertionID
    ) external {
        _initializeRollup(confirmationPeriod, challengePeriod, minimumAssertionPeriod, 100000, initialAssertionID);

        uint256 minimumAmount = rollup.baseStakeAmount();
        uint256 aliceBalance = alice.balance;

        // Let's stake something on behalf of Alice
        uint256 aliceAmountToStake = minimumAmount * 10;

        require(aliceBalance >= aliceAmountToStake, "Increase balance of Alice to proceed");
        _stake(alice, aliceAmountToStake);

        uint256 aliceBalanceInitial = alice.balance;

        amountToWithdraw = _generateRandomUintInRange(1, (aliceAmountToStake - minimumAmount), amountToWithdraw);

        vm.prank(alice);
        rollup.unstake(amountToWithdraw);

        uint256 aliceBalanceFinal = alice.balance;

        assertEq((aliceBalanceFinal - aliceBalanceInitial), amountToWithdraw, "Desired amount could not be withdrawn.");
    }

    function testFuzz_unstakePaused_succeeds(
        uint256 confirmationPeriod,
        uint256 challengePeriod,
        uint256 minimumAssertionPeriod,
        uint256 amountToWithdraw,
        uint256 initialAssertionID
    ) external {
        _initializeRollup(confirmationPeriod, challengePeriod, minimumAssertionPeriod, 100000, initialAssertionID);

        uint256 minimumAmount = rollup.baseStakeAmount();
        uint256 aliceBalance = alice.balance;

        // Let's stake something on behalf of Alice
        uint256 aliceAmountToStake = minimumAmount * 10;

        require(aliceBalance >= aliceAmountToStake, "Increase balance of Alice to proceed");
        _stake(alice, aliceAmountToStake);

        uint256 aliceBalanceInitial = alice.balance;

        amountToWithdraw = _generateRandomUintInRange(1, (aliceAmountToStake - minimumAmount), amountToWithdraw);

        // as the owner, pause
        vm.prank(deployer);
        rollup.pause();

        // as alice, attempt to unstake
        vm.prank(alice);
        rollup.unstake(amountToWithdraw);

        uint256 aliceBalanceFinal = alice.balance;

        assertEq((aliceBalanceFinal - aliceBalanceInitial), amountToWithdraw, "Desired amount could not be withdrawn.");
    }

    function testFuzz_unstake_insufficientStake_reverts(
        uint256 confirmationPeriod,
        uint256 challengePeriod,
        uint256 minimumAssertionPeriod,
        uint256 amountToWithdraw,
        uint256 initialAssertionID
    ) external {
        _initializeRollup(confirmationPeriod, challengePeriod, minimumAssertionPeriod, 100000, initialAssertionID);

        // Alice has not staked yet and therefore, this function should return `false`
        (bool isAliceStaked,,,) = rollup.stakers(alice);
        assertTrue(!isAliceStaked);

        uint256 minimumAmount = rollup.baseStakeAmount();
        uint256 aliceBalance = alice.balance;

        // Let's stake something on behalf of Alice
        uint256 aliceAmountToStake = minimumAmount * 10;

        require(aliceBalance >= aliceAmountToStake, "Increase balance of Alice to proceed");

        _stake(alice, aliceAmountToStake);

        amountToWithdraw =
            _generateRandomUintInRange((aliceAmountToStake - minimumAmount) + 1, type(uint256).max, amountToWithdraw);

        vm.expectRevert(IRollup.InsufficientStake.selector);
        vm.prank(alice);
        rollup.unstake(amountToWithdraw);
    }

    function testFuzz_unstake_stakedOnUnconfirmedAssertion_reverts(uint256 confirmationPeriod, uint256 challengePeriod)
        external
    {
        // Bounding it otherwise, function `newAssertionDeadline()` overflows
        confirmationPeriod = bound(confirmationPeriod, 1, type(uint128).max);
        _initializeRollup(confirmationPeriod, challengePeriod, 1 days, 1 ether, 0);

        uint256 minimumAmount = rollup.baseStakeAmount();
        uint256 aliceBalance = alice.balance;

        // Let's stake something on behalf of Alice
        uint256 aliceAmountToStake = minimumAmount * 10;

        require(aliceBalance >= aliceAmountToStake, "Increase balance of Alice to proceed");
        _stake(alice, aliceAmountToStake);

        // To avoid the MinimumAssertionPeriodNotPassed error, increase block.number
        vm.roll(block.number + rollup.minimumAssertionPeriod());

        _appendTxBatch();

        bytes32 mockStateCommitment = bytes32(hex"00");
        uint256 mockBlockNum = 1;

        vm.prank(alice);
        rollup.createAssertion(mockStateCommitment, mockBlockNum, bytes32(0), 0);

        // The assertionID of alice should change after she called `createAssertion`
        (,, uint256 assertionIDFinal,) = rollup.stakers(address(alice));
        assertEq(assertionIDFinal, 1); // Alice is now staked on assertionID = 1 instead of assertionID = 0.

        // Alice tries to unstake
        vm.prank(alice);
        vm.expectRevert(IRollup.StakedOnUnconfirmedAssertion.selector);
        rollup.unstake(aliceAmountToStake);
    }

    /////////////////////////
    // Advance Stake
    /////////////////////////

    function testFuzz_advanceStake_notStaked_reverts(
        uint256 confirmationPeriod,
        uint256 challengePeriod,
        uint256 minimumAssertionPeriod,
        uint256 baseStakeAmount,
        uint256 assertionID,
        uint256 initialAssertionID
    ) external {
        _initializeRollup(
            confirmationPeriod, challengePeriod, minimumAssertionPeriod, baseStakeAmount, initialAssertionID
        );

        // Alice has not staked yet and therefore, this function should return `false`
        (bool isAliceStaked,,,) = rollup.stakers(alice);
        assertTrue(!isAliceStaked);

        // Since Alice is not staked, function advanceStake should also revert
        vm.expectRevert(IRollup.NotStaked.selector);
        vm.prank(alice);

        rollup.advanceStake(assertionID);
    }

    function testFuzz_advanceStake_assertionOutOfRange_reverts(
        uint256 confirmationPeriod,
        uint256 challengePeriod,
        uint256 minimumAssertionPeriod,
        uint256 assertionID,
        uint256 initialAssertionID
    ) external {
        _initializeRollup(confirmationPeriod, challengePeriod, minimumAssertionPeriod, 1 ether, initialAssertionID);

        uint256 minimumAmount = rollup.baseStakeAmount();
        uint256 aliceBalance = alice.balance;

        // Let's stake something on behalf of Alice
        uint256 aliceAmountToStake = minimumAmount * 10;
        require(aliceBalance >= aliceAmountToStake, "Increase balance of Alice to proceed");

        _stake(alice, aliceAmountToStake);

        vm.expectRevert(IRollup.AssertionOutOfRange.selector);
        vm.prank(alice);
        rollup.advanceStake(assertionID);
    }

    function testFuzz_advanceStake_succeeds(uint256 confirmationPeriod, uint256 challengePeriod) external {
        // Bounding it otherwise, function `newAssertionDeadline()` overflows
        confirmationPeriod = bound(confirmationPeriod, 1, type(uint128).max);

        _initializeRollup(confirmationPeriod, challengePeriod, 1 days, 1 ether, 0);

        uint256 minimumAmount = rollup.baseStakeAmount();

        // Let's stake something on behalf of Bob and Alice
        uint256 aliceBalance = alice.balance;
        uint256 bobBalance = bob.balance;

        uint256 aliceAmountToStake = minimumAmount * 10;
        uint256 bobAmountToStake = minimumAmount * 10;

        require(aliceBalance >= aliceAmountToStake, "Increase balance of Alice to proceed");
        require(bobBalance >= bobAmountToStake, "Increase balance of Bob to proceed");

        _stake(alice, aliceAmountToStake);
        _stake(bob, bobAmountToStake);

        // Let's create a brand new assertion, so that the lastCreatedAssertionID goes up and we can successfully advance stake to the new ID after that

        // Increase the sequencerInbox inboxSize with mock transactions we can assert on.
        _appendTxBatch();

        bytes32 mockStateCommitment = bytes32(hex"00");
        uint256 mockBlockNum = 1;

        // To avoid the MinimumAssertionPeriodNotPassed error, increase block.number
        vm.roll(block.number + rollup.minimumAssertionPeriod());

        assertEq(rollup.lastCreatedAssertionID(), 0, "The lastCreatedAssertionID should be 0 (genesis)");

        vm.prank(alice);
        rollup.createAssertion(mockStateCommitment, mockBlockNum, bytes32(0), 0);

        // A successful assertion should bump the lastCreatedAssertionID to 1.
        assertEq(rollup.lastCreatedAssertionID(), 1, "LastCreatedAssertionID not updated correctly");

        // The assertionID of alice should change after she called `createAssertion`
        (, uint256 aliceAmountStakedFinal, uint256 aliceAssertionIdFinal,) = rollup.stakers(address(alice));

        assertEq(aliceAmountToStake, aliceAmountStakedFinal);
        assertEq(aliceAssertionIdFinal, 1);

        // The assertionID of bob should remain unchanged
        (,, uint256 bobAssertionIdInitial,) = rollup.stakers(address(bob));
        assertEq(bobAssertionIdInitial, 0);

        // Advance stake of the staker
        // Since Alice's stake was already advanced when she called createAssertion, her call to `rollup.advanceStake` should fail
        vm.expectRevert(IRollup.AssertionOutOfRange.selector);
        vm.prank(alice);
        rollup.advanceStake(1);

        // Bob's call to `rollup.advanceStake` should succeed as he is still staked on the previous assertion
        vm.prank(bob);
        rollup.advanceStake(1);

        (, uint256 bobAmountStakedFinal, uint256 bobAssertionIdFinal,) = rollup.stakers(address(alice));

        assertEq(bobAmountToStake, bobAmountStakedFinal);
        assertEq(bobAssertionIdFinal, 1);
    }

    function testFuzz_advanceStake_paused_reverts(uint256 confirmationPeriod, uint256 challengePeriod) external {
        // Bounding it otherwise, function `newAssertionDeadline()` overflows
        confirmationPeriod = bound(confirmationPeriod, 1, type(uint128).max);

        _initializeRollup(confirmationPeriod, challengePeriod, 1 days, 1 ether, 0);

        uint256 minimumAmount = rollup.baseStakeAmount();

        // Let's stake something on behalf of Bob and Alice
        uint256 aliceBalance = alice.balance;
        uint256 bobBalance = bob.balance;

        uint256 aliceAmountToStake = minimumAmount * 10;
        uint256 bobAmountToStake = minimumAmount * 10;

        require(aliceBalance >= aliceAmountToStake, "Increase balance of Alice to proceed");
        require(bobBalance >= bobAmountToStake, "Increase balance of Bob to proceed");

        _stake(alice, aliceAmountToStake);
        _stake(bob, bobAmountToStake);

        // Let's create a brand new assertion, so that the lastCreatedAssertionID goes up and we can successfully advance stake to the new ID after that

        // Increase the sequencerInbox inboxSize with mock transactions we can assert on.
        _appendTxBatch();

        bytes32 mockStateCommitment = bytes32(hex"00");
        uint256 mockBlockNum = 1;

        // To avoid the MinimumAssertionPeriodNotPassed error, increase block.number
        vm.roll(block.number + rollup.minimumAssertionPeriod());

        assertEq(rollup.lastCreatedAssertionID(), 0, "The lastCreatedAssertionID should be 0 (genesis)");

        // run paused then unpause and proceed with test.
        vm.prank(deployer);
        rollup.pause();

        // try as alice
        vm.expectRevert("Pausable: paused");
        vm.prank(alice);
        rollup.createAssertion(mockStateCommitment, mockBlockNum, bytes32(0), 0);

        // unpause and continue setup
        vm.prank(deployer);
        rollup.unpause();

        // try again now that pause is over
        vm.prank(alice);
        rollup.createAssertion(mockStateCommitment, mockBlockNum, bytes32(0), 0);

        // A successful assertion should bump the lastCreatedAssertionID to 1.
        assertEq(rollup.lastCreatedAssertionID(), 1, "LastCreatedAssertionID not updated correctly");

        // The assertionID of alice should change after she called `createAssertion`
        (, uint256 aliceAmountStakedFinal, uint256 aliceAssertionIdFinal,) = rollup.stakers(address(alice));

        assertEq(aliceAmountToStake, aliceAmountStakedFinal);
        assertEq(aliceAssertionIdFinal, 1);

        // The assertionID of bob should remain unchanged
        (,, uint256 bobAssertionIdInitial,) = rollup.stakers(address(bob));
        assertEq(bobAssertionIdInitial, 0);

        vm.prank(deployer);
        rollup.pause();

        // Advance stake of the staker
        // Since Alice's stake was already advanced when she called createAssertion, her call to `rollup.advanceStake` should fail
        vm.expectRevert("Pausable: paused");
        vm.prank(alice);
        rollup.advanceStake(1);

        // Bob's call to `rollup.advanceStake` should succeed as he is still staked on the previous assertion
        vm.expectRevert("Pausable: paused");
        vm.prank(bob);
        rollup.advanceStake(1);

        (, uint256 bobAmountStakedFinal, uint256 bobAssertionIdFinal,) = rollup.stakers(address(alice));

        assertEq(bobAmountToStake, bobAmountStakedFinal);
        assertEq(bobAssertionIdFinal, 1);
    }

    /////////////////////////
    // Challenge Assertion
    /////////////////////////

    function testFuzz_challengeAssertion_wrongOrder_reverts(
        uint256 confirmationPeriod,
        uint256 challengePeriod,
        uint256 minimumAssertionPeriod,
        uint256 defenderAssertionID,
        uint256 challengerAssertionID,
        uint256 initialAssertionID
    ) public {
        // Initializing the rollup
        _initializeRollup(
            confirmationPeriod, challengePeriod, minimumAssertionPeriod, type(uint256).max, initialAssertionID
        );

        defenderAssertionID = bound(defenderAssertionID, challengerAssertionID, type(uint256).max);

        address[2] memory players;
        uint256[2] memory assertionIDs;

        players[0] = defender;
        players[1] = challenger;

        assertionIDs[0] = defenderAssertionID;
        assertionIDs[1] = challengerAssertionID;

        vm.expectRevert(IRollup.WrongOrder.selector);
        vm.prank(alice); // validator only
        rollup.challengeAssertion(players, assertionIDs);
    }

    function testFuzz_challengeAssertion_unproposedAssertion_reverts(
        uint256 confirmationPeriod,
        uint256 challengePeriod,
        uint256 minimumAssertionPeriod,
        uint256 initialAssertionID,
        uint256 challengerAssertionID,
        uint256 defenderAssertionID
    ) public {
        // Initializing the rollup
        initialAssertionID = bound(initialAssertionID, 0, (type(uint256).max - 10));
        _initializeRollup(
            confirmationPeriod, challengePeriod, minimumAssertionPeriod, type(uint256).max, initialAssertionID
        );

        uint256 lastCreatedAssertionID = rollup.lastCreatedAssertionID();

        challengerAssertionID = bound(challengerAssertionID, lastCreatedAssertionID + 1, type(uint256).max);
        defenderAssertionID = bound(defenderAssertionID, 0, challengerAssertionID - 1);

        address[2] memory players;
        uint256[2] memory assertionIDs;

        players[0] = defender;
        players[1] = challenger;

        assertionIDs[0] = defenderAssertionID;
        assertionIDs[1] = challengerAssertionID;

        vm.expectRevert(IRollup.UnproposedAssertion.selector);
        vm.prank(alice); // validator only
        rollup.challengeAssertion(players, assertionIDs);
    }

    function testFuzz_challengeAssertion_assertionAlreadyResolved_reverts(
        uint256 confirmationPeriod,
        uint256 challengePeriod
    ) public {
        // Initializing the rollup
        confirmationPeriod = bound(confirmationPeriod, 1, type(uint128).max);
        _initializeRollup(confirmationPeriod, challengePeriod, 1 days, 1 ether, 0);

        uint256 lastConfirmedAssertionID = rollup.lastConfirmedAssertionID();

        // Let's increase the lastCreatedAssertionID
        {
            uint256 minimumAmount = rollup.baseStakeAmount();
            uint256 aliceBalance = alice.balance;

            // Let's stake something on behalf of Alice
            uint256 aliceAmountToStake = minimumAmount * 10;

            require(aliceBalance >= aliceAmountToStake, "Increase balance of Alice to proceed");
            _stake(alice, aliceAmountToStake);

            _appendTxBatch();

            bytes32 mockStateCommitment = bytes32(hex"00");
            uint256 mockBlockNum = 1;

            // To avoid the MinimumAssertionPeriodNotPassed error, increase block.number
            vm.roll(block.number + rollup.minimumAssertionPeriod());

            vm.prank(alice);
            rollup.createAssertion(mockStateCommitment, mockBlockNum, bytes32(0), 0);
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
        vm.prank(alice); // validator only
        rollup.challengeAssertion(players, assertionIDs);
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
    function _appendTxBatch() internal {
        // txLengths is defined as: Array of lengths of each encoded tx in txBatch
        // txBatch is defined as: Batch of RLP-encoded transactions
        bytes memory txBatch = _helper_createTxBatch_hardcoded();

        // Pranking as the sequencer and calling appendTxBatch
        vm.prank(sequencerAddress);
        // Expect TxBatchAppended event
        vm.expectEmit(true, true, true, true, address(seqIn));
        emit TxBatchAppended();
        seqIn.appendTxBatch(txBatch);
    }

    function _initializeRollup(
        uint256 confirmationPeriod,
        uint256 challengePeriod,
        uint256 minimumAssertionPeriod,
        uint256 baseStakeAmount,
        uint256 initialAssertionID
    ) internal {
        address[] memory validators = new address[](2);
        validators[0] = alice;
        validators[1] = bob;
        _initializeRollup(
            confirmationPeriod,
            challengePeriod,
            minimumAssertionPeriod,
            baseStakeAmount,
            initialAssertionID,
            0, // initialInboxSize
            validators
        );
    }

    function _initializeRollup(
        uint256 confirmationPeriod,
        uint256 challengePeriod,
        uint256 minimumAssertionPeriod,
        uint256 baseStakeAmount,
        uint256 initialAssertionID,
        uint256 initialInboxSize,
        address[] memory validators
    ) internal {
        Config memory cfg = Config(
            sequencerAddress,
            address(seqIn), // sequencerInbox
            address(verifier),
            confirmationPeriod, //confirmationPeriod
            challengePeriod, //challengePeriod
            minimumAssertionPeriod, // minimumAssertionPeriod
            baseStakeAmount, //baseStakeAmount
            validators
        );

        InitialRollupState memory state =
            InitialRollupState(initialAssertionID, initialInboxSize, bytes32(""), bytes32(""));

        bytes memory initializingData = abi.encodeWithSelector(Rollup.initialize.selector, cfg);

        // Deploying the rollup contract as the rollup owner/deployer
        vm.startPrank(deployer);
        Rollup implementationRollup = new Rollup();
        rollup = Rollup(address(new ERC1967Proxy(address(implementationRollup), initializingData)));
        rollup.initializeGenesis(state);
        vm.stopPrank();

        // Check initial validators are in the whitelist
        for (uint256 i = 0; i < validators.length; i++) {
            assertTrue(
                rollup.hasRole(rollup.VALIDATOR_ROLE(), validators[i]), "Expected address to be in validator whitelist"
            );
        }
    }

    function _stake(address staker, uint256 amountToStake) internal {
        // Staker has not staked yet and therefore, this function should return `false`
        (bool isInitiallyStaked,, uint256 assertionIDInitial,) = rollup.stakers(staker);
        assertEq(assertionIDInitial, 0);
        assertTrue(!isInitiallyStaked);

        vm.prank(staker);
        // Calling the staking function as Alice
        //slither-disable-next-line arbitrary-send-eth
        rollup.stake{value: amountToStake}();

        // Staker should now be staked on the lastConfirmedAssertionId
        (bool isStaked, uint256 stakedAmount, uint256 stakedAssertionId,) = rollup.stakers(staker);
        assertEq(
            stakedAssertionId, rollup.lastConfirmedAssertionID(), "Staker is not staked on the lastConfirmedAssertionId"
        );
        assertEq(stakedAmount, amountToStake);
        assertTrue(isStaked);
    }
}
