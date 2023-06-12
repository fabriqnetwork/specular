package stage

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/specularl2/specular/clients/geth/specular/rollup/rpc/eth"
	"github.com/specularl2/specular/clients/geth/specular/rollup/types"
	"github.com/specularl2/specular/clients/geth/specular/utils/fmt"
	"github.com/specularl2/specular/clients/geth/specular/utils/log"
)

type L2ForkchoiceUpdater struct {
	cfg         GenesisConfig
	execBackend ExecutionBackend
	l2Client    L2Client
	l1State     L1State

	l1ForkChoice   ForkChoiceState
	l2ForkChoice   ForkChoiceState
	blockRelations types.BlockRelations
}

func NewL2ForkChoiceUpdater(
	cfg GenesisConfig,
	execBackend ExecutionBackend,
	l2Client L2Client,
	l1State L1State,
) *L2ForkchoiceUpdater {
	return &L2ForkchoiceUpdater{cfg: cfg, execBackend: execBackend, l2Client: l2Client, l1State: l1State}
}

// Always ingest.
func (s *L2ForkchoiceUpdater) hasNext() bool { return false }
func (s *L2ForkchoiceUpdater) next() any     { return nil }

// 1. Ingests a new block relation.
// 2. Gets an updated L1 fork-choice.
// 3. Derives the corresponding L2 fork-choice.
func (s *L2ForkchoiceUpdater) ingest(ctx context.Context, relation types.BlockRelation) error {
	if err := s.ensureForkChoiceInitialized(ctx); err != nil {
		return fmt.Errorf("Failed to ensure fork-choice initialized: %w", err)
	}
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
		log.Trace("No change in l1 fork-choice, skipping.")
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
	fcResponse, err := s.execBackend.ForkchoiceUpdate(&l2Forkchoice)
	if err != nil {
		return fmt.Errorf("Failed to update fork-choice state: %w", err)
	}
	s.l1ForkChoice = updatedL1ForkChoice
	s.l2ForkChoice = l2Forkchoice
	s.l2ForkChoice.HeadBlockHash = *fcResponse.PayloadStatus.LatestValidHash
	return nil
}

func (s *L2ForkchoiceUpdater) recover(ctx context.Context, l1BlockID types.BlockID) error {
	s.blockRelations.MarkReorgedOut(l1BlockID.GetNumber())
	return nil
}

func (s *L2ForkchoiceUpdater) findRecoveryPoint(ctx context.Context) (types.BlockID, error) {
	return types.BlockID{}, nil
}

func (s *L2ForkchoiceUpdater) ensureForkChoiceInitialized(ctx context.Context) error {
	err := s.l2Client.EnsureDialed(ctx)
	if err != nil {
		return RetryableError{fmt.Errorf("failed to create l2 client: %w", err)}
	}
	if s.l2ForkChoice != (ForkChoiceState{}) {
		return nil
	}
	forkChoice, err := GetForkChoice(ctx, s.l2Client, s.cfg)
	if err != nil {
		return fmt.Errorf("Failed to get genesis-floored fork-choice state: %w", err)
	}
	_, err = s.execBackend.ForkchoiceUpdate(&forkChoice)
	if err != nil {
		return fmt.Errorf("Failed to set initial fork-choice state: %w", err)
	}
	s.l2ForkChoice = forkChoice
	log.Info("Initialized fork-choice state.")
	return nil
}

// Gets the current fork-choice state from the L2 client. Defaults to genesis block for safe and finalized.
func GetForkChoice(ctx context.Context, l2Client L2Client, cfg GenesisConfig) (ForkChoiceState, error) {
	forkChoice, err := eth.GetForkChoice(ctx, l2Client)
	if err != nil {
		return ForkChoiceState{}, RetryableError{fmt.Errorf("Failed to get fork-choice state: %w", err)}
	}
	if forkChoice.FinalizedBlockHash == (common.Hash{}) {
		log.Info("Using Genesis block for finalized block.")
		genesisL2Header, err := l2Client.HeaderByNumber(ctx, common.Big0)
		if err != nil {
			return ForkChoiceState{}, RetryableError{fmt.Errorf("Failed to get genesis L2 header: %w", err)}
		}
		forkChoice.FinalizedBlockHash = genesisL2Header.Hash()
	}
	if forkChoice.SafeBlockHash == (common.Hash{}) {
		forkChoice.SafeBlockHash = forkChoice.FinalizedBlockHash
	}
	return forkChoice, nil
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
