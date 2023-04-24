package client

import (
	"context"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
	"github.com/specularl2/specular/clients/geth/specular/rollup/utils"
	"github.com/specularl2/specular/clients/geth/specular/rollup/utils/log"
)

type EthSyncer struct {
	OnNewHandler
	LatestHeaderBroker    *utils.Broker[*types.Header]
	FinalizedHeaderBroker *utils.Broker[*types.Header]
	wg                    sync.WaitGroup
}

type OnNewHandler interface {
	OnLatest(ctx context.Context, header *types.Header) error
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
		FinalizedHeaderBroker: utils.NewBroker[*types.Header](),
	}
}

func (s *EthSyncer) Start(ctx context.Context, client EthPollingClient) {
	if s.LatestHeaderBroker != nil {
		s.subscribeNewHead(ctx, client, s.LatestHeaderBroker, Latest)
		s.LatestHeaderBroker.SubscribeWithCallback(ctx, s.OnLatest)
	}
	if s.FinalizedHeaderBroker != nil {
		s.subscribeNewHead(ctx, client, s.FinalizedHeaderBroker, Finalized)
		s.FinalizedHeaderBroker.SubscribeWithCallback(ctx, s.OnFinalized)
	}
}

func (s *EthSyncer) Stop(ctx context.Context) {
	s.LatestHeaderBroker.Stop()
	s.FinalizedHeaderBroker.Stop()
	s.wg.Wait()
}

// Starts polling for new headers and publishes them to the broker.
func (s *EthSyncer) subscribeNewHead(
	ctx context.Context,
	client EthPollingClient,
	broker *utils.Broker[*types.Header],
	tag BlockTag,
) {
	s.wg.Add(1)
	sub := client.SubscribeNewHeadByPolling(ctx, broker.PubCh(), tag, 10*time.Second, 10*time.Second)
	go func() {
		defer s.wg.Done()
		err := broker.Start(ctx, sub)
		if err != nil {
			log.Error("Failed running header broker", "tag", tag, "err", err)
		}
	}()
}
