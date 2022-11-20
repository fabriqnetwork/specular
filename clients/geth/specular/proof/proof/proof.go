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

package proof

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/specularl2/specular/clients/geth/specular/proof/state"
)

type Proof interface {
	Encode() []byte
}

type OneStepProof struct {
	Proofs []Proof
}

func EmptyProof() *OneStepProof {
	return &OneStepProof{}
}

func (p *OneStepProof) AddProof(proof Proof) {
	p.Proofs = append(p.Proofs, proof)
}

func (p *OneStepProof) Encode() []byte {
	if len(p.Proofs) == 0 {
		// Empty proof!
		return []byte{}
	}
	encodedLen := 0
	encodedProofs := make([][]byte, len(p.Proofs))
	for idx, proof := range p.Proofs {
		encodedProofs[idx] = proof.Encode()
		encodedLen += len(encodedProofs[idx])
	}
	encoded := make([]byte, encodedLen)
	offset := 0
	for _, encodedProof := range encodedProofs {
		copy(encoded[offset:], encodedProof)
		offset += len(encodedProof)
	}
	return encoded
}

func GetBlockInitiationProof(blockState *state.BlockState) *OneStepProof {
	proof := EmptyProof()
	proof.AddProof(BlockStateProofFromBlockState(blockState))
	return proof
}

func GetBlockFinalizationProof() *OneStepProof {
	proof := EmptyProof()
	// To clean up
	return proof
}

func GetTransactionInitaitionProof(tx *types.Transaction) *OneStepProof {
	proof := EmptyProof()
	// To clean up
	return proof
}
