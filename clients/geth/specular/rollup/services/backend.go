package services

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// TODO: generalize
type ExecutionBackend interface {
	SubscribeNewTxsEvent(chan<- core.NewTxsEvent) event.Subscription
	ForkchoiceUpdate(update ForkchoiceState) error
	BuildPayload(payload ExecutionPayload) error
	CommitTransactions(txs []*types.Transaction) error            // TODO: remove
	Prepare(txs []*types.Transaction) TransactionsByPriceAndNonce // TODO: probably remove
}

type ForkchoiceState interface {
	HeadBlockHash() common.Hash
	SafeBlockHash() common.Hash
	FinalizedBlockHash() common.Hash
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
