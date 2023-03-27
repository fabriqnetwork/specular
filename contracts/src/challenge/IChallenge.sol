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

import "../IDAProvider.sol";
import "./verifier/IVerifier.sol";

interface IChallengeResultReceiver {
    /**
     * @notice Completes ongoing challenge. Callback, called by a challenge contract.
     * @param winner Address of winning staker.
     * @param loser Address of losing staker.
     */
    function completeChallenge(address winner, address loser) external;
}

interface IChallenge {
    enum CompletionReason {
        OSP_VERIFIED, // OSP verified by winner.
        TIMEOUT // Loser timed out before completing their round.
    }

    event Completed(address winner, address loser, CompletionReason reason);

    event Bisected(bytes32 challengeState, uint256 challengedSegmentStart, uint256 challengedSegmentLength);

    // Participant called function while it's not their turn.
    error NotYourTurn();
    // Participant did not respond prior before deadline.
    error DeadlineExpired();
    // Caller called function prematurely, before deadline passed.
    error DeadlineNotPassed();
    // Caller called function prematurely, before challenge initialized.
    error NotInitialized();
    // Caller called initialize function more than once.
    error AlreadyInitialized();

    /**
     * @notice Triggers completion of challenge protocol if a responder timed out.
     */
    function timeout() external;

    function currentResponder() external view returns (address);

    function currentResponderTimeLeft() external view returns (uint256);
}

/**
 * Symmetric challenge protocol.
 * @notice Protocol execution:
 * `initialize` (challenger, via Rollup) ->
 * `initializeChallengeLength` (defender) ->
 * `bisectExecution` (challenger, defender -- alternating) ->
 * `verifyOneStepProof` ->
 * `IResultReceiver.completeChallenge`
 */
interface ISymChallenge is IChallenge {
    /**
     * @notice Initializes the length of the challenge. Must be called by defender before bisection rounds begin.
     * @param _numSteps Number of steps executed from the start of the assertion to its end.
     * If this parameter is incorrect, the defender will be slashed (assuming successful execution of the protocol by the challenger).
     */
    function initializeChallengeLength(uint256 _numSteps) external;

    /**
     * @notice Bisects a segment. The challenged segment is defined by: {`challengedSegmentStart`, `challengedSegmentLength`, `bisection[0]`, `oldEndHash`}
     * @param bisection Bisection of challenged segment. Each element is a state hash (see `ChallengeLib.stateHash`).
     * The first element is the last agreed upon state hash. Must be of length MAX_BISECTION_LENGTH for all rounds except the last.
     * In the last round, the bisection segments must be single steps.
     * @param challengedSegmentIndex Index into `prevBisection`. Must be greater than 0 (since the first is agreed upon).
     * @param prevBisection Bisection in the preceding round.
     * @param prevChallengedSegmentStart Offset of the segment challenged in the preceding round (in steps).
     * Note: this is relative to the assertion being challenged (i.e. always between 0 and the initial `numSteps`).
     * @param prevChallengedSegmentLength Length of the segment challenged in the preceding round (in steps).
     */
    function bisectExecution(
        bytes32[] calldata bisection,
        uint256 challengedSegmentIndex,
        bytes32[] calldata prevBisection,
        uint256 prevChallengedSegmentStart,
        uint256 prevChallengedSegmentLength
    ) external;

    /**
     * @notice Verifies one step proof and completes challenge protocol.
     * @param oneStepProof TODO.
     * @param challengedStepIndex Index into `prevBisection`. Must be greater than 0 (since the first is agreed upon).
     * @param prevBisection Bisection in the preceding round. Each segment must now be of length 1 (i.e. a single step).
     * @param prevChallengedSegmentStart Offset of the segment challenged in the preceding round (in steps).
     * Note: this is relative to the assertion being challenged (i.e. always between 0 and the initial `numSteps`).
     * @param prevChallengedSegmentLength Length of the segment challenged in the preceding round (in steps).
     */
    function verifyOneStepProof(
        bytes calldata oneStepProof,
        bytes calldata txInclusionProof,
        VerificationContextLib.RawContext calldata ctx,
        uint256 challengedStepIndex,
        bytes32[] calldata prevBisection,
        uint256 prevChallengedSegmentStart,
        uint256 prevChallengedSegmentLength
    ) external;
}

// Assymetric challenge protocol.
interface IAsymChallenge is IChallenge {
// TODO.
}
