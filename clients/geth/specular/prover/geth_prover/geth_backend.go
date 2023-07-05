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
	oss "github.com/specularl2/specular/clients/geth/specular/prover/state"
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

func (g GethBackend) BlockByHash(ctx context.Context, hash common.Hash) (*types.Block, error) {
	return g.Backend.BlockByHash(ctx, hash)
}

func (g GethBackend) BlockByNumber(ctx context.Context, number rpc.BlockNumber) (*types.Block, error) {
	return g.Backend.BlockByNumber(ctx, number)
}

func (g GethBackend) GetTransaction(ctx context.Context, txHash common.Hash) (*types.Transaction, common.Hash, uint64, uint64, error) {
	return g.Backend.GetTransaction(ctx, txHash)
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

func (g GethBackend) StateAtBlock(ctx context.Context, block *types.Block, reexec uint64, base oss.L2ELClientStateInterface, checkLive, preferDisk bool) (oss.L2ELClientStateInterface, error) {
	geth_state, e := base.(*GethState)
	if e {
		panic("base state is not a GethState")
	}
	s, err := g.Backend.StateAtBlock(ctx, block, reexec, geth_state.StateDB, checkLive, preferDisk)
	return GethState{StateDB: s}, err
}

func (g GethBackend) StateAtTransaction(ctx context.Context, block *types.Block, txIndex int, reexec uint64) (core.Message, oss.L2ELClientBlockContextInterface, oss.L2ELClientStateInterface, error) {
	msg, block_context, state, err := g.Backend.StateAtTransaction(ctx, block, txIndex, reexec)
	return msg, &GethBlockContext{block_context}, GethState{StateDB: state}, err
}

func (g GethBackend) NewEVM(blockCtx oss.L2ELClientBlockContextInterface, txCtx vm.TxContext, statedb oss.L2ELClientStateInterface, chainConfig *params.ChainConfig, config oss.L2ELClientConfig) oss.L2ELClientEVMInterface {
	return &GethEVM{vm.NewEVM(blockCtx.(*GethBlockContext).Context, txCtx, statedb.(GethState).StateDB, chainConfig, vm.Config{Debug: config.Debug, Tracer: config.Tracer.(vm.EVMLogger)})}
}

func (g GethBackend) NewEVMBlockContext(header *types.Header, chain core.ChainContext, author *common.Address) oss.L2ELClientBlockContextInterface {
	return GethBlockContext{core.NewEVMBlockContext(header, chain, author)}
}

func (g GethBackend) ApplyMessage(evm oss.L2ELClientEVMInterface, msg core.Message, gp *core.GasPool) (*core.ExecutionResult, error) {
	return core.ApplyMessage(evm.(*GethEVM).EVM, msg, gp)
}
