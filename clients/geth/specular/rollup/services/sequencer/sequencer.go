package sequencer

import (
	"context"

	"github.com/specularl2/specular/clients/geth/specular/rollup/services/api"
	"github.com/specularl2/specular/clients/geth/specular/utils/fmt"
	"github.com/specularl2/specular/clients/geth/specular/utils/log"
)

// Responsible for:
// - Receiving transactions from users.
// - Ordering transactions.
// - Executing transactions.
// - Disseminating L2 blocks.
type Sequencer struct {
	executor     *executor
	l1DAProvider *batchDisseminator
	// Factory function for creating l2 client.
	// Necessary to create the L2Client *after* API access is made available.
	l2ClientCreatorFn l2ClientCreatorFn
}

type l2ClientCreatorFn func(ctx context.Context) (L2Client, error)

func NewSequencer(
	cfg Config,
	backend ExecutionBackend,
	l2ClientCreatorFn l2ClientCreatorFn,
	l1TxMgr TxManager,
	batchBuilder BatchBuilder,
) *Sequencer {
	return &Sequencer{
		executor:          &executor{cfg, backend, newOrdererByFee(backend)},
		l1DAProvider:      &batchDisseminator{cfg: cfg, batchBuilder: batchBuilder, l1TxMgr: l1TxMgr},
		l2ClientCreatorFn: l2ClientCreatorFn,
	}
}

func (s *Sequencer) Start(ctx context.Context, eg api.ErrGroup) error {
	log.Info("Starting sequencer...")
	l2Client, err := s.l2ClientCreatorFn(ctx)
	if err != nil {
		return fmt.Errorf("Failed to create L2 client: %w", err)
	}
	eg.Go(func() error { return s.executor.start(ctx, l2Client) })
	eg.Go(func() error { return s.l1DAProvider.start(ctx, l2Client) })
	log.Info("Sequencer started")
	return nil
}
