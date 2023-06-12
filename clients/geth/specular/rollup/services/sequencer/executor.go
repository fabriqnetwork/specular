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
	var ticker = time.NewTicker(e.cfg.GetMaxExecutionInterval())
	for {
		select {
		case _ = <-ticker.C:
			// TODO: Begin block with a msg-passing tx.
			var payloadAttrs = engine.NewBuildPayloadAttributes(
				uint64(time.Now().Unix()),
				common.Hash{},
				e.cfg.GetAccountAddr(),
				nil,
				false,
			)
			log.Info("Building new payload...")
			if err := e.backend.BuildPayload(payloadAttrs); err != nil {
				return fmt.Errorf("Failed to commit txs: %w", err)
			}
		case <-ctx.Done():
			log.Info("Aborting.")
			return nil
		}
	}
}
