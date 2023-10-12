package api

import (
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/txpool"
)

// Required interface for interacting with Ethereum instance.
type ExecutionBackend interface {
	BlockChain() *core.BlockChain
	TxPool() *txpool.TxPool
}
