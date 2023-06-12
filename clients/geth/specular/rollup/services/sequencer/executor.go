package sequencer

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/specularl2/specular/clients/geth/specular/rollup/types/engine"
	"github.com/specularl2/specular/clients/geth/specular/utils/fmt"
	"github.com/specularl2/specular/clients/geth/specular/utils/log"
)

// Responsible for executing transactions.
type executor struct {
	cfg     Config
	backend ExecutionBackend
}

// Responsible for ordering transactions (prior to their execution).
// TODO: Support:
// - PBS-style ordering: publicize mempool and call remote engine API.
// - remote ordering + weak DA in single call (systems conflate these roles -- e.g. Espresso)

// This goroutine fetches, orders and executes txs from the tx pool.
// Commits an empty block if no txs are received within an interval
// TODO: handle reorgs in the decentralized sequencer case.
func (e *executor) start(ctx context.Context, l2Client L2Client) error {
	var (
		minTicker        *time.Ticker
		minTickerCh      <-chan time.Time // a nil channel blocks.
		maxTicker        *time.Ticker
		maxTickerCh      <-chan time.Time // a nil channel blocks.
		minBatchInterval = e.cfg.GetMinExecutionInterval()
		maxBatchInterval = e.cfg.GetMaxExecutionInterval()
	)
	if minBatchInterval > 0 {
		minTicker = time.NewTicker(e.cfg.GetMinExecutionInterval())
		minTickerCh = minTicker.C
	}
	if maxBatchInterval > 0 {
		maxTicker = time.NewTicker(e.cfg.GetMaxExecutionInterval())
		maxTickerCh = maxTicker.C
	}
	for {
		select {
		case <-minTickerCh:
			status, err := l2Client.TxPoolStatus(ctx)
			if err != nil {
				return fmt.Errorf("Failed to fetch tx pool status: %w", err)
			}
			// Check if there are pending txs to build a payload from.
			numQueued, numPending := uint64(status["queued"]), uint64(status["pending"])
			log.Trace("Tx pool status.", "#queued", numQueued, "#pending", numPending)
			if numPending > 0 {
				if maxBatchInterval > 0 {
					maxTicker.Reset(maxBatchInterval)
				}
				if err := e.buildPayload(); err != nil {
					return err
				}
			} else {
				log.Info("Nothing to publish.")
			}
		case <-maxTickerCh:
			if minBatchInterval > 0 {
				minTicker.Reset(minBatchInterval)
			}
			// TODO: Begin block with a msg-passing tx.
			if err := e.buildPayload(); err != nil {
				return err
			}
		case <-ctx.Done():
			log.Info("Aborting.")
			return nil
		}
	}
}

func (e *executor) buildPayload() error {
	var payloadAttrs = engine.NewBuildPayloadAttributes(
		uint64(time.Now().Unix()),
		common.Hash{},
		e.cfg.GetAccountAddr(),
		nil,
		false,
	)
	if err := e.backend.BuildPayload(payloadAttrs); err != nil {
		return fmt.Errorf("Failed to commit txs: %w", err)
	}
	log.Info("Built payload.")
	return nil
}
