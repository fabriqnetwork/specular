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
	"github.com/specularl2/specular/clients/geth/specular/proof/proof"
	"github.com/specularl2/specular/clients/geth/specular/proof/prover"
	"github.com/specularl2/specular/clients/geth/specular/proof/state"
	"github.com/specularl2/specular/clients/geth/specular/rollup/utils/fmt"
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
	VMHash            common.Hash
	CumulativeGasUsed *big.Int
	BlockGasUsed      *big.Int
	StateType         state.StateType
	Block             *types.Block
	TransactionIdx    uint64
	StepIdx           uint64
}

func (s *ExecutionState) MarshalJson() ([]byte, error) {
	return json.Marshal(&struct {
		VMHash            common.Hash `json:"vmHash"`
		CumulativeGasUsed *big.Int    `json:"cumulativeGasUsed"`
		BlockGasUsed      *big.Int    `json:"blockGasUsed"`
		StateType         string      `json:"stateType"`
		BlockHash         common.Hash `json:"blockHash"`
		TransactionIdx    uint64      `json:"txnIdx"`
		StepIdx           uint64      `json:"stepIdx"`
	}{
		VMHash:            s.VMHash,
		CumulativeGasUsed: s.CumulativeGasUsed,
		BlockGasUsed:      s.BlockGasUsed,
		StateType:         string(s.StateType),
		BlockHash:         s.Block.Hash(),
		TransactionIdx:    s.TransactionIdx,
		StepIdx:           s.StepIdx,
	})
}

func (s *ExecutionState) Hash() common.Hash {
	return s.VMHash
}

// [GenerateStates] generates execution states across blocks [startNum, endNum)
//
// For example there are 2 blocks: block A with transactions A1, A2; block B without transactions.
//
// The states are (in order):
//  1. the BlockState of the block right before A
//  2. the InterState before A1
//  3. the IntraStates of A1 if A1 is not an EOA transfer
//  4. the InterState before A2 (after A1)
//  5. the IntraStates of A2 if A2 is not an EOA transfer
//  6. the InterState after A2
//  7. the BlockState of the block A
//  8. a dummy InterState
//  9. the BlockState of the block B
func GenerateStates(backend Backend, ctx context.Context, startGasUsed *big.Int, startNum, endNum uint64, config *ProverConfig) ([]*ExecutionState, error) {
	parent, err := backend.BlockByNumber(ctx, rpc.BlockNumber(startNum-1))
	if err != nil {
		return nil, err
	}
	reexec := defaultProveReexec
	if config != nil && config.Reexec != nil {
		reexec = *config.Reexec
	}
	// The statedb here is the state at the end of the parent blcok
	statedb, err := backend.StateAtBlock(ctx, parent, reexec, nil, true, false)
	if err != nil {
		return nil, err
	}
	var (
		states []*ExecutionState
		block  *types.Block
	)
	chainCtx := createChainContext(backend, ctx)
	cumulativeGasUsed := new(big.Int).Set(startGasUsed)

	// Push the block state of the parent block
	parentBlockCtx := core.NewEVMBlockContext(parent.Header(), chainCtx, nil)
	cumulativeGasUsedCopy := new(big.Int).Set(cumulativeGasUsed)
	blockHashTree, err := state.BlockHashTreeFromBlockContext(&parentBlockCtx)
	if err != nil {
		return nil, err
	}
	bs, err := state.BlockStateFromBlock(parentBlockCtx.BlockNumber.Uint64(), statedb, cumulativeGasUsedCopy, blockHashTree)
	if err != nil {
		return nil, err
	}
	states = append(states, &ExecutionState{
		VMHash:            bs.Hash(),
		CumulativeGasUsed: cumulativeGasUsedCopy,
		BlockGasUsed:      common.Big0,
		StateType:         state.BlockStateType,
		Block:             parent,
		TransactionIdx:    0,
		StepIdx:           0,
	})

	for num := startNum; num < endNum; num++ {
		// Preparation of block context
		block, err = backend.BlockByNumber(ctx, rpc.BlockNumber(num))
		if err != nil {
			return nil, err
		}
		if block == nil {
			return nil, fmt.Errorf("block #%d not found", num)
		}
		signer := types.MakeSigner(backend.ChainConfig(), block.Number())
		blockCtx := core.NewEVMBlockContext(block.Header(), chainCtx, nil)
		blockHashTree, err := state.BlockHashTreeFromBlockContext(&blockCtx)
		if err != nil {
			return nil, err
		}
		transactions := block.Transactions()
		receipts, _ := backend.GetReceipts(ctx, block.Hash())
		cumulativeGasUsedBeforeBlock := new(big.Int).Set(cumulativeGasUsed)

		// Trace all the transactions contained within
		for i, tx := range transactions {
			// Prepare the transaction context
			msg, _ := tx.AsMessage(signer, block.BaseFee())
			txContext := core.NewEVMTxContext(msg)
			// Call Prepare to clear out the statedb access list
			statedb.Prepare(tx.Hash(), i)

			// Push the inter state before transaction i
			blockGasUsed := new(big.Int).Sub(cumulativeGasUsed, cumulativeGasUsedBeforeBlock)
			its := state.InterStateFromCaptured(
				block.NumberU64(),
				uint64(i),
				statedb,
				cumulativeGasUsedBeforeBlock,
				blockGasUsed,
				transactions,
				receipts,
				blockHashTree,
			)
			states = append(states, &ExecutionState{
				VMHash:            its.Hash(),
				CumulativeGasUsed: cumulativeGasUsedBeforeBlock,
				BlockGasUsed:      blockGasUsed,
				StateType:         state.InterStateType,
				Block:             block,
				TransactionIdx:    uint64(i),
				StepIdx:           0,
			})

			// Execute transaction i with intra state generator enabled.
			prover := prover.NewIntraStateGenerator(block.NumberU64(), uint64(i), statedb, *its, blockHashTree)
			vmenv := vm.NewEVM(blockCtx, txContext, statedb, backend.ChainConfig(), vm.Config{Debug: true, Tracer: prover, NoBaseFee: true})
			executionResult, err := core.ApplyMessage(vmenv, msg, new(core.GasPool).AddGas(msg.Gas()))
			if err != nil {
				return nil, fmt.Errorf("tracing failed: %w", err)
			}
			generatedStates, err := prover.GetGeneratedStates()
			if err != nil {
				return nil, fmt.Errorf("tracing failed: %w", err)
			}
			for idx, s := range generatedStates {
				states = append(states, &ExecutionState{
					VMHash:            s.VMHash,
					CumulativeGasUsed: cumulativeGasUsedBeforeBlock,
					BlockGasUsed:      new(big.Int).Add(blockGasUsed, new(big.Int).SetUint64(tx.Gas()-s.Gas)),
					StateType:         state.IntraStateType,
					Block:             block,
					TransactionIdx:    uint64(i),
					StepIdx:           uint64(idx + 1),
				})
			}
		}

		// Push the inter state after all transactions
		blockGasUsed := new(big.Int).Sub(cumulativeGasUsed, cumulativeGasUsedBeforeBlock)
		its := state.InterStateFromCaptured(
			block.NumberU64(),
			uint64(len(transactions)),
			statedb,
			cumulativeGasUsedBeforeBlock,
			blockGasUsed,
			transactions,
			receipts,
			blockHashTree,
		)
		states = append(states, &ExecutionState{
			VMHash:            its.Hash(),
			CumulativeGasUsed: cumulativeGasUsedBeforeBlock,
			BlockGasUsed:      blockGasUsed,
			StateType:         state.InterStateType,
			Block:             block,
			TransactionIdx:    uint64(len(block.Transactions())),
			StepIdx:           0,
		})

		// Get next statedb to skip simulating block finalization
		// Here the statedb is the state at the end of the current block
		// (i.e. start of the new block)
		statedb, err = backend.StateAtBlock(ctx, block, reexec, statedb, true, false)
		if err != nil {
			return nil, err
		}

		// Push the block state of the current finalized block
		cumulativeGasUsedCopy = new(big.Int).Set(cumulativeGasUsed)
		bs, err := state.BlockStateFromBlock(block.NumberU64(), statedb, cumulativeGasUsedCopy, blockHashTree)
		if err != nil {
			return nil, err
		}
		states = append(states, &ExecutionState{
			VMHash:            bs.Hash(),
			CumulativeGasUsed: cumulativeGasUsedCopy,
			BlockGasUsed:      common.Big0,
			StateType:         state.BlockStateType,
			Block:             block,
			TransactionIdx:    0,
			StepIdx:           0,
		})
	}
	return states, nil
}

// [GenerateProof] serves as an entrypoint for one-step proof generation.
// There are 6 types of one-step proofs:
//  1. BlockState -> InterState: block initiation
//  2. InterState -> IntraState: transaction initiation (contract call or creation)
//  3. InterState -> InterState: EOA transfer transaction
//  4. IntraState -> IntraState: one-step EVM execution (require tracing)
//  5. IntraState -> InterState: transaction finalization (require tracing)
//  6. InterState -> BlockState: block finalization
func GenerateProof(backend Backend, ctx context.Context, startState *ExecutionState, config *ProverConfig) (*proof.OneStepProof, error) {
	if startState.Block == nil {
		return nil, fmt.Errorf("bad start state")
	}
	transactions := startState.Block.Transactions()
	if startState.TransactionIdx > uint64(len(transactions)) {
		return nil, fmt.Errorf("bad start state")
	}

	reexec := defaultProveReexec
	if config != nil && config.Reexec != nil {
		reexec = *config.Reexec
	}

	// Type 1: block initiation or Type 6: block finalization
	if startState.StateType == state.BlockStateType || (startState.StateType == state.InterStateType && startState.TransactionIdx == uint64(len(transactions))) {
		statedb, err := backend.StateAtBlock(ctx, startState.Block, reexec, nil, true, false)
		if err != nil {
			return nil, err
		}
		chainCtx := createChainContext(backend, ctx)
		vmctx := core.NewEVMBlockContext(startState.Block.Header(), chainCtx, nil)
		blockHashTree, err := state.BlockHashTreeFromBlockContext(&vmctx)
		if err != nil {
			return nil, err
		}
		if startState.StateType == state.BlockStateType {
			// Type 1: block initiation
			bs, err := state.BlockStateFromBlock(startState.Block.NumberU64(), statedb, startState.CumulativeGasUsed, blockHashTree)
			if err != nil {
				return nil, err
			}
			return proof.GetBlockInitiationProof(bs)
		} else {
			// Type 6: block finalization
			receipts, _ := backend.GetReceipts(ctx, startState.Block.Hash())
			its := state.InterStateFromCaptured(
				startState.Block.NumberU64(),
				startState.TransactionIdx,
				statedb,
				startState.CumulativeGasUsed,
				startState.BlockGasUsed,
				transactions,
				receipts,
				blockHashTree,
			)
			return proof.GetBlockFinalizationProof(its)
		}
	}

	// Prepare block and transaction context
	msg, vmctx, statedb, err := backend.StateAtTransaction(ctx, startState.Block, int(startState.TransactionIdx), reexec)
	if err != nil {
		return nil, err
	}
	txContext := core.NewEVMTxContext(msg)
	receipts, _ := backend.GetReceipts(ctx, startState.Block.Hash())
	blockHashTree, err := state.BlockHashTreeFromBlockContext(&vmctx)
	if err != nil {
		return nil, err
	}

	// Prepare the inter state before transaction for the prover
	its := state.InterStateFromCaptured(
		startState.Block.NumberU64(),
		startState.TransactionIdx,
		statedb,
		startState.CumulativeGasUsed,
		startState.BlockGasUsed,
		transactions,
		receipts,
		blockHashTree,
	)

	transaction := transactions[startState.TransactionIdx]

	if startState.StateType == state.InterStateType {
		// Type 2: transaction initiation or Type 3: EOA transfer transaction
		return proof.GetTransactionInitaitionProof(backend.ChainConfig(), &vmctx, transaction, &txContext, its, statedb)
	}
	// Type 4: one-step EVM execution or Type 5: transaction finalization. Both require tracing.

	// Set up the prover
	prover := prover.NewProver(
		startState.VMHash,
		startState.StepIdx,
		backend.ChainConfig().Rules(vmctx.BlockNumber, vmctx.Random != nil),
		startState.Block.NumberU64(),
		startState.TransactionIdx,
		transaction.Gas(),
		transaction.GasPrice(),
		statedb,
		*its,
		blockHashTree,
		receipts[startState.TransactionIdx],
	)
	// Run the transaction with prover enabled.
	vmenv := vm.NewEVM(vmctx, txContext, statedb, backend.ChainConfig(), vm.Config{Debug: true, Tracer: prover, NoBaseFee: true})
	// Call Prepare to clear out the statedb access list
	txHash := transactions[startState.TransactionIdx].Hash()
	statedb.Prepare(txHash, int(startState.TransactionIdx))
	_, err = core.ApplyMessage(vmenv, msg, new(core.GasPool).AddGas(msg.Gas()))
	if err != nil {
		return nil, fmt.Errorf("tracing failed: %w", err)
	}
	return prover.GetProof()
}
