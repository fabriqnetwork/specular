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
	"github.com/holiman/uint256"
)

func (t *MerkleTree) GetRangeProof(start, length uint64) []common.Hash {
	// Add elementNum to get the index in the tree
	start += t.elementNum
	end := start + length - 1
	proof := make([]common.Hash, 0)
	for start != 1 {
		if start&1 == 1 {
			proof = append(proof, t.tree[start-1])
		}
		if end&1 == 0 {
			proof = append(proof, t.tree[end+1])
		}
		start >>= 1
		end >>= 1
	}
	return proof
}

func (t *MerkleTree) GetProof(offset uint64) []common.Hash {
	return t.GetRangeProof(offset, 1)
}

func VerifyProof(root common.Hash, elementNum, start uint64, elements []*uint256.Int, proof []common.Hash) bool {
	elementsAsHash := make([]common.Hash, len(elements))
	for i := 0; i < len(elements); i++ {
		elementsAsHash[i] = elements[i].Bytes32()
	}
	// Add elementNum to get the index in the tree
	start += elementNum
	end := start + uint64(len(elements)) - 1
	proofOffset := 0
	for start != 1 {
		hashRange := end - start + 1
		if start&1 == 1 {
			elementsAsHash[0] = hashNode(proof[proofOffset], elementsAsHash[0])
			proofOffset++
			for i := uint64(1); i < hashRange/2; i++ {
				elementsAsHash[i] = hashNode(elementsAsHash[2*i-1], elementsAsHash[2*i])
			}
			last := hashRange / 2
			if end&1 == 0 {
				elementsAsHash[last] = hashNode(elementsAsHash[2*last-1], proof[proofOffset])
				proofOffset++
			} else if last > 0 {
				elementsAsHash[last] = hashNode(elementsAsHash[2*last-1], elementsAsHash[2*last])
			}
		} else {
			for i := uint64(0); i < hashRange/2; i++ {
				elementsAsHash[i] = hashNode(elementsAsHash[2*i], elementsAsHash[2*i+1])
			}
			if end&1 == 0 {
				last := hashRange / 2
				elementsAsHash[last] = hashNode(elementsAsHash[2*last], proof[proofOffset])
				proofOffset++
			}
		}
		start >>= 1
		end >>= 1
	}
	return root == elementsAsHash[0]
}
