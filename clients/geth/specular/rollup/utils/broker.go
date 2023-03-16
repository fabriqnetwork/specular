package utils

import (
	"context"

	"github.com/ethereum/go-ethereum/event"
	"github.com/specularl2/specular/clients/geth/specular/rollup/utils/log"
)

type Broker[T any] struct {
	PubCh   chan T
	subCh   chan chan T
	unsubCh chan chan T
	stopCh  chan struct{}
}

func NewBroker[T any]() *Broker[T] {
	return &Broker[T]{
		PubCh:   make(chan T, 1),
		subCh:   make(chan chan T, 1),
		unsubCh: make(chan chan T, 1),
		stopCh:  make(chan struct{}),
	}
}

// TODO: remove dependency on `event.Subscription`
func (b *Broker[T]) Start(ctx context.Context, sub event.Subscription) error {
	subs := map[chan T]struct{}{}
	for {
		select {
		case msg := <-b.PubCh:
			for msgCh := range subs {
				// Note: msgCh is buffered, non-blocking send protects broker
				select {
				case msgCh <- msg:
				default:
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
			return ctx.Err()
		}
	}
}

func (b *Broker[T]) Stop() {
	close(b.stopCh)
}

func (b *Broker[T]) Subscribe() chan T {
	msgCh := make(chan T, 1)
	b.subCh <- msgCh
	return msgCh
}

func (b *Broker[T]) Unsubscribe(msgCh chan T) {
	b.unsubCh <- msgCh
}

func (b *Broker[T]) Publish(msg T) {
	b.PubCh <- msg
}

// Creates and publishes events to a channel mapped from `inCh`.
func SubscribeMapped[T any, U any](
	ctx context.Context,
	broker *Broker[T],
	mapFn func(context.Context, T) ([]U, error),
) <-chan U {
	inCh := broker.Subscribe()
	outCh := make(chan U, 128)
	go func() {
		defer broker.Unsubscribe(inCh)
		defer close(outCh)
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
			return
		}
		return
	}()
	return outCh
}
