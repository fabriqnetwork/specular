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
	l2Client     L2Client
}

func NewSequencer(
	cfg Config,
	backend ExecutionBackend,
	l2Client L2Client,
	l1TxMgr TxManager,
	batchBuilder BatchBuilder,
) *Sequencer {
	return &Sequencer{
		executor:     &executor{cfg, backend},
		l1DAProvider: &batchDisseminator{cfg: cfg, batchBuilder: batchBuilder, l1TxMgr: l1TxMgr},
		l2Client:     l2Client,
	}
}

func (s *Sequencer) Start(ctx context.Context, eg api.ErrGroup) error {
	log.Info("Starting sequencer...")
	if err := s.l2Client.EnsureDialed(ctx); err != nil {
		return fmt.Errorf("Failed to create L2 client: %w", err)
	}
	eg.Go(func() error { return s.executor.start(ctx, s.l2Client) })
	eg.Go(func() error { return s.l1DAProvider.start(ctx, s.l2Client) })
	log.Info("Sequencer started")
	return nil
}
