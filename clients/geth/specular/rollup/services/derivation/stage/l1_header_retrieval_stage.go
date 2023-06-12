package stage

import (
	"context"
	"fmt"
	"math/big"

	"github.com/specularl2/specular/clients/geth/specular/rollup/types"
	"github.com/specularl2/specular/clients/geth/specular/utils/log"
)

// Note: source stage (no prev stage).
type L1HeaderRetrievalStage struct {
	currL1BlockID types.BlockID
	l1Client      L1Client
}

func (s *L1HeaderRetrievalStage) Pull(ctx context.Context) (types.BlockID, error) {
	nextHeader, err := s.l1Client.HeaderByNumber(ctx, big.NewInt(0).SetUint64(s.currL1BlockID.GetNumber()+1))
	if err != nil {
		return types.EmptyBlockID, RetryableError{fmt.Errorf("failed to get next L1 header: %w", err)}
	}
	nextL1BlockID := types.NewBlockIDFromHeader(nextHeader)
	log.Info("Retrieved next L1 header.", "id", nextL1BlockID)
	if nextHeader.ParentHash != s.currL1BlockID.GetHash() {
		return types.EmptyBlockID, RecoverableError{
			fmt.Errorf("next L1 block %s not linked to %s", nextL1BlockID, s.currL1BlockID),
		}
	}
	s.currL1BlockID = nextL1BlockID
	return s.currL1BlockID, nil
}

func (s *L1HeaderRetrievalStage) Recover(ctx context.Context, l1BlockID types.BlockID) error {
	s.currL1BlockID = l1BlockID
	return nil
}
