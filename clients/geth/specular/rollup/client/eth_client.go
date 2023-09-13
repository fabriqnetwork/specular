package client

import (
	"context"

	"github.com/avast/retry-go/v4"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/specularl2/specular/clients/geth/specular/rollup/utils/fmt"
)

type BlockTag string

// L1 block tag values
// https://ethereum.github.io/execution-apis/api-documentation/
const (
	// Lowest numbered block client has available
	Earliest = "earliest"
	// The most recent crypto-economically secure block,
	// cannot be re-orged outside of manual intervention driven by community coordination
	Finalized = "finalized"
	// The most recent block that is safe from re-orgs under honest majority and certain synchronicity assumptions
	Safe = "safe"
	// The most recent block in the canonical chain observed by the client,
	// this block may be re-orged out of the canonical chain even under healthy/normal conditions
	Latest = "latest"
	// A sample next block built by the client on top of `latest` and
	// containing the set of transactions usually taken from local mempool.
	Pending = "pending"
)

type EthClient struct {
	*ethclient.Client
	C *rpc.Client
}

func NewEthClient(c *rpc.Client) *EthClient {
	return &EthClient{ethclient.NewClient(c), c}
}

func (c *EthClient) HeaderByTag(ctx context.Context, tag BlockTag) (*types.Header, error) {
	var header *types.Header
	err := c.C.CallContext(ctx, &header, "eth_getBlockByNumber", tag, false)
	if err == nil && header == nil {
		err = ethereum.NotFound
	}
	return header, err
}

func DialWithRetry(ctx context.Context, endpoint string, retryOpts []retry.Option) (*EthClient, error) {
	var client *EthClient
	err := retry.Do(func() error {
		rpcClient, err := rpc.DialContext(ctx, endpoint)
		client = NewEthClient(rpcClient)
		return err
	}, retryOpts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to node, err: %w", err)
	}
	return client, nil
}
