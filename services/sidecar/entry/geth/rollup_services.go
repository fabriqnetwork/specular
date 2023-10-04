package geth

import (
	"context"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/specularl2/specular/clients/geth/specular/proof"
	"github.com/specularl2/specular/clients/geth/specular/rollup"
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

// TODO: remove this when EL and CL client are decoupled.
// Creates rollup services and registers them as Geth services.
func RegisterGethRollupServices(
	stack Node,
	cliCtx *cli.Context,
	eth api.ExecutionBackend,
	proofBackend proof.Backend,
) error {
	log.Info("Parsing system config...")
	cfg, err := services.ParseSystemConfig(cliCtx)
	if err != nil {
		return fmt.Errorf("failed to parse system config: %w", err)
	}
	vmConfig := eth.BlockChain().GetVMConfig()
	type hookCfg struct {
		services.L2Config
		services.SequencerConfig
	}
	vmConfig.SpecularEVMPreTransferHook = MakeSpecularEVMPreTransferHook(hookCfg{cfg.L2(), cfg.Sequencer()})

	services, err := rollup.CreateRollupServices(stack.AccountManager(), eth, proofBackend, cfg)
	if err != nil {
		return fmt.Errorf("failed to create rollup services: %w", err)
	}
	log.Info("Registering services...")
	for _, service := range services {
		eg, ctx := errgroup.WithContext(context.Background())
		stack.RegisterLifecycle(&gethRollupService{ctx, eg, service})
	}
	return nil
}

type gethRollupService struct {
	StartCtx context.Context // Can't be passed to `Start` so we have to store it here.
	Eg       *errgroup.Group // Group used to manage service goroutines.
	service  api.Service     // The rollup service to start.
}

// Ensures that gethRollupService implements the node.Lifecycle interface.
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
