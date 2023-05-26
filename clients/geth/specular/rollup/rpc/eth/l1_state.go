package eth

import (
	"context"
	"fmt"

	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/specularl2/specular/clients/geth/specular/rollup/types"
	"github.com/specularl2/specular/clients/geth/specular/rollup/utils"
	"github.com/specularl2/specular/clients/geth/specular/rollup/utils/log"
)

// Thread-safe. Tracks the latest, last safe and last finalized L1 headers received.
type L1State struct {
	// Thread-safe map from BlockTag to last corresponding BlockID.
	headers utils.Map[BlockTag, types.BlockID]
}

var _ OnNewHandler = (*L1State)(nil)

func NewL1State() *L1State { return &L1State{} }

func (s *L1State) Head() types.BlockID      { return s.headers.Load(Latest) }
func (s *L1State) Safe() types.BlockID      { return s.headers.Load(Safe) }
func (s *L1State) Finalized() types.BlockID { return s.headers.Load(Finalized) }

func (s *L1State) OnLatest(_ context.Context, header *ethTypes.Header) error {
	prev := s.headers.LoadAndStore(Latest, types.NewBlockIDFromHeader(header))
	if header.Number.Uint64() <= prev.GetNumber() {
		log.Warn(
			"Received old latest L1 header",
			"number", header.Number, "hash", header.Hash(),
			"prev_number", prev.GetNumber(), "prev_hash", prev.GetHash(),
		)
	}
	return nil
}

func (s *L1State) OnSafe(_ context.Context, header *ethTypes.Header) error {
	prev := s.Safe()
	if header.Number.Uint64() < prev.GetNumber() {
		// Assuming L1 safety, this should only happen due to network issues / slow nodes.
		log.Warn(
			"Received old safe L1 header; ignoring.",
			"recvd_number", header.Number, "recvd_hash", header.Hash(),
		)
		return nil
	} else if header.Number.Uint64() == prev.GetNumber() {
		if header.Hash() != prev.GetHash() {
			return fmt.Errorf("received two safe headers for block_num=%d with hashes: %s and %s", prev.GetNumber(), header.Hash(), prev.GetHash())
		}
	}
	s.headers.Store(Safe, types.NewBlockIDFromHeader(header))
	return nil
}

func (s *L1State) OnFinalized(_ context.Context, header *ethTypes.Header) error {
	prev := s.Finalized()
	if header.Number.Uint64() < prev.GetNumber() {
		// Assuming L1 safety, this should only happen due to network issues / slow nodes.
		log.Warn(
			"Received old finalized L1 header; ignoring.",
			"recvd_number", header.Number, "recvd_hash", header.Hash(),
		)
		return nil
	} else if header.Number.Uint64() == prev.GetNumber() {
		if header.Hash() != prev.GetHash() {
			return fmt.Errorf(
				"received two finalized headers for block_num=%d with hashes: %s and %s",
				prev.GetNumber(), header.Hash(), prev.GetHash(),
			)
		}
	}
	s.headers.Store(Finalized, types.NewBlockIDFromHeader(header))
	return nil
}
