package derivation

import (
	"context"
	"errors"
	"math"
	"math/big"
	"time"

	"github.com/avast/retry-go/v4"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/specularl2/specular/clients/geth/specular/rollup/rpc/eth"
	"github.com/specularl2/specular/clients/geth/specular/rollup/services/api"
	"github.com/specularl2/specular/clients/geth/specular/rollup/services/derivation/stage"
	"github.com/specularl2/specular/clients/geth/specular/rollup/types"
	"github.com/specularl2/specular/clients/geth/specular/utils/backoff"
	"github.com/specularl2/specular/clients/geth/specular/utils/fmt"
	"github.com/specularl2/specular/clients/geth/specular/utils/log"
)

type Driver struct {
	cfg           Config
	terminalStage TerminalStageOps
	retryOpts     []retry.Option
}

type Config interface {
	GetStepInterval() time.Duration
	GetRetryDelay() time.Duration
	GetNumAttempts() uint
}

type TerminalStageOps interface {
	stage.StageOps[any]
	FindRecoveryPoint(ctx context.Context) (types.BlockID, error)
}

type EthClient interface {
	BlockNumber(ctx context.Context) (uint64, error)
	HeaderByNumber(ctx context.Context, number *big.Int) (*ethTypes.Header, error)
	HeaderByTag(ctx context.Context, tag eth.BlockTag) (*ethTypes.Header, error)
}

// Driver state machine states.
const (
	driverStateHealthy   = iota // Driver is making progress.
	driverStateUnhealthy        // Driver is recovering from a failure.
)

func NewDriver(cfg Config, terminalStage TerminalStageOps) *Driver {
	var (
		backoffStrat = backoff.Exponential(float64(cfg.GetRetryDelay().Milliseconds()), math.MaxFloat64)
		// Retry RetryableErrors with exponential backoff up to cfg.NumAttempts() times.
		retryOpts = []retry.Option{
			retry.OnRetry(func(n uint, err error) { log.Warn("Retrying after delay...", "attempt#", n+1, "error", err) }),
			retry.RetryIf(func(err error) bool { return errors.As(err, &stage.RetryableError{}) }),
			retry.Attempts(cfg.GetNumAttempts()),
			retry.Delay(cfg.GetRetryDelay() * time.Second),
			retry.DelayType(func(n uint, _ error, _ *retry.Config) time.Duration { return backoffStrat.Duration(n) }),
		}
	)
	return &Driver{cfg: cfg, terminalStage: terminalStage, retryOpts: retryOpts}
}

func (d *Driver) Start(ctx context.Context, eg api.ErrGroup) error {
	eg.Go(func() error {
		log.Crit("Driver failed.", "error", d.drive(ctx)) // TODO: get rid of crit
		return d.drive(ctx)
	})
	log.Info("Driver started.")
	return nil
}

func (d *Driver) drive(ctx context.Context) error {
	d.retryOpts = append(d.retryOpts, retry.Context(ctx))
	var (
		driverState = driverStateHealthy
		pullFn      = func() error { return d.pull(ctx) }
		recoverFn   = func() error { return d.recover(ctx) }
		ticker      = time.NewTicker(d.cfg.GetStepInterval())
	)
	defer ticker.Stop()
	// TODO: consider async control w/ channels.
	for {
		select {
		case <-ticker.C:
			// Determine step type.
			var stepFn func() error
			switch driverState {
			case driverStateHealthy:
				stepFn = pullFn
			case driverStateUnhealthy:
				stepFn = recoverFn
			}
			// Perform step.
			err := retry.Do(stepFn, d.retryOpts...)
			// Process result.
			if err == nil {
				// Success (possibly after multiple retries).
				driverState = driverStateHealthy
				log.Info("Successful step.")
			} else if errors.As(err, &stage.RecoverableError{}) {
				// Note: this error type is only expected to be returned by `Pull`.
				log.Warn("Failed to advance, attempting recovery.", "error", err)
				driverState = driverStateUnhealthy
			} else {
				// Unrecoverable error or all attempts failed.
				return fmt.Errorf("Failed to advance unrecoverably: %w", err)
			}
		case <-ctx.Done():
			return nil
		}
	}
}

func (d *Driver) pull(ctx context.Context) error {
	_, err := d.terminalStage.Pull(ctx)
	return err
}

func (d *Driver) recover(ctx context.Context) error {
	recoveryPoint, err := d.terminalStage.FindRecoveryPoint(ctx)
	if err != nil {
		return err
	}
	return d.terminalStage.Recover(ctx, recoveryPoint)
}
