package client

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

const DefaultRetries = 3

type EthClient struct {
	*ethclient.Client
	C *rpc.Client
}

func NewEthClient(c *rpc.Client) *EthClient {
	return &EthClient{ethclient.NewClient(c), c}
}

func DialWithRetry(ctx context.Context, endpoint string, numAttempts uint) (*EthClient, error) {
	var client *EthClient
	var err error
	retryOpts := []retry.Option{
		retry.Context(ctx),
		retry.Attempts(numAttempts),
		retry.Delay(5 * time.Second),
		retry.LastErrorOnly(true),
		retry.OnRetry(func(n uint, err error) {
			log.Error("Failed to connect to node", "endpoint", endpoint, "attempt", n, "err", err)
		}),
	}
	err = retry.Do(func() error {
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
		lastHeaderNumber := big.NewInt(-1)
		poll := func() {
			reqCtx, cancel := context.WithTimeout(ctx, requestTimeout)
			header, err := c.HeaderByTag(reqCtx, tag)
			cancel()
			if err != nil {
				log.Warn("Failed to poll for latest L1 block header", "err", err)
				return
			}
			if header.Number.Cmp(lastHeaderNumber) <= 0 {
				log.Warn("Polled header is not new", "number", header.Number, "newest", lastHeaderNumber)
				return
			}
			headCh <- header
			lastHeaderNumber = header.Number
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
