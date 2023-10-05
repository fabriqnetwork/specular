package utils

import (
	"context"

	"github.com/ethereum/go-ethereum/event"
	"github.com/specularL2/specular/services/sidecar/utils/log"
)

type Broker[T any] struct {
	PubCh   chan T        // Input to broker
	subCh   chan chan T   // Subscribes to broker
	unsubCh chan chan T   // Unsubscribes from broker
	stopCh  chan struct{} // Stops broker
}

func NewBroker[T any]() *Broker[T] {
	return &Broker[T]{
		PubCh:   make(chan T, 1),
		subCh:   make(chan chan T, 1),
		unsubCh: make(chan chan T, 1),
		stopCh:  make(chan struct{}),
	}
}

var chanSize = 8

// TODO: remove dependency on `event.Subscription`
func (b *Broker[T]) Start(ctx context.Context, sub event.Subscription) error {
	subs := map[chan T]struct{}{}
	for {
		select {
		case msg := <-b.PubCh:
			for msgCh := range subs {
				// Note: msgCh is buffered, non-blocking send protects broker
				// default:
				select {
				case msgCh <- msg:
				}
			}
		case msgCh := <-b.subCh:
			subs[msgCh] = struct{}{}
		case msgCh := <-b.unsubCh:
			delete(subs, msgCh)
		case <-b.stopCh:
			return nil
		case err := <-sub.Err():
			log.Warn("Subscription error, stopping broker", "err", sub.Err())
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

func (b *Broker[T]) Subscribe() chan T {
	msgCh := make(chan T, chanSize)
	b.subCh <- msgCh
	return msgCh
}

// Subscribes to a new channel mapped from `inCh` (one-to-one).
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
	b.PubCh <- msg
}

// Creates and publishes events to a channel mapped from `inCh` (one-to-many).
func SubscribeMappedToMany[T any, U any](
	ctx context.Context,
	broker *Broker[T],
	mapFn func(context.Context, T) ([]U, error),
) <-chan U {
	inCh := broker.Subscribe()
	outCh := make(chan U, chanSize*chanSize)
	go func() {
		defer broker.Unsubscribe(inCh)
		defer close(outCh)
		for {
			select {
			case head := <-inCh:
				out, err := mapFn(ctx, head)
				if err != nil {
					log.Error("Failed to map, err: %w", err)
					return
				}
				for _, event := range out {
					outCh <- event
				}
			case <-ctx.Done():
				log.Info("Aborting.")
				return
			}
		}
	}()
	return outCh
}
