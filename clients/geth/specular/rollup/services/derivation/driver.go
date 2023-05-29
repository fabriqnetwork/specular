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
	"github.com/specularl2/specular/clients/geth/specular/rollup/services/derivation/stage"
	"github.com/specularl2/specular/clients/geth/specular/rollup/types"
	"github.com/specularl2/specular/clients/geth/specular/rollup/utils/backoff"
	"github.com/specularl2/specular/clients/geth/specular/rollup/utils/fmt"
	"github.com/specularl2/specular/clients/geth/specular/rollup/utils/log"
	"golang.org/x/sync/errgroup"
)

type Driver struct {
	BaseService   // TODO: remove?
	cfg           Config
	terminalStage TerminalStageOps
	retryOpts     []retry.Option
}

type BaseService interface {
	Start() context.Context
	Stop() error
	ErrGroup() *errgroup.Group
}

type Config interface {
	StepInterval() time.Duration
	RetryDelay() time.Duration
	NumAttempts() uint
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

func NewDriver(
	base BaseService,
	cfg Config,
	terminalStage TerminalStageOps,
) *Driver {
	var (
		backoffStrat = backoff.Exponential(float64(cfg.RetryDelay().Milliseconds()), math.MaxFloat64)
		retryOpts    = []retry.Option{
			retry.OnRetry(func(n uint, err error) { log.Warn("Retrying...", "attempt#", n+1, "error", err) }),
			retry.RetryIf(func(err error) bool { return errors.As(err, &stage.RetryableError{}) }),
			retry.Attempts(cfg.NumAttempts()),
			retry.Delay(cfg.RetryDelay() * time.Second),
			retry.DelayType(func(n uint, _ error, _ *retry.Config) time.Duration { return backoffStrat.Duration(n) }),
		}
	)
	return &Driver{BaseService: base, cfg: cfg, terminalStage: terminalStage, retryOpts: retryOpts}
}

func (d *Driver) Start() error {
	ctx := d.BaseService.Start()
	d.ErrGroup().Go(func() error { return d.drive(ctx) })
	log.Info("Driver started.")
	return nil
}

func (d *Driver) drive(ctx context.Context) error {
	d.retryOpts = append(d.retryOpts, retry.Context(ctx))
	var (
		driverState = driverStateHealthy
		pullFn      = func() error { return d.pull(ctx) }
		recoverFn   = func() error { return d.recover(ctx) }
		ticker      = time.NewTicker(d.cfg.StepInterval())
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
