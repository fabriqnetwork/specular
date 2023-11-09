package bridge

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/specularL2/specular/services/sidecar/bindings"
)

type BridgeClient struct {
	*bindings.ISequencerInbox
	*bindings.IRollup
}

type ProtocolConfig interface {
	GetSequencerInboxAddr() common.Address
	GetRollupAddr() common.Address
}

func NewBridgeClient(backend bind.ContractBackend, cfg ProtocolConfig) (*BridgeClient, error) {
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
