package eth

import (
	"context"
	"math/big"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/specularl2/specular/clients/geth/specular/rollup/utils/fmt"
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
		log.Error("Failed attempt", "attempt", n, "err", err)
	}),
}

type EthClient struct {
	*ethclient.Client
	C *rpc.Client
}

func NewEthClient(c *rpc.Client) *EthClient { return &EthClient{ethclient.NewClient(c), c} }

func DialWithRetry(ctx context.Context, endpoint string, retryOpts []retry.Option) (*EthClient, error) {
	if retryOpts == nil {
		retryOpts = DefaultRetryOpts
	}
	retryOpts = append(retryOpts, retry.Context(ctx))
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

func (c *EthClient) SubscribeNewHeadByPolling(
	ctx context.Context,
	headCh chan<- *types.Header,
	tag BlockTag,
	interval time.Duration,
	requestTimeout time.Duration,
) event.Subscription {
	return event.NewSubscription(func(unsub <-chan struct{}) error {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		poll := func() error {
			reqCtx, cancel := context.WithTimeout(ctx, requestTimeout)
			header, err := c.HeaderByTag(reqCtx, tag)
			cancel()
			if err != nil {
				log.Warn("Failed to poll for latest L1 block header", "err", err)
				return err
			}
			headCh <- header
			return nil
		}
		poll()
		for {
			select {
			case <-ticker.C:
				poll()
			case <-ctx.Done():
				return ctx.Err()
			case <-unsub:
				return nil
			}
		}
	})
}

func (c *EthClient) ResubscribeErrNewHead(ctx context.Context, headCh chan<- *types.Header) event.Subscription {
	sub := event.ResubscribeErr(
		time.Second*10,
		func(ctx context.Context, err error) (event.Subscription, error) {
			if err != nil {
				log.Warn("Error in NewHead subscription, resubscribing", "err", err)
			}
			return c.SubscribeNewHead(ctx, headCh)
		},
	)
	return sub
}
