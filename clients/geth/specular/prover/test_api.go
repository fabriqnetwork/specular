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
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/rpc"
	oss "github.com/specularl2/specular/clients/geth/specular/prover/state"
	prover_types "github.com/specularl2/specular/clients/geth/specular/prover/types"
)

func (api *ProverAPI) GenerateProofForTest(ctx context.Context, hash common.Hash, cumulativeGasUsed, blockGasUsed *big.Int, step uint64, config *ProverConfig) (json.RawMessage, error) {
	transaction, blockHash, blockNumber, index, err := api.backend.GetTransaction(ctx, hash)
	if err != nil {
		return nil, err
	}
	// It shouldn't happen in practice.
	if blockNumber == 0 {
		return nil, errors.New("genesis is not traceable")
	}
	reexec := defaultProveReexec
	if config != nil && config.Reexec != nil {
		reexec = *config.Reexec
	}
	block, err := api.backend.BlockByNumber(ctx, rpc.BlockNumber(blockNumber))
	if err != nil {
		return nil, err
	}
	if block == nil {
		return nil, fmt.Errorf("block #%d not found", blockNumber)
	}
	msg, vmctx, statedb, err := api.backend.StateAtTransaction(ctx, block, int(index), reexec)
	if err != nil {
		return nil, err
	}
	txContext := core.NewEVMTxContext(msg)
	receipts, err := api.backend.GetReceipts(ctx, blockHash)
	if err != nil {
		return nil, err
	}
	blockHashTree, err := oss.BlockHashTreeFromBlockContext(vmctx)
	if err != nil {
		return nil, err
	}
	its := oss.InterStateFromCaptured(
		blockNumber,
		index,
		statedb,
		cumulativeGasUsed,
		blockGasUsed,
		block.Transactions(),
		receipts,
		blockHashTree,
	)
	prover := NewTestProver(
		step,
		transaction,
		&txContext,
		receipts[index],
		api.backend.ChainConfig().Rules(vmctx.BlockNumber(), vmctx.Random() != nil),
		blockNumber,
		index,
		statedb,
		*its,
		blockHashTree,
	)
	vmenv := api.backend.NewEVM(vmctx, txContext, statedb, api.backend.ChainConfig(), prover_types.L2ELClientConfig{Debug: true, Tracer: prover})
	// TODO: under geth situation we don't need to pass in the block hash
	statedb.Prepare(hash, common.Hash{}, int(index))
	_, err = api.backend.ApplyMessage(vmenv, msg, new(core.GasPool).AddGas(msg.Gas()))
	if err != nil {
		return nil, fmt.Errorf("tracing failed: %w", err)
	}
	return prover.GetResult()
}
