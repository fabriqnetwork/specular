// SPDX-License-Identifier: Apache-2.0

/*
 * Modifications Copyright 2022, Specular contributors
 *
 * This file was changed in accordance to Apache License, Version 2.0.
 *
 * Copyright 2021, Offchain Labs, Inc.
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

pragma solidity ^0.8.0;

import {Initializable} from "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import {UUPSUpgradeable} from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import {OwnableUpgradeable} from "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import {PausableUpgradeable} from "@openzeppelin/contracts-upgradeable/security/PausableUpgradeable.sol";
import {AccessControlDefaultAdminRulesUpgradeable} from
    "@openzeppelin/contracts-upgradeable/access/AccessControlDefaultAdminRulesUpgradeable.sol";
import {Hashing} from "./libraries/Hashing.sol";

import "./challenge/IChallenge.sol";
import "./challenge/SymChallenge.sol";
import "./challenge/ChallengeLib.sol";
import "./libraries/Errors.sol";
import "./IDAProvider.sol";
import "./IRollup.sol";

abstract contract RollupBase is
    IRollup,
    IChallengeResultReceiver,
    Initializable,
    UUPSUpgradeable,
    PausableUpgradeable,
    AccessControlDefaultAdminRulesUpgradeable
{
    // Access role identifiers
    bytes32 public constant VALIDATOR_ROLE = keccak256("VALIDATOR_ROLE");

    // Config parameters
    uint256 public confirmationPeriod; // number of L1 blocks
    uint256 public challengePeriod; // number of L1 blocks
    uint256 public minimumAssertionPeriod; // number of L1 blocks
    uint256 public baseStakeAmount; // number of stake tokens

    address public vault;
    IDAProvider public daProvider;
    IVerifier public verifier;

    function __RollupBase_init() internal onlyInitializing {
        __AccessControlDefaultAdminRules_init(3 days, msg.sender);
        __UUPSUpgradeable_init();
        __Pausable_init();
    }
}

contract Rollup is RollupBase {
    modifier validatorOnly() {
        _checkRole(VALIDATOR_ROLE);
        _;
    }

    modifier stakedOnly() {
        if (!isStaked(msg.sender)) {
            revert NotStaked();
        }
        _;
    }

    // Assertion state
    uint256 public lastResolvedAssertionID;
    uint256 public lastConfirmedAssertionID;
    uint256 public lastCreatedAssertionID;
    mapping(uint256 => Assertion) public assertions; // mapping from assertionID to assertion
    mapping(uint256 => AssertionState) private assertionState; // mapping from assertionID to assertion state

    // Staking state
    uint256 public numStakers; // current total number of stakers
    mapping(address => Staker) public stakers; // mapping from staker addresses to corresponding stakers
    mapping(address => uint256) public withdrawableFunds; // mapping from addresses to withdrawable funds (won in challenge)
    Zombie[] public zombies; // stores stakers that lost a challenge

    function initialize() public initializer {
        __RollupBase_init();
        // set values to 0, makes sure unpause checks work correctly
        verifier = IVerifier(address(0));
        daProvider = IDAProvider(address(0));
        vault = address(0);
        // initialize in a paused state to prevent core interactions until necessary values are set
        pause();
    }

    function initializeGenesis(InitialRollupState calldata _initialRollupState) external onlyRole(DEFAULT_ADMIN_ROLE) {
        require(lastCreatedAssertionID == 0, "Rollup: genesis already initialized");

        lastResolvedAssertionID = _initialRollupState.assertionID;
        lastConfirmedAssertionID = _initialRollupState.assertionID;
        lastCreatedAssertionID = _initialRollupState.assertionID;

        bytes32 initialStateCommitment =
            Hashing.createStateCommitmentV0(_initialRollupState.l2BlockHash, _initialRollupState.l2StateRoot);

        assertions[_initialRollupState.assertionID] = Assertion(
            initialStateCommitment,
            _initialRollupState.l2BlockNum, // blockNum (genesis)
            0, // parentID
            block.number, // deadline (unchallengeable)
            block.number, // proposal time
            0, // numStakers
            0 // childInboxSize
        );
        emit AssertionCreated(lastCreatedAssertionID, msg.sender, initialStateCommitment);
    }

    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }

    function pause() public onlyRole(DEFAULT_ADMIN_ROLE) {
        _pause();
    }

    function unpause() public onlyRole(DEFAULT_ADMIN_ROLE) {
        // verify we have values set to allow contract to unpause
        if (verifier == IVerifier(address(0)) || daProvider == IDAProvider(address(0)) || vault == address(0)) {
            revert ZeroAddress();
        }
        _unpause();
    }

    modifier hasConsensus() {
        require(lastConfirmedAssertionID == lastCreatedAssertionID, "Rollup: no consensus");
        _;
    }

    function _authorizeUpgrade(address) internal override onlyRole(DEFAULT_ADMIN_ROLE) whenPaused hasConsensus {}

    /// @inheritdoc IRollup
    function currentRequiredStake() public view override returns (uint256) {
        return baseStakeAmount;
    }

    /// @inheritdoc IRollup
    function confirmedBlockNum() public view override returns (uint256) {
        return assertions[lastConfirmedAssertionID].blockNum;
    }

    /// @inheritdoc IRollup
    function getStaker(address addr) external view override returns (Staker memory) {
        return stakers[addr];
    }

    /// @inheritdoc IRollup
    function getAssertion(uint256 assertionID) external view override returns (Assertion memory) {
        return assertions[assertionID];
    }

    /// @inheritdoc IRollup
    function getLastConfirmedAssertionID() external view override returns (uint256) {
        return lastConfirmedAssertionID;
    }

    /// @inheritdoc IRollup
    function isStakedOnAssertion(uint256 assertionID, address stakerAddress) external view override returns (bool) {
        return assertionState[assertionID].stakers[stakerAddress];
    }

    /// @inheritdoc IRollup
    function requireFirstUnresolvedAssertionIsConfirmable() public view override {
        if (lastResolvedAssertionID >= lastCreatedAssertionID) {
            revert NoUnresolvedAssertion();
        }
        uint256 firstUnresolvedID = lastResolvedAssertionID + 1;
        Assertion storage firstUnresolved = assertions[firstUnresolvedID];
        // (1) confirmation period has passed
        if (block.number < firstUnresolved.deadline) {
            revert ConfirmationPeriodPending();
        }
        // (2) predecessor has been confirmed
        if (firstUnresolved.parent != lastConfirmedAssertionID) {
            revert InvalidParent();
        }
        // (3) at least one staker is staked on the assertion.
        if (firstUnresolved.numStakers != countStakedZombies(firstUnresolvedID) + 1) {
            revert NoStaker();
        }
    }

    /// @inheritdoc IRollup
    function requireFirstUnresolvedAssertionIsRejectable(address stakerAddress) public view override {
        if (lastResolvedAssertionID >= lastCreatedAssertionID) {
            revert NoUnresolvedAssertion();
        }

        uint256 firstUnresolvedAssertionID = lastResolvedAssertionID + 1;
        Assertion storage firstUnresolvedAssertion = assertions[firstUnresolvedAssertionID];

        // First case - parent of first unresolved is last confirmed (`if` condition below). e.g.
        // [1] <- [3]           | valid chain ([1] is last confirmed, [3] is stakerAddress's unresolved assertion)
        //  ^---- [2]           | invalid chain ([2] is firstUnresolved)
        // Second case (trivial) - parent of first unresolved is not last confirmed. i.e.:
        //   parent is previous rejected, e.g.
        //   [1] <- [4]           | valid chain ([1] is last confirmed, [4] is stakerAddress's unresolved assertion)
        //   [2] <- [3]           | invalid chain ([3] is firstUnresolved)
        //   OR
        //   parent is previous confirmed, e.g.
        //   [1] <- [2] <- [4]    | valid chain ([2] is last confirmed, [4] is stakerAddress's unresolved assertion)
        //    ^---- [3]           | invalid chain ([3] is firstUnresolved)
        if (firstUnresolvedAssertion.parent == lastConfirmedAssertionID) {
            // 1a. confirmation period has passed.
            if (block.number < firstUnresolvedAssertion.deadline) {
                revert ConfirmationPeriodPending();
            }

            // 1b. at least one staker exists (on a sibling)
            // - stakerAddress is indeed a staker
            if (!stakers[stakerAddress].isStaked) {
                revert NotStaked();
            }
            // - staker's assertion can't be a ancestor of firstUnresolved (because staker's assertion is also unresolved)
            if (stakers[stakerAddress].assertionID < firstUnresolvedAssertionID) {
                revert AssertionAlreadyResolved();
            }
            AssertionState storage firstUnresolvedAssertionState = assertionState[firstUnresolvedAssertionID];

            // - staker's assertion can't be a descendant of firstUnresolved (because staker has never staked on firstUnresolved)
            if (firstUnresolvedAssertionState.stakers[stakerAddress]) {
                revert StakerStakedOnTarget();
            }
            // If a staker is staked on an assertion that is neither an ancestor nor a descendant of firstUnresolved, it must be a sibling, QED
            // 1c. no staker is staked on this assertion
            if (firstUnresolvedAssertion.numStakers != countStakedZombies(firstUnresolvedAssertionID)) {
                revert StakersPresent();
            }
        }
    }

    /// @inheritdoc IRollup
    function setConfig(Config calldata _config) external override onlyRole(DEFAULT_ADMIN_ROLE) {
        if (_config.vault == address(0) || _config.daProvider == address(0) || _config.verifier == address(0)) {
            revert ZeroAddress();
        }

        vault = _config.vault;
        daProvider = IDAProvider(_config.daProvider);
        verifier = IVerifier(_config.verifier);

        confirmationPeriod = _config.confirmationPeriod;
        challengePeriod = _config.challengePeriod;
        minimumAssertionPeriod = _config.minimumAssertionPeriod;
        baseStakeAmount = _config.baseStakeAmount;

        // Initialize role based access control
        for (uint256 i = 0; i < _config.validators.length; i++) {
            grantRole(VALIDATOR_ROLE, _config.validators[i]);
        }
        emit ConfigChanged();
    }

    /// @inheritdoc IRollup
    function setVault(address newVault) external override onlyRole(DEFAULT_ADMIN_ROLE) {
        if (newVault == address(0)) {
            revert ZeroAddress();
        }
        vault = newVault;
        emit ConfigChanged();
    }

    /// @inheritdoc IRollup
    function setDAProvider(address newDAProvider) external override onlyRole(DEFAULT_ADMIN_ROLE) {
        if (lastCreatedAssertionID > lastResolvedAssertionID) {
            revert InvalidConfigChange();
        }
        daProvider = IDAProvider(newDAProvider);
        emit ConfigChanged();
    }

    /// @inheritdoc IRollup
    function setVerifier(address newVerifier) external override onlyRole(DEFAULT_ADMIN_ROLE) {
        if (newVerifier == address(0)) {
            revert ZeroAddress();
        }
        verifier = IVerifier(newVerifier);
        emit ConfigChanged();
    }

    /// @inheritdoc IRollup
    function setConfirmationPeriod(uint256 newPeriod) external override onlyRole(DEFAULT_ADMIN_ROLE) {
        if (lastCreatedAssertionID > lastResolvedAssertionID) {
            revert InvalidConfigChange();
        }
        confirmationPeriod = newPeriod;
        emit ConfigChanged();
    }

    /// @inheritdoc IRollup
    function setChallengePeriod(uint256 newPeriod) external override onlyRole(DEFAULT_ADMIN_ROLE) {
        if (lastCreatedAssertionID > lastResolvedAssertionID) {
            revert InvalidConfigChange();
        }
        challengePeriod = newPeriod;
        emit ConfigChanged();
    }

    /// @inheritdoc IRollup
    function setMinimumAssertionPeriod(uint256 newPeriod) external override onlyRole(DEFAULT_ADMIN_ROLE) {
        if (lastCreatedAssertionID > lastResolvedAssertionID) {
            revert InvalidConfigChange();
        }
        minimumAssertionPeriod = newPeriod;
        emit ConfigChanged();
    }

    /// @inheritdoc IRollup
    function setBaseStakeAmount(uint256 newAmount) external override onlyRole(DEFAULT_ADMIN_ROLE) {
        if (lastCreatedAssertionID != lastResolvedAssertionID || newAmount > baseStakeAmount) {
            revert InvalidConfigChange();
        }
        baseStakeAmount = newAmount;
        emit ConfigChanged();
    }

    /// @inheritdoc IRollup
    function addValidator(address validator) external override onlyRole(DEFAULT_ADMIN_ROLE) {
        if (hasRole(VALIDATOR_ROLE, validator)) {
            revert RoleAlreadyGranted();
        }
        grantRole(VALIDATOR_ROLE, validator);
        emit ConfigChanged();
    }

    /// @inheritdoc IRollup
    function removeValidator(address validator) external override onlyRole(DEFAULT_ADMIN_ROLE) {
        if (!hasRole(VALIDATOR_ROLE, validator)) {
            revert NoRoleToRevoke();
        }
        revokeRole(VALIDATOR_ROLE, validator);
        emit ConfigChanged();
    }

    /// @inheritdoc IRollup
    function removeOwnValidatorRole() external override {
        if (!hasRole(VALIDATOR_ROLE, msg.sender)) {
            revert NoRoleToRevoke();
        }
        renounceRole(VALIDATOR_ROLE, msg.sender);
        emit ConfigChanged();
    }

    /// @inheritdoc IRollup
    function stake() external payable override validatorOnly {
        if (isStaked(msg.sender)) {
            stakers[msg.sender].amountStaked += msg.value;
        } else {
            if (msg.value < baseStakeAmount) {
                revert InsufficientStake();
            }
            stakers[msg.sender] = Staker(true, msg.value, 0, address(0));
            numStakers++;
            stakeOnAssertion(msg.sender, lastConfirmedAssertionID);
        }
    }

    /// @inheritdoc IRollup
    function unstake(uint256 stakeAmount) external override {
        requireStaked(msg.sender);
        // Require that staker is staked on a confirmed assertion.
        Staker storage staker = stakers[msg.sender];
        if (staker.assertionID > lastConfirmedAssertionID) {
            revert StakedOnUnconfirmedAssertion();
        }
        if (stakeAmount > staker.amountStaked - currentRequiredStake()) {
            revert InsufficientStake();
        }
        staker.amountStaked -= stakeAmount;
        // Note: we don't need to modify assertion state because you can only unstake from a confirmed assertion.
        (bool success,) = msg.sender.call{value: stakeAmount}("");
        if (!success) revert TransferFailed();
    }

    /// @inheritdoc IRollup
    function removeStake(address stakerAddress) external override validatorOnly {
        requireStaked(stakerAddress);
        // Require that staker is staked on a confirmed assertion.
        Staker storage staker = stakers[stakerAddress];
        if (staker.assertionID > lastConfirmedAssertionID) {
            revert StakedOnUnconfirmedAssertion();
        }

        uint256 stakerAmountStaked = staker.amountStaked;

        // Note: we don't need to modify assertion state because you can only unstake from a confirmed assertion.
        deleteStaker(stakerAddress);

        //slither-disable-next-line arbitrary-send-eth
        (bool success,) = stakerAddress.call{value: stakerAmountStaked}("");
        if (!success) revert TransferFailed();
    }

    /// @inheritdoc IRollup
    function advanceStake(uint256 assertionID) external override stakedOnly whenNotPaused validatorOnly {
        Staker storage staker = stakers[msg.sender];
        if (assertionID <= staker.assertionID || assertionID > lastCreatedAssertionID) {
            revert AssertionOutOfRange();
        }
        // TODO: allow arbitrary descendant of current staked assertionID, not just child.
        if (staker.assertionID != assertions[assertionID].parent) {
            revert ParentAssertionUnstaked();
        }
        stakeOnAssertion(msg.sender, assertionID);
    }

    /// @inheritdoc IRollup
    function withdraw() external override {
        uint256 withdrawableFund = withdrawableFunds[msg.sender];
        withdrawableFunds[msg.sender] = 0;
        (bool success,) = msg.sender.call{value: withdrawableFund}("");
        if (!success) revert TransferFailed();
    }

    /// @inheritdoc IRollup
    function createAssertion(bytes32 stateCommitment, uint256 blockNum, bytes32 l1BlockHash, uint256 l1BlockNumber)
        external
        override
        stakedOnly
        whenNotPaused
        validatorOnly
    {
        if (l1BlockHash != bytes32(0) && blockhash(l1BlockNumber) != l1BlockHash) {
            // This check allows the validator to propose an output based on a given L1 block,
            // without fear that it will be reorged out.
            // It will also revert if the blockheight provided is more than 256 blocks behind the
            // chain tip (as the hash will return as zero). This does open the door to a griefing
            // attack in which the validator's submission is censored until the block is no longer
            // retrievable, if the validator is experiencing this attack it can simply leave out the
            // blockhash value, and delay submission until it is confident that the L1 block is
            // finalized.
            revert MismatchingL1Blockhashes();
        }
        uint256 parentID = stakers[msg.sender].assertionID;
        Assertion storage parent = assertions[parentID];
        // Require that enough time has passed since the last assertion.
        if (block.number - parent.proposalTime < minimumAssertionPeriod) {
            revert MinimumAssertionPeriodNotPassed();
        }
        // Require that the assertion at least includes one block
        if (blockNum <= parent.blockNum) {
            revert EmptyAssertion();
        }

        // Initialize assertion.
        lastCreatedAssertionID++;
        emit AssertionCreated(lastCreatedAssertionID, msg.sender, stateCommitment);
        createAssertionHelper(lastCreatedAssertionID, stateCommitment, blockNum, parentID, newAssertionDeadline());

        // Update stake.
        stakeOnAssertion(msg.sender, lastCreatedAssertionID);
    }

    /// @inheritdoc IRollup
    function challengeAssertion(address[2] calldata players, uint256[2] calldata assertionIDs)
        external
        override
        validatorOnly
        returns (address)
    {
        uint256 defenderAssertionID = assertionIDs[0];
        uint256 parentID = assertions[defenderAssertionID].parent;
        uint256 challengerAssertionID = assertionIDs[1];

        // Require IDs ordered and in-range.
        if (defenderAssertionID >= challengerAssertionID) {
            revert WrongOrder();
        }
        if (challengerAssertionID > lastCreatedAssertionID) {
            revert UnproposedAssertion();
        }
        if (lastConfirmedAssertionID >= defenderAssertionID) {
            revert AssertionAlreadyResolved();
        }
        // Require that players have attested to sibling assertions.
        if (parentID != assertions[challengerAssertionID].parent) {
            revert NotSiblings();
        }
        // Require that neither player is currently engaged in a challenge.
        address defender = players[0];
        address challenger = players[1];
        requireUnchallengedStaker(defender);
        requireUnchallengedStaker(challenger);

        // TODO: Calculate upper limit for allowed node proposal time.

        // Initialize challenge.
        SymChallenge challenge = new SymChallenge();
        address challengeAddr = address(challenge);
        bytes32 startStateCommitment = assertions[parentID].stateCommitment;
        bytes32 endStateDefenderCommitment = assertions[defenderAssertionID].stateCommitment;
        bytes32 endStateChallengerCommitment = assertions[challengerAssertionID].stateCommitment;
        stakers[challenger].currentChallenge = challengeAddr;
        stakers[defender].currentChallenge = challengeAddr;
        emit AssertionChallenged(defenderAssertionID, challengeAddr);
        challenge.initialize(
            defender,
            challenger,
            verifier,
            daProvider,
            IChallengeResultReceiver(address(this)),
            startStateCommitment,
            endStateDefenderCommitment,
            endStateChallengerCommitment,
            challengePeriod
        );
        return challengeAddr;
    }

    /// @inheritdoc IRollup
    function confirmFirstUnresolvedAssertion() external override validatorOnly {
        requireFirstUnresolvedAssertionIsConfirmable();
        // Confirm assertion.
        // delete assertions[lastConfirmedAssertionID];
        lastResolvedAssertionID++;
        lastConfirmedAssertionID = lastResolvedAssertionID;
        emit AssertionConfirmed(lastResolvedAssertionID);
    }

    /// @inheritdoc IRollup
    function rejectFirstUnresolvedAssertion(address stakerAddress) external override validatorOnly {
        requireFirstUnresolvedAssertionIsRejectable(stakerAddress);
        // Reject assertion.
        lastResolvedAssertionID++;
        emit AssertionRejected(lastResolvedAssertionID);
        delete assertions[lastResolvedAssertionID];
    }

    /// @inheritdoc IChallengeResultReceiver
    function completeChallenge(address winner, address loser) external override {
        address challenge = getChallenge(winner, loser);
        if (msg.sender != challenge) {
            revert NotChallengeManager(msg.sender, challenge);
        }

        uint256 remainingLoserStake = stakers[loser].amountStaked;
        uint256 winnerStake = stakers[winner].amountStaked;
        if (remainingLoserStake > winnerStake) {
            // If loser has a higher stake than the winner, refund the difference.
            // Loser gets deleted anyways, so maybe unnecessary to set amountStaked.
            stakers[loser].amountStaked = winnerStake;
            withdrawableFunds[loser] += remainingLoserStake - winnerStake;
            remainingLoserStake = winnerStake;
        }
        // Reward the winner with half the remaining stake
        uint256 amountWon = remainingLoserStake / 2;
        stakers[winner].amountStaked += amountWon; // why +stake instead of +withdrawable?
        stakers[winner].currentChallenge = address(0);
        // Credit the other half to the vault address
        withdrawableFunds[vault] += remainingLoserStake - amountWon;
        // Turning loser into zombie renders the loser's remaining stake inaccessible.
        uint256 assertionID = stakers[loser].assertionID;
        deleteStaker(loser);
        // Track as zombie so we can account for it during assertion resolution.
        zombies.push(Zombie(loser, assertionID));
    }

    function isStaked(address addr) private view returns (bool) {
        return stakers[addr].isStaked;
    }

    /**
     * @notice Updates staker and assertion metadata.
     * @param stakerAddress Address of existing staker.
     * @param assertionID ID of existing assertion to stake on.
     */
    function stakeOnAssertion(address stakerAddress, uint256 assertionID) private {
        stakers[stakerAddress].assertionID = assertionID;
        assertions[assertionID].numStakers++;
        assertionState[assertionID].stakers[stakerAddress] = true;
        emit StakerStaked(stakerAddress, assertionID);
    }

    /**
     * @notice Creates a new assertion. See `Assertion` documentation.
     */
    function createAssertionHelper(
        uint256 assertionID,
        bytes32 stateCommitment,
        uint256 blockNum,
        uint256 parentID,
        uint256 deadline
    ) private {
        Assertion storage parentAssertion = assertions[parentID];
        AssertionState storage parentAssertionState = assertionState[parentID];
        // Siblings must have same inbox size.
        uint256 siblingConstraint = parentAssertion.childBlockNum;
        if (siblingConstraint == 0) {
            // Set the constraint for future siblings.
            parentAssertion.childBlockNum = blockNum;
        } else if (blockNum != siblingConstraint) {
            // Enforce the constraint if it's set.
            revert InvalidInboxSize();
        } else if (parentAssertionState.childStateCommitments[stateCommitment]) {
            revert DuplicateAssertion();
        }
        parentAssertionState.childStateCommitments[stateCommitment] = true;
        assertions[assertionID] = Assertion(
            stateCommitment,
            blockNum,
            parentID,
            deadline,
            block.number, // proposal time
            0, // numStakers
            0 // childInboxSize
        );
    }

    /**
     * @notice Deletes the staker from global state. Does not touch assertion staker state.
     * @param stakerAddress Address of the staker to delete
     */
    function deleteStaker(address stakerAddress) private {
        numStakers--;
        delete stakers[stakerAddress];
    }

    /**
     * @notice Checks to see whether the two stakers are in the same challenge
     * @param staker1Address Address of the first staker
     * @param staker2Address Address of the second staker
     * @return Address of the challenge that the two stakers are in
     */
    function getChallenge(address staker1Address, address staker2Address) private view returns (address) {
        Staker storage staker1 = stakers[staker1Address];
        Staker storage staker2 = stakers[staker2Address];
        address challenge = staker1.currentChallenge;
        if (challenge == address(0)) {
            revert NotInChallenge();
        }
        if (challenge != staker2.currentChallenge) {
            revert InDifferentChallenge(challenge, staker2.currentChallenge);
        }
        return challenge;
    }

    function newAssertionDeadline() private view returns (uint256) {
        // TODO: account for prev assertion, gas
        return block.number + confirmationPeriod;
    }

    // *****************
    // zombie processing
    // *****************

    /**
     * @notice Removes any zombies whose latest stake is earlier than the first unresolved assertion.
     * @dev Uses pop() instead of delete to prevent gaps, although order is not preserved
     */
    // function removeOldZombies() private {
    // }

    /**
     * @notice Counts the number of zombies staked on an assertion.
     * @dev O(n), where n is # of zombies (but is expected to be small).
     * This function could be uncallable if there are too many zombies. However,
     * removeOldZombies() can be used to remove any zombies that exist so that this
     * will then be callable.
     * @param assertionID The assertion on which to count staked zombies
     * @return The number of zombies staked on the assertion
     */
    function countStakedZombies(uint256 assertionID) private view returns (uint256) {
        uint256 numStakedZombies = 0;
        for (uint256 i = 0; i < zombies.length; i++) {
            if (assertionState[assertionID].stakers[zombies[i].stakerAddress]) {
                numStakedZombies++;
            }
        }
        return numStakedZombies;
    }

    // ************
    // requirements
    // ************

    function requireStaked(address stakerAddress) private view {
        if (!isStaked(stakerAddress)) {
            revert NotStaked();
        }
    }

    function requireUnchallengedStaker(address stakerAddress) private view {
        requireStaked(stakerAddress);
        if (stakers[stakerAddress].currentChallenge != address(0)) {
            revert ChallengedStaker();
        }
    }
}
