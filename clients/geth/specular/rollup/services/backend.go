package services

import (
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// TODO: generalize
type ExecutionBackend interface {
	SubscribeNewTxsEvent(chan<- core.NewTxsEvent) event.Subscription
	Prepare(txs []*types.Transaction) TransactionsByPriceAndNonce
	// TODO: dedup
	CommitTransactions(txs []*types.Transaction) error
	CommitPayload(payload ExecutionPayload) error
}

type ExecutionPayload interface {
	BlockNumber() uint64
	Timestamp() uint64
	Txs() [][]byte
}

type TransactionsByPriceAndNonce interface {
	Peek() *types.Transaction
	Pop()
}
