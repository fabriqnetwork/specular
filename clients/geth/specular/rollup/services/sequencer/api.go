package sequencer

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/beacon"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/specularl2/specular/clients/geth/specular/rollup/rpc/eth"
	"github.com/specularl2/specular/clients/geth/specular/rollup/services/api"
	"github.com/specularl2/specular/clients/geth/specular/rollup/types"
	"github.com/specularl2/specular/clients/geth/specular/rollup/types/da"
)

type Config interface {
	GetAccountAddr() common.Address
	GetMinExecutionInterval() time.Duration
	GetMaxExecutionInterval() time.Duration
	GetSequencingInterval() time.Duration
}

type ExecutionBackend interface {
	ForkchoiceUpdate(update *ForkChoiceState) (*ForkChoiceResponse, error)
	BuildPayload(attrs api.BuildPayloadAttributes) error
}

type ForkChoiceState = beacon.ForkchoiceStateV1
type ForkChoiceResponse = beacon.ForkChoiceResponse

type BatchBuilder interface {
	Append(block da.DerivationBlock, header da.HeaderRef) error
	LastAppended() types.BlockID
	Build() (*da.BatchAttributes, error)
	Advance()
	Reset(lastAppended types.BlockID)
}

type headerRef interface {
	GetHash() common.Hash
	GetParentHash() common.Hash
}

type batchAttributes interface {
	FirstL2BlockNumber() *big.Int
}

type TxManager interface {
	AppendTxBatch(
		ctx context.Context,
		contexts []*big.Int,
		txLengths []*big.Int,
		firstL2BlockNumber *big.Int,
		txBatch []byte,
	) (*ethTypes.Receipt, error)
}

type L2Client interface {
	EnsureDialed(ctx context.Context) error
	BlockNumber(ctx context.Context) (uint64, error)
	BlockByNumber(ctx context.Context, number *big.Int) (*ethTypes.Block, error)
	HeaderByTag(ctx context.Context, tag eth.BlockTag) (*ethTypes.Header, error)
	TransactionByHash(ctx context.Context, hash common.Hash) (*ethTypes.Transaction, bool, error)
}
