package indexer

import (
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/specularl2/specular/clients/geth/specular/prover"
	"github.com/specularl2/specular/clients/geth/specular/rollup/client"
	"github.com/specularl2/specular/clients/geth/specular/rollup/services"
)

type Indexer struct{ *services.BaseService }

func New(eth services.Backend, proofBackend prover.L2ELClientBackend, l1Client client.L1BridgeClient, cfg *services.Config) (*Indexer, error) {
	base, err := services.NewBaseService(eth, proofBackend, l1Client, cfg)
	if err != nil {
		return nil, err
	}
	return &Indexer{BaseService: base}, nil
}

func (i *Indexer) Start() error {
	log.Info("Starting indexer...")
	ctx, err := i.BaseService.Start()
	if err != nil {
		return err
	}
	i.Wg.Add(1)
	go i.SyncLoop(ctx, i.Config.L1RollupGenesisBlock, nil)
	log.Info("Indexer started")
	return nil
}

func (i *Indexer) APIs() []rpc.API { return []rpc.API{} }
