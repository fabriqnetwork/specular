package services

import (
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
	"github.com/specularl2/specular/clients/geth/specular/rollup/types/data"
)

// TODO: generalize
type ExecutionBackend interface {
	SubscribeNewTxsEvent(chan<- core.NewTxsEvent) event.Subscription
	Prepare(txs []*types.Transaction) *types.TransactionsByPriceAndNonce
	// TODO: dedup
	CommitTransactions(txs []*types.Transaction) error
	CommitBlock(block *data.DerivationBlock) error
}
