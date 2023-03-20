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
import "./verifier/IVerifier.sol";
import "./ChallengeBase.sol";
import "../IDAProvider.sol";
import "../libraries/Errors.sol";

contract AsymChallenge is ChallengeBase, IAsymChallenge {
    // See `ChallengeLib.computeBisectionHash` for the format of this commitment.
    bytes32 public bisectionHash;

    uint256 private constant MAX_BISECTION_DEGREE = 2;

    modifier onlyDefender() {
        if (msg.sender != challenger) {
            revert NotYourTurn();
        }
        _;
    }

    modifier onlyChallenger() {
        if (msg.sender != challenger) {
            revert NotYourTurn();
        }
        _;
    }

    function initialize(
        address _defender,
        address _challenger,
        bytes32 _bisectionHash,
        IVerifier _verifier,
        IDAProvider _daProvider,
        IChallengeResultReceiver _resultReceiver
    ) external {
        if (turn != Turn.NoChallenge) {
            revert AlreadyInitialized();
        }
        if (_defender == address(0) || _challenger == address(0)) {
            revert ZeroAddress();
        }
        defender = _defender;
        challenger = _challenger;
        _bisectionHash = bisectionHash;
        verifier = _verifier;
        daProvider = _daProvider;
        resultReceiver = _resultReceiver;

        turn = Turn.Defender;
        lastMoveBlock = block.number;
        // TODO(ujval): initialize timeout
        defenderTimeLeft = 10;
        challengerTimeLeft = 10;
    }
}
