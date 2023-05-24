package stage

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/specularl2/specular/clients/geth/specular/rollup/l2types"
	"github.com/specularl2/specular/clients/geth/specular/rollup/utils/log"
)

type L2ForkchoiceUpdater struct {
	execBackend ExecutionBackend
	l1State     L1State

	l1ForkChoice   forkChoiceState
	l2ForkChoice   forkChoiceState
	blockRelations l2types.BlockRelations
}

func NewL2ForkChoiceUpdater(execBackend ExecutionBackend, l1State L1State) *L2ForkchoiceUpdater {
	return &L2ForkchoiceUpdater{execBackend: execBackend, l1State: l1State}
}

// Always ingest.
func (s *L2ForkchoiceUpdater) hasNext() bool { return false }
func (s *L2ForkchoiceUpdater) next() any     { return nil }

// 1. Ingests a new block relation.
// 2. Gets an updated L1 fork-choice.
// 3. Derives the corresponding L2 fork-choice.
func (s *L2ForkchoiceUpdater) ingest(ctx context.Context, relation l2types.BlockRelation) error {
	// TODO: handle old block relations (older than the fc marking)
	if relation != (l2types.BlockRelation{}) {
		err := s.blockRelations.Append(relation)
		if err != nil {
			return fmt.Errorf("Failed to append block relation: %w", err)
		}
	}
	// Get latest L1 fork-choice.
	updatedL1ForkChoice := forkChoiceState{s.l1State.Head(), s.l1State.Safe(), s.l1State.Finalized()}
	// Skip if no change.
	if updatedL1ForkChoice == s.l1ForkChoice {
		log.Info("No change in l1 fork-choice, skipping.")
		return nil
	}
	// Derive l2 fork-choice from l1 fork-choice.
	safeL2BlockID := s.blockRelations.MarkSafe(uint64(updatedL1ForkChoice.headID.Number()))
	finalizedL2BlockID := s.blockRelations.MarkFinal(uint64(updatedL1ForkChoice.finalizedID.Number()))
	l2Forkchoice := forkChoiceState{
		headID:      l2types.NewBlockID(0, common.Hash{}), // TODO
		safeID:      safeL2BlockID,
		finalizedID: finalizedL2BlockID,
	}
	// Skip if no change.
	if l2Forkchoice == s.l2ForkChoice {
		log.Info("No change in l2 fork-choice, skipping.")
		return nil
	}
	// Update fork-choice.
	err := s.execBackend.ForkchoiceUpdate(&l2Forkchoice)
	if err != nil {
		return fmt.Errorf("Failed to update fork-choice state: %w", err)
	}
	s.l1ForkChoice = updatedL1ForkChoice
	s.l2ForkChoice = l2Forkchoice
	return nil
}

func (s *L2ForkchoiceUpdater) recover(ctx context.Context, l1BlockID l2types.BlockID) error {
	s.blockRelations.MarkReorgedOut(l1BlockID.Number())
	return nil
}

func (s *L2ForkchoiceUpdater) findRecoveryPoint(ctx context.Context) (l2types.BlockID, error) {
	return l2types.BlockID{}, nil
}

// func (d *Driver[T]) recover(ctx context.Context) error {
// 	lastPlausibleBlockID, err := findLatestPlausibleL1BlockID(ctx, d.l1Client)
// 	if err != nil {
// 		err = fmt.Errorf("Failed to find plausible L1 block ID for recovery: %w", err)
// 	} else {
// 		err = d.terminalStage.Recover(ctx, lastPlausibleBlockID)
// 	}
// 	return err
// }

// Finds the "plausible" tip of the canonical chain, walking back from the latest L1 head.
// Plausible <=> block number is known and canonical, OR unknown
// func findLatestPlausibleL1BlockID(ctx context.Context, l1Client EthClient) (l2types.BlockID, error) {
// 	latest, err := l1Client.HeaderByTag(ctx, client.Latest)
// 	if err != nil {
// 		return l2types.BlockID{}, RetryableError{fmt.Errorf("Could not get latest L1 block header: %w", err)}
// 	}
// 	return l2types.NewBlockIDFromHeader(header), nil
// }

type forkChoiceState struct {
	headID      l2types.BlockID
	safeID      l2types.BlockID
	finalizedID l2types.BlockID
}

func (fcs *forkChoiceState) HeadBlockHash() common.Hash      { return fcs.headID.Hash() }
func (fcs *forkChoiceState) SafeBlockHash() common.Hash      { return fcs.safeID.Hash() }
func (fcs *forkChoiceState) FinalizedBlockHash() common.Hash { return fcs.finalizedID.Hash() }
