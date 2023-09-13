package api

import (
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
)

// Required interface for interacting with Ethereum instance.
type ExecutionBackend interface {
	BlockChain() *core.BlockChain
	TxPool() *core.TxPool
	StateAtBlock(block *types.Block, reexec uint64, base *state.StateDB, checkLive bool, preferDisk bool) (statedb *state.StateDB, err error)
}
