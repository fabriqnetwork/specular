package bridge

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/specularl2/specular/clients/geth/specular/bindings"
	"github.com/specularl2/specular/clients/geth/specular/rollup/rpc/eth"
)

type BridgeClient struct {
	*eth.EthClient
	*bindings.ISequencerInbox
	*bindings.IRollup
}

type L1Config interface {
	GetEndpoint() string
	GetSequencerInboxAddr() common.Address
	GetRollupAddr() common.Address
}

func NewBridgeClient(client *eth.EthClient, cfg L1Config) (*BridgeClient, error) {
	inbox, err := bindings.NewISequencerInbox(cfg.GetSequencerInboxAddr(), client)
	if err != nil {
		return nil, err
	}
	rollup, err := bindings.NewIRollup(cfg.GetRollupAddr(), client)
	if err != nil {
		return nil, err
	}
	return &BridgeClient{EthClient: client, ISequencerInbox: inbox, IRollup: rollup}, nil
}

func DialWithRetry(ctx context.Context, cfg L1Config) (*BridgeClient, error) {
	l1Client, err := eth.DialWithRetry(ctx, cfg.GetEndpoint(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to dial L1 client: %v", err)
	}
	return NewBridgeClient(l1Client, cfg)
}
