package stage

import (
	"context"
	"fmt"
	"math/big"

	"github.com/specularl2/specular/clients/geth/specular/rollup/l2types"
	"github.com/specularl2/specular/clients/geth/specular/rollup/utils"
)

type L1HeaderRetrievalStage struct {
	currL1BlockID l2types.BlockID
	l1Client      EthClient
}

func (s *L1HeaderRetrievalStage) Step(ctx context.Context) (l2types.BlockID, error) {
	nextHeader, err := s.l1Client.HeaderByNumber(ctx, big.NewInt(0).SetUint64(s.currL1BlockID.Number()+1))
	if err != nil {
		return l2types.BlockID{}, &RetryableError{fmt.Errorf("failed to get next L1 header: %w", err)}
	}
	if nextHeader.ParentHash != s.currL1BlockID.Hash() {
		return l2types.BlockID{}, &utils.L1ReorgDetectedError{Msg: "received parent hash does not match current L1 block hash"}
	}
	s.currL1BlockID = l2types.NewBlockIDFromHeader(nextHeader)
	return l2types.NewBlockIDFromHeader(nextHeader), nil
}

func (s *L1HeaderRetrievalStage) Recover(ctx context.Context, l1BlockID l2types.BlockID) error {
	s.currL1BlockID = l1BlockID
	return nil
}
