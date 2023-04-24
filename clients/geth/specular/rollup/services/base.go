package services

import (
	"context"
	"sync"

	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rpc"
)

type BaseService struct {
	cancel context.CancelFunc
	Wg     sync.WaitGroup
}

// Starts the rollup service.
func (b *BaseService) Start() (context.Context, error) {
	ctx, cancel := context.WithCancel(context.Background())
	b.cancel = cancel
	return ctx, nil
}

func (b *BaseService) Stop() error {
	log.Info("Stopping service...")
	b.cancel()
	b.Wg.Wait()
	log.Info("Service stopped.")
	return nil
}

func (i *BaseService) APIs() []rpc.API { return []rpc.API{} }
