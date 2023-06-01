package api

import (
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/beacon"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Defines interface between Specular services and the underlying execution backend.
// TODO: generalize to better support clients other than Geth.
type ExecutionBackend interface {
	SubscribeNewTxsEvent(chan<- core.NewTxsEvent) event.Subscription
	ForkchoiceUpdate(update *ForkChoiceState) (*ForkChoiceResponse, error)
	BuildPayload(payload ExecutionPayload) error
	CommitTransactions(txs []*types.Transaction) error // TODO: remove
	Order(txs []*types.Transaction) TransactionQueue   // TODO: remove
}

type ForkChoiceState = beacon.ForkchoiceStateV1
type ForkChoiceResponse = beacon.ForkChoiceResponse

type ExecutionPayload interface {
	BlockNumber() uint64
	Timestamp() uint64
	Txs() [][]byte
}

type TransactionQueue interface {
	Peek() *types.Transaction
	Pop()
}
