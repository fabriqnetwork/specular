package stage

import (
	"context"

	"github.com/specularl2/specular/clients/geth/specular/rollup/types"
)

// Note: each Stage is itself a `StageOps`.
// This struct ensures that `processor.process` is only called on the same cached input
// multiple times if the stage returned a retryable error.
// Generic parameters:
// `T`: Stage input type.
// `U`: Stage output type.
type Stage[T, U any] struct {
	prev      StageOps[T]          // note: can't be of type Stage due to golang type system
	processor stageProcessor[T, U] // processor for this stage.
	cached    T                    // cached output from previous stage
	isCached  bool                 // necessary for non-comparable types
}

// `T`: Stage processor input type (via `ingest`).
// `U`: Stage processor output type (via `next`).
type stageProcessor[T, U any] interface {
	hasNext() bool
	next() U
	ingest(ctx context.Context, prev T) error
	recover(ctx context.Context, l1BlockID types.BlockID) error
}

func NewStage[T, U any](prev StageOps[T], processor stageProcessor[T, U]) *Stage[T, U] {
	return &Stage[T, U]{prev: prev, processor: processor}
}

func (s *Stage[T, U]) Prev() StageOps[T] { return s.prev }

func (s *Stage[T, U]) Pull(ctx context.Context) (out U, err error) {
	// If there's nothing queued up, pull from the previous stage.
	if s.processor.hasNext() {
		return s.processor.next(), nil
	}
	err = s.ensurePulled(ctx)
	if err != nil {
		return out, nil
	}
	err = s.processor.ingest(ctx, s.cached)
	if err != nil {
		return out, err
	}
	// Successfully processed cached input,
	// so we should clear it to fetch the next one.
	s.clearPulled()
	// Note: `ingest` does not guarantee that `next` will return a non-zero/nil value.
	if s.processor.hasNext() {
		return s.processor.next(), nil
	}
	return out, nil
}

func (s *Stage[T, U]) Recover(ctx context.Context, l1BlockID types.BlockID) error {
	err := s.processor.recover(ctx, l1BlockID)
	if err != nil {
		return err
	}
	s.clearPulled()
	return s.prev.Recover(ctx, l1BlockID)
}

func (s *Stage[T, U]) ensurePulled(ctx context.Context) error {
	if !s.isCached {
		out, err := s.Prev().Pull(ctx)
		if err != nil {
			return err
		}
		s.cached = out
		s.isCached = true
	}
	return nil
}

func (s *Stage[T, U]) clearPulled() {
	var empty T
	s.cached = empty
	s.isCached = false
}
