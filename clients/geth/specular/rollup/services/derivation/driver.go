package derivation

import (
	"context"
	"errors"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/specularl2/specular/clients/geth/specular/rollup/l2types"
	"github.com/specularl2/specular/clients/geth/specular/rollup/rpc/client"
	"github.com/specularl2/specular/clients/geth/specular/rollup/services/derivation/stage"
	"github.com/specularl2/specular/clients/geth/specular/rollup/utils/backoff"
	"github.com/specularl2/specular/clients/geth/specular/rollup/utils/fmt"
	"github.com/specularl2/specular/clients/geth/specular/rollup/utils/log"
)

type Driver[T any] struct {
	cfg           DriverConfig
	terminalStage stage.Stage[T]
	l1Client      EthClient
	backoffStrat  backoff.Strategy
}

type DriverConfig interface {
	NumAttempts() int
}

type EthClient interface {
	BlockNumber(ctx context.Context) (uint64, error)
	HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error)
	HeaderByTag(ctx context.Context, tag client.BlockTag) (*types.Header, error)
}

func (d *Driver[T]) Drive(ctx context.Context) error {
	var ticker = time.NewTicker(1)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			err := d.Step(ctx)
			if err != nil {
				return err
			}
		case <-ctx.Done():
			return nil
		}
	}
}

// TODO: consider async actions.
// Attempts a pipeline step (up to numAttempts).
func (d *Driver[T]) Step(ctx context.Context) error {
	numAttempts := d.cfg.NumAttempts()
	var i int
	for ; i < numAttempts; i++ {
		if isDone(ctx) {
			return nil
		}
		// Wait for a delay if this is a retry.
		if i > 0 {
			delay := d.backoffStrat.Duration(i)
			log.Info("Retrying after delay...", "num_attempt", i+1, "delay (s)", delay.Seconds())
			time.Sleep(delay)
		}
		// Try step.
		_, err := d.terminalStage.Step(ctx)
		// Process result.
		if err == nil {
			return nil
		} else if errors.As(err, &stage.RetryableError{}) {
			log.Warn("Retryable error encountered", "error", err)
			continue
		} else if errors.As(err, &stage.RecoverableError{}) {
			safeBlockID, err := d.getSafeL1BlockID(ctx)
			if err != nil {
				return fmt.Errorf("Driver could not fetch safe L1 block number to recover from re-org: %w", err)
			}
			err = d.terminalStage.Recover(ctx, safeBlockID)
			if err != nil {
				return fmt.Errorf("Driver could not recover from re-org: %w", err)
			}
		} else {
			return fmt.Errorf("Driver pipeline fatally failed: %w", err)
		}
	}
	if i == numAttempts {
		return fmt.Errorf("Driver pipeline failed after %d attempts", numAttempts)
	}
	return nil
}

// TODO: get (potentially unsafe) tip on safe chain

func (d *Driver[T]) getSafeL1BlockID(ctx context.Context) (l2types.BlockID, error) {
	header, err := d.l1Client.HeaderByTag(ctx, client.Safe)
	if err != nil {
		return l2types.BlockID{}, fmt.Errorf("Could not get safe L1 block header: %w", err)
	}
	return l2types.NewBlockIDFromHeader(header), nil
}

func isDone(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}
