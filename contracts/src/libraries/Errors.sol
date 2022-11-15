// SPDX-License-Identifier: Apache-2.0

/*
 * Copyright 2022, Specular contributors
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

/// @dev Thrown when unauthorized (!rollup) address calls an only-rollup function
/// @param sender Address of the caller
/// @param rollup The rollup address authorized to call this function
error NotRollup(address sender, address rollup);

/// @dev Thrown when unauthorized (!challenge) address calls an only-challenge function
/// @param sender Address of the caller
/// @param challenge The challenge address authorized to call this function
error NotChallenge(address sender, address challenge);

/// @dev Thrown when unauthorized (!sequencer) address calls an only-sequencer function
/// @param sender Address of the caller
/// @param sequencer The sequencer address authorized to call this function
error NotSequencer(address sender, address sequencer);

/// @dev Thrown when function is called with a zero address argument
error ZeroAddress();

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

/// @dev Thrown when the staker is currently in Challenge
error ChallengedStaker();

/// @dev Thrown when the staker is not in a challenge
error NotInChallenge();

/// @dev Thrown when the two stakers are in different challenge
/// @param staker1Challenge challenge address of staker 1
/// @param staker2Challenge challenge address of staker 2
error InDifferentChallenge(address staker1Challenge, address staker2Challenge);

/// @dev Thrown when a sender tries to create assertion before the minimum assertion time period
error MinimumAssertionPeriodNotPassed();

/// @dev Thrown when the L2 gas used by the assertion is more the max allowed limit.
error MaxGasLimitExceeded();

/// @dev Thrown when parent's statehash is not equal to the start state(or previous state)/
error PreviousStateHash();

/// @dev Thrown when a sender tries to create assertion without any tx.
error EmptyAssertion();

/// @dev Thrown when the requested assertion read past the end of current Inbox.
error InboxReadLimitExceeded();

/// @dev Thrown when there is no unresolved assertion
error NoUnresolvedAssertion();

/// @dev Thrown when the challenge period has not passed
error ChallengePeriodPending();

/// @dev Thrown when the assertion is already resolved
error AssertionAlreadyResolved();

/// @dev Thrown staker's assertion is descendant of firstUnresolved assertion
error StakerStakedOnTarget();

/// @dev Thrown when there are staker's present on the assertion
error StakersPresent();

/// @dev Thrown when the assertion's parent is not the last confirmed assertion
error InvalidParent();

/// @dev Thrown when all the stakers are not staked
error NotAllStaked();

/// @dev Thrown when there are zero stakers
error NoStaker();

/// @dev Thrown when the challenger and defender didn't attest to sibling assertions
error DifferentParent();

/// @dev Thrown when the challenge assertion Id is not ordered or in range.
error WrongOrder();

/// @dev Thrown when the challenger tries to challenge an unproposed assertion
error UnproposedAssertion();

/// @dev Thrown when sequencer tries to append an empty batch
error EmptyBatch();

/// @dev Thrown when the given tx inlcusion proof has incorrect accumulator or batch no.
error IncorrectAccOrBatch();
