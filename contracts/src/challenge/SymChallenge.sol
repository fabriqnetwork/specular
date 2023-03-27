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

import "./IChallenge.sol";
import "./ChallengeBase.sol";
import "./ChallengeLib.sol";
import "./verifier/IVerifier.sol";
import "../IDAProvider.sol";
import "../libraries/DeserializationLib.sol";
import "../libraries/Errors.sol";

contract SymChallenge is ChallengeBase, ISymChallenge {
    // Previous state consistency could not be verified against `bisectionHash`.
    error PreviousStateInconsistent();
    // Tx context consistency could not be verified against ground truth.
    error TxContextInconsistent();

    uint256 private constant MAX_BISECTION_DEGREE = 2;

    // See `ChallengeLib.computeBisectionHash` for the format of this commitment.
    bytes32 public bisectionHash;
    // Initial state used to initialize bisectionHash (write-once).
    bytes32 private startStateHash;
    bytes32 private endStateHash;

    /**
     * @notice Ensures challenge has been initialized.
     */
    modifier postInitialization() {
        if (bisectionHash != 0) {
            revert NotInitialized();
        }
        _;
    }

    /**
     * @notice Initializes contract.
     * @param _defender Defending party.
     * @param _challenger Challenging party. Challenger starts.
     * @param _verifier Address of the verifier contract.
     * @param _daProvider DA provider.
     * @param _resultReceiver Address of contract that will receive the outcome (via callback `completeChallenge`).
     * @param _startStateHash Bisection root being challenged.
     * @param _endStateHash Bisection root being challenged.
     */
    function initialize(
        address _defender,
        address _challenger,
        IVerifier _verifier,
        IDAProvider _daProvider,
        IChallengeResultReceiver _resultReceiver,
        bytes32 _startStateHash,
        bytes32 _endStateHash,
        uint256 challengePeriod
    ) external {
        if (turn != Turn.NoChallenge) {
            revert AlreadyInitialized();
        }
        if (_defender == address(0) || _challenger == address(0)) {
            revert ZeroAddress();
        }
        defender = _defender;
        challenger = _challenger;
        verifier = _verifier;
        daProvider = _daProvider;
        resultReceiver = _resultReceiver;
        startStateHash = _startStateHash;
        endStateHash = _endStateHash;

        turn = Turn.Defender;
        lastMoveBlock = block.number;
        defenderTimeLeft = challengePeriod;
        challengerTimeLeft = challengePeriod;
    }

    function initializeChallengeLength(uint256 _numSteps) external override onlyOnTurn {
        if (bisectionHash != 0) {
            revert AlreadyInitialized();
        }
        require(_numSteps > 0, "INVALID_NUM_STEPS");
        bisectionHash = ChallengeLib.initialBisectionHash(startStateHash, endStateHash, _numSteps);
        // TODO: consider emitting a different event?
        emit Bisected(bisectionHash, 0, _numSteps);
    }

    function bisectExecution(
        bytes32[] calldata bisection,
        uint256 challengedSegmentIndex,
        bytes32[] calldata prevBisection,
        uint256 prevChallengedSegmentStart,
        uint256 prevChallengedSegmentLength
    ) external override onlyOnTurn postInitialization {
        // Verify provided prev bisection.
        bytes32 prevHash =
            ChallengeLib.computeBisectionHash(prevBisection, prevChallengedSegmentStart, prevChallengedSegmentLength);
        if (prevHash != bisectionHash) {
            revert PreviousStateInconsistent();
        }
        require(challengedSegmentIndex > 0 && challengedSegmentIndex < prevBisection.length, "INVALID_INDEX");
        // Require agreed upon start state hash and disagreed upon end state hash.
        require(bisection[0] == prevBisection[challengedSegmentIndex - 1], "INVALID_START");
        require(bisection[bisection.length - 1] != prevBisection[challengedSegmentIndex], "INVALID_END");

        // Compute segment start/length.
        uint256 challengedSegmentStart = prevChallengedSegmentStart;
        uint256 challengedSegmentLength = prevChallengedSegmentLength;
        if (prevBisection.length > 2) {
            // prevBisection.length == 2 means first round
            uint256 firstSegmentLength =
                ChallengeLib.firstSegmentLength(prevChallengedSegmentLength, MAX_BISECTION_DEGREE);
            uint256 otherSegmentLength =
                ChallengeLib.otherSegmentLength(prevChallengedSegmentLength, MAX_BISECTION_DEGREE);
            challengedSegmentLength = challengedSegmentIndex == 1 ? firstSegmentLength : otherSegmentLength;

            if (challengedSegmentIndex > 1) {
                challengedSegmentStart += firstSegmentLength + otherSegmentLength * (challengedSegmentIndex - 2);
            }
        }
        require(challengedSegmentLength > 1, "TOO_SHORT");

        // Require that bisection has the correct length. This is only ever less than BISECTION_DEGREE at the last bisection.
        uint256 target = challengedSegmentLength < MAX_BISECTION_DEGREE ? challengedSegmentLength : MAX_BISECTION_DEGREE;
        require(bisection.length == target + 1, "CUT_COUNT");

        // Compute new challenge state.
        bisectionHash = ChallengeLib.computeBisectionHash(bisection, challengedSegmentStart, challengedSegmentLength);
        emit Bisected(bisectionHash, challengedSegmentStart, challengedSegmentLength);
    }

    function verifyOneStepProof(
        bytes calldata oneStepProof,
        bytes calldata txInclusionProof,
        VerificationContextLib.RawContext calldata ctx,
        uint256 challengedStepIndex,
        bytes32[] calldata prevBisection,
        uint256 prevChallengedSegmentStart,
        uint256 prevChallengedSegmentLength
    ) external override onlyOnTurn {
        // Verify provided prev bisection.
        bytes32 prevHash =
            ChallengeLib.computeBisectionHash(prevBisection, prevChallengedSegmentStart, prevChallengedSegmentLength);
        if (prevHash != bisectionHash) {
            revert PreviousStateInconsistent();
        }
        require(challengedStepIndex > 0 && challengedStepIndex < prevBisection.length, "INVALID_INDEX");
        // Require that this is the last round.
        require(prevChallengedSegmentLength / MAX_BISECTION_DEGREE <= 1, "BISECTION_INCOMPLETE");
        {
            // Verify tx inclusion.
            daProvider.verifyTxInclusion(ctx.encodedTx, txInclusionProof);
            // Verify tx context consistency.
            // TODO: leaky abstraction (assumes `txInclusionProof` structure).
            (, bytes32 txContextHash) = DeserializationLib.deserializeBytes32(txInclusionProof, 0);
            if (VerificationContextLib.txContextHash(ctx) != txContextHash) {
                revert TxContextInconsistent();
            }
        }
        // Verify OSP.
        bytes32 endHash = verifier.verifyOneStepProof(prevBisection[challengedStepIndex - 1], ctx, oneStepProof);
        // Require that the end state differs from the counterparty's.
        if (endHash != prevBisection[challengedStepIndex]) {
            _currentWin(CompletionReason.OSP_VERIFIED);
        }
    }
}
