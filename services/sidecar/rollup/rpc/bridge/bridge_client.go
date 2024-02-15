package bridge

import (
	"context"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"

	"github.com/specularL2/specular/bindings-go/bindings"
	"github.com/specularL2/specular/services/sidecar/utils/fmt"
)

type BridgeClient struct {
	*bindings.ISequencerInbox
	*bindings.IRollup
}

type ProtocolConfig interface {
	GetSequencerInboxAddr() common.Address
	GetRollupAddr() common.Address
}

type UnsatisfiedCondition = string

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

func (c *BridgeClient) RequireFirstUnresolvedAssertionIsConfirmable(ctx context.Context) (*UnsatisfiedCondition, error) {
	err := c.IRollup.RequireFirstUnresolvedAssertionIsConfirmable(&bind.CallOpts{Pending: false, Context: ctx})
	return processRollupError(err)
}

func (c *BridgeClient) RequireFirstUnresolvedAssertionIsRejectable(ctx context.Context, address common.Address) (*UnsatisfiedCondition, error) {
	err := c.IRollup.RequireFirstUnresolvedAssertionIsRejectable(&bind.CallOpts{Pending: false, Context: ctx}, address)
	return processRollupError(err)
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

func (c *BridgeClient) IsStakedOnAssertion(ctx context.Context, assertionID *big.Int, address common.Address) (bool, error) {
	return c.IRollup.IsStakedOnAssertion(&bind.CallOpts{Pending: false, Context: ctx}, assertionID, address)
}

func processRollupError(err error) (*UnsatisfiedCondition, error) {
	if err == nil {
		return nil, nil
	}
	var dataErr rpc.DataError
	if errors.As(err, &dataErr) {
		unpackedErr, err := UnpackRollupError(dataErr)
		if err != nil {
			return nil, fmt.Errorf("failed to unpack rollup error: %w", err)
		}
		return &unpackedErr.Name, nil
	}
	return nil, fmt.Errorf("failed call with unknown error: %w", err)
}
