package sequencer

import (
	"context"
	"errors"
	"io"
	"math/big"
	"time"

	"github.com/specularl2/specular/clients/geth/specular/rollup/geth"
	"github.com/specularl2/specular/clients/geth/specular/rollup/rpc/eth"
	"github.com/specularl2/specular/clients/geth/specular/rollup/types"
	"github.com/specularl2/specular/clients/geth/specular/rollup/utils"
	"github.com/specularl2/specular/clients/geth/specular/rollup/utils/fmt"
	"github.com/specularl2/specular/clients/geth/specular/rollup/utils/log"
)

// Disseminates batches of L2 blocks via L1.
type batchDisseminator struct {
	cfg          Config
	batchBuilder BatchBuilder
	l1TxMgr      TxManager
	l2Client     L2Client
}

type unexpectedSystemStateError struct{ msg string }

func (e *unexpectedSystemStateError) Error() string {
	return fmt.Sprintf("service in unexpected state: %s", e.msg)
}

func newBatchDisseminator(
	cfg Config,
	batchBuilder BatchBuilder,
	l1TxMgr TxManager,
) *batchDisseminator {
	return &batchDisseminator{cfg: cfg, batchBuilder: batchBuilder, l1TxMgr: l1TxMgr}
}

func (d *batchDisseminator) start(ctx context.Context, l2Client L2Client) error {
	d.l2Client = l2Client
	// Start with latest safe state.
	d.revertToSafe()
	var ticker = time.NewTicker(d.cfg.SequencingInterval())
	defer ticker.Stop()
	d.step(ctx)
	for {
		select {
		case <-ticker.C:
			d.step(ctx)
		case <-ctx.Done():
			log.Info("Aborting.")
			return nil
		}
	}
}

func (d *batchDisseminator) step(ctx context.Context) {
	if err := d.appendToBuilder(ctx); err != nil {
		log.Error("Failed to append to batch builder", "error", err)
		if errors.As(err, &utils.L2ReorgDetectedError{}) {
			log.Error("Reorg detected, reverting to safe state.")
			d.revertToSafe()
			return
		}
	}
	if err := d.sequenceBatches(ctx); err != nil {
		log.Error("Failed to sequence batch", "error", err)
	}
}

func (d *batchDisseminator) revertToSafe() error {
	safeHeader, err := d.l2Client.HeaderByTag(context.Background(), eth.Safe)
	if err != nil {
		return fmt.Errorf("Failed to get safe header: %w", err)
	}
	d.batchBuilder.Reset(types.NewBlockIDFromHeader(safeHeader))
	return nil
}

// Appends blocks to batch builder.
func (d *batchDisseminator) appendToBuilder(ctx context.Context) error {
	start, end, err := d.unsafeBlockRange(ctx)
	if err != nil {
		return fmt.Errorf("Failed to get l2 block number: %w", err)
	}
	// TODO: this check might not be necessary
	if start > end {
		return &utils.L2ReorgDetectedError{Msg: "Last appended l2 block number exceeds l2 chain head"}
	}
	for i := start; i <= end; i++ {
		block, err := d.l2Client.BlockByNumber(ctx, big.NewInt(0).SetUint64(i))
		if err != nil {
			return fmt.Errorf("Failed to get block: %w", err)
		}
		txs, err := geth.EncodeRLP(block.Transactions())
		if err != nil {
			return fmt.Errorf("Failed to encode txs: %w", err)
		}
		dBlock := types.NewDerivationBlock(block.NumberU64(), block.Time(), txs)
		err = d.batchBuilder.Append(dBlock, types.NewBlockRefFromHeader(block.Header()))
		if err != nil {
			return fmt.Errorf("Failed to append block (num=%s): %w", i, err)
		}
	}
	return nil
}

// Determines first and last unsafe block numbers.
func (d *batchDisseminator) unsafeBlockRange(ctx context.Context) (uint64, uint64, error) {
	lastAppended := d.batchBuilder.LastAppended()
	start := lastAppended.GetNumber() + 1 // TODO: fix assumption
	safe, err := d.l2Client.HeaderByTag(ctx, eth.Safe)
	if err != nil {
		return 0, 0, fmt.Errorf("Failed to get l2 safe header: %w", err)
	}
	if safe.Number.Uint64() > lastAppended.GetNumber() {
		// This should currently not be possible. TODO: handle restart case?
		return 0, 0, &unexpectedSystemStateError{msg: "Safe header exceeds last appended header"}
	}
	// else if safe.Number.Uint64() < lastAppended {
	// 	// Hasn't caught up yet... OR re-org.
	// 	// TODO: allow continue to appending blocks.
	// 	return 0, 0, &utils.L2ReorgDetectedError{Msg: "Safe header is behind last appended header"}
	// }
	end, err := d.l2Client.BlockNumber(ctx)
	if err != nil {
		return 0, 0, fmt.Errorf("Failed to get most recent l2 block number: %w", err)
	}
	return start, end, nil
}

// Sequences batches until batch builder runs out (or signal from `ctx`).
func (d *batchDisseminator) sequenceBatches(ctx context.Context) error {
	for {
		// Non-blocking ctx check.
		select {
		case <-ctx.Done():
			return nil
		default:
			if err := d.sequenceBatch(ctx); err != nil {
				if errors.Is(err, io.EOF) {
					log.Info("No more batches to sequence")
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
func (d *batchDisseminator) sequenceBatch(ctx context.Context) error {
	// Construct tx data.
	batchAttrs, err := d.batchBuilder.Build()
	if err != nil {
		return fmt.Errorf("Failed to build batch: %w", err)
	}
	receipt, err := d.l1TxMgr.AppendTxBatch(ctx, batchAttrs.contexts, batchAttrs.txLengths, batchAttrs.firstL2BlockNumber, batchAttrs.txs)
	if err != nil {
		return fmt.Errorf("Failed to send batch transaction: %w", err)
	}
	log.Info("Sequenced batch successfully", "tx hash", receipt.TxHash)
	d.batchBuilder.Advance()
	return nil
}
