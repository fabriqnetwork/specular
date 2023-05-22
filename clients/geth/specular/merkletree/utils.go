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
)

func roundUpToPowerOf2(num uint64) uint64 {
	if bits.OnesCount64(num) == 1 {
		return num
	}
	return 1 << (64 - bits.LeadingZeros64(num))
}

type hcachekey struct {
	left, right common.Hash
}

var hcache = make(map[hcachekey]common.Hash)

func hashNode(left, right common.Hash) common.Hash {
	if h, ok := hcache[hcachekey{left, right}]; ok {
		return h
	}
	h := crypto.Keccak256Hash(left[:], right[:])
	hcache[hcachekey{left, right}] = h
	return h
}
