package driver

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum/core/beacon"
	"github.com/specularl2/specular/clients/geth/specular/rollup/services/api"
)

type Config interface {
	GetStepInterval() time.Duration
	GetRetryDelay() time.Duration
	GetNumAttempts() uint
}

type Sequencer interface {
	Step(ctx context.Context) error
	Plan() time.Duration
}

type ExecutionBackend interface {
	ForkchoiceUpdate(ctx context.Context, update *ForkChoiceState) (*ForkChoiceResponse, error)
	BuildPayload(ctx context.Context, attrs api.BuildPayloadAttributes) (*BuildPayloadResponse, error)
}

type ForkChoiceState = beacon.ForkchoiceStateV1
type ForkChoiceResponse = beacon.ForkChoiceResponse
type BuildPayloadResponse = beacon.ForkChoiceResponse
