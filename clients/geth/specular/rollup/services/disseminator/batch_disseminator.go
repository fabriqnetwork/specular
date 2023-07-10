package disseminator

import (
	"bytes"
	"context"
	"errors"
	"io"
	"math/big"
	"time"

	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/specularl2/specular/clients/geth/specular/rollup/rpc/eth"
	"github.com/specularl2/specular/clients/geth/specular/rollup/services/api"
	"github.com/specularl2/specular/clients/geth/specular/rollup/types"
	"github.com/specularl2/specular/clients/geth/specular/rollup/types/da"
	"github.com/specularl2/specular/clients/geth/specular/utils/fmt"
	"github.com/specularl2/specular/clients/geth/specular/utils/log"
)

// Disseminates batches of L2 blocks via L1.
type BatchDisseminator struct {
	cfg          Config
	batchBuilder BatchBuilder
	l1TxMgr      TxManager
	l2Client     L2Client
}

type unexpectedSystemStateError struct{ msg string }

func (e unexpectedSystemStateError) Error() string {
	return fmt.Sprintf("service in unexpected state: %s", e.msg)
}

type L2ReorgDetectedError struct{ err error }

func (e L2ReorgDetectedError) Error() string { return e.err.Error() }

func NewBatchDisseminator(
	cfg Config,
	batchBuilder BatchBuilder,
	l1TxMgr TxManager,
	l2Client L2Client,
) *BatchDisseminator {
	return &BatchDisseminator{cfg: cfg, batchBuilder: batchBuilder, l1TxMgr: l1TxMgr, l2Client: l2Client}
}

func (s *BatchDisseminator) Start(ctx context.Context, eg api.ErrGroup) error {
	log.Info("Starting batch disseminator...")
	if err := s.l2Client.EnsureDialed(ctx); err != nil {
		return fmt.Errorf("Failed to create L2 client: %w", err)
	}
	eg.Go(func() error { return s.start(ctx) })
	log.Info("Sequencer started")
	return nil
}

func (d *BatchDisseminator) start(ctx context.Context) error {
	// Start with latest safe state.
	d.revertToFinalized()
	var ticker = time.NewTicker(d.cfg.GetDisseminationInterval())
	defer ticker.Stop()
	d.step(ctx)
	for {
		select {
		case <-ticker.C:
			if err := d.step(ctx); err != nil {
				if errors.As(err, &unexpectedSystemStateError{}) {
					return fmt.Errorf("Aborting: %w", err)
				}
				log.Errorf("Failed to step: %w", err)
			}
		case <-ctx.Done():
			log.Info("Aborting.")
			return nil
		}
	}
}

// TODO: document
func (d *BatchDisseminator) step(ctx context.Context) error {
	if err := d.appendToBuilder(ctx); err != nil {
		if errors.As(err, &L2ReorgDetectedError{}) {
			log.Error("Reorg detected, reverting to safe state.", "error", err)
			d.revertToFinalized()
		}
		return fmt.Errorf("Failed to append to batch builder: %w", err)
	}
	if err := d.sequenceBatches(ctx); err != nil {
		return fmt.Errorf("Failed to sequence batches: %w", err)
	}
	return nil
}

// TODO: document
func (d *BatchDisseminator) revertToFinalized() error {
	finalizedHeader, err := d.l2Client.HeaderByTag(context.Background(), eth.Finalized)
	if err != nil {
		return fmt.Errorf("Failed to get last finalized header: %w", err)
	}
	d.batchBuilder.Reset(types.NewBlockIDFromHeader(finalizedHeader))
	return nil
}

// Appends blocks to batch builder.
func (d *BatchDisseminator) appendToBuilder(ctx context.Context) error {
	start, end, err := d.pendingL2BlockRange(ctx)
	if err != nil {
		return fmt.Errorf("Failed to get l2 block number: %w", err)
	}
	if start > end {
		log.Info("No pending blocks to append", "start", start, "end", end)
		return nil
	}
	log.Info("Appending blocks to builder", "start", start, "end", end)
	for i := start; i <= end; i++ {
		block, err := d.l2Client.BlockByNumber(ctx, big.NewInt(0).SetUint64(i))
		if err != nil {
			return fmt.Errorf("Failed to get block: %w", err)
		}
		txs, err := encodeRLP(block.Transactions())
		if err != nil {
			return fmt.Errorf("Failed to encode txs: %w", err)
		}
		dBlock := da.NewDerivationBlock(block.NumberU64(), block.Time(), txs)
		err = d.batchBuilder.Append(dBlock, types.NewBlockRefFromHeader(block.Header()))
		log.Info("Appended block to builder", "block", block.NumberU64(), "#txs", len(txs))
		if err != nil {
			if errors.As(err, &da.InvalidBlockError{}) {
				return L2ReorgDetectedError{err}
			}
			return fmt.Errorf("Failed to append block (num=%d): %w", i, err)
		}
	}
	return nil
}

// Determines first and last unsafe block numbers.
func (d *BatchDisseminator) pendingL2BlockRange(ctx context.Context) (uint64, uint64, error) {
	var (
		lastAppended = d.batchBuilder.LastAppended()
		start        uint64
	)
	safe, err := d.l2Client.HeaderByTag(ctx, eth.Safe)
	if err != nil {
		return 0, 0, fmt.Errorf("Failed to get l2 safe header: %w", err)
	}
	log.Info("Retrieved safe head", "number", safe.Number, "hash", safe.Hash)
	if lastAppended == types.EmptyBlockID {
		// First time running; use safe (assumes local chain fork choice is in sync...)
		start = safe.Number.Uint64() + 1
	} else if safe.Number.Uint64() > lastAppended.GetNumber() {
		// This should currently not be possible (single sequencer). TODO: handle restart case?
		return 0, 0, &unexpectedSystemStateError{msg: "Safe header exceeds last appended header"}
	} else {
		// Normal case.
		start = lastAppended.GetNumber() + 1 // TODO: fix assumption
	}
	end, err := d.l2Client.BlockNumber(ctx)
	if err != nil {
		return 0, 0, fmt.Errorf("Failed to get most recent l2 block number: %w", err)
	}
	return start, end, nil
}

// Sequences batches until batch builder runs out (or signal from `ctx`).
func (d *BatchDisseminator) sequenceBatches(ctx context.Context) error {
	for {
		// Non-blocking ctx check.
		select {
		case <-ctx.Done():
			return nil
		default:
			if err := d.sequenceBatch(ctx); err != nil {
				if errors.Is(err, io.EOF) {
					log.Info("No pending batches to sequence")
					return nil
				}
				return fmt.Errorf("Failed to sequence batch: %w", err)
			}
		}
	}
}

// Fetches a batch from batch builder and sequences to L1.
// Blocking call until batch is sequenced and N confirmations received.
// Note: this does not guarantee safety (re-org resistance) but should make re-orgs less likely.
func (d *BatchDisseminator) sequenceBatch(ctx context.Context) error {
	// Construct tx data.
	batchAttrs, err := d.batchBuilder.Build()
	if err != nil {
		return fmt.Errorf("Failed to build batch: %w", err)
	}
	receipt, err := d.l1TxMgr.AppendTxBatch(
		ctx,
		batchAttrs.Contexts(),
		batchAttrs.TxLengths(),
		batchAttrs.FirstL2BlockNumber(),
		batchAttrs.TxBatch(),
	)
	if err != nil {
		return fmt.Errorf("Failed to send batch transaction: %w", err)
	}
	log.Info("Sequenced batch to L1", "tx_hash", receipt.TxHash, "l1Block#", receipt.BlockNumber)
	d.batchBuilder.Advance()
	return nil
}

func encodeRLP(txs ethTypes.Transactions) ([][]byte, error) {
	var encodedTxs [][]byte
	for _, tx := range txs {
		var txBuf bytes.Buffer
		if err := tx.EncodeRLP(&txBuf); err != nil {
			return nil, err
		}
		encodedTxs = append(encodedTxs, txBuf.Bytes())
	}
	return encodedTxs, nil
}
