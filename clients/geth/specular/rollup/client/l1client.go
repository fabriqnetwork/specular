package client

import (
	"context"

	"github.com/avast/retry-go/v4"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/specularl2/specular/clients/geth/specular/bindings"
	"github.com/specularl2/specular/clients/geth/specular/rollup/rpc/eth"
)

// TODO: delete this file
type L1BridgeClient interface {
	TransactionByHash(ctx context.Context, hash common.Hash) (*types.Transaction, bool, error)
	HeaderByTag(ctx context.Context, tag eth.BlockTag) (*types.Header, error)
	BlockNumber(ctx context.Context) (uint64, error)
	FilterTxBatchAppendedEvents(opts *bind.FilterOpts) (*bindings.ISequencerInboxTxBatchAppendedIterator, error)
}

type EthBridgeClient struct {
	*eth.EthClient
	inbox *bindings.ISequencerInboxSession
}

func NewEthBridgeClient(
	ctx context.Context,
	l1Endpoint string,
	sequencerInboxAddress common.Address,
	retryOpts []retry.Option,
) (*EthBridgeClient, error) {
	client, err := eth.DialWithRetry(ctx, l1Endpoint, retryOpts...)
	if err != nil {
		return nil, err
	}
	inbox, err := bindings.NewISequencerInbox(sequencerInboxAddress, client)
	if err != nil {
		return nil, err
	}
	callOpts := bind.CallOpts{Pending: true, Context: ctx}
	inboxSession := &bindings.ISequencerInboxSession{Contract: inbox, CallOpts: callOpts}
	return &EthBridgeClient{EthClient: client, inbox: inboxSession}, nil
}

func (c *EthBridgeClient) FilterTxBatchAppendedEvents(
	opts *bind.FilterOpts,
) (*bindings.ISequencerInboxTxBatchAppendedIterator, error) {
	return c.inbox.Contract.FilterTxBatchAppended(opts)
}
