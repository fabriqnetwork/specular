package disseminator

import (
	"bytes"
	"context"
	"errors"
	"io"
	"math/big"
	"time"

	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/specularl2/specular/clients/geth/specular/rollup/derivation"
	"github.com/specularl2/specular/clients/geth/specular/rollup/rpc/eth"
	"github.com/specularl2/specular/clients/geth/specular/rollup/services/api"
	"github.com/specularl2/specular/clients/geth/specular/rollup/types"
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
	return fmt.Sprintf("service entered unexpected state: %s", e.msg)
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
	d.rollback()
	var ticker = time.NewTicker(d.cfg.GetDisseminationInterval())
	defer ticker.Stop()
	d.step(ctx)
	for {
		select {
		case <-ticker.C:
			if err := d.step(ctx); err != nil {
				if errors.As(err, &unexpectedSystemStateError{}) {
					return fmt.Errorf("aborting: %w", err)
				}
				log.Errorf("Failed to step: %w", err)
			}
		case <-ctx.Done():
			log.Info("Aborting.")
			return nil
		}
	}
}

// Attempts to (incrementally) build a batch and disseminate it via L1.
func (d *BatchDisseminator) step(ctx context.Context) error {
	if err := d.appendToBuilder(ctx); err != nil {
		if errors.As(err, &L2ReorgDetectedError{}) {
			log.Error("Reorg detected, reverting to safe state.", "error", err)
			d.rollback()
		}
		return fmt.Errorf("failed to append to batch builder: %w", err)
	}
	if err := d.disseminateBatches(ctx); err != nil {
		return fmt.Errorf("failed to sequence batches: %w", err)
	}
	return nil
}

// Rolls back the disseminator state to the last safe L2 header.
func (d *BatchDisseminator) rollback() error {
	// TODO: use eth.Safe once Engine API is enabled.
	head, err := d.l2Client.HeaderByTag(context.Background(), eth.Latest)
	if err != nil {
		return fmt.Errorf("failed to get last finalized header: %w", err)
	}
	log.Info("Rolling back disseminator to checkpoint", "l2Block#", head.Number)
	d.batchBuilder.Reset(types.NewBlockIDFromHeader(head))
	return nil
}

// Appends L2 blocks to batch builder.
func (d *BatchDisseminator) appendToBuilder(ctx context.Context) error {
	start, end, err := d.pendingL2BlockRange(ctx)
	if err != nil {
		return fmt.Errorf("failed to get l2 block number: %w", err)
	}
	if start > end {
		log.Info("No pending blocks to append", "start", start, "end", end)
		return nil
	}
	log.Info("Appending blocks to builder", "start", start, "end", end)
	for i := start; i <= end; i++ {
		block, err := d.l2Client.BlockByNumber(ctx, big.NewInt(0).SetUint64(i))
		if err != nil {
			return fmt.Errorf("failed to get block: %w", err)
		}
		txs, err := encodeRLP(block.Transactions())
		if err != nil {
			return fmt.Errorf("failed to encode txs: %w", err)
		}
		dBlock := derivation.NewDerivationBlock(block.NumberU64(), block.Time(), txs)
		err = d.batchBuilder.Append(dBlock, types.NewBlockRefFromHeader(block.Header()))
		log.Info("Appended block to builder", "block", block.NumberU64(), "#txs", len(txs))
		if err != nil {
			if errors.As(err, &derivation.InvalidBlockError{}) {
				return L2ReorgDetectedError{err}
			}
			return fmt.Errorf("failed to append block (num=%d): %w", i, err)
		}
	}
	return nil
}

// Determines first and last unsafe block numbers.
// Typically, we start from the last appended block number + 1, and end at the current unsafe head.
func (d *BatchDisseminator) pendingL2BlockRange(ctx context.Context) (uint64, uint64, error) {
	var (
		lastAppended = d.batchBuilder.LastAppended()
		start        = lastAppended.GetNumber() + 1 // TODO: fix assumption
	)
	// TODO: uncomment the following cases after enabling Engine API.
	// safe, err := d.l2Client.HeaderByTag(ctx, eth.Safe)
	// if err != nil {
	// 	return 0, 0, fmt.Errorf("failed to get l2 safe header: %w", err)
	// }
	// log.Info("Retrieved safe head", "number", safe.Number, "hash", safe.Hash)
	// if lastAppended == types.EmptyBlockID {
	// 	// First time running; use safe (assumes local chain fork choice is in sync...)
	// 	start = safe.Number.Uint64() + 1
	// } else if safe.Number.Uint64() > lastAppended.GetNumber() {
	// 	// This should currently not be possible (single sequencer). TODO: handle restart case?
	// 	return 0, 0, &unexpectedSystemStateError{msg: "Safe header exceeds last appended header"}
	// }
	end, err := d.l2Client.BlockNumber(ctx)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get most recent l2 block number: %w", err)
	}
	return start, end, nil
}

// Disseminates batches until batch builder runs out (or signal from `ctx`).
func (d *BatchDisseminator) disseminateBatches(ctx context.Context) error {
	for {
		// Non-blocking ctx check.
		select {
		case <-ctx.Done():
			return nil
		default:
			if err := d.disseminateBatch(ctx); err != nil {
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
func (d *BatchDisseminator) disseminateBatch(ctx context.Context) error {
	// Construct tx data.
	batchAttrs, err := d.batchBuilder.Build()
	if err != nil {
		return fmt.Errorf("failed to build batch: %w", err)
	}
	receipt, err := d.l1TxMgr.AppendTxBatch(
		ctx,
		batchAttrs.Contexts(),
		batchAttrs.TxLengths(),
		batchAttrs.FirstL2BlockNumber(),
		batchAttrs.TxBatch(),
	)
	if err != nil {
		return fmt.Errorf("failed to send batch transaction: %w", err)
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
