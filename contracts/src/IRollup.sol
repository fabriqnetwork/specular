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

interface IRollup {

    event ConfigurationChanged();

    event AssertionCreated(uint256 assertionID, address asserterAddr, bytes32 vmHash);

    event AssertionChallenged(uint256 assertionID, address challengeAddr);

    event AssertionConfirmed(uint256 assertionID);

    event AssertionRejected(uint256 assertionID);

    event StakerStaked(address stakerAddr, uint256 assertionID);

    // TODO: Include errors thrown in function documentation.

    /// @dev Thrown when the new config parameter is invalid (configuration methods).
    error InvalidConfigChange();

    /// @dev Thrown when assertion creation requested with invalid inbox size.
    error InvalidInboxSize();

    /// @dev Thrown when assertion is a duplicate of an existing one.
    error DuplicateAssertion();

    /// @dev Thrown when address that have not staked any token calls a only-staked function
    error NotStaked();

    /// @dev Thrown when the function is called with Insufficient Stake
    error InsufficientStake();

    /// @dev Thrown when the caller is staked on unconfirmed assertion.
    error StakedOnUnconfirmedAssertion();

    /// @dev Thrown when transfer fails
    error TransferFailed();

    /// @dev Thrown when a staker tries to advance stake to invalid assertionId.
    error AssertionOutOfRange();

    /// @dev Thrown when a staker tries to advance stake to non-child assertion
    error ParentAssertionUnstaked();

    /// @dev Thrown when a sender tries to create assertion before the minimum assertion time period
    error MinimumAssertionPeriodNotPassed();

    /// @dev Thrown when parent's statehash is not equal to the start state(or previous state)/
    error PreviousStateHash();

    /// @dev Thrown when a sender tries to create assertion without any tx.
    error EmptyAssertion();

    /// @dev Thrown when the requested assertion read past the end of current Inbox.
    error InboxReadLimitExceeded();

    /// @dev Thrown when the challenge assertion Id is not ordered or in range.
    error WrongOrder();

    /// @dev Thrown when the challenger tries to challenge an unproposed assertion
    error UnproposedAssertion();

    /// @dev Thrown when the assertion is already resolved
    error AssertionAlreadyResolved();

    /// @dev Thrown when there is no unresolved assertion
    error NoUnresolvedAssertion();

    /// @dev Thrown when the confirmation period has not passed
    error ConfirmationPeriodPending();

    /// @dev Thrown when the challenger and defender didn't attest to sibling assertions
    error NotSiblings();

    /// @dev Thrown when the assertion's parent is not the last confirmed assertion
    error InvalidParent();

    /// @dev Thrown when the staker is not in a challenge
    error NotInChallenge();

    /// @dev Thrown when the two stakers are in different challenge
    /// @param staker1Challenge challenge address of staker 1
    /// @param staker2Challenge challenge address of staker 2
    error InDifferentChallenge(address staker1Challenge, address staker2Challenge);

    /// @dev Thrown when the staker is currently in Challenge
    error ChallengedStaker();

    /// @dev Thrown when all the stakers are not staked
    error NotAllStaked();

    /// @dev Thrown staker's assertion is descendant of firstUnresolved assertion
    error StakerStakedOnTarget();

    /// @dev Thrown when there are staker's present on the assertion
    error StakersPresent();

    /// @dev Thrown when there are zero stakers
    error NoStaker();

    struct Staker {
        bool isStaked;
        uint256 amountStaked;
        uint256 assertionID; // latest staked assertion ID
        address currentChallenge; // address(0) if none
    }

    struct Assertion {
        bytes32 stateHash; // Hash of execution state associated with assertion. Currently equiv to `vmHash`.
        uint256 inboxSize; // Inbox size this assertion advanced to
        uint256 parent; // Parent assertion ID
        uint256 deadline; // Dispute deadline (L1 block number)
        uint256 proposalTime; // L1 block number at which assertion was proposed
        // Staking state
        uint256 numStakers; // total number of stakers that have ever staked on this assertion. increasing only.
        // Child state
        uint256 childInboxSize; // child assertion inbox state
    }

    // *** Getters ***

    /**
     * @param addr Staker address.
     * @return Staker corresponding to address.
     */
    function getStaker(address addr) external view returns (Staker memory);

    /**
     * @param assertionID Assertion ID.
     * @return Assertion corresponding to ID.
     */
    function getAssertion(uint256 assertionID) external view returns (Assertion memory);

    /**
     * @return Whether or not the staker is staked on the assertion.
     */
    function isStakedOnAssertion(uint256 assertionID, address stakerAddress) external view returns (bool);

    /**
     * @return The current required stake amount.
     */
    function currentRequiredStake() external view returns (uint256);

    /**
     * @return confirmedInboxSize size of inbox confirmed
     */
    function confirmedInboxSize() external view returns (uint256);

    // *** Configuration ***

    /**
     * @notice Sets a new DA provider
     * @param newDAProvider New DA provider
     */
    function setDAProvider(address newDAProvider) external;

    /**
     * @notice Sets a new confirmation period
     * @param newPeriod New confirmation period
     */
    function setConfirmationPeriod(uint256 newPeriod) external;

    /**
     * @notice Sets a new challenge period
     * @param newPeriod New challenge period
     */
    function setChallengePeriod(uint256 newPeriod) external;

    /**
     * @notice Sets a new minimum assertion period
     * @param newPeriod New minimum assertion period
     */
    function setMinimumAssertionPeriod(uint256 newPeriod) external;

    /**
     * @notice Sets a new base stake amount.
     * @param newAmount New base stake amount; this can currently only be decreased.
     */
    function setBaseStakeAmount(uint256 newAmount) external;

    // *** State mutation ***

    /**
     * @notice Deposits stake on staker's current assertion (or the last confirmed assertion if not currently staked).
     * @dev Currently uses Ether to stake; Must be > than defined threshold if this is a new stake.
     */
    function stake() external payable;

    /**
     * @notice Withdraws stakeAmount from staker's stake if assertion it is staked on is confirmed.
     * @param stakeAmount Token amount to withdraw. Must be <= sender's current stake minus the current required stake.
     */
    function unstake(uint256 stakeAmount) external;

    /**
     * @notice Removes stakerAddress from the set of stakers and withdraws the full stake amount to stakerAddress.
     * This can be called by anyone since it is currently necessary to keep the chain progressing.
     * @param stakerAddress Address of staker for which to unstake.
     */
    function removeStake(address stakerAddress) external;

    /**
     * @notice Advances msg.sender's existing stake to assertionID.
     * @param assertionID ID of assertion to advance stake to. Currently this must be a child of the current assertion.
     * TODO: generalize to arbitrary descendants.
     */
    function advanceStake(uint256 assertionID) external;

    /**
     * @notice Withdraws all of msg.sender's withdrawable funds.
     */
    function withdraw() external;

    /**
     * @notice Creates a new DA representing the rollup state after executing a block of transactions (sequenced in SequencerInbox).
     * Block is represented by all transactions in range [prevInboxSize, inboxSize]. The latest staked DA of the sender
     * is considered to be the predecessor. Moves sender stake onto the new DA.
     *
     * Emits: `AssertionCreated` and `StakerStaked` events.
     *
     * @param vmHash New VM hash.
     * @param inboxSize Size of inbox corresponding to assertion (number of transactions).
     */
    function createAssertion(bytes32 vmHash, uint256 inboxSize) external;

    /**
     * @notice Initiates a dispute between a defender and challenger on an unconfirmed DA.
     * @param players Defender (first) and challenger (second) addresses. Must be staked on DAs on different branches.
     * @param assertionIDs Assertion IDs of the players engaged in the challenge. The first ID should be the earlier-created and is the one being challenged.
     * @return Newly created challenge contract address.
     */
    function challengeAssertion(address[2] calldata players, uint256[2] calldata assertionIDs)
        external
        returns (address);

    /**
     * @notice Confirms first unresolved assertion. Assertion is confirmed if and only if:
     * (1) there is at least one staker, and
     * (2) challenge period has passed, and
     * (3) predecessor has been confirmed, and
     * (4) all stakers are staked on the assertion.
     */
    function confirmFirstUnresolvedAssertion() external;

    /**
     * @notice Rejects first unresolved assertion. Assertion is rejected if and only if:
     * (1) all of the following are true:
     * (a) challenge period has passed, and
     * (b) at least one staker exists, and
     * (c) no staker remains staked on the assertion (all have been destroyed).
     * OR
     * (2) predecessor has been rejected
     * @param stakerAddress Address of a staker staked on a different branch to the first unresolved assertion.
     * If the first unresolved assertion's parent is confirmed, this parameter is used to establish that a staker exists
     * on a different branch of the assertion chain. This parameter is ignored when the parent of the first unresolved
     * assertion is not the last confirmed assertion.
     */
    function rejectFirstUnresolvedAssertion(address stakerAddress) external;
}
