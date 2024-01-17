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

package prover

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"

	"github.com/specularL2/specular/services/sidecar/proof/proof"
)

type PlaceHolderProof struct{}

func (p *PlaceHolderProof) Encode() []byte {
	return []byte{'p', 'r', 'o', 'o', 'f'}
}

type OneStepProver struct {
	// Config
	target common.Hash
	step   uint64
}

func NewProver(target common.Hash, step uint64) *OneStepProver {
	return &OneStepProver{target: target, step: step}
}

func (l *OneStepProver) CaptureTxStart(gasLimit uint64) {}

func (l *OneStepProver) CaptureTxEnd(restGas uint64) {}

func (l *OneStepProver) CaptureStart(env *vm.EVM, from common.Address, to common.Address, create bool, input []byte, gas uint64, value *big.Int) {
}

func (l *OneStepProver) CaptureState(pc uint64, op vm.OpCode, gas, cost uint64, scope *vm.ScopeContext, rData []byte, depth int, err error) {
}

func (l *OneStepProver) CaptureEnter(typ vm.OpCode, from common.Address, to common.Address, input []byte, gas uint64, value *big.Int) {
}

func (l *OneStepProver) CaptureExit(output []byte, gasUsed uint64, err error) {
}

func (l *OneStepProver) CaptureFault(pc uint64, op vm.OpCode, gas, cost uint64, scope *vm.ScopeContext, depth int, err error) {
}

func (l *OneStepProver) CaptureEnd(output []byte, gasUsed uint64, err error) {
}

func (l *OneStepProver) GetProof() (*proof.OneStepProof, error) {
	proof := proof.EmptyProof()
	proof.AddProof(&PlaceHolderProof{})
	return proof, nil
}
