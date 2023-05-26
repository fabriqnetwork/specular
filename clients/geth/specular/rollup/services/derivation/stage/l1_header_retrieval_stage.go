package stage

import (
	"context"
	"fmt"
	"math/big"

	"github.com/specularl2/specular/clients/geth/specular/rollup/types"
	"github.com/specularl2/specular/clients/geth/specular/rollup/utils"
)

// Note: source stage (no prev stage).
type L1HeaderRetrievalStage struct {
	currL1BlockID types.BlockID
	l1Client      EthClient
}

func (s *L1HeaderRetrievalStage) Pull(ctx context.Context) (types.BlockID, error) {
	nextHeader, err := s.l1Client.HeaderByNumber(ctx, big.NewInt(0).SetUint64(s.currL1BlockID.GetNumber()+1))
	if err != nil {
		return types.BlockID{}, &RetryableError{fmt.Errorf("failed to get next L1 header: %w", err)}
	}
	if nextHeader.ParentHash != s.currL1BlockID.GetHash() {
		return types.BlockID{}, &utils.L1ReorgDetectedError{Msg: "received parent hash does not match current L1 block hash"}
	}
	s.currL1BlockID = types.NewBlockIDFromHeader(nextHeader)
	return s.currL1BlockID, nil
}

func (s *L1HeaderRetrievalStage) Recover(ctx context.Context, l1BlockID types.BlockID) error {
	s.currL1BlockID = l1BlockID
	return nil
}
