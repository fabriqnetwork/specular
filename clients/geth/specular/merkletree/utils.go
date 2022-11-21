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
	"math/bits"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/holiman/uint256"
)

func roundUpToPowerOf2(num uint64) uint64 {
	if bits.OnesCount64(num) == 1 {
		return num
	}
	return 1 << (64 - bits.LeadingZeros64(num))
}

func uint64To32Bytes(num uint64) [32]byte {
	return uint256.NewInt(num).Bytes32()
}

func hashLeaf(element *uint256.Int) common.Hash {
	elementBytes := element.Bytes32()
	return crypto.Keccak256Hash(PREFIX, elementBytes[:])
}

func hashNode(left, right common.Hash) common.Hash {
	return crypto.Keccak256Hash(left[:], right[:])
}

func hashMixedRoot(elementCount uint64, root common.Hash) common.Hash {
	elementCountBytes := uint64To32Bytes(elementCount)
	return crypto.Keccak256Hash(elementCountBytes[:], root[:])
}

func boolSetToUint256(set []bool) *uint256.Int {
	result := uint256.NewInt(0)
	one := uint256.NewInt(1)
	tmp := uint256.NewInt(0)
	for i, v := range set {
		if v {
			tmp.Lsh(one, uint(i))
			result.Or(result, tmp)
		}
	}
	return result
}

func boolSetToBytes32(set []bool) [32]byte {
	return boolSetToUint256(set).Bytes32()
}

func minimalCombinedProofIndex(elementCount uint64) uint64 {
	for shifts := uint64(0); shifts < 64; shifts++ {
		if elementCount&1 > 0 {
			return (elementCount & 0xFFFFFFFFFFFFFFFE) << shifts
		}
		elementCount >>= 1
	}
	return 0
}
