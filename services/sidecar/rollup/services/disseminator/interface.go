package disseminator

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/beacon/engine"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/specularL2/specular/services/sidecar/rollup/rpc/eth"
	"github.com/specularL2/specular/services/sidecar/rollup/types"
)

type Config interface{ GetDisseminationInterval() time.Duration }

type ForkChoiceState = engine.ForkchoiceStateV1
type ForkChoiceResponse = engine.ForkChoiceResponse
type BuildPayloadResponse = engine.ForkChoiceResponse

type BatchBuilder interface {
	Enqueue(block *ethTypes.Block) error
	LastEnqueued() types.BlockID
	Build() ([]byte, error)
	Advance()
	Reset(lastEnqueued types.BlockID)
}

type TxManager interface {
	AppendTxBatch(ctx context.Context, batch []byte) (*ethTypes.Receipt, error)
}

type L2Client interface {
	EnsureDialed(ctx context.Context) error
	BlockNumber(ctx context.Context) (uint64, error)
	BlockByNumber(ctx context.Context, number *big.Int) (*ethTypes.Block, error)
	HeaderByTag(ctx context.Context, tag eth.BlockTag) (*ethTypes.Header, error)
}
