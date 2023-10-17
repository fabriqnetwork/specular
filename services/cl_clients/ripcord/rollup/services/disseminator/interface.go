package disseminator

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/beacon/engine"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/specularl2/specular/services/cl_clients/ripcord/rollup/derivation"
	"github.com/specularl2/specular/services/cl_clients/ripcord/rollup/rpc/eth"
	"github.com/specularl2/specular/services/cl_clients/ripcord/rollup/types"
)

type Config interface{ GetDisseminationInterval() time.Duration }

type ForkChoiceState = engine.ForkchoiceStateV1
type ForkChoiceResponse = engine.ForkChoiceResponse
type BuildPayloadResponse = engine.ForkChoiceResponse

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
		txBatchVersion *big.Int,
		txBatch []byte,
	) (*ethTypes.Receipt, error)
}

type L2Client interface {
	EnsureDialed(ctx context.Context) error
	BlockNumber(ctx context.Context) (uint64, error)
	BlockByNumber(ctx context.Context, number *big.Int) (*ethTypes.Block, error)
	HeaderByTag(ctx context.Context, tag eth.BlockTag) (*ethTypes.Header, error)
}
