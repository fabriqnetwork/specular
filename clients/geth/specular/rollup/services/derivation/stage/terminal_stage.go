package stage

import (
	"context"

	"github.com/specularl2/specular/clients/geth/specular/rollup/types"
)

// Last stage in a pipeline. No output.
// Note: the terminal stage is responsible for finding the recovery point.
type TerminalStage[T any] struct{ *Stage[T, any] }

type terminalStageProcessor[T any] interface {
	stageProcessor[T, any]
	findRecoveryPoint(ctx context.Context) (types.BlockID, error)
}

func NewTerminalStage[T any](prev StageOps[T], processor terminalStageProcessor[T]) *TerminalStage[T] {
	return &TerminalStage[T]{&Stage[T, any]{prev: prev, processor: processor}}
}

func (s *TerminalStage[T]) FindRecoveryPoint(ctx context.Context) (types.BlockID, error) {
	return s.processor.(terminalStageProcessor[T]).findRecoveryPoint(ctx)
}
