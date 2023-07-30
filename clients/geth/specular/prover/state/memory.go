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

package state

import (
	"encoding/binary"

	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	"github.com/specularl2/specular/clients/geth/specular/merkletree"
	prover_types "github.com/specularl2/specular/clients/geth/specular/prover/types"
)

const MaxMemoryCell uint64 = 1 << 17

type Memory struct {
	content []byte
	size    uint64
	tree    *merkletree.MerkleTree
}

func bytesToMemoryCells(values []byte) []*uint256.Int {
	size := len(values)
	bytesToPad := 32 - size%32
	if bytesToPad != 32 {
		values = append(values, make([]byte, bytesToPad)...)
	}
	cells := make([]*uint256.Int, len(values)/32)
	for idx := range cells {
		cells[idx] = new(uint256.Int).SetBytes(values[idx*32 : (idx+1)*32])
	}
	return cells
}

func NewMemoryFromBytesWithCapacity(values []byte, capacity uint64) *Memory {
	cells := bytesToMemoryCells(values)
	return &Memory{
		content: values,
		size:    uint64(len(values)),
		tree:    merkletree.NewWithCapacity(cells, capacity),
	}
}

func NewMemoryFromBytes(values []byte) *Memory {
	cells := bytesToMemoryCells(values)
	return &Memory{
		content: values,
		size:    uint64(len(values)),
		tree:    merkletree.New(cells),
	}
}

func MemoryFromEVMMemory(mem prover_types.L2ELClientMemoryInterface) *Memory {
	return NewMemoryFromBytesWithCapacity(mem.Data(), MaxMemoryCell)
}

func (m *Memory) Size() uint64 {
	return m.size
}

func (m *Memory) Data() []byte {
	return m.content[:m.size]
}

func (m *Memory) CellNum() uint64 {
	return m.tree.ElementCount()
}

func (m *Memory) Empty() bool {
	return m.size == 0
}

func (m *Memory) Cell(i uint64) *uint256.Int {
	return m.tree.GetElement(i)
}

func (m *Memory) Root() common.Hash {
	return m.tree.GetRoot()
}

func (m *Memory) GetProof(start, length uint64) []common.Hash {
	return m.tree.GetRangeProof(start, length)
}

func (m *Memory) EncodeState() []byte {
	encoded := make([]byte, 8)
	binary.BigEndian.PutUint64(encoded, m.size)
	if m.size != 0 {
		encoded = append(encoded, m.Root().Bytes()...)
	}
	return encoded
}
