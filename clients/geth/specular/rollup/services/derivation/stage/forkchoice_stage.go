package stage

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/specularl2/specular/clients/geth/specular/rollup/l2types"
)

type L2ForkchoiceUpdateStage struct {
	prev        Stage[l2types.BlockRelation]
	execBackend ExecutionBackend
	l1State     L1State

	blockRelations l2types.BlockRelations
}

func NewL2ForkChoiceUpdateStage(
	prev Stage[l2types.BlockRelation],
	execBackend ExecutionBackend,
	l1State L1State,
) *L2ForkchoiceUpdateStage {
	return &L2ForkchoiceUpdateStage{
		prev:        prev,
		execBackend: execBackend,
		l1State:     l1State,
	}
}

func (s *L2ForkchoiceUpdateStage) Step(ctx context.Context) (interface{}, error) {
	relation, err := s.prev.Step(ctx)
	if err != nil {
		return nil, err
	}
	if relation != (l2types.BlockRelation{}) {
		err := s.blockRelations.Append(relation)
		if err != nil {
			return nil, fmt.Errorf("Failed to append block relation: %w", err)
		}
	}
	// Get latest L1 fork-choice.
	updatedL1ForkChoice := forkChoiceState{s.l1State.Head(), s.l1State.Safe(), s.l1State.Finalized()}
	// TODO: handle no forkChoice change
	safeL2BlockID := s.blockRelations.MarkSafe(uint64(updatedL1ForkChoice.safeID.Number()))
	finalizedL2BlockID := s.blockRelations.MarkFinal(uint64(updatedL1ForkChoice.finalizedID.Number()))
	l2Forkchoice := forkChoiceState{
		headID:      l2types.NewBlockID(0, common.Hash{}),
		safeID:      safeL2BlockID,
		finalizedID: finalizedL2BlockID,
	}
	err = s.execBackend.ForkchoiceUpdate(&l2Forkchoice)
	if err != nil {
		return nil, fmt.Errorf("Failed to update forkchoice state: %w", err)
	}
	return nil, nil
}

func (s *L2ForkchoiceUpdateStage) Recover(ctx context.Context, l1BlockID l2types.BlockID) error {
	s.blockRelations.MarkReorgedOut(l1BlockID.Number())
	return nil
}

type forkChoiceState struct {
	headID      l2types.BlockID
	safeID      l2types.BlockID
	finalizedID l2types.BlockID
}

func (fcs *forkChoiceState) HeadBlockHash() common.Hash      { return fcs.headID.Hash() }
func (fcs *forkChoiceState) SafeBlockHash() common.Hash      { return fcs.safeID.Hash() }
func (fcs *forkChoiceState) FinalizedBlockHash() common.Hash { return fcs.finalizedID.Hash() }
