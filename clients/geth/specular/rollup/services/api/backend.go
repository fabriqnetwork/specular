package api

import (
	"context"

	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/txpool"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth/tracers"
)

// Required interface for interacting with Ethereum instance.
type ExecutionBackend interface {
	BlockChain() *core.BlockChain
	TxPool() *txpool.TxPool
	StateAtBlock(ctx context.Context, block *types.Block, reexec uint64, base *state.StateDB, checkLive bool, preferDisk bool) (statedb *state.StateDB, release tracers.StateReleaseFunc, err error)
}
