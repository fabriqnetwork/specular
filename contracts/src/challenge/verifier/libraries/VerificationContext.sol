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

import "../../../ISequencerInbox.sol";
import "../../../libraries/DeserializationLib.sol";
import "../../../libraries/BytesLib.sol";

import "./MemoryLib.sol";
import "./EVMTypesLib.sol";

library VerificationContext {
    using BytesLib for bytes;
    using EVMTypesLib for EVMTypesLib.Transaction;

    struct RawContext {
        bytes encodedTx;
        // Transaction context.
        address coinbase;
        uint256 number;
        uint256 timestamp;
    }

    function txContextHash(RawContext calldata ctx) internal pure returns (bytes32) {
        return keccak256(abi.encodePacked(ctx.coinbase, ctx.number, ctx.timestamp));
    }

    struct Context {
        address coinbase;
        uint256 timestamp;
        uint256 number;
        address origin;
        EVMTypesLib.Transaction transaction;
        bytes32 inputRoot;
        bytes32 txHash;
    }

    function fromRaw(RawContext memory raw) internal pure returns (Context memory ctx) {
        ctx.coinbase = raw.coinbase;
        ctx.number = raw.number;
        ctx.timestamp = raw.timestamp;
        ctx.transaction = EVMTypesLib.decodeTransaction(raw.encodedTx);
    }

    function getCoinbase(Context memory ctx) internal pure returns (address) {
        return ctx.coinbase;
    }

    function getTimestamp(Context memory ctx) internal pure returns (uint256) {
        return ctx.timestamp;
    }

    function getBlockNumber(Context memory ctx) internal pure returns (uint256) {
        return ctx.number;
    }

    function getDifficulty(Context memory) internal pure returns (uint256) {
        return 1;
    }

    function getGasLimit(Context memory) internal pure returns (uint64) {
        return 80000000;
    }

    function getChainID(Context memory) internal pure returns (uint256) {
        return 13527;
    }

    // Transaction
    function getOrigin(Context memory ctx) internal pure returns (address) {
        return ctx.origin;
    }

    function getRecipient(Context memory ctx) internal pure returns (address) {
        return ctx.transaction.to;
    }

    function getValue(Context memory ctx) internal pure returns (uint256) {
        return ctx.transaction.value;
    }

    function getGas(Context memory ctx) internal pure returns (uint64) {
        return ctx.transaction.gas;
    }

    function getGasPrice(Context memory ctx) internal pure returns (uint256) {
        return ctx.transaction.gasPrice;
    }

    function getInput(Context memory ctx) internal pure returns (bytes memory) {
        return ctx.transaction.data;
    }

    function getInputSize(Context memory ctx) internal pure returns (uint64) {
        return uint64(ctx.transaction.data.length);
    }

    function getInputRoot(Context memory ctx) internal pure returns (bytes32) {
        if (ctx.inputRoot == 0x0) {
            ctx.inputRoot = MemoryLib.getMemoryRoot(ctx.transaction.data);
        }
        return ctx.inputRoot;
    }

    function getTxHash(Context memory ctx) internal pure returns (bytes32) {
        if (ctx.txHash == 0x0) {
            ctx.txHash = ctx.transaction.hashTransaction();
        }
        return ctx.transaction.hashTransaction();
    }
}
