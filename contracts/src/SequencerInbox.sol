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

import "./ISequencerInbox.sol";
import "./libraries/DeserializationLib.sol";
import "./libraries/Errors.sol";

contract SequencerInbox is ISequencerInbox, Initializable, UUPSUpgradeable, OwnableUpgradeable, PausableUpgradeable {
    // Current txBatch serialization version
    uint8 public constant currentTxBatchVersion = 0;

    address public sequencerAddress;

    function initialize(address _sequencerAddress) public initializer {
        if (_sequencerAddress == address(0)) {
            revert ZeroAddress();
        }
        sequencerAddress = _sequencerAddress;
        __Ownable_init();
        __Pausable_init();
        __UUPSUpgradeable_init();
    }

    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }

    function pause() public onlyOwner {
        _pause();
    }

    function unpause() public onlyOwner {
        _unpause();
    }

    function _authorizeUpgrade(address) internal override onlyOwner whenPaused {}

    /// @inheritdoc ISequencerInbox
    function appendTxBatch(bytes calldata txBatchData) external override whenNotPaused {
        if (msg.sender != sequencerAddress) {
            revert NotSequencer(msg.sender, sequencerAddress);
        }
        if (txBatchData.length == 0) {
            revert TxBatchDataUnderflow();
        }
        uint8 txBatchVersion = uint8(txBatchData[0]);
        if (txBatchVersion != currentTxBatchVersion) {
            revert TxBatchVersionIncorrect();
        }
        emit TxBatchAppended();
    }

    // TODO post EIP-4844: KZG proof verification
    // https://eips.ethereum.org/EIPS/eip-4844#point-evaluation-precompile

    function verifyTxInclusion(bytes calldata, bytes calldata) external pure override {
        revert ProofVerificationFailed();
    }
}
