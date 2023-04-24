package state

import (
	"context"
	"sync"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/specularl2/specular/clients/geth/specular/rollup/comms/client"
)

// State of L2 chain, to be derived from L1 state.
type L2State struct {
	head      *types.Header
	safe      *types.Header
	finalized *types.Header

	// L1 state to derive L2 state from.
	l1State ReadOnlyL1State

	headLock      *sync.RWMutex
	safeLock      *sync.RWMutex
	finalizedLock *sync.RWMutex
}

var _ client.OnNewHandler = (*L2State)(nil)

func NewL2State() *L2State {
	return &L2State{
		headLock:      &sync.RWMutex{},
		safeLock:      &sync.RWMutex{},
		finalizedLock: &sync.RWMutex{},
	}
}

type ReadOnlyL1State interface {
	Head() *types.Header
	Safe() *types.Header
	Finalized() *types.Header
}

func (s *L2State) Head() *types.Header {
	s.headLock.RLock()
	defer s.headLock.RUnlock()
	return s.head
}

func (s *L2State) Safe() *types.Header {
	s.safeLock.RLock()
	defer s.safeLock.RUnlock()
	return s.safe
}

func (s *L2State) Finalized() *types.Header {
	s.finalizedLock.RLock()
	defer s.finalizedLock.RUnlock()
	return s.finalized
}

func (s *L2State) OnLatest(_ context.Context, header *types.Header) error {
	s.headLock.RLock()
	defer s.headLock.RUnlock()
	s.head = header
	// TODO: check l1 state to mark safe/finalize.
	// This requires tracking which L1 block corresponds to which L2 block.
	return nil
}

// Must be derived from L1 state directly since Engine API not used.
func (s *L2State) OnSafe(_ context.Context, header *types.Header) error {
	panic("should not be called")
}

// Must be derived from L1 state directly since Engine API not used.
func (s *L2State) OnFinalized(_ context.Context, header *types.Header) error {
	panic("should not be called")
}
