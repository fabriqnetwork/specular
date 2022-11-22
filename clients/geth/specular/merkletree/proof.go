// Copyright 2022, Specular contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package merkletree

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"github.com/holiman/uint256"
)

func (t *MerkleTree) generateSingleProof(index uint64) []common.Hash {
	leafCount := uint64(len(t.tree)) >> 1
	var decommitments []common.Hash
	for i := leafCount + index; i > 1; i >>= 1 {
		if i&1 == 0 {
			if i-1 < uint64(len(t.tree)) {
				decommitments = append(decommitments, t.tree[i-1])
			}
		} else {
			if i+1 < uint64(len(t.tree)) {
				decommitments = append(decommitments, t.tree[i+1])
			}
		}
	}
	proof := []common.Hash{uint256.NewInt(t.ElementCount()).Bytes32()}
	for i := len(decommitments) - 1; i >= 0; i-- {
		if decommitments[i] != (common.Hash{}) {
			proof = append(proof, decommitments[i])
		}
	}
	return proof
}

func (t *MerkleTree) generateMultiProof(indices []uint64) []common.Hash {
	elementCount := t.ElementCount()
	proof := []common.Hash{uint256.NewInt(elementCount).Bytes32()}
	known := make([]bool, len(t.tree))
	relevant := make([]bool, len(t.tree))
	leafcount := uint64(len(t.tree)) >> 1
	var flags, orders, skips []bool
	var decommitments []common.Hash

	for _, index := range indices {
		known[leafcount+index] = true
		relevant[(leafcount+index)>>1] = true
	}
	for i := leafcount - 1; i > 0; i-- {
		leftChildIndex := i << 1
		left := known[leftChildIndex]
		right := known[leftChildIndex+1]
		sibling := t.tree[leftChildIndex]
		if left {
			sibling = t.tree[leftChildIndex+1]
		}
		if left != right {
			decommitments = append(decommitments, sibling)
		}
		if relevant[i] {
			flags = append(flags, left == right)
			skips = append(skips, sibling == (common.Hash{}))
			orders = append(orders, left)
			relevant[i>>1] = true
		}
		known[i] = left || right
	}
	if len(flags) > 255 {
		log.Error("too many flags")
	}
	stopMask := uint256.NewInt(1)
	stopMask.Lsh(stopMask, uint(len(flags)))
	flagBits := boolSetToUint256(flags)
	flagBits.Or(flagBits, stopMask)
	proof = append(proof, flagBits.Bytes32())
	skipBits := boolSetToUint256(skips)
	skipBits.Or(skipBits, stopMask)
	proof = append(proof, skipBits.Bytes32())
	proof = append(proof, boolSetToBytes32(orders))
	for _, decommitment := range decommitments {
		if decommitment != (common.Hash{}) {
			proof = append(proof, decommitment)
		}
	}
	return proof
}

func (t *MerkleTree) GenerateProof(indices []uint64) []common.Hash {
	if len(indices) == 1 {
		return t.generateSingleProof(indices[0])
	}
	return t.generateMultiProof(indices)
}

func (t *MerkleTree) GenerateAppendProof() []common.Hash {
	elementCount := t.ElementCount()
	proof := []common.Hash{uint256.NewInt(elementCount).Bytes32()}
	if elementCount == 0 {
		return proof
	}
	leafCount := uint64(len(t.tree)) >> 1
	var decommitments []common.Hash
	for i := leafCount + elementCount; i > 1; i >>= 1 {
		if i&1 == 0 || i == 2 {
			decommitments = append(decommitments, t.tree[i-1])
		}
	}
	for i := len(decommitments) - 1; i >= 0; i-- {
		if decommitments[i] != (common.Hash{}) {
			proof = append(proof, decommitments[i])
		}
	}
	return proof
}

func (t *MerkleTree) GenerateCombinedProof(indices []uint64) ([]common.Hash, error) {
	if len(t.tree) == 0 {
		return nil, ErrEmptyTree
	}
	if indices[len(indices)-1] < minimalCombinedProofIndex(t.ElementCount()) {
		return nil, ErrIndexTooSmallForCombinedProof
	}
	return t.GenerateProof(indices), nil
}
