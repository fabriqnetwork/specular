package sequencer

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/specularl2/specular/clients/geth/specular/rollup/utils/fmt"
	"github.com/specularl2/specular/clients/geth/specular/rollup/utils/log"
)

type executor struct {
	cfg      SequencerServiceConfig
	backend  ExecutionBackend
	l2Client L2Client
}

// This goroutine fetches txs from txpool and executes them (immediately when received).
// Commits an empty block if no txs are received within an interval
// TODO: handle reorgs in the decentralized sequencer case.
// TODO: commit a msg-passing tx in empty block.
func (e *executor) start(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	// Watch transactions in TxPool
	txsCh := make(chan core.NewTxsEvent, 4096)
	txsSub := e.backend.SubscribeNewTxsEvent(txsCh)
	defer txsSub.Unsubscribe()
	ticker := time.NewTicker(e.cfg.Sequencer().ExecutionInterval)
	for {
		select {
		case _ = <-ticker.C:
			if err := e.backend.CommitTransactions([]*types.Transaction{}); err != nil {
				log.Crit("Failed to commit empty block", "err", err)
			}
			log.Info("No new txs; committed empty block")
		case ev := <-txsCh:
			log.Info("Received txsCh event", "txs", len(ev.Txs))
			sanitizedTxs, err := e.sanitize(ctx, e.backend.Prepare(ev.Txs))
			if err != nil {
				log.Crit("Failed to sanitize txs", "err", err)
			}
			err = e.backend.CommitTransactions(sanitizedTxs)
			if err != nil {
				log.Crit("Failed to commit txsCh event ", "err", err)
			}
			log.Info("Committed txs", "num_txs", len(sanitizedTxs))
			ticker.Reset(e.cfg.Sequencer().ExecutionInterval)
		case <-ctx.Done():
			log.Info("Aborting.")
			return
		}
	}
}

func (e *executor) sanitize(
	ctx context.Context,
	sortedTxs *types.TransactionsByPriceAndNonce,
) ([]*types.Transaction, error) {
	var sanitizedTxs []*types.Transaction
	for {
		tx := sortedTxs.Peek()
		if tx == nil {
			break
		}
		err := e.validateTx(ctx, tx)
		if errors.Is(err, &txValidationError{}) {
			log.Warn("Dropping tx", "tx", tx.Hash(), "err", err)
			sortedTxs.Pop()
			continue
		} else if err != nil {
			return nil, fmt.Errorf("Sanitization failed: %w", err)
		}
		sanitizedTxs = append(sanitizedTxs, tx)
		sortedTxs.Pop()
	}
	return sanitizedTxs, nil
}

func (e *executor) validateTx(ctx context.Context, tx *types.Transaction) error {
	// Check if tx exists on the L2 chain (TODO: is this really necessary)
	prevTx, _, err := e.l2Client.TransactionByHash(ctx, tx.Hash())
	if err != nil {
		return fmt.Errorf("Failed to query for tx by hash: %w", err)
	}
	if prevTx != nil {
		return &txValidationError{"tx already exists on-chain"}
	}
	return nil
}
