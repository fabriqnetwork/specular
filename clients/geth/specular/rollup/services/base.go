package services

import (
	"context"

	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/specularl2/specular/clients/geth/specular/rollup/utils/fmt"
	"golang.org/x/sync/errgroup"
)

type BaseService struct {
	StartCtx context.Context // Can't be passed to `Start` so we have to store it here.
	Eg       *errgroup.Group // Group used to manage service goroutines.
}

// Starts the rollup service.
func (b *BaseService) Start() context.Context {
	ctx := b.StartCtx
	// Initialize if not already initialized.
	if b.Eg == nil {
		b.Eg, ctx = errgroup.WithContext(b.StartCtx)
	}
	return ctx
}

func (b *BaseService) Stop() error {
	log.Info("Stopping service...")
	b.Eg.Go(func() error { return fmt.Errorf("Force-stopping service.") })
	// Ignore error (we raised it in the above goroutine).
	_ = b.Eg.Wait()
	log.Info("Service stopped.")
	return nil
}

func (b *BaseService) ErrGroup() *errgroup.Group { return b.Eg }

func (b *BaseService) APIs() []rpc.API { return []rpc.API{} }
