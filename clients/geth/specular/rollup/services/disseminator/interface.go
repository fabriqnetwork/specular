package disseminator

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/core/beacon"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/specularl2/specular/clients/geth/specular/rollup/rpc/eth"
	"github.com/specularl2/specular/clients/geth/specular/rollup/types"
	"github.com/specularl2/specular/clients/geth/specular/rollup/types/da"
)

type Config interface{ GetDisseminationInterval() time.Duration }

type ForkChoiceState = beacon.ForkchoiceStateV1
type ForkChoiceResponse = beacon.ForkChoiceResponse
type BuildPayloadResponse = beacon.ForkChoiceResponse

type BatchBuilder interface {
	Append(block da.DerivationBlock, header da.HeaderRef) error
	LastAppended() types.BlockID
	Build() (*da.BatchAttributes, error)
	Advance()
	Reset(lastAppended types.BlockID)
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
}
