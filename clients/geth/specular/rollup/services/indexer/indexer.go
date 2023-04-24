package indexer

import (
	"github.com/ethereum/go-ethereum/log"
	"github.com/specularl2/specular/clients/geth/specular/rollup/services"
)

type Indexer struct {
	*services.BaseService
	cfg    IndexerServiceConfig
	syncer *services.Syncer
}

type IndexerServiceConfig interface{ L1() *services.L1Config }

func NewIndexer(cfg IndexerServiceConfig, syncer *services.Syncer) *Indexer {
	return &Indexer{BaseService: &services.BaseService{}, syncer: syncer}
}

func (i *Indexer) Start() error {
	log.Info("Starting indexer...")
	ctx, err := i.BaseService.Start()
	if err != nil {
		return err
	}
	i.Wg.Add(1)
	go func() {
		defer i.Wg.Done()
		i.syncer.SyncLoop(ctx, i.cfg.L1().L1RollupGenesisBlock, nil)
	}()
	log.Info("Indexer started")
	return nil
}
