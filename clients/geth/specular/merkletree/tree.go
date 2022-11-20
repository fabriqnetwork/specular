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

var PREFIX []byte = []byte{0x00}

type MerkleTree struct {
	elements []*uint256.Int
	tree     []common.Hash
}

func New(elements []*uint256.Int) *MerkleTree {
	if len(elements) == 0 {
		return &MerkleTree{}
	}
	balancedLeafCount := roundUpToPowerOf2(uint64(len(elements)))
	tree := make([]common.Hash, balancedLeafCount+uint64(len(elements)))
	for i := uint64(0); i < uint64(len(elements)); i++ {
		tree[balancedLeafCount+i] = hashLeaf(elements[i])
	}

	lowerBound := balancedLeafCount
	upperBound := balancedLeafCount + uint64(len(elements)) - 1

	for i := balancedLeafCount - 1; i > 0; i-- {
		index := i << 1
		if index > upperBound {
			continue
		}
		if index <= lowerBound {
			lowerBound >>= 1
			upperBound >>= 1
		}
		if index == upperBound {
			tree[i] = tree[index]
			continue
		}
		tree[i] = hashNode(tree[index], tree[index+1])
	}
	tree[0] = hashMixedRoot(uint64(len(elements)), tree[1])

	return &MerkleTree{
		elements: elements,
		tree:     tree,
	}
}

func (m *MerkleTree) Root() common.Hash {
	if len(m.elements) == 0 {
		return common.Hash{}
	}
	return m.tree[0]
}

func (m *MerkleTree) ElementCount() uint64 {
	return uint64(len(m.elements))
}

func (m *MerkleTree) Element(index uint64) *uint256.Int {
	return m.elements[index]
}
