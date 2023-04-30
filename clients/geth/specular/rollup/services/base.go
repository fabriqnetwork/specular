package services

import (
	"context"

	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/specularl2/specular/clients/geth/specular/rollup/utils/fmt"
	"golang.org/x/sync/errgroup"
)

type BaseService struct {
	Eg *errgroup.Group
}

// Starts the rollup service.
func (b *BaseService) Start() context.Context {
	eg, ctx := errgroup.WithContext(context.Background())
	b.Eg = eg
	return ctx
}

func (b *BaseService) Stop() error {
	log.Info("Stopping service...")
	b.Eg.Go(func() error { return fmt.Errorf("Force-stopping service.") })
	b.Eg.Wait()
	log.Info("Service stopped.")
	return nil
}

func (i *BaseService) APIs() []rpc.API { return []rpc.API{} }
