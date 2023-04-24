package state

import (
	"context"
	"sync"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/specularl2/specular/clients/geth/specular/rollup/comms/client"
	"github.com/specularl2/specular/clients/geth/specular/rollup/utils/log"
)

// Tracks the latest and last finalized L1 headers received.
type L1State struct {
	head      *types.Header
	safe      *types.Header
	finalized *types.Header

	headLock      *sync.RWMutex
	safeLock      *sync.RWMutex
	finalizedLock *sync.RWMutex
}

var _ client.OnNewHandler = (*L1State)(nil)

func NewL1State() *L1State {
	return &L1State{
		headLock:      &sync.RWMutex{},
		finalizedLock: &sync.RWMutex{},
	}
}

func (s *L1State) Head() *types.Header {
	s.headLock.RLock()
	defer s.headLock.RUnlock()
	return s.head
}

func (s *L1State) Safe() *types.Header {
	s.safeLock.RLock()
	defer s.safeLock.RUnlock()
	return s.safe
}

func (s *L1State) Finalized() *types.Header {
	s.finalizedLock.RLock()
	defer s.finalizedLock.RUnlock()
	return s.finalized
}

func (s *L1State) OnLatest(_ context.Context, header *types.Header) error {
	s.headLock.Lock()
	defer s.headLock.Unlock()
	if header.Number.Uint64() < s.head.Number.Uint64() {
		log.Warn(
			"Received old L1 header",
			"number", header.Number, "hash", header.Hash(),
			"prev", s.head.Number, "hash", s.head.Hash(),
		)
	}
	s.head = header
	return nil
}

func (s *L1State) OnSafe(_ context.Context, header *types.Header) error {
	s.safeLock.Lock()
	defer s.safeLock.Unlock()
	if header.Number.Uint64() < s.safe.Number.Uint64() {
		log.Crit(
			"Received old safe L1 header (should rarely happen)",
			"number", header.Number, "hash", header.Hash(),
		)
	}
	// if header.Number.Uint64() > s.Head().Number.Uint64() {
	// 	log.Crit(
	// 		"Received safe L1 header that is newer than latest L1 header (should never happen)",
	// 		"number", header.Number, "hash", header.Hash(),
	// 		"latest", s.Head().Number, "hash", s.Head().Hash(),
	// 	)
	// }
	s.safe = header
	return nil
}

func (s *L1State) OnFinalized(_ context.Context, header *types.Header) error {
	s.finalizedLock.Lock()
	defer s.finalizedLock.Unlock()
	if header.Number.Uint64() < s.finalized.Number.Uint64() {
		log.Crit(
			"Received old finalized L1 header (should basically never happen)",
			"number", header.Number, "hash", header.Hash(),
		)
	}
	s.finalized = header
	return nil
}
