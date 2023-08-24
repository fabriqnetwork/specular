package client

import (
	"context"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/specularl2/specular/clients/geth/specular/rollup/utils"
	"github.com/specularl2/specular/clients/geth/specular/rollup/utils/log"
)

type L1Syncer struct {
	Latest    *types.Header
	Finalized *types.Header

	LatestHeaderBroker    *utils.Broker[*types.Header]
	FinalizedHeaderBroker *utils.Broker[*types.Header]

	wg sync.WaitGroup
}

func NewL1Syncer(ctx context.Context, l1Client L1BridgeClient) *L1Syncer {
	// This is the source subscription from which we broadcast to other subscribers.
	syncer := L1Syncer{
		LatestHeaderBroker:    utils.NewBroker[*types.Header](),
		FinalizedHeaderBroker: utils.NewBroker[*types.Header](),
	}
	syncer.wg.Add(2)
	latestSub := l1Client.SubscribeNewHeadByPolling(ctx, syncer.LatestHeaderBroker.PubCh, Latest, 10*time.Second, 10*time.Second)
	go func() {
		defer syncer.wg.Done()
		err := syncer.LatestHeaderBroker.Start(ctx, latestSub)
		if err != nil {
			log.Error("Failed running latest head broker", "err", err)
		}
	}()
	finalizedSub := l1Client.SubscribeNewHeadByPolling(ctx, syncer.FinalizedHeaderBroker.PubCh, Finalized, 10*time.Second, 10*time.Second)
	go func() {
		defer syncer.wg.Done()
		err := syncer.FinalizedHeaderBroker.Start(ctx, finalizedSub)
		if err != nil {
			log.Error("Failed running finalized head broker", "err", err)
		}
	}()
	return &syncer
}

func (s *L1Syncer) Start(ctx context.Context) {
	s.LatestHeaderBroker.SubscribeWithCallback(ctx, s.onLatest)
	s.FinalizedHeaderBroker.SubscribeWithCallback(ctx, s.onFinalized)
	for s.Latest == nil {
		log.Info("Waiting for L1 latest header...")
		time.Sleep(2 * time.Second)
	}
	log.Info("Latest header received", "number", s.Latest.Number)
	// TODO: Wait for finalized head? Might cause issues in tests.
}

func (s *L1Syncer) onLatest(ctx context.Context, header *types.Header) error {
	s.Latest = header
	return nil
}

func (s *L1Syncer) onFinalized(ctx context.Context, header *types.Header) error {
	s.Finalized = header
	return nil
}
