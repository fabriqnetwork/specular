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
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/specularl2/specular/clients/geth/specular/proof/state"
)

type GeneratedIntraState struct {
	VMHash common.Hash
	Gas    uint64
}

type IntraStateGenerator struct {
	// Context (read-only)
	committedGlobalState vm.StateDB
	startInterState      *state.InterState
	blockHashTree        *state.BlockHashTree

	// Global
	env             *vm.EVM
	counter         int
	states          []GeneratedIntraState
	err             error
	done            bool
	selfDestructSet *state.SelfDestructSet
	logSeries       *state.LogSeries
	accessListTrie  *state.AccessListTrie

	// Current Call Frame
	callFlag       state.CallFlag
	lastState      *state.IntraState
	lastDepthState state.OneStepState
	input          *state.Memory
	out            uint64
	outSize        uint64
	selfDestructed bool
}

func NewIntraStateGenerator(
	committedGlobalState vm.StateDB,
	interState state.InterState,
	blockHashTree *state.BlockHashTree,
) *IntraStateGenerator {
	return &IntraStateGenerator{
		committedGlobalState: committedGlobalState,
		startInterState:      &interState,
		blockHashTree:        blockHashTree,
	}
}

func (l *IntraStateGenerator) CaptureTxStart(gasLimit uint64) {}

func (l *IntraStateGenerator) CaptureTxEnd(restGas uint64) {}

func (l *IntraStateGenerator) CaptureStart(env *vm.EVM, from common.Address, to common.Address, create bool, input []byte, gas uint64, value *big.Int) {
	l.env = env
	l.counter = 0
	if create {
		l.callFlag = state.CALLFLAG_CREATE
	} else {
		l.callFlag = state.CALLFLAG_CALL
	}
	l.input = state.NewMemoryFromBytes(input)
	l.accessListTrie = state.NewAccessListTrie()
	l.logSeries = state.EmptyLogSeries()
	l.selfDestructSet = state.NewSelfDestructSet()
	l.startInterState.GlobalState = env.StateDB.Copy() // This state includes gas-buying and nonce-increment
	l.lastDepthState = l.startInterState
	// log.Info("Capture Start", "from", from, "to", to)
}

// CaptureState will be called before the opcode execution
// vmerr is for stack validation and gas validation
// the execution error is captured in CaptureFault
func (l *IntraStateGenerator) CaptureState(pc uint64, op vm.OpCode, gas, cost uint64, scope *vm.ScopeContext, rData []byte, depth int, vmerr error) {
	if l.done {
		// Something went wrong during tracing, exit early
		return
	}
	// Construct intra state
	s := state.StateFromCaptured(
		l.committedGlobalState,
		l.selfDestructSet,
		l.logSeries,
		l.blockHashTree,
		l.accessListTrie,
		l.env,
		l.lastDepthState,
		l.callFlag,
		l.input,
		l.out, l.outSize, pc,
		op,
		gas, cost,
		scope,
		rData,
		depth,
	)
	l.states = append(l.states, GeneratedIntraState{s.Hash(), gas})
	l.lastState = s
	l.counter += 1

	if op == vm.CALL || op == vm.CALLCODE {
		l.out = scope.Stack.Back(5).Uint64()
		l.outSize = scope.Stack.Back(6).Uint64()
	} else if op == vm.DELEGATECALL || op == vm.STATICCALL {
		l.out = scope.Stack.Back(4).Uint64()
		l.outSize = scope.Stack.Back(5).Uint64()
	}
}

func (l *IntraStateGenerator) CaptureEnter(typ vm.OpCode, from common.Address, to common.Address, input []byte, gas uint64, value *big.Int) {
	if l.done {
		// Something went wrong during tracing, exit early
		return
	}
	if typ == vm.SELFDESTRUCT {
		// This enter is for the selfdestruct, record the address
		l.selfDestructed = true
		l.selfDestructSet = l.selfDestructSet.Add(from)
		return
	}
	l.callFlag = state.OpCodeToCallFlag(typ)
	l.lastDepthState = l.lastState.StateAsLastDepth(l.callFlag)
	l.input = state.NewMemoryFromBytes(input)
}

func (l *IntraStateGenerator) CaptureExit(output []byte, gasUsed uint64, vmerr error) {
	if l.done {
		// Something went wrong during tracing, exit early
		return
	}
	if l.selfDestructed {
		// This exit is for selfdestruct
		l.selfDestructed = false
		return
	}
	// TODO: next line seems unnecessary because CaptureEnd will be instantly called
	// if depth of the last state is 1
	if l.lastState.Depth > 1 {
		lastDepthState := l.lastDepthState.(*state.IntraState)
		l.callFlag = lastDepthState.CallFlag
		l.out = lastDepthState.Out
		l.outSize = lastDepthState.OutSize
		l.input = lastDepthState.InputData
		l.lastDepthState = lastDepthState.LastDepthState
		if vmerr != nil {
			// Call reverted, so revert the logs, selfdestructs, and access list changes
			l.logSeries = lastDepthState.LogSeries
			l.selfDestructSet = lastDepthState.SelfDestructSet
			l.accessListTrie = lastDepthState.AccessListTrie
		}
	}
}

// CaptureFault will be called when the stack/gas validation is passed but
// the execution failed. The current call will immediately be reverted.
// The error is handled in CaptureExit so nothing to do here.
func (l *IntraStateGenerator) CaptureFault(pc uint64, op vm.OpCode, gas, cost uint64, scope *vm.ScopeContext, depth int, err error) {
}

func (l *IntraStateGenerator) CaptureEnd(output []byte, gasUsed uint64, t time.Duration, err error) {
	// State generation finished, mark it as done
	l.done = true
}

func (l *IntraStateGenerator) GetGeneratedStates() ([]GeneratedIntraState, error) {
	if !l.done {
		return nil, fmt.Errorf("states generation not finished")
	}
	return l.states, l.err
}
