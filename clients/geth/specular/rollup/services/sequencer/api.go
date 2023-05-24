package sequencer

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/beacon"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
	"github.com/specularl2/specular/clients/geth/specular/rollup/l2types"
	"github.com/specularl2/specular/clients/geth/specular/rollup/rpc/client"
	"github.com/specularl2/specular/clients/geth/specular/rollup/services"
)

type SequencerServiceConfig interface {
	Sequencer() *services.SequencerConfig
	L1() *services.L1Config
}

type ExecutionBackend interface {
	SubscribeNewTxsEvent(chan<- core.NewTxsEvent) event.Subscription
	ForkchoiceUpdate(update *ForkChoiceState) (*ForkChoiceResponse, error)
	BuildPayload(payload services.ExecutionPayload) error
	CommitTransactions(txs []*types.Transaction) error                     // TODO: remove
	Prepare(txs []*types.Transaction) services.TransactionsByPriceAndNonce // TODO: remove
}

type ForkChoiceState = beacon.ForkchoiceStateV1
type ForkChoiceResponse = beacon.ForkChoiceResponse

type BatchBuilder interface {
	Append(block l2types.DerivationBlock, header Header) error
	LastAppended() l2types.BlockID
	Build() (*batchAttributes, error)
	Advance()
	Reset(lastAppended l2types.BlockID)
}

type TxManager interface {
	AppendTxBatch(
		ctx context.Context,
		contexts []*big.Int,
		txLengths []*big.Int,
		txs []byte,
	) (*types.Receipt, error)
}

type L2Client interface {
	BlockNumber(ctx context.Context) (uint64, error)
	BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error)
	HeaderByTag(ctx context.Context, tag client.BlockTag) (*types.Header, error)
	TransactionByHash(ctx context.Context, hash common.Hash) (*types.Transaction, bool, error)
	Close()
}
