package indexer

import (
	bind "github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/specularl2/specular/clients/geth/specular/proof"
	"github.com/specularl2/specular/clients/geth/specular/rollup/services"
)

func RegisterService(stack *node.Node, eth services.Backend, proofBackend proof.Backend, cfg *services.Config, auth *bind.TransactOpts) {
	indexer, err := New(eth, proofBackend, cfg, auth)
	if err != nil {
		log.Crit("Failed to register the Rollup service", "err", err)
	}
	stack.RegisterLifecycle(indexer)
	// stack.RegisterAPIs(indexer.APIs())
	log.Info("Indexer registered")
}

type Indexer struct {
	*services.BaseService

	newBatchCh chan struct{}
}

func New(eth services.Backend, proofBackend proof.Backend, cfg *services.Config, auth *bind.TransactOpts) (*Indexer, error) {
	base, err := services.NewBaseService(eth, proofBackend, cfg, auth)
	if err != nil {
		return nil, err
	}
	s := &Indexer{
		BaseService: base,
		newBatchCh:  make(chan struct{}, 1),
	}
	return s, nil
}

func (i *Indexer) newBatchConsumeLoop() {
	defer i.Wg.Done()
	for {
		select {
		case <-i.newBatchCh:
			continue
		case <-i.Ctx.Done():
			return
		}
	}
}

func (i *Indexer) Start() error {
	i.BaseService.Start(false)

	i.Wg.Add(2)
	go i.newBatchConsumeLoop()
	go i.SyncLoop(i.newBatchCh)
	log.Info("Indexer started")
	return nil
}

func (i *Indexer) Stop() error {
	log.Info("Indexer stopped")
	i.Cancel()
	i.Wg.Wait()
	return nil
}

func (i *Indexer) APIs() []rpc.API {
	// TODO: sequencer APIs
	return []rpc.API{}
}
