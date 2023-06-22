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

import "../../../libraries/BytesLib.sol";

library MerkleRangeLib {
    using BytesLib for bytes;

    // Counts number of set bits (1's) in 32-bit unsigned integer
    function bitCount32(uint32 n) internal pure returns (uint32) {
        n = n - ((n >> 1) & 0x55555555);
        n = (n & 0x33333333) + ((n >> 2) & 0x33333333);

        return (((n + (n >> 4)) & 0xF0F0F0F) * 0x1010101) >> 24;
    }

    // Round 32-bit unsigned integer up to the nearest power of 2
    function roundUpToPowerOf2(uint32 n) internal pure returns (uint32) {
        if (bitCount32(n) == 1) return n;

        n |= n >> 1;
        n |= n >> 2;
        n |= n >> 4;
        n |= n >> 8;
        n |= n >> 16;

        return n + 1;
    }

    /**
     * @notice The merkle hash function
     */
    function hashNode(bytes32 left, bytes32 right) internal pure returns (bytes32) {
        // TODO: use assembly to save gas?
        return keccak256(abi.encodePacked(left, right));
    }

    /**
     * @notice Calculate the merkle tree root from leaves
     */
    function buildMerkleRoot(bytes32[] memory elements) internal pure returns (bytes32) {
        if (elements.length == 0) {
            return bytes32(0);
        }
        uint256 realLeafCount = elements.length;
        // Invariant: realLeafCount * 2 > leafCount
        uint256 leafCount = roundUpToPowerOf2(uint32(realLeafCount));
        while (leafCount > 1) {
            for (uint256 i = 0; i < leafCount / 2; i++) {
                bytes32 left = bytes32(0);
                if (i * 2 < realLeafCount) {
                    left = elements[i * 2];
                }
                bytes32 right = bytes32(0);
                if (i * 2 + 1 < realLeafCount) {
                    right = elements[i * 2 + 1];
                }
                elements[i] = hashNode(left, right);
            }
            leafCount = leafCount / 2;
        }
        return elements[0];
    }

    /**
     * @notice Calculate the merkle tree root using the range proof
     * @param elementNum The number of elements in the tree
     * @param start The index of the first element
     * @param elements The elements in the range
     * @param proof The range proof
     * @return root The merkle tree root
     */
    function getNewRootFromRangeProof(
        uint256 elementNum,
        uint256 start,
        bytes32[] memory elements,
        bytes32[] memory proof
    ) internal pure returns (bytes32) {
        start = start + elementNum;
        uint256 end = start + elements.length - 1;
        uint256 proofOffset = 0;
        while (start != 1) {
            uint256 hashRange = end - start + 1;
            if (start & 1 == 1) {
                elements[0] = hashNode(proof[proofOffset], elements[0]);
                proofOffset += 1;
                for (uint256 i = 1; i < hashRange / 2; i++) {
                    elements[i] = hashNode(elements[i * 2 - 1], elements[i * 2]);
                }
                uint256 last = hashRange / 2;
                if (end & 1 == 0) {
                    elements[last] = hashNode(elements[last * 2 - 1], proof[proofOffset]);
                    proofOffset += 1;
                } else {
                    elements[last] = hashNode(elements[last * 2 - 1], elements[last * 2]);
                }
            } else {
                for (uint256 i = 0; i < hashRange / 2; i++) {
                    elements[i] = hashNode(elements[i * 2], elements[i * 2 + 1]);
                }
                if (end & 1 == 0) {
                    uint256 last = hashRange / 2;
                    elements[last] = hashNode(elements[last * 2], proof[proofOffset]);
                    proofOffset += 1;
                }
            }
            start = start >> 1;
            end = end >> 1;
        }
        return elements[0];
    }

    /**
     * @notice Verify the range proof given a merkle tree root
     * @param root The merkle tree root
     * @param elementNum The number of elements in the tree
     * @param start The index of the first element
     * @param elements The elements in the range
     * @param proof The range proof
     * @return result The verification result, true if root is the correct merkle tree root
     */
    function verifyRangeProof(
        bytes32 root,
        uint256 elementNum,
        uint256 start,
        bytes32[] memory elements,
        bytes32[] memory proof
    ) internal pure returns (bool) {
        return root == getNewRootFromRangeProof(elementNum, start, elements, proof);
    }
}
