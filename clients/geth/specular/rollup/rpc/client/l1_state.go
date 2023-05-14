package client

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/specularl2/specular/clients/geth/specular/rollup/l2types"
	"github.com/specularl2/specular/clients/geth/specular/rollup/utils"
	"github.com/specularl2/specular/clients/geth/specular/rollup/utils/log"
)

// Thread-safe. Tracks the latest, safe and last finalized L1 headers received.
type L1State struct {
	// Thread-safe map from block tag to block ID.
	headers utils.Map[BlockTag, l2types.BlockID]
}

var _ OnNewHandler = (*L1State)(nil)

func NewL1State() *L1State { return &L1State{} }

func (s *L1State) Head() l2types.BlockID      { return s.headers.Load(Latest) }
func (s *L1State) Safe() l2types.BlockID      { return s.headers.Load(Safe) }
func (s *L1State) Finalized() l2types.BlockID { return s.headers.Load(Finalized) }

func (s *L1State) OnLatest(_ context.Context, header *types.Header) error {
	prev := s.headers.LoadAndStore(Latest, l2types.NewBlockIDFromHeader(header))
	if header.Number.Uint64() <= prev.Number() {
		log.Warn(
			"Received old latest L1 header",
			"number", header.Number, "hash", header.Hash(),
			"prev_number", prev.Number(), "prev_hash", prev.Hash(),
		)
	}
	return nil
}

func (s *L1State) OnSafe(_ context.Context, header *types.Header) error {
	prev := s.Safe()
	if header.Number.Uint64() < prev.Number() {
		// Assuming L1 safety, this should only happen due to network issues / slow nodes.
		log.Warn(
			"Received old safe L1 header; ignoring.",
			"recvd_number", header.Number, "recvd_hash", header.Hash(),
		)
		return nil
	} else if header.Number.Uint64() == prev.Number() {
		if header.Hash() != prev.Hash() {
			return fmt.Errorf("received two safe headers for block_num=%d with hashes: %s and %s", prev.Number(), header.Hash(), prev.Hash())
		}
	}
	s.headers.Store(Safe, l2types.NewBlockIDFromHeader(header))
	return nil
}

func (s *L1State) OnFinalized(_ context.Context, header *types.Header) error {
	prev := s.Finalized()
	if header.Number.Uint64() < prev.Number() {
		// Assuming L1 safety, this should only happen due to network issues / slow nodes.
		log.Warn(
			"Received old finalized L1 header; ignoring.",
			"recvd_number", header.Number, "recvd_hash", header.Hash(),
		)
		return nil
	} else if header.Number.Uint64() == prev.Number() {
		if header.Hash() != prev.Hash() {
			return fmt.Errorf("received two finalized headers for block_num=%d with hashes: %s and %s", prev.Number(), header.Hash(), prev.Hash())
		}
	}
	s.headers.Store(Finalized, l2types.NewBlockIDFromHeader(header))
	return nil
}
