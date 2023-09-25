package indexer

import (
	"context"

	"github.com/ethereum/go-ethereum/log"
	"github.com/specularl2/specular/clients/geth/specular/proof"
	"github.com/specularl2/specular/clients/geth/specular/rollup/client"
	"github.com/specularl2/specular/clients/geth/specular/rollup/services"
	"github.com/specularl2/specular/clients/geth/specular/rollup/services/api"
)

// TODO: delete.
type Indexer struct{ *services.BaseService }

func New(eth api.ExecutionBackend, proofBackend proof.Backend, l1Client client.L1BridgeClient, cfg services.BaseConfig) (*Indexer, error) {
	base, err := services.NewBaseService(eth, proofBackend, l1Client, cfg)
	if err != nil {
		return nil, err
	}
	return &Indexer{BaseService: base}, nil
}

func (i *Indexer) Start(ctx context.Context, eg api.ErrGroup) error {
	log.Info("Starting indexer...")
	err := i.BaseService.Start(ctx, eg)
	if err != nil {
		return err
	}
	i.Wg.Add(1)
	go i.SyncLoop(ctx, i.Config.GetRollupGenesisBlock(), nil)
	log.Info("Indexer started")
	return nil
}
