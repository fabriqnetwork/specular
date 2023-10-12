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
)

type GeneratedState struct {
	VMHash common.Hash
	Gas    uint64
}

type StateGenerator struct {
	// Global
	states []GeneratedState
}

func NewStateGenerator() *StateGenerator {
	return &StateGenerator{}
}

func (l *StateGenerator) CaptureTxStart(gasLimit uint64) {}

func (l *StateGenerator) CaptureTxEnd(restGas uint64) {}

func (l *StateGenerator) CaptureStart(env *vm.EVM, from common.Address, to common.Address, create bool, input []byte, gas uint64, value *big.Int) {
}

func (l *StateGenerator) CaptureState(pc uint64, op vm.OpCode, gas, cost uint64, scope *vm.ScopeContext, rData []byte, depth int, err error) {
	l.states = append(l.states, GeneratedState{common.Hash{}, gas})
}

func (l *StateGenerator) CaptureEnter(typ vm.OpCode, from common.Address, to common.Address, input []byte, gas uint64, value *big.Int) {
}

func (l *StateGenerator) CaptureExit(output []byte, gasUsed uint64, err error) {
}

func (l *StateGenerator) CaptureFault(pc uint64, op vm.OpCode, gas, cost uint64, scope *vm.ScopeContext, depth int, err error) {
}

func (l *StateGenerator) CaptureEnd(output []byte, gasUsed uint64, err error) {
}

func (l *StateGenerator) GetGeneratedStates() ([]GeneratedState, error) {
	return l.states, nil
}
