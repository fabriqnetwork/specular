package geth_prover

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/eth"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rpc"
	prover_types "github.com/specularl2/specular/clients/geth/specular/prover/types"
)

type GethBackend struct {
	Backend *eth.EthAPIBackend
}

func (g GethBackend) HeaderByHash(ctx context.Context, hash common.Hash) (*types.Header, error) {
	return g.Backend.HeaderByHash(ctx, hash)
}

func (g GethBackend) HeaderByNumber(ctx context.Context, number rpc.BlockNumber) (*types.Header, error) {
	return g.Backend.HeaderByNumber(ctx, number)
}

func (g GethBackend) BlockByHash(ctx context.Context, hash common.Hash) (prover_types.Block, error) {
	block, err := g.Backend.BlockByHash(ctx, hash)
	return &GethBlock{block}, err
}

func (g GethBackend) BlockByNumber(ctx context.Context, number rpc.BlockNumber) (prover_types.Block, error) {
	block, err := g.Backend.BlockByNumber(ctx, number)
	return &GethBlock{block}, err
}

func (g GethBackend) GetTransaction(ctx context.Context, txHash common.Hash) (prover_types.Transaction, common.Hash, uint64, uint64, error) {
	tx, bhash, bnum, idx, err := g.Backend.GetTransaction(ctx, txHash)
	return &GethTransaction{tx}, bhash, bnum, idx, err
}

func (g GethBackend) GetReceipts(ctx context.Context, hash common.Hash) (types.Receipts, error) {
	return g.Backend.GetReceipts(ctx, hash)
}

func (g GethBackend) RPCGasCap() uint64 {
	return g.Backend.RPCGasCap()
}

func (g GethBackend) ChainConfig() *params.ChainConfig {
	return g.Backend.ChainConfig()
}

func (g GethBackend) Engine() consensus.Engine {
	return g.Backend.Engine()
}

func (g GethBackend) ChainDb() ethdb.Database {
	return g.Backend.ChainDb()
}

func (g GethBackend) StateAtBlock(ctx context.Context, block prover_types.Block, reexec uint64, base prover_types.L2ELClientStateInterface, checkLive, preferDisk bool) (prover_types.L2ELClientStateInterface, error) {
	geth_state, cast_failed := base.(*GethState)
	if cast_failed {
		panic("base state is not a GethState")
	}
	geth_block, cast_failed := block.(*GethBlock)
	if cast_failed {
		panic("block is not a GethBlock")
	}
	s, err := g.Backend.StateAtBlock(ctx, geth_block.Block, reexec, geth_state.StateDB, checkLive, preferDisk)
	return GethState{StateDB: s}, err
}

func (g GethBackend) StateAtTransaction(ctx context.Context, block prover_types.Block, txIndex int, reexec uint64) (core.Message, prover_types.L2ELClientBlockContextInterface, prover_types.L2ELClientStateInterface, error) {
	geth_block, cast_failed := block.(*GethBlock)
	if cast_failed {
		panic("block is not a GethBlock")
	}
	msg, block_context, state, err := g.Backend.StateAtTransaction(ctx, geth_block.Block, txIndex, reexec)
	return msg, &GethBlockContext{block_context}, GethState{StateDB: state}, err
}

func (g GethBackend) NewEVM(blockCtx prover_types.L2ELClientBlockContextInterface, txCtx vm.TxContext, statedb prover_types.L2ELClientStateInterface, chainConfig *params.ChainConfig, config prover_types.L2ELClientConfig) prover_types.L2ELClientEVMInterface {
	return &GethEVM{vm.NewEVM(blockCtx.(*GethBlockContext).Context, txCtx, statedb.(GethState).StateDB, chainConfig, vm.Config{Debug: config.Debug, Tracer: config.Tracer.(vm.EVMLogger)})}
}

func (g GethBackend) NewEVMBlockContext(header *types.Header, chain core.ChainContext, author *common.Address) prover_types.L2ELClientBlockContextInterface {
	return GethBlockContext{core.NewEVMBlockContext(header, chain, author)}
}

func (g GethBackend) ApplyMessage(evm prover_types.L2ELClientEVMInterface, msg core.Message, gp *core.GasPool) (*core.ExecutionResult, error) {
	return core.ApplyMessage(evm.(*GethEVM).EVM, msg, gp)
}
