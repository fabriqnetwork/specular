package eth

import (
	"context"
	"fmt"

	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/specularL2/specular/services/sidecar/rollup/types"
	"github.com/specularL2/specular/services/sidecar/utils"
	"github.com/specularL2/specular/services/sidecar/utils/log"
)

// Thread-safe. Tracks the latest, last safe and last finalized L1 headers received.
type EthState struct {
	// Thread-safe map from BlockTag to last corresponding BlockID.
	headers utils.Map[BlockTag, types.BlockID]
}

func NewEthState() *EthState { return &EthState{} }

func (s *EthState) Head() types.BlockID      { return s.headers.Load(Latest) }
func (s *EthState) Safe() types.BlockID      { return s.headers.Load(Safe) }
func (s *EthState) Finalized() types.BlockID { return s.headers.Load(Finalized) }
func (s *EthState) Tips() (types.BlockID, types.BlockID, types.BlockID) {
	return s.Head(), s.Safe(), s.Finalized()
}

func (s *EthState) OnLatest(_ context.Context, header *ethTypes.Header) error {
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

func (s *EthState) OnSafe(_ context.Context, header *ethTypes.Header) error {
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

func (s *EthState) OnFinalized(_ context.Context, header *ethTypes.Header) error {
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
