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

    event ConfigurationChanged();

    Verifier verifier = new Verifier();

    function setUp() public virtual {
        utils = new Utils();
        users = utils.createUsers(4);

        alice = users[0];
        bob = users[1];
        deployer = users[2];
        sequencerAddress = users[3];
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
        bytes memory seqInInitData = abi.encodeWithSignature(
            "initialize(address)",
            sequencerAddress
        );
        vm.startPrank(deployer);
        implementationSequencer = new SequencerInbox();
        seqIn = SequencerInbox(
            address(
                new ERC1967Proxy(
                    address(implementationSequencer),
                    seqInInitData
                )
            )
        );
        vm.stopPrank();

        // Making sure that the proxy returns the correct proxy owner and sequencerAddress
        address sequencerInboxDeployer = seqIn.owner();
        assertEq(sequencerInboxDeployer, deployer);

        address fetchedSequencerAddress = seqIn.sequencerAddress();
        assertEq(fetchedSequencerAddress, sequencerAddress);
    }

    // test changes for parameters: confirmationPeriod, challengePeriod, minimumAssertionPeriod, baseStakeAmount
    function test_parametersChangedWhen_callerIsDeployer() public {
        bytes memory initializingData = abi.encodeWithSelector(
            Rollup.initialize.selector,
            sequencerAddress, // vault
            address(seqIn),
            address(verifier),
            0, //confirmationPeriod
            0, //challengePeriod
            0, // minimumAssertionPeriod
            20, //baseStakeAmount, it can only decrease so high initial value to test
            0, // initialAssertionID
            0, // initialInboxSize
            bytes32("")
        );

        vm.startPrank(deployer);

        Rollup implementationRollup = new Rollup(); // implementation contract
        rollup = Rollup(
            address(
                new ERC1967Proxy(
                    address(implementationRollup),
                    initializingData
                )
            )
        ); // The rollup contract (proxy, not implementation should have been initialized by now)

        assertEq(rollup.confirmationPeriod(), 0);
        vm.expectEmit(false, false, false, true);
        emit ConfigurationChanged();
        rollup.setConfirmationPeriod(10);
        assertEq(rollup.confirmationPeriod(), 10);

        assertEq(rollup.challengePeriod(), 0);
        vm.expectEmit(false, false, false, true);
        emit ConfigurationChanged();
        rollup.setChallengePeriod(10);
        assertEq(rollup.challengePeriod(), 10);

        assertEq(rollup.minimumAssertionPeriod(), 0);
        vm.expectEmit(false, false, false, true);
        emit ConfigurationChanged();
        rollup.setMinimumAssertionPeriod(10);
        assertEq(rollup.minimumAssertionPeriod(), 10);

        assertEq(rollup.baseStakeAmount(), 20);
        vm.expectEmit(false, false, false, true);
        emit ConfigurationChanged();
        rollup.setBaseStakeAmount(10);
        assertEq(rollup.baseStakeAmount(), 10);

        vm.stopPrank();
    }

    function test_parametersNotChangedWhen_callerIsNotDeployer() public {
        bytes memory initializingData = abi.encodeWithSelector(
            Rollup.initialize.selector,
            sequencerAddress, // vault
            address(seqIn),
            address(verifier),
            0, //confirmationPeriod
            0, //challengePeriod
            0, // minimumAssertionPeriod
            20, //baseStakeAmount,
            0, // initialAssertionID
            0, // initialInboxSize
            bytes32("")
        );

        vm.startPrank(deployer);

        Rollup implementationRollup = new Rollup(); // implementation contract
        rollup = Rollup(
            address(
                new ERC1967Proxy(
                    address(implementationRollup),
                    initializingData
                )
            )
        ); // The rollup contract (proxy, not implementation should have been initialized by now)

        vm.stopPrank();

        // now, try with sequencer address
        vm.startPrank(sequencerAddress);

        assertEq(rollup.confirmationPeriod(), 0);
        vm.expectRevert("Ownable: caller is not the owner");
        rollup.setConfirmationPeriod(10);
        assertEq(rollup.confirmationPeriod(), 0);

        assertEq(rollup.challengePeriod(), 0);
        vm.expectRevert("Ownable: caller is not the owner");
        rollup.setChallengePeriod(10);
        assertEq(rollup.challengePeriod(), 0);

        assertEq(rollup.minimumAssertionPeriod(), 0);
        vm.expectRevert("Ownable: caller is not the owner");
        rollup.setMinimumAssertionPeriod(10);
        assertEq(rollup.minimumAssertionPeriod(), 0);

        assertEq(rollup.baseStakeAmount(), 20);
        vm.expectRevert("Ownable: caller is not the owner");
        rollup.setBaseStakeAmount(10);
        vm.expectRevert("Ownable: caller is not the owner");
        rollup.setBaseStakeAmount(30);
        assertEq(rollup.baseStakeAmount(), 20);

        vm.stopPrank();
    }

    function test_baseStakeChangeFailsIf_setToBeHigher() public {
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
        rollup = Rollup(
            address(
                new ERC1967Proxy(
                    address(implementationRollup),
                    initializingData
                )
            )
        ); // The rollup contract (proxy, not implementation should have been initialized by now)

        assertEq(rollup.baseStakeAmount(), 0);
        vm.expectRevert("Cannot increase base stake amount");
        rollup.setBaseStakeAmount(10);
        assertEq(rollup.baseStakeAmount(), 0);

        vm.stopPrank();
    }
}
