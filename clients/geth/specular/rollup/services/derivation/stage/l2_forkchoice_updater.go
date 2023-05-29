package stage

import (
	"context"
	"fmt"

	"github.com/specularl2/specular/clients/geth/specular/rollup/types"
	"github.com/specularl2/specular/clients/geth/specular/rollup/utils/log"
)

type L2ForkchoiceUpdater struct {
	execBackend ExecutionBackend
	l1State     L1State

	l1ForkChoice   ForkChoiceState
	l2ForkChoice   ForkChoiceState
	blockRelations types.BlockRelations
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
func (s *L2ForkchoiceUpdater) ingest(ctx context.Context, relation types.BlockRelation) error {
	// TODO: handle 'old' block relations (older than the fc marking)
	// TODO: update fork choice for latest
	if relation != (types.BlockRelation{}) {
		err := s.blockRelations.Append(relation)
		if err != nil {
			return fmt.Errorf("Failed to append block relation: %w", err)
		}
		log.Info("Appended block relation.", "l1", relation.L1BlockID.Number, "l2", relation.L2BlockID.Number)
	}
	// Get latest L1 fork-choice.
	var (
		l1Head              = s.l1State.Head()
		l1Safe              = s.l1State.Safe()
		l1Finalized         = s.l1State.Finalized()
		updatedL1ForkChoice = ForkChoiceState{
			HeadBlockHash:      l1Head.GetHash(),
			SafeBlockHash:      l1Safe.GetHash(),
			FinalizedBlockHash: l1Finalized.GetHash(),
		}
	)
	// Skip if no change.
	if updatedL1ForkChoice == s.l1ForkChoice {
		log.Info("No change in l1 fork-choice, skipping.")
		return nil
	}
	// Derive l2 fork-choice from l1 fork-choice.
	var (
		safeL2BlockID      = s.blockRelations.MarkSafe(uint64(l1Head.GetNumber()))
		finalizedL2BlockID = s.blockRelations.MarkFinal(uint64(l1Finalized.GetNumber()))
		l2Forkchoice       = ForkChoiceState{
			HeadBlockHash:      s.l2ForkChoice.HeadBlockHash,
			SafeBlockHash:      safeL2BlockID.GetHash(),
			FinalizedBlockHash: finalizedL2BlockID.GetHash(),
		}
	)
	// Skip if no change.
	if l2Forkchoice == s.l2ForkChoice {
		log.Info("No change in l2 fork-choice, skipping.")
		return nil
	}
	// Update fork-choice.
	response, err := s.execBackend.ForkchoiceUpdate(&l2Forkchoice)
	if err != nil {
		return fmt.Errorf("Failed to update fork-choice state: %w", err)
	}
	s.l1ForkChoice = updatedL1ForkChoice
	s.l2ForkChoice = l2Forkchoice
	s.l2ForkChoice.HeadBlockHash = *response.PayloadStatus.LatestValidHash
	return nil
}

func (s *L2ForkchoiceUpdater) recover(ctx context.Context, l1BlockID types.BlockID) error {
	s.blockRelations.MarkReorgedOut(l1BlockID.GetNumber())
	return nil
}

func (s *L2ForkchoiceUpdater) findRecoveryPoint(ctx context.Context) (types.BlockID, error) {
	return types.BlockID{}, nil
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
// func findLatestPlausibleL1BlockID(ctx context.Context, l1Client EthClient) (types.BlockID, error) {
// 	latest, err := l1Client.HeaderByTag(ctx, client.Latest)
// 	if err != nil {
// 		return types.BlockID{}, RetryableError{fmt.Errorf("Could not get latest L1 block header: %w", err)}
// 	}
// 	return types.NewBlockIDFromHeader(header), nil
// }
