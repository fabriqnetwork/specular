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

import {Utils} from "./utils/Utils.sol";
import {MockToken} from "./utils/MockToken.sol";

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";

import "../src/ISequencerInbox.sol";
import "../src/libraries/Errors.sol";

import {Verifier} from "../src/challenge/verifier/Verifier.sol";
import {Rollup} from "../src/Rollup.sol";
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
        users = utils.createUsers(6);

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

    function test_initializeRollup_valuesAfterInit() external {
        Rollup _tempRollup = new Rollup();

        uint256 confirmationPeriod = _generateRandomUint();
        uint256 challengePeriod = _generateRandomUint();
        uint256 minimumAssertionPeriod = _generateRandomUint();
        uint256 maxGasPerAssertion = _generateRandomUint();
        uint256 baseStakeAmount = _generateRandomUint();

        emit log_named_uint("confirmationPeriod", confirmationPeriod);
        emit log_named_uint("CP", challengePeriod);
        emit log_named_uint("BSA", baseStakeAmount);

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
    }

    /////////////////////////
    // Auxillary Functions
    /////////////////////////

    function _generateRandomUint() internal returns (uint256) {
        ++randomNonce;
        return uint256(keccak256(abi.encodePacked(block.timestamp, randomNonce)));
    }
}
