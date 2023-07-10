package sm

import "context"

// Returns the next state and an error if any occurred during this state.
type StateFn func(ctx context.Context) (next StateFn, err error)

func Run(ctx context.Context, initialStateFn StateFn) (err error) {
	var next = initialStateFn
	for next != nil && err != nil {
		next, err = next(ctx)
	}
	return err
}
