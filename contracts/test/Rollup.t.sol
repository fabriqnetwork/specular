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
import {Utils} from "./utils/Utils.sol";
import {IRollup} from "../src/IRollup.sol";
import {Verifier} from "../src/challenge/verifier/Verifier.sol";
import {Rollup} from "../src/Rollup.sol";
import {SequencerInbox} from "../src/SequencerInbox.sol";
import {RLPEncodedTransactionsUtil} from "./utils/RLPEncodedTransactions.sol";

contract RollupBaseSetup is Test, RLPEncodedTransactionsUtil {
    Utils internal utils;
    address payable[] internal users;

    address internal alice;
    address internal bob;
    address internal deployer;
    address internal sequencerAddress;

    Verifier verifier = new Verifier();

    function setUp() public virtual {
        utils = new Utils();
        users = utils.createUsers(4);

        deployer = users[0];
        sequencerAddress = users[1];
        alice = users[2];
        bob = users[3];
    }
}

contract RollupTest is RollupBaseSetup {
    Rollup public rollup;
    uint256 randomNonce;
    SequencerInbox public seqIn;
    SequencerInbox public implementationSequencer;

    address defender = makeAddr("defender");
    address challenger = makeAddr("challenger");

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
    }

    function test_fuzz_zeroValues(address _vault, address _sequencerInboxAddress, address _verifier) external {
        vm.assume(_vault >= address(0));
        vm.assume(_sequencerInboxAddress >= address(0));
        vm.assume(_verifier >= address(0));
        bytes memory initializingData = abi.encodeWithSelector(
            Rollup.initialize.selector,
            _vault, // vault
            _sequencerInboxAddress,
            _verifier,
            0, //confirmationPeriod
            0, //challengePeriod
            0, // minimumAssertionPeriod
            0, //baseStakeAmount,
            0, // initialAssertionID
            0, // initialInboxSize
            bytes32("")
        );
        if (_vault == address(0) || _sequencerInboxAddress == address(0) || _verifier == address(0)) {
            vm.startPrank(deployer);

            Rollup implementationRollup = new Rollup(); // implementation contract

            vm.expectRevert(ZeroAddress.selector);
            rollup = Rollup(address(new ERC1967Proxy(address(implementationRollup), initializingData)));
        }
    }

    function test_initializeRollup_cannotReinitializeRollup() external {
        bytes memory initializingData = abi.encodeWithSelector(
            Rollup.initialize.selector,
            sequencerAddress, // vault
            address(seqIn),
            address(verifier),
            0, //confirmationPeriod
            0, //challengePeriod
            0, // minimumAssertionPeriod
            0, //baseStakeAmount,
            0, // initialAssertionID
            0, // initialInboxSize
            bytes32("")
        );

        vm.startPrank(deployer);

        Rollup implementationRollup = new Rollup(); // implementation contract
        rollup = Rollup(address(new ERC1967Proxy(address(implementationRollup), initializingData))); // The rollup contract (proxy, not implementation should have been initialized by now)

        // Trying to call initialize for the second time
        vm.expectRevert("Initializable: contract is already initialized");

        rollup.initialize(
            sequencerAddress, // vault
            address(seqIn),
            address(verifier),
            0, //confirmationPeriod
            0, //challengePeriod
            0, // minimumAssertionPeriod
            0, //baseStakeAmount,
            0, // initialAssertionID
            0, // initialInboxSize
            bytes32("")
        );
    }

    function test_initializeRollup_valuesAfterInit(
        uint256 confirmationPeriod,
        uint256 challengePeriod,
        uint256 minimumAssertionPeriod,
        uint256 baseStakeAmount,
        uint256 initialInboxSize,
        uint256 initialAssertionID
    ) external {
        {
            bytes memory initializingData = abi.encodeWithSelector(
                Rollup.initialize.selector,
                sequencerAddress,
                address(seqIn), // sequencerInbox
                address(verifier),
                confirmationPeriod, //confirmationPeriod
                challengePeriod, //challengePeriod
                minimumAssertionPeriod, // minimumAssertionPeriod
                baseStakeAmount, //baseStakeAmount
                initialAssertionID,
                initialInboxSize,
                bytes32("") //initialVMHash
            );

            vm.startPrank(deployer);

            Rollup implementationRollup = new Rollup(); // implementation contract
            rollup = Rollup(address(new ERC1967Proxy(address(implementationRollup), initializingData))); // The rollup contract (proxy, not implementation should have been initialized by now)

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

    /////////////////////////
    // Challenge Assertion
    /////////////////////////

    function test_challengeAssertion_wrongOrderOfAssertionIDs(
        uint256 confirmationPeriod,
        uint256 challengePeriod,
        uint256 minimumAssertionPeriod,
        uint256 defenderAssertionID,
        uint256 challengerAssertionID,
        uint256 initialAssertionID,
        uint256 initialInboxSize
    ) public {
        // Initializing the rollup
        _initializeRollup(
            confirmationPeriod,
            challengePeriod,
            minimumAssertionPeriod,
            type(uint256).max,
            initialAssertionID,
            initialInboxSize
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
        uint256 initialAssertionID,
        uint256 initialInboxSize
    ) public {
        // Initializing the rollup
        initialAssertionID = bound(initialAssertionID, 0, (type(uint256).max - 10));
        _initializeRollup(
            confirmationPeriod,
            challengePeriod,
            minimumAssertionPeriod,
            type(uint256).max,
            initialAssertionID,
            initialInboxSize
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
        uint256 challengePeriod
    ) public {
        // Initializing the rollup
        confirmationPeriod = bound(confirmationPeriod, 1, type(uint128).max);
        _initializeRollup(confirmationPeriod, challengePeriod, 1 days, 1 ether, 0, 5);

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

            // To avoid the MinimumAssertionPeriodNotPassed error, increase block.number
            vm.warp(block.timestamp + 50 days);
            vm.roll(block.number + (50 * 86400) / 20);

            assertEq(rollup.lastCreatedAssertionID(), 0, "The lastCreatedAssertionID should be 0 (genesis)");
            (,, uint256 assertionIDInitial,) = rollup.stakers(address(alice));

            assertEq(assertionIDInitial, 0);

            vm.prank(alice);
            rollup.createAssertion(mockVmHash, mockInboxSize);

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
        uint256 confirmationPeriod,
        uint256 challengePeriod,
        uint256 minimumAssertionPeriod,
        uint256 baseStakeAmount,
        uint256 initialAssertionID,
        uint256 initialInboxSize
    ) internal {
        bytes memory initializingData = abi.encodeWithSelector(
            Rollup.initialize.selector,
            sequencerAddress,
            address(seqIn), // sequencerInbox
            address(verifier),
            confirmationPeriod, //confirmationPeriod
            challengePeriod, //challengePeriod
            minimumAssertionPeriod, // minimumAssertionPeriod
            baseStakeAmount, //baseStakeAmount
            initialAssertionID,
            initialInboxSize,
            bytes32("") //initialVMHash
        );

        // Deploying the rollup contract as the rollup owner/deployer
        vm.startPrank(deployer);
        Rollup implementationRollup = new Rollup();
        rollup = Rollup(address(new ERC1967Proxy(address(implementationRollup), initializingData)));
        vm.stopPrank();
    }
}
