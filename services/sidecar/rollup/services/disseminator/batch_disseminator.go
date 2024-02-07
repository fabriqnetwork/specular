package disseminator

import (
	"context"
	"errors"
	"io"
	"math/big"
	"time"

	"github.com/specularL2/specular/services/sidecar/rollup/derivation"
	"github.com/specularL2/specular/services/sidecar/rollup/rpc/eth"
	"github.com/specularL2/specular/services/sidecar/rollup/types"
	"github.com/specularL2/specular/services/sidecar/utils/fmt"
	"github.com/specularL2/specular/services/sidecar/utils/log"
)

// Disseminates batches of L2 blocks via L1.
type BatchDisseminator struct {
	cfg          Config
	batchBuilder BatchBuilder
	l1TxMgr      TxManager
	l1State      *eth.EthState // Expected to generally be kept in sync with L1 chain.
	l2Client     L2Client
}

type recoverableSystemStateError struct{ msg string }

func (e recoverableSystemStateError) Error() string {
	return fmt.Sprintf("service entered unexpected state: %s", e.msg)
}

type L2ReorgDetectedError struct{ err error }

func (e L2ReorgDetectedError) Error() string { return e.err.Error() }

func NewBatchDisseminator(
	cfg Config,
	batchBuilder BatchBuilder,
	l1TxMgr TxManager,
	l1State *eth.EthState,
	l2Client L2Client,
) *BatchDisseminator {
	return &BatchDisseminator{cfg, batchBuilder, l1TxMgr, l1State, l2Client}
}

func (s *BatchDisseminator) Start(ctx context.Context, eg ErrGroup) error {
	log.Info("Starting disseminator...")
	if err := s.l2Client.EnsureDialed(ctx); err != nil {
		return fmt.Errorf("failed to create L2 client: %w", err)
	}
	eg.Go(func() error { return s.start(ctx) })
	log.Info("Disseminator started")
	return nil
}

func (d *BatchDisseminator) start(ctx context.Context) error {
	// Start with latest safe state.
	if err := d.rollback(ctx); err != nil {
		return err
	}
	var ticker = time.NewTicker(d.cfg.GetDisseminationInterval())
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if err := d.step(ctx); err != nil {
				log.Errorf("Failed to step: %w", err)
				if errors.As(err, &recoverableSystemStateError{}) {
					d.rollback(ctx)
					// return fmt.Errorf("aborting: %w", err)
				}
			}
		case <-ctx.Done():
			log.Info("Aborting.")
			return nil
		}
	}
}

// Attempts to (incrementally) build a batch and disseminate it via L1.
func (d *BatchDisseminator) step(ctx context.Context) error {
	start, end, safe, err := d.pendingL2BlockRange(ctx)
	if err != nil {
		return fmt.Errorf("failed to get l2 block number: %w", err)
	}
	if err := d.appendToBuilder(ctx, start, end); err != nil {
		if errors.As(err, &L2ReorgDetectedError{}) {
			log.Error("Reorg detected, reverting to safe state.", "error", err)
			if err := d.rollback(ctx); err != nil {
				return err
			}
		}
		return fmt.Errorf("failed to append to batch builder: %w", err)
	}
	if err := d.disseminateBatches(ctx, end-safe); err != nil {
		return fmt.Errorf("failed to sequence batches: %w", err)
	}
	return nil
}

// Rolls back the disseminator state to the last safe L2 header.
func (d *BatchDisseminator) rollback(ctx context.Context) error {
	head, err := d.l2Client.HeaderByTag(ctx, eth.Safe)
	if err != nil {
		return fmt.Errorf("failed to get last safe header: %w", err)
	}
	log.Info("Rolling back disseminator to checkpoint", "l2Block#", head.Number)
	d.batchBuilder.Reset(types.NewBlockIDFromHeader(head))
	return nil
}

// Appends L2 blocks to batch builder.
func (d *BatchDisseminator) appendToBuilder(ctx context.Context, start uint64, end uint64) error {
	if start > end {
		log.Info("No pending blocks to append", "start", start, "end", end)
		return nil
	}
	log.Info("Enqueuing blocks to builder", "start", start, "end", end)
	for i := start; i <= end; i++ {
		block, err := d.l2Client.BlockByNumber(ctx, big.NewInt(0).SetUint64(i))
		if err != nil {
			return fmt.Errorf("failed to get block: %w", err)
		}
		if err := d.batchBuilder.Enqueue(block); err != nil {
			if errors.As(err, &derivation.InvalidBlockError{}) {
				return L2ReorgDetectedError{err}
			}
			return fmt.Errorf("failed to enqueue block (num=%d): %w", i, err)
		}
		log.Info("Enqueued block at builder", "block#", block.NumberU64(), "#txs", len(block.Transactions()))
	}
	return nil
}

// Determines first and last unsafe block numbers.
// Typically, we start from the last appended block number + 1, and end at the current unsafe head.
func (d *BatchDisseminator) pendingL2BlockRange(ctx context.Context) (uint64, uint64, uint64, error) {
	var (
		lastEnqueued = d.batchBuilder.LastEnqueued()
		start        = lastEnqueued.GetNumber() + 1 // TODO: fix assumption
	)
	safe, err := d.l2Client.HeaderByTag(ctx, eth.Safe)
	safeBlockNum := safe.Number.Uint64()
	if err != nil {
		return 0, 0, 0, fmt.Errorf("failed to get l2 safe header: %w", err)
	}
	log.Info("Retrieved safe head", "number", safe.Number, "hash", safe.Hash)
	if lastEnqueued == types.EmptyBlockID {
		// First time running; use safe (assumes local chain fork-choice is in sync...)
		start = safeBlockNum + 1
	} else if safeBlockNum > lastEnqueued.GetNumber() {
		// This should currently not be possible (single sequencer).
		return 0, 0, 0, &recoverableSystemStateError{
			msg: fmt.Sprintf("safe header exceeds last appended header (safe=%d, last=%d)", safeBlockNum, lastEnqueued.GetNumber()),
		}
	}
	end, err := d.l2Client.BlockNumber(ctx)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("failed to get most recent l2 block number: %w", err)
	}
	return start, end, safeBlockNum, nil
}

// Disseminates batches until batch builder runs out (or signal from `ctx`).
func (d *BatchDisseminator) disseminateBatches(ctx context.Context, currentLag uint64) error {
	for {
		// Non-blocking ctx check.
		select {
		case <-ctx.Done():
			log.Info("Done disseminating batches")
			return nil
		default:
			if err := d.disseminateBatch(ctx, currentLag); err != nil {
				if errors.Is(err, io.EOF) {
					log.Info("No pending batches to sequence")
					return nil
				}
				return fmt.Errorf("failed to sequence batch: %w", err)
			}
		}
	}
}

// Fetches a batch from batch builder and disseminates it via L1.
// Blocking call until batch is sequenced and N confirmations received.
// Note: this does not guarantee safety (re-org resistance) but should make re-orgs less likely.
func (d *BatchDisseminator) disseminateBatch(ctx context.Context, currentLag uint64) error {
	// Construct tx data.
	data, err := d.batchBuilder.Build(d.l1State.Head(), currentLag)
	if err != nil {
		return fmt.Errorf("failed to build batch: %w", err)
	}
	receipt, err := d.l1TxMgr.AppendTxBatch(ctx, data)
	if err != nil {
		return fmt.Errorf("failed to send batch transaction: %w", err)
	}
	log.Info("Sequenced batch to L1", "size", len(data), "tx_hash", receipt.TxHash, "l1Block#", receipt.BlockNumber)
	d.batchBuilder.Advance()
	return nil
}
