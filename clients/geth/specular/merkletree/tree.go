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
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
)

// Given a list of elements (2^n), construct a merkle tree with size 2^(n+1)
// tree[0]: temp
// tree[1]: root
// in depth d, the first element is tree[2^d], the last element is tree[2^(d+1)-1]
// elements start from tree[2^n] to tree[2^(n+1)-1]

type MerkleTree struct {
	elementNum uint64
	nonEmpty   uint64
	tree       []common.Hash
}

func (t *MerkleTree) buildTree() {
	for i := t.elementNum - 1; i > 0; i-- {
		t.tree[i] = hashNode(t.tree[i<<1], t.tree[i<<1+1])
	}
}

func New(elements []*uint256.Int) *MerkleTree {
	elementNum := roundUpToPowerOf2(uint64(len(elements)))
	tree := make([]common.Hash, elementNum*2)
	for i := uint64(0); i < uint64(len(elements)); i++ {
		tree[elementNum+i] = elements[i].Bytes32()
	}
	t := &MerkleTree{
		elementNum: elementNum,
		nonEmpty:   uint64(len(elements)),
		tree:       tree,
	}
	t.buildTree()
	return t
}

func NewWithCapacity(elements []*uint256.Int, capacity uint64) *MerkleTree {
	elementNum := roundUpToPowerOf2(uint64(len(elements)))
	tree := make([]common.Hash, capacity*2)
	for i := uint64(0); i < uint64(len(elements)); i++ {
		tree[elementNum+i] = elements[i].Bytes32()
	}
	t := &MerkleTree{
		elementNum: elementNum,
		nonEmpty:   uint64(len(elements)),
		tree:       tree,
	}
	t.buildTree()
	return t
}

func NewWithPrealloc(elements []uint256.Int, prealloc []common.Hash) (*MerkleTree, error) {
	elementNum := roundUpToPowerOf2(uint64(len(elements)))
	if len(prealloc) < int(elementNum*2) {
		return nil, errors.New("prealloc size is too small")
	}
	tree := prealloc
	for i := uint64(0); i < uint64(len(elements)); i++ {
		tree[elementNum+i] = elements[i].Bytes32()
	}
	t := &MerkleTree{
		elementNum: elementNum,
		nonEmpty:   uint64(len(elements)),
		tree:       tree,
	}
	t.buildTree()
	return t, nil
}

func (t *MerkleTree) Rebuild(elements []*uint256.Int) error {
	elementNum := roundUpToPowerOf2(uint64(len(elements)))
	if elementNum != t.elementNum {
		return errors.New("element number mismatch")
	}
	for i := uint64(0); i < uint64(len(elements)); i++ {
		t.tree[elementNum+i] = elements[i].Bytes32()
	}
	for i := uint64(len(elements)); i < t.nonEmpty; i++ {
		t.tree[elementNum+i] = common.Hash{}
	}
	t.elementNum = elementNum
	t.nonEmpty = uint64(len(elements))
	t.buildTree()
	return nil
}

func (t *MerkleTree) GetElement(offset uint64) *uint256.Int {
	element := &uint256.Int{}
	element.SetBytes(t.tree[t.elementNum+offset][:])
	return element
}

func (t *MerkleTree) ElementCount() uint64 {
	return t.elementNum
}

func (t *MerkleTree) GetRoot() common.Hash {
	return t.tree[1]
}
