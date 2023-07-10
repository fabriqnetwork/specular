package sm

import (
	"context"
	"math"
	"time"
)

type BasicStateFn struct {
	fn                     retryableFn
	onSuccessStateFn       StateFn
	onErrorStateFnSelector onErrorStateFnSelector
}

type RetryableStateFn struct {
	BasicStateFn
	delayFn     func(n uint) time.Duration
	maxAttempts uint
	t           *time.Timer
	numAttempts uint
}

type Option func(*RetryableStateFn)
type retryableFn = func(context.Context) error
type onErrorStateFnSelector = func(n uint, err error) StateFn

// By default,
// On success: terminal state.
// On error: retry all errors. Calls fn.
// Delay: no delay.
func NewRetryableStateFn(fn retryableFn, options ...Option) *RetryableStateFn {
	r := &RetryableStateFn{
		BasicStateFn: BasicStateFn{
			fn:                     fn,
			onSuccessStateFn:       nil, // terminal state by default
			onErrorStateFnSelector: nil, // set below (retry all errors)
		},
		delayFn:     func(n uint) time.Duration { return 0 }, // no delay
		maxAttempts: math.MaxUint,                            // retry forever
		t:           time.NewTimer(0),                        // signal to start immediately
	}
	// By default: retry all errors by iterating on the signal.
	r.onErrorStateFnSelector = func(n uint, err error) StateFn { <-r.C(); return r.Do }
	// Override defaults with options.
	r.Reconfigure(options...)
	return r
}

// Signals when Do should be called again, if at all.
func (r *RetryableStateFn) C() <-chan time.Time { return r.t.C }

func (r *RetryableStateFn) Do(ctx context.Context) (StateFn, error) {
	if err := r.fn(ctx); err != nil {
		r.numAttempts++
		var onErrorStateFn = r.onErrorStateFnSelector(r.numAttempts, err)
		if onErrorStateFn == nil {
			return nil, err // terminal state
		}
		if r.numAttempts < r.maxAttempts {
			r.t.Reset(r.delayFn(r.numAttempts))
		}
		return onErrorStateFn, nil
	}
	r.numAttempts = 0
	return r.onSuccessStateFn, nil
}

func (r *RetryableStateFn) Reconfigure(options ...Option) {
	for _, o := range options {
		o(r)
	}
}

// Configurable options.

func WithFn(fn retryableFn) Option {
	return func(r *RetryableStateFn) { r.fn = fn }
}

func WithDelay(delay time.Duration) Option {
	return func(r *RetryableStateFn) {
		// Drain the channel if necessary.
		if !r.t.Stop() {
			<-r.t.C
		}
		r.t.Reset(delay)
	}
}

func WithDelayFn(delayFn func(n uint) time.Duration) Option {
	return func(r *RetryableStateFn) { r.delayFn = delayFn }
}

func WithOnErrorStateFnSelector(onErrStateFnSelector onErrorStateFnSelector) Option {
	return func(r *RetryableStateFn) { r.onErrorStateFnSelector = onErrStateFnSelector }
}

func WithOnSuccessStateFn(onSuccessStateFn StateFn) Option {
	return func(r *RetryableStateFn) { r.onSuccessStateFn = onSuccessStateFn }
}

func WithAttempts(maxAttempts uint) Option {
	return func(r *RetryableStateFn) { r.maxAttempts = maxAttempts }
}
