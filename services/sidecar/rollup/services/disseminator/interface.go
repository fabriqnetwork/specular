package disseminator

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/core/beacon"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/specularL2/specular/services/sidecar/rollup/derivation"
	"github.com/specularL2/specular/services/sidecar/rollup/rpc/eth"
	"github.com/specularL2/specular/services/sidecar/rollup/types"
)

type Config interface{ GetDisseminationInterval() time.Duration }

type ForkChoiceState = beacon.ForkchoiceStateV1
type ForkChoiceResponse = beacon.ForkChoiceResponse
type BuildPayloadResponse = beacon.ForkChoiceResponse

type BatchBuilder interface {
	Append(block derivation.DerivationBlock, header derivation.HeaderRef) error
	LastAppended() types.BlockID
	Build() (*derivation.BatchAttributes, error)
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
