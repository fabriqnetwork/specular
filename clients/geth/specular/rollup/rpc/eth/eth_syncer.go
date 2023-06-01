package eth

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
	"github.com/specularl2/specular/clients/geth/specular/utils"
	"golang.org/x/sync/errgroup"
)

// TODO: move to config
const (
	EthSlotInterval  = 12 * time.Second
	EthEpochInterval = 6*time.Minute + 24*time.Second
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

// Starts a subscription in a separate goroutine for each commitment level.
func (s *EthSyncer) Start(ctx context.Context, client EthPollingClient) {
	s.subscribeNewHead(ctx, client, Latest, s.LatestHeaderBroker, s.OnLatest, EthSlotInterval)
	s.subscribeNewHead(ctx, client, Safe, s.SafeHeaderBroker, s.OnSafe, EthEpochInterval)
	s.subscribeNewHead(ctx, client, Finalized, s.FinalizedHeaderBroker, s.OnFinalized, EthEpochInterval)
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
	tag BlockTag,
	broker *utils.Broker[*types.Header],
	fn func(context.Context, *types.Header) error,
	pollInterval time.Duration,
) {
	sub := client.SubscribeNewHeadByPolling(ctx, broker.PubCh(), tag, pollInterval, 10*time.Second)
	s.eg.Go(func() error { return broker.Start(ctx, sub) })
	broker.SubscribeWithCallback(ctx, fn)
}
