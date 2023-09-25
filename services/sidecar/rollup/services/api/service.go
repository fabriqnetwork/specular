package api

import "context"

type Service interface {
	// Starts the service (non-blocking).
	// Long-running goroutines must be scheduled via `eg`, using `ctx`.
	Start(ctx context.Context, eg ErrGroup) error
}

type ErrGroup interface{ Go(f func() error) }
