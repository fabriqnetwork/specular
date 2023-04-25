package sequencer

import (
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/specularl2/specular/clients/geth/specular/rollup/services"
	"github.com/specularl2/specular/clients/geth/specular/rollup/utils/fmt"
	"github.com/specularl2/specular/clients/geth/specular/rollup/utils/log"
)

type Sequencer struct {
	*services.BaseService
	executor     *executor
	l1DAProvider *batchDisseminator
}

type unexpectedStateError struct{ msg string }

func (e *unexpectedStateError) Error() string {
	return fmt.Sprintf("service in unexpected state: %s", e.msg)
}

type txValidationError struct{ msg string }

func (e *txValidationError) Error() string {
	return fmt.Sprintf("Tx validation failed: %s", e.msg)
}

func NewSequencer(
	cfg SequencerServiceConfig,
	backend ExecutionBackend,
	l2Client L2Client,
	l1TxMgr TxManager,
	batchBuilder BatchBuilder,
) *Sequencer {
	return &Sequencer{
		BaseService:  &services.BaseService{},
		executor:     &executor{cfg, backend, &ordererByFee{backend, l2Client}, l2Client},
		l1DAProvider: &batchDisseminator{cfg: cfg, batchBuilder: batchBuilder, l2Client: l2Client, l1TxMgr: l1TxMgr},
	}
}

func (s *Sequencer) Start() error {
	log.Info("Starting sequencer...")
	ctx, err := s.BaseService.Start()
	if err != nil {
		return fmt.Errorf("Failed to start base service: %w", err)
	}
	// We assume a single sequencer (us) for now, so we don't
	// need to sync sequenced transactions up.
	s.Wg.Add(2)
	go s.executor.start(ctx, &s.Wg)
	go s.l1DAProvider.start(ctx, &s.Wg)
	log.Info("Sequencer started")
	return nil
}

func (s *Sequencer) APIs() []rpc.API {
	// TODO: sequencer APIs
	return []rpc.API{}
}
