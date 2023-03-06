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
        {
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
        }

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
