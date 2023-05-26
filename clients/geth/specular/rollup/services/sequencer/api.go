package sequencer

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/beacon"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
	"github.com/specularl2/specular/clients/geth/specular/rollup/rpc/eth"
	"github.com/specularl2/specular/clients/geth/specular/rollup/services"
	"github.com/specularl2/specular/clients/geth/specular/rollup/types"
)

type SequencerServiceConfig interface {
	Sequencer() *services.SequencerConfig
	L1() *services.L1Config
}

type ExecutionBackend interface {
	SubscribeNewTxsEvent(chan<- core.NewTxsEvent) event.Subscription
	ForkchoiceUpdate(update *ForkChoiceState) (*ForkChoiceResponse, error)
	BuildPayload(payload services.ExecutionPayload) error
	CommitTransactions(txs []*ethTypes.Transaction) error          // TODO: remove
	Prepare(txs []*ethTypes.Transaction) services.TransactionQueue // TODO: remove
}

type ForkChoiceState = beacon.ForkchoiceStateV1
type ForkChoiceResponse = beacon.ForkChoiceResponse

type BatchBuilder interface {
	Append(block types.DerivationBlock, header Header) error
	LastAppended() types.BlockID
	Build() (*batchAttributes, error)
	Advance()
	Reset(lastAppended types.BlockID)
}

type TxManager interface {
	AppendTxBatch(
		ctx context.Context,
		contexts []*big.Int,
		txLengths []*big.Int,
		firstL2BlockNumber *big.Int,
		txs []byte,
	) (*ethTypes.Receipt, error)
}

type L2Client interface {
	BlockNumber(ctx context.Context) (uint64, error)
	BlockByNumber(ctx context.Context, number *big.Int) (*ethTypes.Block, error)
	HeaderByTag(ctx context.Context, tag eth.BlockTag) (*ethTypes.Header, error)
	TransactionByHash(ctx context.Context, hash common.Hash) (*ethTypes.Transaction, bool, error)
	Close()
}
