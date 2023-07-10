package stage

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/beacon"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/specularl2/specular/clients/geth/specular/rollup/rpc/eth"
	"github.com/specularl2/specular/clients/geth/specular/rollup/services/api"
	"github.com/specularl2/specular/clients/geth/specular/rollup/services/derivation/engine"
	"github.com/specularl2/specular/clients/geth/specular/rollup/types"
	"github.com/specularl2/specular/clients/geth/specular/utils"
)

// Represents a stage in a pipeline.
// Generic parameters:
// `T`: Stage output type.
type StageOps[T any] interface {
	// Possible errors returned:
	// - RetryableError: Indicates caller should retry step.
	// - RecoverableError: Indicates caller should perform recovery.
	// - Unrecoverable fatal error (i.e. any other type): Unexpected. Indicates caller should not retry.
	Pull(ctx context.Context) (T, error)
	// Recovers from a re-org to the given L1 block number.
	Recover(ctx context.Context, l1BlockID types.BlockID) error
}

type EngineManager interface {
	ForkchoiceUpdate(ctx context.Context, update *ForkChoiceState) (*ForkChoiceResponse, error)
	StartPayload(
		ctx context.Context,
		parent types.L2BlockRef,
		attrs *engine.PayloadAttributes,
		jobType engine.JobType,
	) error
	ConfirmPayload(ctx context.Context) (*engine.ExecutionPayload, error)
	CurrentBuildJob() *engine.EngineBuildJob
}

type ExecutionBackend interface {
	ForkchoiceUpdate(ctx context.Context, update *ForkChoiceState) (*ForkChoiceResponse, error)
	BuildPayload(ctx context.Context, payloadAttrs api.BuildPayloadAttributes) (*BuildPayloadResponse, error)
}

type ForkChoiceState = beacon.ForkchoiceStateV1
type ForkChoiceResponse = beacon.ForkChoiceResponse
type BuildPayloadResponse = beacon.ForkChoiceResponse

type L1Client interface {
	HeaderByNumber(ctx context.Context, number *big.Int) (*ethTypes.Header, error)
	BlockByHash(ctx context.Context, hash common.Hash) (*ethTypes.Block, error)
}

type L2Client interface {
	EnsureDialed(ctx context.Context) error
	BlockNumber(ctx context.Context) (uint64, error)
	HeaderByNumber(ctx context.Context, number *big.Int) (*ethTypes.Header, error)
	HeaderByTag(ctx context.Context, tag eth.BlockTag) (*ethTypes.Header, error)
}

type L1State interface {
	Head() types.BlockID
	Safe() types.BlockID
	Finalized() types.BlockID
}

type RollupState interface {
	OnAssertionCreated(ctx context.Context, l1BlockID types.BlockID, tx *ethTypes.Transaction) error
	OnAssertionConfirmed(ctx context.Context, l1BlockID types.BlockID, tx *ethTypes.Transaction) error
	OnAssertionRejected(ctx context.Context, l1BlockID types.BlockID, tx *ethTypes.Transaction) error
}

type DerivationConfig interface {
	L1Config
	GenesisConfig
}

type L1Config interface {
	GetChainID() uint64
	GetSequencerInboxAddr() common.Address
	GetRollupAddr() common.Address
}

type GenesisConfig interface {
	GetGenesisL1BlockID() types.BlockID
}

type ErrorType uint

const (
	Retryable ErrorType = iota
	Recoverable
	// All other errors are considered fatal.
)

func NewRetryableError(err error) *utils.CategorizedError[ErrorType] {
	return &utils.CategorizedError[ErrorType]{Cat: Retryable, Err: err}
}

func NewRecoverableError(err error) *utils.CategorizedError[ErrorType] {
	return &utils.CategorizedError[ErrorType]{Cat: Recoverable, Err: err}
}

// Aliases for comparison (errors.Is).
var (
	RetryableError   = NewRetryableError(nil)
	RecoverableError = NewRecoverableError(nil)
)
