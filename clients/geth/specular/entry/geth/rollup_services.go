package geth

import (
	"context"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/specularl2/specular/clients/geth/specular/proof"
	"github.com/specularl2/specular/clients/geth/specular/rollup/services"
	"github.com/specularl2/specular/clients/geth/specular/rollup/services/api"

	"github.com/specularl2/specular/clients/geth/specular/utils/fmt"
	"github.com/specularl2/specular/clients/geth/specular/utils/log"
	"github.com/urfave/cli/v2"
	"golang.org/x/sync/errgroup"
)

type Node interface {
	RegisterLifecycle(lc node.Lifecycle)
	AccountManager() *accounts.Manager
}

// Creates rollup services and registers them as Geth services.
func RegisterGethRollupServices(
	stack Node,
	cliCtx *cli.Context,
	eth GethBackend,
	proofBackend proof.Backend,
) error {
	cfg := services.ParseSystemConfig(cliCtx)
	execBackend, err := NewExecutionBackend(eth, cfg.Sequencer().AccountAddr())
	if err != nil {
		return fmt.Errorf("Failed to create geth execution backend: %w", err)
	}
	services, err := services.CreateRollupServices(stack.AccountManager(), execBackend, proofBackend, cfg)
	if err != nil {
		return fmt.Errorf("Failed to create rollup services: %w", err)
	}
	eg, ctx := errgroup.WithContext(context.Background())
	for _, service := range services {
		stack.RegisterLifecycle(&gethRollupService{ctx, eg, service})
	}
	return nil
}

type gethRollupService struct {
	StartCtx context.Context // Can't be passed to `Start` so we have to store it here.
	Eg       *errgroup.Group // Group used to manage service goroutines.
	service  api.Service     // The rollup service to start.
}

var _ node.Lifecycle = (*gethRollupService)(nil)

// Starts the rollup service.
func (b *gethRollupService) Start() error {
	ctx := b.StartCtx
	// Initialize if not already initialized.
	if b.Eg == nil {
		b.Eg, ctx = errgroup.WithContext(b.StartCtx)
	}
	return b.service.Start(ctx, b.Eg)
}

func (b *gethRollupService) Stop() error {
	log.Info("Stopping service...")
	b.Eg.Go(func() error { return fmt.Errorf("Force-stopping service.") })
	_ = b.Eg.Wait() // Ignore error (we raised it in the above goroutine).
	log.Info("Service stopped.")
	return nil
}

// TODO: support APIs.
func (b *gethRollupService) APIs() []rpc.API { return []rpc.API{} }
