package client

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
	"github.com/specularl2/specular/clients/geth/specular/rollup/utils"
	"golang.org/x/sync/errgroup"
)

type EthSyncer struct {
	OnNewHandler
	LatestHeaderBroker    *utils.Broker[*types.Header]
	SafeHeaderBroker      *utils.Broker[*types.Header]
	FinalizedHeaderBroker *utils.Broker[*types.Header]
	eg                    errgroup.Group
}

type OnNewHandler interface {
	OnLatest(ctx context.Context, header *types.Header) error
	OnSafe(ctx context.Context, header *types.Header) error
	OnFinalized(ctx context.Context, header *types.Header) error
}

type EthPollingClient interface {
	SubscribeNewHeadByPolling(
		ctx context.Context,
		headCh chan<- *types.Header,
		tag BlockTag,
		interval time.Duration,
		requestTimeout time.Duration,
	) event.Subscription
}

func NewEthSyncer(handler OnNewHandler) *EthSyncer {
	return &EthSyncer{
		OnNewHandler:          handler,
		LatestHeaderBroker:    utils.NewBroker[*types.Header](),
		SafeHeaderBroker:      utils.NewBroker[*types.Header](),
		FinalizedHeaderBroker: utils.NewBroker[*types.Header](),
	}
}

func (s *EthSyncer) Start(ctx context.Context, client EthPollingClient) {
	s.subscribeNewHead(ctx, client, s.LatestHeaderBroker, Latest)
	s.LatestHeaderBroker.SubscribeWithCallback(ctx, s.OnLatest)
	s.subscribeNewHead(ctx, client, s.SafeHeaderBroker, Safe)
	s.SafeHeaderBroker.SubscribeWithCallback(ctx, s.OnSafe)
	s.subscribeNewHead(ctx, client, s.FinalizedHeaderBroker, Finalized)
	s.FinalizedHeaderBroker.SubscribeWithCallback(ctx, s.OnFinalized)
}

func (s *EthSyncer) Stop(ctx context.Context) {
	s.LatestHeaderBroker.Stop()
	s.FinalizedHeaderBroker.Stop()
	s.eg.Wait()
}

// Starts polling for new headers and publishes them to the broker.
func (s *EthSyncer) subscribeNewHead(
	ctx context.Context,
	client EthPollingClient,
	broker *utils.Broker[*types.Header],
	tag BlockTag,
) {
	sub := client.SubscribeNewHeadByPolling(ctx, broker.PubCh(), tag, 10*time.Second, 10*time.Second)
	s.eg.Go(func() error { return broker.Start(ctx, sub) })
}
