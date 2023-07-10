package driver

import (
	"context"
	"errors"
	"math"
	"time"

	"github.com/specularl2/specular/clients/geth/specular/rollup/services/api"
	"github.com/specularl2/specular/clients/geth/specular/rollup/services/derivation/stage"
	"github.com/specularl2/specular/clients/geth/specular/rollup/types"
	"github.com/specularl2/specular/clients/geth/specular/utils/backoff"
	"github.com/specularl2/specular/clients/geth/specular/utils/log"
	"github.com/specularl2/specular/clients/geth/specular/utils/sm"
)

type Driver struct {
	cfg               Config
	derivationStateFn *retryableDerivationStateFn
	sequencerStateFn  *retryableSequencerStateFn
}

// TODO: cleanup sm usage in driver.

// State fn which interrupts the caller via a signal to C, to call Do().
// Signals to C are driven internally.
// type InterruptingStateFn interface {
// 	C() <-chan time.Time
// 	Do(ctx context.Context) (sm.StateFn, error)
// }
// Interrupting state fn with externally polled planning.
// Signals to C should be driven via calls to Plan.
// type Sequencer interface {
// 	InterruptingStateFn
// 	Plan() time.Duration
// }

type TerminalStageOps interface {
	stage.StageOps[any]
	FindRecoveryPoint(ctx context.Context) (types.BlockID, error)
}

func NewDriver(
	cfg Config,
	terminalStage TerminalStageOps,
	sequencer Sequencer,
) *Driver {
	var d = Driver{cfg: cfg}
	d.derivationStateFn = NewRetryableDerivationStateFn(
		cfg,
		terminalStage,
		d.handleDerivationSuccess, // on success, continue
		d.waitForSignal,           // on retry, repeat the previous call (i.e. either pull or recovery)
		d.handleRecoverableError,  // on recoverable error, attempt recovery
	)
	if sequencer != nil {
		d.sequencerStateFn = NewRetryableSequencerStateFn(
			sequencer,
			d.waitForSignal,          // on success, continue
			d.handleRecoverableError, // on recoverable error, attempt recovery
		)
	}
	return &d
}

// Starts the driver's state machine in a goroutine.
func (d *Driver) Start(ctx context.Context, eg api.ErrGroup) error {
	eg.Go(func() error { return sm.Run(ctx, d.waitForSignal) })
	log.Info("Driver started.")
	return nil
}

// Driver state fn:
// 1. Schedule another derivation step (after the configured delay).
// 2. Wait for the next signal.
func (d *Driver) handleDerivationSuccess(ctx context.Context) (sm.StateFn, error) {
	log.Trace("Successful derivation step.")
	d.derivationStateFn.Reconfigure(sm.WithFn(d.derivationStateFn.pull), sm.WithDelay(d.cfg.GetStepInterval()))
	return d.waitForSignal(ctx)
}

// Driver state fn:
// 1. Schedule a recovery attempt.
// 2. Wait for the next signal.
func (d *Driver) handleRecoverableError(ctx context.Context) (sm.StateFn, error) {
	d.derivationStateFn.Reconfigure(sm.WithFn(d.derivationStateFn.recover))
	return d.waitForSignal(ctx)
}

// Driver state fn: wait for the next signal, to perform a derivation or sequencing step.
// Returns the next state fn to execute (derivation or sequencing).
func (d *Driver) waitForSignal(ctx context.Context) (sm.StateFn, error) {
	if d.sequencerStateFn != nil {
		// TODO: disable when drift too large
		// TODO: only when mismatched heads
		d.sequencerStateFn.Reconfigure(sm.WithDelay(d.sequencerStateFn.Plan()))
	}
	select {
	case <-d.derivationStateFn.C():
		return d.derivationStateFn.Do, nil
	case <-d.sequencerStateFn.C():
		return d.sequencerStateFn.Do, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

type retryableDerivationStateFn struct {
	*sm.RetryableStateFn
	terminalStage TerminalStageOps
}

// Retry retryable errors with exponential backoff up to cfg.NumAttempts() times.
// Yields control after success and between retries / recovery (expecting to get control back).
func NewRetryableDerivationStateFn(
	cfg Config,
	terminalStage TerminalStageOps,
	successRecvFn sm.StateFn,
	retryRecvFn sm.StateFn,
	recoveryRecvFn sm.StateFn,
) *retryableDerivationStateFn {
	var (
		s            = &retryableDerivationStateFn{terminalStage: terminalStage}
		backoffStrat = backoff.Exponential(float64(cfg.GetRetryDelay().Milliseconds()), math.MaxFloat64)
	)
	// initial state fn operation is s.pull -- this may change on error.
	s.RetryableStateFn = sm.NewRetryableStateFn(
		s.pull,
		sm.WithOnSuccessStateFn(successRecvFn),
		sm.WithOnErrorStateFnSelector(func(n uint, err error) sm.StateFn {
			if errors.As(err, &stage.RetryableError) {
				log.Warn("Retrying after delay...", "attempt#", n+1, "error", err)
				return retryRecvFn // retry same fn
			} else if errors.As(err, &stage.RecoverableError) {
				// Note: this error type is only expected to be returned by `s.pull`.
				log.Warn("Failed to advance derivation, attempting recovery.", "error", err)
				return recoveryRecvFn // try recovery
			}
			return nil // terminate for unknown errors
		}),
		sm.WithAttempts(cfg.GetNumAttempts()),
		sm.WithDelayFn(func(n uint) time.Duration { return backoffStrat.Duration(n) }),
	)
	return s
}

func (s *retryableDerivationStateFn) pull(ctx context.Context) error {
	_, err := s.terminalStage.Pull(ctx)
	return err
}

func (s *retryableDerivationStateFn) recover(ctx context.Context) error {
	recoveryPoint, err := s.terminalStage.FindRecoveryPoint(ctx)
	if err != nil {
		return err
	}
	return s.terminalStage.Recover(ctx, recoveryPoint)
}

type retryableSequencerStateFn struct {
	*sm.RetryableStateFn
	sequencer Sequencer
}

func NewRetryableSequencerStateFn(
	sequencer Sequencer,
	successRecvFn sm.StateFn,
	recoveryRecvFn sm.StateFn,
) *retryableSequencerStateFn {
	retryable := sm.NewRetryableStateFn(
		sequencer.Step,
		sm.WithOnSuccessStateFn(successRecvFn),
		sm.WithOnErrorStateFnSelector(func(n uint, err error) sm.StateFn {
			if errors.As(err, &stage.RecoverableError) {
				log.Warn("Failed to advance sequencer, attempting recovery.", "error", err)
				return recoveryRecvFn
			}
			return nil
		}),
	)
	return &retryableSequencerStateFn{retryable, sequencer}
}

func (s *retryableSequencerStateFn) Plan() time.Duration { return s.sequencer.Plan() }
