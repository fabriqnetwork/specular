package utils

import (
	"context"
	"time"

	"github.com/specularl2/specular/clients/geth/specular/rollup/utils/log"
)

// Creates and publishes batched events to a channel from `inCh` (one-to-many).
// Publishes batch every `minBatchInterval` if batch is non-empty; otherwise every `maxBatchInterval`.
func SubscribeBatched[T any](
	ctx context.Context,
	inCh chan T,
	minBatchInterval time.Duration,
	maxBatchInterval time.Duration,
) chan []T {
	outCh := make(chan []T)
	var minTicker *time.Ticker // a nil channel blocks.
	if minBatchInterval > 0 {
		minTicker = time.NewTicker(minBatchInterval)
	}
	maxTicker := time.NewTicker(maxBatchInterval)
	var batch []T
	go func() {
		defer close(outCh)
		for {
			select {
			case <-minTicker.C:
				// Only publish batch if it's non-empty.
				if len(batch) > 0 {
					outCh <- batch
					batch = []T{}
					maxTicker.Reset(maxBatchInterval)
				}
			case <-maxTicker.C:
				// Publish batch even if it's empty.
				outCh <- batch
				batch = []T{}
				minTicker.Reset(minBatchInterval)
			case ev := <-inCh:
				batch = append(batch, ev)
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
	go func() {
		defer close(outCh)
		mapCh(ctx, inCh, outCh, mapFn)
	}()
	return outCh
}

// Creates and publishes events to a channel mapped from `inCh` (one-to-one).
func SubscribeMappedToMany[T any, U any](
	ctx context.Context,
	inCh chan T,
	mapFn func(context.Context, T) ([]U, error),
) <-chan U {
	outCh := make(chan U)
	go func() {
		defer close(outCh)
		mapChToMany(ctx, inCh, outCh, mapFn)
	}()
	return outCh
}

// Maps events from `inCh` to `outCh` (one-to-one).
func mapCh[T any, U any](
	ctx context.Context,
	inCh chan T,
	outCh chan U,
	mapFn func(context.Context, T) (U, error),
) {
	for {
		select {
		case ev := <-inCh:
			out, err := mapFn(ctx, ev)
			if err != nil {
				log.Error("Failed to map, err: %w", err)
				return
			}
			outCh <- out
		case <-ctx.Done():
			log.Info("Aborting.")
			return
		}
	}
}

// Maps events from `inCh` to `outCh` (one-to-many).
func mapChToMany[T any, U any](
	ctx context.Context,
	inCh chan T,
	outCh chan U,
	mapFn func(context.Context, T) ([]U, error),
) {
	for {
		select {
		case ev := <-inCh:
			out, err := mapFn(ctx, ev)
			if err != nil {
				log.Error("Failed to map, err: %w", err)
				return
			}
			for _, ev := range out {
				outCh <- ev
			}
		case <-ctx.Done():
			log.Info("Aborting.")
			return
		}
	}
}
