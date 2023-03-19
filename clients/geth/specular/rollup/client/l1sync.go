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
	// Watch L1 blockchain for confirmation period
	latestHeaderBroker := utils.NewBroker[*types.Header]()
	// This is the source subscription from which we broadcast to other subscribers.
	finalizedHeaderBroker := utils.NewBroker[*types.Header]()
	syncer := L1Syncer{
		LatestHeaderBroker:    latestHeaderBroker,
		FinalizedHeaderBroker: finalizedHeaderBroker,
	}
	syncer.wg.Add(2)
	latestSub := l1Client.SubscribeNewHeadByPolling(ctx, latestHeaderBroker.PubCh, Latest, 10*time.Second, 10*time.Second)
	go func() {
		defer syncer.wg.Done()
		err := latestHeaderBroker.Start(ctx, latestSub)
		if err != nil {
			log.Error("Failed running latest head broker", "err", err)
		}
	}()
	finalizedSub := l1Client.SubscribeNewHeadByPolling(ctx, latestHeaderBroker.PubCh, Finalized, 10*time.Second, 10*time.Second)
	go func() {
		defer syncer.wg.Done()
		err := finalizedHeaderBroker.Start(ctx, finalizedSub)
		if err != nil {
			log.Error("Failed running finalized head broker", "err", err)
		}
	}()
	return &syncer
}

func (s *L1Syncer) Start(ctx context.Context) {
	s.LatestHeaderBroker.SubscribeWithCallback(ctx, func(ctx context.Context, head *types.Header) error { return s.onLatest(head) })
	s.FinalizedHeaderBroker.SubscribeWithCallback(ctx, func(ctx context.Context, head *types.Header) error { return s.onFinalized(head) })
	// TODO: cleanup.
	for s.Latest == nil {
		log.Info("Waiting for L1 head...")
		log.Info("test \n test", "err", "gg no \n re")
		time.Sleep(4 * time.Second)
	}
	log.Info("Latest received", "number", s.Latest.Number)
	for s.Finalized == nil {
		log.Info("Waiting for L1 head...\n")
		time.Sleep(4 * time.Second)
	}
	log.Info("Finalized received", "number", s.Finalized.Number)
}

func (s *L1Syncer) onLatest(header *types.Header) error {
	s.Latest = header
	return nil
}

func (s *L1Syncer) onFinalized(header *types.Header) error {
	s.Finalized = header
	return nil
}
