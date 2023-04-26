package sequencer

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
	"github.com/specularl2/specular/clients/geth/specular/rollup/comms/client"
	"github.com/specularl2/specular/clients/geth/specular/rollup/comms/txmgr"
	"github.com/specularl2/specular/clients/geth/specular/rollup/services"
	"github.com/specularl2/specular/clients/geth/specular/rollup/types/data"
)

type SequencerServiceConfig interface {
	Sequencer() *services.SequencerConfig
	L1() *services.L1Config
}

type ExecutionBackend interface {
	SubscribeNewTxsEvent(chan<- core.NewTxsEvent) event.Subscription
	ForkchoiceUpdate(update services.ForkchoiceState) error
	BuildPayload(payload services.ExecutionPayload) error
	CommitTransactions(txs []*types.Transaction) error                     // TODO: remove
	Prepare(txs []*types.Transaction) services.TransactionsByPriceAndNonce // TODO: probably remove
}

type BatchBuilder interface {
	Append(block data.DerivationBlock, header data.Header) error
	LastAppended() uint64
	Build() ([]byte, error)
	Advance()
	Reset(lastAppendedBlockNumber uint64, lastAppendedBlockHash common.Hash)
}

type TxManager interface {
	Send(ctx context.Context, candidate txmgr.TxCandidate) (*types.Receipt, error)
}

type L2Client interface {
	BlockNumber(ctx context.Context) (uint64, error)
	BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error)
	HeaderByTag(ctx context.Context, tag client.BlockTag) (*types.Header, error)
	TransactionByHash(ctx context.Context, hash common.Hash) (*types.Transaction, bool, error)
	Close()
}
