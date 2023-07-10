package sequencer

import (
	"context"
	"errors"
	"time"

	"github.com/specularl2/specular/clients/geth/specular/rollup/services/derivation/engine"
	"github.com/specularl2/specular/clients/geth/specular/rollup/services/derivation/stage"
	"github.com/specularl2/specular/clients/geth/specular/utils"
)

type PlanningSequencer struct {
	simpleSequencer

	blockInterval time.Duration
	targetTs      time.Time
}

type simpleSequencer interface {
	Step(ctx context.Context) *SequencerError
	EngineManager() EngineManager
}

func NewPlanningSequencer(sequencer simpleSequencer, blockInterval time.Duration) *PlanningSequencer {
	return &PlanningSequencer{
		simpleSequencer: sequencer,
		blockInterval:   blockInterval,
		targetTs:        time.Now(),
	}
}

func (s *PlanningSequencer) Step(ctx context.Context) error {
	err := s.simpleSequencer.Step(ctx)
	if err == nil {
		return nil
	}
	// Update target time depending on the error.
	// TODO: payload building cancellation upon error.
	switch err.Category() {
	case busy:
		// approximates the worst-case time it takes to build a block, to reattempt sequencing after.
		s.updateTarget(time.Duration(s.blockInterval) * time.Second)
	case starting:
		if errors.Is(err, stage.RecoverableError) {
			s.updateTarget(s.blockInterval) // hold off from sequencing for a full block
			return err
		} else if errors.Is(err, stage.RetryableError) {
			s.updateTarget(time.Second)
		} else {
			s.updateTarget(time.Second)
			return err // unclassified errors are currently considered fatal.
		}
	case confirming:
		if errors.Is(err, stage.RecoverableError) {
			s.updateTarget(s.blockInterval) // hold off from sequencing for a full block
			return err
		} else if errors.Is(err, stage.RetryableError) {
			s.updateTarget(time.Second)
		} else {
			s.updateTarget(time.Second)
			return err // unclassified errors are currently considered fatal.
		}
	}
	return nil
}

// Returns the desired delay till the next step call.
func (p *PlanningSequencer) Plan() time.Duration {
	var isBuilding = p.EngineManager().IsBuilding()
	// If the engine is busy building safe blocks (and thus changing the head that we would sync on top of),
	// then give it time to sync up.
	if isBuilding && p.EngineManager().CurrentBuildJob().Type() == engine.Safe {
		// approximates the worst-case time it takes to build a block, to reattempt sequencing after.
		return time.Second * time.Duration(p.blockInterval)
	}
	var (
		isBuildingConsistently = p.EngineManager().IsBuildingConsistently()
		now                    = time.Now()
	)
	// We may have to wait till the next sequencing action, e.g. upon an error.
	// If the head changed we need to respond and will not delay the sequencing.
	if delay := p.targetTs.Sub(now); delay > 0 && isBuildingConsistently {
		return delay
	}
	var (
		head          = p.EngineManager().L2State().Head()
		blockTime     = time.Duration(p.blockInterval) * time.Second
		payloadTime   = time.Unix(int64(head.GetTime()+uint64(p.blockInterval.Seconds())), 0)
		remainingTime = payloadTime.Sub(now)
		margin        time.Duration
	)
	// If we started building a block already, and if that work is still consistent,
	// then we would like to finish it by sealing the block.
	if isBuildingConsistently {
		// If we started building already, then we will schedule the sealing.
		margin = sealingDurationEstimate
	} else {
		// If we did not yet start building, then we will schedule the start.
		// if we have too much time, then wait before starting the build, otherwise start instantly.
		margin = blockTime
	}
	return utils.Max(remainingTime-margin, 0)
}

func (s *PlanningSequencer) updateTarget(delay time.Duration) { s.targetTs = time.Now().Add(delay) }
