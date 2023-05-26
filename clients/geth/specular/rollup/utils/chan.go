package utils

import (
	"context"
	"time"

	"github.com/specularl2/specular/clients/geth/specular/rollup/utils/log"
)

// Creates and publishes batched events to a channel from `inCh` (one-to-many).
// Publishes batch every `minBatchInterval` if batch is non-empty; otherwise every `maxBatchInterval`.
// Setting either to 0 disables batching at that interval.
func SubscribeBatched[T any](
	ctx context.Context,
	inCh chan T,
	minBatchInterval time.Duration,
	maxBatchInterval time.Duration,
) chan []T {
	var (
		minTicker   *time.Ticker
		minTickerCh <-chan time.Time // a nil channel blocks.
		maxTicker   *time.Ticker
		maxTickerCh <-chan time.Time // a nil channel blocks.
	)
	if minBatchInterval > 0 {
		minTicker = time.NewTicker(minBatchInterval)
		minTickerCh = minTicker.C
	}
	if maxBatchInterval > 0 {
		maxTicker = time.NewTicker(maxBatchInterval)
		maxTickerCh = maxTicker.C
	}
	var (
		outCh = make(chan []T)
		batch []T
	)
	go func() {
		defer close(outCh)
		for {
			select {
			case <-minTickerCh:
				// Only publish batch if it's non-empty.
				if len(batch) > 0 {
					outCh <- batch
					batch = []T{}
					if maxBatchInterval > 0 {
						maxTicker.Reset(maxBatchInterval)
					}
				}
			case <-maxTickerCh:
				// Publish batch even if it's empty.
				outCh <- batch
				batch = []T{}
				if minBatchInterval > 0 {
					minTicker.Reset(minBatchInterval)
				}
			case ev := <-inCh:
				if minBatchInterval == 0 {
					outCh <- []T{ev}
				} else {
					batch = append(batch, ev)
				}
			case <-ctx.Done():
				if len(batch) > 0 {
					log.Info("Dropping batch", "size", len(batch))
				}
				return
			}
		}
	}()
	return outCh
}

// Creates and publishes events to a channel mapped from `inCh` (one-to-one).
func SubscribeMapped[T any, U any](
	ctx context.Context,
	inCh chan T,
	mapFn func(context.Context, T) (U, error),
) <-chan U {
	outCh := make(chan U)
	go func() { defer close(outCh); mapCh(ctx, inCh, outCh, mapFn) }()
	return outCh
}

// Creates and publishes events to a channel mapped from `inCh` (one-to-one).
func SubscribeMappedToMany[T any, U any](
	ctx context.Context,
	inCh chan T,
	mapFn func(context.Context, T) ([]U, error),
) <-chan U {
	outCh := make(chan U)
	go func() { defer close(outCh); mapChToMany(ctx, inCh, outCh, mapFn) }()
	return outCh
}

// Maps events from `inCh` to `outCh` (one-to-one).
func mapCh[T any, U any](
	ctx context.Context,
	inCh chan T,
	outCh chan U,
	mapFn func(context.Context, T) (U, error),
) {
	applyCh(
		ctx,
		inCh,
		func(ctx context.Context, ev T) (err error) {
			out, err := mapFn(ctx, ev)
			if err == nil {
				outCh <- out
			}
			return
		},
	)
}

// Maps events from `inCh` to `outCh` (one-to-many).
func mapChToMany[T any, U any](
	ctx context.Context,
	inCh chan T,
	outCh chan U,
	mapFn func(context.Context, T) ([]U, error),
) {
	applyCh(
		ctx,
		inCh,
		func(ctx context.Context, ev T) (err error) {
			out, err := mapFn(ctx, ev)
			if err == nil {
				for _, ev := range out {
					outCh <- ev
				}
			}
			return
		},
	)
}

// Applies a function to each event in `inCh`.
func applyCh[T any](
	ctx context.Context,
	inCh chan T,
	callbackFn func(context.Context, T) error,
) {
	for {
		select {
		case head := <-inCh:
			err := callbackFn(ctx, head)
			if err != nil {
				log.Error("Failed to map", "error", err)
				return
			}
		case <-ctx.Done():
			log.Info("Aborting.")
			return
		}
	}
}
