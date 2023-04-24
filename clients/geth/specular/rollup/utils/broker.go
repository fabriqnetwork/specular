package utils

import (
	"context"

	"github.com/specularl2/specular/clients/geth/specular/rollup/utils/log"
)

// Example usage:
// broker := NewBroker[uint64]()
// broker.Publish(123)
// broker.PubCh() <- 123 // Equivalent behavior as above (used when abstracting away broker from publisher)
// recv := broker.Subscribe(ctx)
// broker.Start(ctx, interrupt)
type Broker[T any] struct {
	pubCh   chan T        // Input to broker
	subCh   chan chan T   // Subscribes to broker
	unsubCh chan chan T   // Unsubscribes from broker
	stopCh  chan struct{} // Stops broker

	isBlocking bool // Whether broadcast should block on receivers.
}

func NewBroker[T any]() *Broker[T] {
	return &Broker[T]{
		pubCh:      make(chan T, 1),
		subCh:      make(chan chan T, 1),
		unsubCh:    make(chan chan T, 1),
		stopCh:     make(chan struct{}),
		isBlocking: true,
	}
}

var chanSize = 8

type Interrupt interface {
	Err() <-chan error
}

func (b *Broker[T]) PubCh() chan<- T {
	return b.pubCh
}

// Starts broadcasting messages received on `pubCh`, until interrupted.
// Note: Interruptor used to support `event.Subscription`
func (b *Broker[T]) Start(ctx context.Context, interrupt Interrupt) error {
	subs := map[chan T]struct{}{}
	for {
		select {
		case msg := <-b.pubCh:
			// Broadcasts `msg`. Note: this is a blocking operation.
			if b.isBlocking {
				for msgCh := range subs {
					msgCh <- msg
				}
			} else {
				for msgCh := range subs {
					select {
					case msgCh <- msg:
					default:
						log.Info("Dropping msg; subscriber not ready", "msg", msg)
					}
				}
			}
		case msgCh := <-b.subCh:
			// `msgCh` receives broadcasts as of the next iteration.
			subs[msgCh] = struct{}{}
		case msgCh := <-b.unsubCh:
			// `msgCh` stops receiving broadcasts as of the next iteration.
			delete(subs, msgCh)
		case <-b.stopCh:
			// Stops broadcasting.
			log.Info("Stopping broker.")
			return nil
		case err := <-interrupt.Err():
			log.Warn("Interrupted by error, stopping broker", "err", interrupt.Err())
			return err
		case <-ctx.Done():
			log.Info("Aborting.")
			return ctx.Err()
		}
	}
}

func (b *Broker[T]) Stop() {
	close(b.stopCh)
}

// Subscribes to a new channel.
func (b *Broker[T]) Subscribe() chan T {
	msgCh := make(chan T, chanSize)
	b.subCh <- msgCh
	return msgCh
}

// Subscribes to a new channel, mapped from `pubCh` (one-to-one).
func (b *Broker[T]) SubscribeWithCallback(
	ctx context.Context,
	callbackFn func(context.Context, T) error,
) {
	inCh := b.Subscribe()
	go func() {
		defer b.Unsubscribe(inCh)
		for {
			select {
			case head := <-inCh:
				err := callbackFn(ctx, head)
				if err != nil {
					log.Error("Failed triggering callback, err: %w", err)
					return
				}
			case <-ctx.Done():
				log.Info("Aborting.")
				return
			}
		}
	}()
}

func (b *Broker[T]) Unsubscribe(msgCh chan T) {
	b.unsubCh <- msgCh
}

func (b *Broker[T]) Publish(msg T) {
	b.pubCh <- msg
}

// Creates and publishes events to a channel mapped from broker (one-to-many).
func SubscribeBrokerMappedToMany[T any, U any](
	ctx context.Context,
	broker *Broker[T],
	mapFn func(context.Context, T) ([]U, error),
) <-chan U {
	inCh := broker.Subscribe()
	outCh := make(chan U, chanSize*chanSize)
	go func() {
		defer broker.Unsubscribe(inCh)
		defer close(outCh)
		mapChToMany(ctx, inCh, outCh, mapFn)
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
		case head := <-inCh:
			out, err := mapFn(ctx, head)
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
		case head := <-inCh:
			out, err := mapFn(ctx, head)
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
