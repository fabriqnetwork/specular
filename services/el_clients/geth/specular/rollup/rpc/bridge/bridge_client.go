package bridge

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/specularl2/specular/clients/geth/specular/bindings"
	"github.com/specularl2/specular/clients/geth/specular/rollup/client"
)

type BridgeClient struct {
	*client.EthClient
	*bindings.ISequencerInbox
	*bindings.IRollup
}

type L1Config interface {
	Endpoint() string
	SequencerInboxAddr() common.Address
	RollupAddr() common.Address
}

func NewBridgeClient(client *client.EthClient, cfg L1Config) (*BridgeClient, error) {
	inbox, err := bindings.NewISequencerInbox(cfg.SequencerInboxAddr(), client)
	if err != nil {
		return nil, err
	}
	rollup, err := bindings.NewIRollup(cfg.RollupAddr(), client)
	if err != nil {
		return nil, err
	}
	return &BridgeClient{EthClient: client, ISequencerInbox: inbox, IRollup: rollup}, nil
}

func DialWithRetry(ctx context.Context, cfg L1Config) (*BridgeClient, error) {
	l1Client, err := client.DialWithRetry(ctx, cfg.Endpoint(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to dial L1 client: %v", err)
	}
	return NewBridgeClient(l1Client, cfg)
}
