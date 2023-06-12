package eth

import (
	"context"
	"errors"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/beacon"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/log"
	"github.com/specularl2/specular/clients/geth/specular/utils/fmt"
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

type ForkChoiceState = beacon.ForkchoiceStateV1

func GetForkChoice(ctx context.Context, client ethClient) (ForkChoiceState, error) {
	// Get latest head
	l2Head, err := client.HeaderByTag(ctx, Latest)
	if err != nil {
		return ForkChoiceState{}, fmt.Errorf("Failed to get latest L2 block header: %w", err)
	}
	// Get safe head
	var l2SafeHash common.Hash
	l2Safe, err := client.HeaderByTag(ctx, Safe)
	if err == nil {
		l2SafeHash = l2Safe.Hash()
	} else if !errors.Is(err, ethereum.NotFound) {
		return ForkChoiceState{}, fmt.Errorf("Failed to get safe L2 block header: %w", err)
	}
	// Get finalized head
	var l2FinalizedHash common.Hash
	l2Finalized, err := client.HeaderByTag(ctx, Finalized)
	if err == nil {
		l2FinalizedHash = l2Finalized.Hash()
	} else if !errors.Is(err, ethereum.NotFound) {
		return ForkChoiceState{}, fmt.Errorf("Failed to get finalized L2 block header: %w", err)
	}
	return ForkChoiceState{
		HeadBlockHash:      l2Head.Hash(),
		SafeBlockHash:      l2SafeHash,
		FinalizedBlockHash: l2FinalizedHash,
	}, nil
}

type LazyEthClient struct {
	*EthClient
	endpoint  string
	retryOpts []retry.Option
}

func NewLazyDialedEthClient(endpoint string, retryOpts []retry.Option) *LazyEthClient {
	return &LazyEthClient{endpoint: endpoint, retryOpts: retryOpts}
}

func (c *LazyEthClient) EnsureDialed(ctx context.Context) error {
	if c.EthClient != nil {
		return nil
	}
	client, err := DialWithRetry(ctx, c.endpoint, c.retryOpts)
	if err != nil {
		return fmt.Errorf("failed to connect to node, err: %w", err)
	}
	c.EthClient = client
	return nil
}
