package services

import (
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/beacon"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// TODO: generalize
type ExecutionBackend interface {
	SubscribeNewTxsEvent(chan<- core.NewTxsEvent) event.Subscription
	ForkchoiceUpdate(update *ForkChoiceState) (*ForkChoiceResponse, error)
	BuildPayload(payload ExecutionPayload) error
	CommitTransactions(txs []*types.Transaction) error // TODO: remove
	Prepare(txs []*types.Transaction) TransactionQueue // TODO: remove
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
