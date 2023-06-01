package derivation

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/specularl2/specular/clients/geth/specular/bindings"
	"github.com/specularl2/specular/clients/geth/specular/rollup/rpc/bridge"
	"github.com/specularl2/specular/clients/geth/specular/rollup/types"
	"github.com/specularl2/specular/clients/geth/specular/utils/fmt"
)

// Tracks L2 state as a function of synced L1 state.
// Mirrors `Rollup.sol`.
type RollupState struct {
	l1Client *bridge.BridgeClient

	lastResolvedAssertionID  *big.Int
	lastConfirmedAssertionID *big.Int
	lastCreatedAssertionID   *big.Int
	// Assertion state
	assertions      map[*big.Int]*Assertion
	assertionsState map[*big.Int]*AssertionState
	// Staking state
	numStakers *big.Int
	stakers    map[common.Address]*Staker
}

type Staker struct{ bindings.IRollupStaker }
type Assertion struct{ bindings.IRollupAssertion }
type AssertionState struct{ stakers map[common.Address]bool }

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
	l1BlockID types.BlockID,
	tx *ethTypes.Transaction,
) error {
	receipt, err := s.l1Client.TransactionReceipt(ctx, tx.Hash())
	if err != nil {
		return fmt.Errorf("failed to get receipt for `createAssertion` tx %v: %w", tx.Hash(), err)
	}
	var assertionID = big.NewInt(receipt.BlockNumber.Int64()) // TODO
	_, err = s.GetAssertion(ctx, assertionID)
	if err != nil {
		return fmt.Errorf("failed to get assertion %v: %w", assertionID, err)
	}
	s.lastCreatedAssertionID = assertionID
	return nil
}

func (s *RollupState) OnAssertionConfirmed(
	ctx context.Context,
	l1BlockID types.BlockID,
	tx *ethTypes.Transaction,
) error {
	// s.lastResolvedAssertionID
	return nil
}

func (s *RollupState) OnAssertionRejected(
	ctx context.Context,
	l1BlockID types.BlockID,
	tx *ethTypes.Transaction,
) error {
	return nil
}

func (s *RollupState) OnReorg(l1BlockID types.BlockID) {
	return
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
