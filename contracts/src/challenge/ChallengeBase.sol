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
import "./ChallengeLib.sol";

abstract contract ChallengeBase is IChallenge {
    enum Turn {
        NoChallenge,
        Challenger,
        Defender
    }

    Turn public turn;

    IVerifierEntry internal verifier;
    IDAProvider internal daProvider;
    IChallengeResultReceiver internal resultReceiver;

    // Challenge state
    address public defender;
    address public challenger;
    uint256 public lastMoveBlock;
    uint256 public defenderTimeLeft;
    uint256 public challengerTimeLeft;

    /**
     * @notice Pre-condition: `msg.sender` is correct and still has time remaining.
     * Post-condition: `turn` changes and `lastMoveBlock` set to current `block.number`.
     */
    modifier onlyOnTurn() {
        if (msg.sender != currentResponder()) {
            revert NotYourTurn();
        }
        if (block.number - lastMoveBlock > currentResponderTimeLeft()) {
            revert DeadlineExpired();
        }

        _;

        if (turn == Turn.Challenger) {
            challengerTimeLeft = challengerTimeLeft - (block.number - lastMoveBlock);
            turn = Turn.Defender;
        } else if (turn == Turn.Defender) {
            defenderTimeLeft = defenderTimeLeft - (block.number - lastMoveBlock);
            turn = Turn.Challenger;
        }
        lastMoveBlock = block.number;
    }

    function timeout() external override {
        if (block.number - lastMoveBlock <= currentResponderTimeLeft()) {
            revert DeadlineNotPassed();
        }
        if (turn == Turn.Defender) {
            _challengerWin(CompletionReason.TIMEOUT);
        } else {
            _asserterWin(CompletionReason.TIMEOUT);
        }
    }

    function currentResponder() public view override returns (address) {
        if (turn == Turn.Defender) {
            return defender;
        } else if (turn == Turn.Challenger) {
            return challenger;
        } else {
            revert NotInitialized();
        }
    }

    function currentResponderTimeLeft() public view override returns (uint256) {
        if (turn == Turn.Defender) {
            return defenderTimeLeft;
        } else if (turn == Turn.Challenger) {
            return challengerTimeLeft;
        } else {
            revert NotInitialized();
        }
    }

    function _currentWin(CompletionReason reason) internal {
        if (turn == Turn.Defender) {
            _asserterWin(reason);
        } else {
            _challengerWin(reason);
        }
    }

    function _asserterWin(CompletionReason reason) internal {
        emit Completed(defender, challenger, reason);
        resultReceiver.completeChallenge(defender, challenger); // safeSelfDestruct(msg.sender);
    }

    function _challengerWin(CompletionReason reason) internal {
        emit Completed(challenger, defender, reason);
        resultReceiver.completeChallenge(challenger, defender); // safeSelfDestruct(msg.sender);
    }
}
