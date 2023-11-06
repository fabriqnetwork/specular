package bridge

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/specularL2/specular/services/sidecar/bindings"
	"github.com/specularL2/specular/services/sidecar/rollup/rpc/eth"
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
