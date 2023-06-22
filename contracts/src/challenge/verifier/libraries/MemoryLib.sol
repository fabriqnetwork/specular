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

import "../../../libraries/DeserializationLib.sol";
import "../../../libraries/BytesLib.sol";
import "./OneStepProof.sol";
import "./Params.sol";
import "./MerkleRangeLib.sol";

library MemoryLib {
    using BytesLib for bytes;

    function calcCellNum(uint64 offset, uint64 length) internal pure returns (uint64) {
        return (offset + length + 31) / 32 - offset / 32;
    }

    function getMemoryRoot(bytes memory content) internal pure returns (bytes32) {
        uint64 cellNum = MemoryLib.calcCellNum(0, uint64(content.length));
        bytes32[] memory elements = new bytes32[](cellNum);
        for (uint256 i = 0; i < cellNum - 1; i++) {
            elements[i] = content.toBytes32(i * 32);
        }
        elements[cellNum - 1] = content.toBytes32Pad((cellNum - 1) * 32);
        return MerkleRangeLib.buildMerkleRoot(elements);
    }

    function decodeAndVerifyMemoryReadProof(
        OneStepProof.StateProof memory stateProof,
        bytes calldata encoded,
        uint64 offset,
        uint64 memoryOffset,
        uint64 memoryReadLength
    ) internal pure returns (uint64, bytes memory) {
        if (stateProof.memSize == 0 || memoryReadLength == 0) {
            return (offset, new bytes(memoryReadLength));
        }
        uint64 startCell = memoryOffset / 32;
        uint64 cellNum = calcCellNum(memoryOffset, memoryReadLength);
        // uint64 memoryCell = calcCellNum(0, stateProof.memSize);
        OneStepProof.MemoryReadProof memory readProof;
        OneStepProof.MemoryMerkleProof memory merkleProof;
        (offset, readProof) = OneStepProof.decodeMemoryReadProof(encoded, offset, cellNum);
        (offset, merkleProof) = OneStepProof.decodeMemoryMerkleProof(encoded, offset);
        require(
            MerkleRangeLib.verifyRangeProof(
                stateProof.memRoot, Params.MAX_MEMORY_SIZE, startCell, readProof.cells, merkleProof.proof
            ),
            "Bad Memory Proof"
        );
        bytes memory readContent = abi.encodePacked(readProof.cells).slice(memoryOffset % 32, memoryReadLength);
        return (offset, readContent);
    }

    // For input and return data
    function decodeAndVerifyROMReadProof(
        bytes32 romRoot,
        uint256 romSize,
        bytes calldata encoded,
        uint64 offset,
        uint64 memoryOffset,
        uint64 memoryReadLength
    ) internal pure returns (uint64, bytes memory) {
        if (memoryReadLength == 0) {
            return (offset, new bytes(memoryReadLength));
        }
        uint256 romTreeSize = MerkleRangeLib.roundUpToPowerOf2(uint32((romSize + 31) / 32));
        uint64 startCell = memoryOffset / 32;
        uint64 cellNum = calcCellNum(memoryOffset, memoryReadLength);
        // uint64 memoryCell = calcCellNum(0, stateProof.memSize);
        OneStepProof.MemoryReadProof memory readProof;
        OneStepProof.MemoryMerkleProof memory merkleProof;
        (offset, readProof) = OneStepProof.decodeMemoryReadProof(encoded, offset, cellNum);
        (offset, merkleProof) = OneStepProof.decodeMemoryMerkleProof(encoded, offset);
        require(
            MerkleRangeLib.verifyRangeProof(romRoot, romTreeSize, startCell, readProof.cells, merkleProof.proof),
            "Bad Memory Proof"
        );
        bytes memory readContent = abi.encodePacked(readProof.cells).slice(memoryOffset % 32, memoryReadLength);
        return (offset, readContent);
    }

    function decodeAndVerifyMemoryWriteProof(
        OneStepProof.StateProof memory stateProof,
        bytes calldata encoded,
        uint64 offset,
        uint64 memoryOffset,
        uint64 memoryWriteLength
    ) internal pure returns (uint64, bytes32, bytes memory) {
        if (memoryWriteLength == 0) {
            return (offset, stateProof.memRoot, new bytes(0));
        }
        if (stateProof.memSize == 0) {
            // Don't call decodeMemoryWriteProof if memory is empty
            // Instead, update memory root and size directly
            revert();
        }
        uint64 startCell = memoryOffset / 32;
        uint64 cellNum = calcCellNum(memoryOffset, memoryWriteLength);
        // uint64 memoryCell = calcCellNum(0, stateProof.memSize);
        OneStepProof.MemoryWriteProof memory writeProof;
        OneStepProof.MemoryMerkleProof memory merkleProof;
        (offset, writeProof) = OneStepProof.decodeMemoryWriteProof(encoded, offset, cellNum);
        (offset, merkleProof) = OneStepProof.decodeMemoryMerkleProof(encoded, offset);
        require(
            MerkleRangeLib.verifyRangeProof(
                stateProof.memRoot, Params.MAX_MEMORY_SIZE, startCell, writeProof.cells, merkleProof.proof
            ),
            "Bad Memory Read Proof"
        );
        bytes32 newRoot = MerkleRangeLib.getNewRootFromRangeProof(
            Params.MAX_MEMORY_SIZE, startCell, writeProof.updatedCells, merkleProof.proof
        );
        bytes memory writeContent =
            abi.encodePacked(writeProof.updatedCells).slice(memoryOffset % 32, memoryWriteLength);
        return (offset, newRoot, writeContent);
    }
}
