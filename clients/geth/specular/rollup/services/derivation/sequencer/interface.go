package sequencer

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/specularl2/specular/clients/geth/specular/rollup/services/derivation/engine"
	"github.com/specularl2/specular/clients/geth/specular/rollup/types"
)

type Config interface {
	IsEnabled() bool
	GetAccountAddr() common.Address
}

type EngineManager interface {
	StartPayload(
		ctx context.Context,
		parent types.L2BlockRef,
		attrs *engine.PayloadAttributes,
		jobType engine.JobType,
	) error
	ConfirmPayload(ctx context.Context) (*engine.ExecutionPayload, error)
	CurrentBuildJob() *engine.EngineBuildJob
	IsBuilding() bool
	IsBuildingConsistently() bool
	L2State() *engine.EthState
}
