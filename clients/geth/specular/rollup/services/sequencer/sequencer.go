package sequencer

import (
	"context"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/specularl2/specular/clients/geth/specular/rollup/utils/fmt"
	"github.com/specularl2/specular/clients/geth/specular/rollup/utils/log"
)

type Sequencer struct {
	BaseService
	executor     *executor
	l1DAProvider *batchDisseminator
	// Factory function for creating l2 client.
	// Necessary to create the L2Client *after* API access is made available.
	l2ClientCreatorFn l2ClientCreatorFn
}

type l2ClientCreatorFn func(ctx context.Context) (L2Client, error)

func NewSequencer(
	base BaseService,
	cfg Config,
	backend ExecutionBackend,
	l2ClientCreatorFn l2ClientCreatorFn,
	l1TxMgr TxManager,
	batchBuilder BatchBuilder,
) *Sequencer {
	return &Sequencer{
		BaseService:       base,
		executor:          &executor{cfg, backend, newOrdererByFee(backend)},
		l1DAProvider:      &batchDisseminator{cfg: cfg, batchBuilder: batchBuilder, l1TxMgr: l1TxMgr},
		l2ClientCreatorFn: l2ClientCreatorFn,
	}
}

func (s *Sequencer) Start() error {
	log.Info("Starting sequencer...")
	ctx := s.BaseService.Start()
	l2Client, err := s.l2ClientCreatorFn(ctx)
	if err != nil {
		return fmt.Errorf("Failed to create L2 client: %w", err)
	}
	s.ErrGroup().Go(func() error { return s.executor.start(ctx, l2Client) })
	s.ErrGroup().Go(func() error { return s.l1DAProvider.start(ctx, l2Client) })
	log.Info("Sequencer started")
	return nil
}

// TODO: sequencer APIs
func (s *Sequencer) APIs() []rpc.API { return []rpc.API{} }
