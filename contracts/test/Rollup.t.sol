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

import '@openzeppelin/contracts/token/ERC20/IERC20.sol';

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
    
    function setUp() public virtual override {
        BaseSetup.setUp();

        SequencerInbox _impl = new SequencerInbox();
        bytes memory data = abi.encodeWithSelector(SequencerInbox.initialize.selector, sequencer);
        address admin = address(47); // Random admin
        TransparentUpgradeableProxy proxy = new TransparentUpgradeableProxy(address(_impl), admin, data);

        seqIn = SequencerInbox(address(proxy));
    }

    function testFail_initializeRollup_ownerAddressZero() external {
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
        address proxyAdmin = makeAddr("Proxy Admin");
        TransparentUpgradeableProxy proxy = new TransparentUpgradeableProxy(
            address(_tempRollup), 
            proxyAdmin, 
            initializingData
        );

        rollup = Rollup(address(proxy));
    }

    function testFail_initializeRollup_verifierAddressZero() external {
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
        address proxyAdmin = makeAddr("Proxy Admin");
        TransparentUpgradeableProxy proxy = new TransparentUpgradeableProxy(
            address(_tempRollup), 
            proxyAdmin, 
            initializingData
        );

        rollup = Rollup(address(proxy));
    }

    function testFail_initializeRollup_sequencerInboxAddressZero() external {
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
}