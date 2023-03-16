package client

import (
	"context"
	"reflect"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/specularl2/specular/clients/geth/specular/rollup/utils"
	"github.com/specularl2/specular/clients/geth/specular/rollup/utils/fmt"
)

// Interface used to represent Iterators from `bindings`
type Iterable interface {
	Next() bool
	Error() error
	Close() error
}

func SubscribeHeaderMapped[T any, U Iterable](
	ctx context.Context,
	broker *utils.Broker[*types.Header],
	fn func(*bind.FilterOpts) (U, error),
	start uint64,
) <-chan T {
	return utils.SubscribeMapped(ctx, broker, toMapperFn[T](fn, start))
}

func toMapperFn[T any, U Iterable](
	filterFn func(*bind.FilterOpts) (U, error),
	start uint64,
) func(context.Context, *types.Header) ([]T, error) {
	var last *types.Header
	return func(ctx context.Context, header *types.Header) ([]T, error) {
		end := header.Number.Uint64()
		opts := &bind.FilterOpts{
			Start:   start,
			End:     &end,
			Context: ctx,
		}
		iter, err := filterFn(opts)
		if err != nil {
			return nil, fmt.Errorf("Failed to filter, err: %w", err)
		}
		var mapped []T
		for iter.Next() {
			// TODO: remove this hack
			mapped = append(mapped, reflect.ValueOf(iter).FieldByName("Event").Interface().(T))
		}
		if iter.Error() != nil {

		}
		last = header
		start = last.Number.Uint64()
		return mapped, nil
	}
}
