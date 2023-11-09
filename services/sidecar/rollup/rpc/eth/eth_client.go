package eth

import (
	"context"
	"math/big"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/specularL2/specular/services/sidecar/utils/fmt"
	"github.com/specularL2/specular/services/sidecar/utils/log"
)

const (
	defaultRetryAttempts = 3
	defaultRetryDelay    = 5 * time.Second
)

var DefaultRetryOpts = []retry.Option{
	retry.Attempts(defaultRetryAttempts),
	retry.Delay(defaultRetryDelay),
	retry.LastErrorOnly(true),
	retry.OnRetry(func(n uint, err error) {
		log.Errorf("Failed attempt: %w", err, "attempt", n)
	}),
}

type EthClient struct {
	*ethclient.Client
	C *rpc.Client
}

func NewEthClient(c *rpc.Client) *EthClient { return &EthClient{ethclient.NewClient(c), c} }

func DialWithRetry(ctx context.Context, endpoint string, retryOpts ...retry.Option) (*EthClient, error) {
	if retryOpts == nil {
		retryOpts = DefaultRetryOpts
	}
	retryOpts = append(retryOpts, retry.Context(ctx))
	var client *EthClient
	err := retry.Do(func() error {
		log.Info("Dialing...", "endpoint", endpoint)
		rpcClient, err := rpc.DialContext(ctx, endpoint)
		client = NewEthClient(rpcClient)
		return err
	}, retryOpts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to node: %w", err)
	}
	return client, nil
}

func (c *EthClient) HeaderByTag(ctx context.Context, tag BlockTag) (*types.Header, error) {
	var header *types.Header
	err := c.C.CallContext(ctx, &header, "eth_getBlockByNumber", tag, false)
	if err == nil && header == nil {
		err = ethereum.NotFound
	}
	return header, err
}

func (c *EthClient) SuggestGasTipCap(ctx context.Context) (*big.Int, error) {
	gasTipCap, err := c.Client.SuggestGasTipCap(ctx)
	if err != nil {
		// This is a workaround for Hardhat and any backend that doesn't support eth_maxPriorityFeePerGas.
		log.Warn("Failed to get gas tip cap by eth_maxPriorityFeePerGas", "err", err)
		return c.SuggestGasPrice(ctx)
	}
	return gasTipCap, nil
}

func (c *EthClient) TxPoolStatus(ctx context.Context) (map[string]hexutil.Uint, error) {
	var status map[string]hexutil.Uint
	err := c.C.CallContext(ctx, &status, "txpool_status")
	return status, err
}
