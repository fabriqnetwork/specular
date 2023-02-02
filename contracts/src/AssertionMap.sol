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

pragma solidity ^0.8.0;

import "./libraries/Errors.sol";

// Exists only to reduce size of Rollup contract (maybe revert since Rollup fits under optimized compilation).
contract AssertionMap {
    error ChildInboxSizeMismatch();

    error SiblingStateHashExists();

    struct Assertion {
        bytes32 stateHash; // Hash of execution state associated with assertion (see `RollupLib.stateHash`)
        uint256 inboxSize; // Inbox size this assertion advanced to
        uint256 parent; // Parent assertion ID
        uint256 deadline; // Confirmation deadline (L1 block number)
        uint256 proposalTime; // L1 block number at which assertion was proposed
        // Staking state
        uint256 numStakers; // total number of stakers that have ever staked on this assertion. increasing only.
        mapping(address => bool) stakers; // all stakers that have ever staked on this assertion.
        // Child state
        uint256 childInboxSize; // child assertion inbox state
        mapping(bytes32 => bool) childStateHashes; // child assertion vm hashes
    }

    mapping(uint256 => Assertion) public assertions;
    address public rollupAddress;

    modifier rollupOnly() {
        if (msg.sender != rollupAddress) {
            revert NotRollup(msg.sender, rollupAddress);
        }
        _;
    }

    constructor(address _rollupAddress) {
        if (_rollupAddress == address(0)) {
            revert ZeroAddress();
        }
        rollupAddress = _rollupAddress;
    }

    function getStateHash(uint256 assertionID) external view returns (bytes32) {
        return assertions[assertionID].stateHash;
    }

    function getInboxSize(uint256 assertionID) external view returns (uint256) {
        return assertions[assertionID].inboxSize;
    }

    function getParentID(uint256 assertionID) external view returns (uint256) {
        return assertions[assertionID].parent;
    }

    function getDeadline(uint256 assertionID) external view returns (uint256) {
        return assertions[assertionID].deadline;
    }

    function getProposalTime(uint256 assertionID) external view returns (uint256) {
        return assertions[assertionID].proposalTime;
    }

    function getNumStakers(uint256 assertionID) external view returns (uint256) {
        return assertions[assertionID].numStakers;
    }

    function isStaker(uint256 assertionID, address stakerAddress) external view returns (bool) {
        return assertions[assertionID].stakers[stakerAddress];
    }

    event DebuggerString(string value);

    function createAssertion(
        uint256 assertionID,
        bytes32 stateHash,
        uint256 inboxSize,
        uint256 parentID,
        uint256 deadline
    ) external rollupOnly {
        emit DebuggerString("0");

        Assertion storage assertion = assertions[assertionID];
        Assertion storage parentAssertion = assertions[parentID];
        // Child assertions must have same inbox size

        emit DebuggerString("1");

        uint256 parentChildInboxSize = parentAssertion.childInboxSize;

        emit DebuggerString("2");

        if (parentChildInboxSize == 0) {
            parentAssertion.childInboxSize = inboxSize;
            emit DebuggerString("3");
        } else {
            emit DebuggerString("4");
            if (inboxSize != parentChildInboxSize) {
                revert ChildInboxSizeMismatch();
            }
        }

        emit DebuggerString("5");

        if (parentAssertion.childStateHashes[stateHash]) {
            revert SiblingStateHashExists();
        }

        emit DebuggerString("6");

        parentAssertion.childStateHashes[stateHash] = true;

        emit DebuggerString("7");

        assertion.stateHash = stateHash;
        assertion.inboxSize = inboxSize;
        assertion.parent = parentID;
        assertion.deadline = deadline;
        assertion.proposalTime = block.number;

        emit DebuggerString("7");
    }

    function stakeOnAssertion(uint256 assertionID, address stakerAddress) external rollupOnly {
        Assertion storage assertion = assertions[assertionID];
        assertion.stakers[stakerAddress] = true;
        assertion.numStakers++;
    }

    function deleteAssertion(uint256 assertionID) external rollupOnly {
        delete assertions[assertionID];
    }
}
