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
	"context"
	"encoding/json"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/specularl2/specular/services/el_clients/geth/specular/proof/proof"
	"github.com/specularl2/specular/services/el_clients/geth/specular/proof/prover"
	"github.com/specularl2/specular/services/el_clients/geth/specular/rollup/utils/fmt"
)

const (
	// defaultProveReexec is the number of blocks the prover is willing to go back
	// and reexecute to produce missing historical state necessary to run a specific
	// trace.
	defaultProveReexec = uint64(128)
)

type ProverConfig struct {
	Reexec *uint64
}

type ExecutionState struct {
	VMHash         common.Hash
	Block          *types.Block
	TransactionIdx uint64
	StepIdx        uint64
}

func (s *ExecutionState) MarshalJson() ([]byte, error) {
	return json.Marshal(&struct {
		VMHash         common.Hash `json:"vmHash"`
		L2GasUsed      *big.Int    `json:"l2GasUsed"`
		BlockHash      common.Hash `json:"blockHash"`
		TransactionIdx uint64      `json:"txnIdx"`
		StepIdx        uint64      `json:"stepIdx"`
	}{
		VMHash:         s.VMHash,
		BlockHash:      s.Block.Hash(),
		TransactionIdx: s.TransactionIdx,
		StepIdx:        s.StepIdx,
	})
}

func (s *ExecutionState) Hash() common.Hash {
	return s.VMHash
}

// This function generates execution states across blocks [startNum, endNum)
// For example there are 2 transactions: a, b
// The states are: inter-state before a, intra-states in a, inter-state before b (after a), intra-states in b, inter-state after b
func GenerateStates(
	backend Backend,
	ctx context.Context,
	startNum uint64,
	endNum uint64,
	config *ProverConfig,
) ([]*ExecutionState, error) {
	parent, err := backend.BlockByNumber(ctx, rpc.BlockNumber(startNum-1))
	if err != nil {
		return nil, err
	}
	reexec := defaultProveReexec
	if config != nil && config.Reexec != nil {
		reexec = *config.Reexec
	}
	statedb, err := backend.StateAtBlock(ctx, parent, reexec, nil, true, false)
	if err != nil {
		return nil, err
	}
	var (
		states []*ExecutionState
		block  *types.Block
	)
	for num := startNum; num < endNum; num++ {
		block, err = backend.BlockByNumber(ctx, rpc.BlockNumber(num))
		if err != nil {
			return nil, err
		}
		if block == nil {
			return nil, fmt.Errorf("block #%d not found", num)
		}
		signer := types.MakeSigner(backend.ChainConfig(), block.Number())
		blockCtx := core.NewEVMBlockContext(block.Header(), createChainContext(backend, ctx), nil)
		// Trace all the transactions contained within
		for i, tx := range block.Transactions() {
			// Push inter-state hash
			states = append(states, &ExecutionState{
				VMHash:         statedb.IntermediateRoot(backend.ChainConfig().IsEIP158(block.Number())),
				Block:          block,
				TransactionIdx: uint64(i),
				StepIdx:        0,
			})
			msg, _ := tx.AsMessage(signer, block.BaseFee())
			txContext := core.NewEVMTxContext(msg)
			prover := prover.NewStateGenerator()
			// Run the transaction with tracing enabled.
			vmenv := vm.NewEVM(blockCtx, txContext, statedb, backend.ChainConfig(), vm.Config{Debug: true, Tracer: prover, NoBaseFee: true})
			// Call Prepare to clear out the statedb access list
			statedb.Prepare(tx.Hash(), i)
			_, err := core.ApplyMessage(vmenv, msg, new(core.GasPool).AddGas(msg.Gas()))
			if err != nil {
				return nil, fmt.Errorf("tracing failed: %w", err)
			}
			generatedStates, err := prover.GetGeneratedStates()
			if err != nil {
				return nil, fmt.Errorf("tracing failed: %w", err)
			}
			for idx, s := range generatedStates {
				states = append(states, &ExecutionState{
					VMHash:         s.VMHash,
					Block:          block,
					TransactionIdx: uint64(i),
					StepIdx:        uint64(idx + 1),
				})
			}
		}
		// Get next statedb if we are not at the last block
		if num < endNum-1 {
			statedb, err = backend.StateAtBlock(ctx, block, reexec, statedb, true, false)
			if err != nil {
				return nil, err
			}
		}
	}
	states = append(states, &ExecutionState{
		VMHash:         block.Root(),
		Block:          block,
		TransactionIdx: uint64(len(block.Transactions())),
		StepIdx:        0,
	})
	return states, nil
}

func GenerateProof(backend Backend, ctx context.Context, startState *ExecutionState, config *ProverConfig) (*proof.OneStepProof, error) {
	if startState.Block == nil {
		return nil, fmt.Errorf("bad start state")
	}
	if startState.TransactionIdx >= uint64(len(startState.Block.Transactions())) {
		return nil, fmt.Errorf("bad start state")
	}
	reexec := defaultProveReexec
	if config != nil && config.Reexec != nil {
		reexec = *config.Reexec
	}
	msg, vmctx, statedb, err := backend.StateAtTransaction(ctx, startState.Block, int(startState.TransactionIdx), reexec)
	if err != nil {
		return nil, err
	}
	txContext := core.NewEVMTxContext(msg)
	prover := prover.NewProver(startState.VMHash, startState.StepIdx)
	// Run the transaction with tracing enabled.
	vmenv := vm.NewEVM(vmctx, txContext, statedb, backend.ChainConfig(), vm.Config{Debug: true, Tracer: prover, NoBaseFee: true})
	// Call Prepare to clear out the statedb access list
	txHash := startState.Block.Transactions()[startState.TransactionIdx].Hash()
	statedb.Prepare(txHash, int(startState.TransactionIdx))
	_, err = core.ApplyMessage(vmenv, msg, new(core.GasPool).AddGas(msg.Gas()))
	if err != nil {
		return nil, fmt.Errorf("tracing failed: %w", err)
	}
	return prover.GetProof()
}
