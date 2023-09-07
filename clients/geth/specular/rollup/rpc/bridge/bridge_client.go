package bridge

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/specularl2/specular/clients/geth/specular/bindings"
	"github.com/specularl2/specular/clients/geth/specular/rollup/rpc/eth"
)

type BridgeClient struct {
	*bindings.ISequencerInbox
	*bindings.IRollup
}

type L1Config interface {
	GetEndpoint() string
	GetSequencerInboxAddr() common.Address
	GetRollupAddr() common.Address
}

func NewBridgeClient(backend bind.ContractBackend, cfg L1Config) (*BridgeClient, error) {
	inbox, err := bindings.NewISequencerInbox(cfg.GetSequencerInboxAddr(), backend)
	if err != nil {
		return nil, err
	}
	rollup, err := bindings.NewIRollup(cfg.GetRollupAddr(), backend)
	if err != nil {
		return nil, err
	}
	return &BridgeClient{ISequencerInbox: inbox, IRollup: rollup}, nil
}

func DialWithRetry(ctx context.Context, cfg L1Config) (*BridgeClient, error) {
	l1Client, err := eth.DialWithRetry(ctx, cfg.GetEndpoint(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to dial L1 client: %v", err)
	}
	return NewBridgeClient(l1Client, cfg)
}

func (c *BridgeClient) RequireFirstUnresolvedAssertionIsConfirmable(ctx context.Context) error {
	return c.IRollup.RequireFirstUnresolvedAssertionIsConfirmable(&bind.CallOpts{Pending: false, Context: ctx})
}

func (c *BridgeClient) GetStaker(ctx context.Context, addr common.Address) (bindings.IRollupStaker, error) {
	return c.IRollup.GetStaker(&bind.CallOpts{Pending: false, Context: ctx}, addr)
}

func (c *BridgeClient) GetAssertion(ctx context.Context, assertionID *big.Int) (bindings.IRollupAssertion, error) {
	return c.IRollup.GetAssertion(&bind.CallOpts{Pending: false, Context: ctx}, assertionID)
}

func (c *BridgeClient) GetLastConfirmedAssertionID(ctx context.Context) (*big.Int, error) {
	return c.IRollup.GetLastConfirmedAssertionID(&bind.CallOpts{Pending: false, Context: ctx})
}

func (c *BridgeClient) GetRequiredStakeAmount(ctx context.Context) (*big.Int, error) {
	return c.IRollup.CurrentRequiredStake(&bind.CallOpts{Pending: false, Context: ctx})
}

// Returns the last assertion ID that was validated *by us*.
// func (c *BridgeClient) GetLastValidatedAssertionID(opts *bind.FilterOpts) (*big.Int, error) {
// 	iter, err := c.IRollup.FilterStakerStaked(opts)
// 	if err != nil {
// 		return nil, fmt.Errorf("Failed to filter through `StakerStaked` events to get last validated assertion ID, err: %w", err)
// 	}
// 	lastValidatedAssertionID := common.Big0
// 	for iter.Next() {
// 		// Note: the second condition should always hold true if the iterator iterates in time order.
// 		if iter.Event.StakerAddr == c.transactOpts.From && iter.Event.AssertionID.Cmp(lastValidatedAssertionID) == 1 {
// 			log.Debug("StakerStaked event found", "staker", iter.Event.StakerAddr, "assertionID", iter.Event.AssertionID)
// 			lastValidatedAssertionID = iter.Event.AssertionID
// 		}
// 	}
// 	if iter.Error() != nil {
// 		return nil, fmt.Errorf("Failed to iterate through validated assertion IDs, err: %w", iter.Error())
// 	}
// 	if lastValidatedAssertionID.Cmp(common.Big0) == 0 {
// 		return nil, fmt.Errorf("No validated assertion IDs found")
// 	}
// 	return lastValidatedAssertionID, nil
// }

// func (c *BridgeClient) GetGenesisAssertionCreated(opts *bind.FilterOpts) (*bindings.IRollupAssertionCreated, error) {
// 	// We could probably do this from initialization calldata too.
// 	iter, err := c.IRollup.FilterAssertionCreated(opts)
// 	if err != nil {
// 		return nil, fmt.Errorf("Failed to filter through `AssertionCreated` events to get genesis assertion ID, err: %w", err)
// 	}
// 	if iter.Next() {
// 		return iter.Event, nil
// 	}
// 	if iter.Error() != nil {
// 		return nil, fmt.Errorf("No genesis `AssertionCreated` event found, err: %w", iter.Error())
// 	}
// 	return nil, fmt.Errorf("No genesis `AssertionCreated` event found")
// }
