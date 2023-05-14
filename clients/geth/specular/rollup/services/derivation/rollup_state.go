package derivation

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/specularl2/specular/clients/geth/specular/bindings"
	"github.com/specularl2/specular/clients/geth/specular/rollup/l2types"
	"github.com/specularl2/specular/clients/geth/specular/rollup/rpc/bridge"
	"github.com/specularl2/specular/clients/geth/specular/rollup/utils/fmt"
)

// Tracks L2 state as a function of synced L1 state.
type RollupState struct {
	stakers         map[common.Address]*Staker
	assertions      map[*big.Int]*Assertion
	assertionsState map[*big.Int]*AssertionState
	assertionsAux   map[*big.Int]*AssertionAux

	l1Client *bridge.BridgeClient
}

type Staker struct{ bindings.IRollupStaker }
type Assertion struct{ bindings.IRollupAssertion }
type AssertionState struct{ stakers map[common.Address]bool }
type AssertionAux struct {
	vmHash    common.Hash
	inboxSize *big.Int
}

func NewRollupState(l1Client *bridge.BridgeClient) *RollupState {
	return &RollupState{l1Client: l1Client}
}

func (s *RollupState) GetStaker(ctx context.Context, stakerAddr common.Address) (*Staker, error) {
	staker := s.stakers[stakerAddr]
	if staker != nil {
		return staker, nil
	}
	opts := bind.CallOpts{Pending: false, Context: ctx}
	recvdStaker, err := s.l1Client.GetStaker(&opts, stakerAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to get staker %v: %w", stakerAddr, err)
	}
	s.stakers[stakerAddr] = &Staker{IRollupStaker: recvdStaker}
	return staker, nil
}

// Note: the assertion struct itself is immutable.
func (s *RollupState) GetAssertion(ctx context.Context, assertionID *big.Int) (*Assertion, error) {
	assertion := s.assertions[assertionID]
	if assertion != nil {
		return assertion, nil
	}
	opts := bind.CallOpts{Pending: false, Context: ctx}
	recvdAssertion, err := s.l1Client.GetAssertion(&opts, assertionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get assertion %v: %w", assertionID, err)
	}
	s.assertions[assertionID] = &Assertion{IRollupAssertion: recvdAssertion}
	return assertion, nil
}

func (s *RollupState) OnAssertionCreated(
	ctx context.Context,
	l1BlockID l2types.BlockID,
	tx *types.Transaction,
) error {
	return nil
}

func (s *RollupState) OnAssertionConfirmed(
	ctx context.Context,
	l1BlockID l2types.BlockID,
	tx *types.Transaction,
) error {
	return nil
}

func (s *RollupState) OnAssertionRejected(
	ctx context.Context,
	l1BlockID l2types.BlockID,
	tx *types.Transaction,
) error {
	return nil
}

// func (r *RollupState) StartSync(ctx context.Context, l1Client, l2Client client.EthPollingClient) {
// 	for r.L1State.Head() == nil {
// 		log.Info("Waiting for L1 latest header...")
// 		time.Sleep(100 * time.Millisecond)
// 	}
// 	log.Info("Latest header received", "number", r.L1State.Head().Number)
// 	for r.L1State.Finalized() == nil {
// 		log.Info("Waiting for L1 finalized header...")
// 		time.Sleep(100 * time.Millisecond)
// 	}
// 	log.Info("Finalized header received", "number", r.L1State.Head().Number)
// }
