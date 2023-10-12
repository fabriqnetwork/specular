package eth

import (
	"context"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
	"github.com/specularL2/specular/services/sidecar/utils/fmt"
	"github.com/specularL2/specular/services/sidecar/utils/log"
)

type ethClient interface {
	HeaderByTag(ctx context.Context, tag BlockTag) (*types.Header, error)
}

func SubscribeNewHeadByPolling(
	ctx context.Context,
	client ethClient,
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
			header, err := client.HeaderByTag(reqCtx, tag)
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

type LazyEthClient struct {
	*EthClient
	endpoint  string
	retryOpts []retry.Option
}

func NewLazilyDialedEthClient(endpoint string, retryOpts ...retry.Option) *LazyEthClient {
	return &LazyEthClient{endpoint: endpoint, retryOpts: retryOpts}
}

func (c *LazyEthClient) EnsureDialed(ctx context.Context) error {
	if c.EthClient != nil {
		return nil
	}
	client, err := DialWithRetry(ctx, c.endpoint, c.retryOpts...)
	if err != nil {
		return fmt.Errorf("failed to connect to node: %w", err)
	}
	c.EthClient = client
	return nil
}
