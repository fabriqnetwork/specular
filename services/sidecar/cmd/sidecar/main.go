package main

import (
	"context"
	"crypto/ecdsa"
	"math"
	"math/big"
	"os"

	"github.com/urfave/cli/v2"
	"golang.org/x/sync/errgroup"

	"github.com/ethereum/go-ethereum/accounts"
	bind "github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/external"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/specularL2/specular/services/sidecar/rollup/derivation"
	"github.com/specularL2/specular/services/sidecar/rollup/rpc/bridge"
	"github.com/specularL2/specular/services/sidecar/rollup/rpc/eth"
	"github.com/specularL2/specular/services/sidecar/rollup/rpc/eth/txmgr"
	"github.com/specularL2/specular/services/sidecar/rollup/services"
	"github.com/specularL2/specular/services/sidecar/rollup/services/disseminator"
	"github.com/specularL2/specular/services/sidecar/rollup/services/validator"
	"github.com/specularL2/specular/services/sidecar/utils/fmt"
	"github.com/specularL2/specular/services/sidecar/utils/log"
)

type serviceCfg interface {
	GetAccountAddr() common.Address
	GetSecretKey() *ecdsa.PrivateKey
	GetClefEndpoint() string
	GetTxMgrCfg() txmgr.Config
}

func main() {
	app := &cli.App{
		Name:   "sidecar",
		Usage:  "launch a validator and/or disseminator",
		Action: startService,
	}
	app.Flags = services.CLIFlags()
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func startService(cliCtx *cli.Context) error {
	log.Info("Reading configuration")
	cfg, err := services.ParseSystemConfig(cliCtx)
	if err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}

	var (
		disseminator *disseminator.BatchDisseminator
		validator    *validator.Validator
		eg, ctx      = errgroup.WithContext(context.Background())
	)
	if cfg.Sequencer().GetIsEnabled() {
		disseminator, err = createDisseminator(context.Background(), cfg)
		if err != nil {
			return fmt.Errorf("failed to create disseminator: %w", err)
		}
		if err := disseminator.Start(ctx, eg); err != nil {
			return fmt.Errorf("failed to start disseminator: %w", err)
		}
	}
	if cfg.Validator().GetIsEnabled() {
		validator, err = createValidator(context.Background(), cfg)
		if err != nil {
			return fmt.Errorf("failed to create validator: %w", err)
		}
		if err := validator.Start(ctx, eg); err != nil {
			return fmt.Errorf("failed to start validator: %w", err)
		}
	}
	log.Info("Services running.")
	if err := eg.Wait(); err != nil {
		return fmt.Errorf("service failed while running: %w", err)
	}
	log.Info("Services stopped.")
	return nil
}

func createDisseminator(
	ctx context.Context,
	cfg *services.SystemConfig,
) (*disseminator.BatchDisseminator, error) {
	l1TxMgr, err := createTxManager(ctx, "disseminator", cfg.L1(), cfg.Sequencer())
	if err != nil {
		return nil, fmt.Errorf("failed to initialize l1 tx manager: %w", err)
	}
	batchBuilder, err := derivation.NewBatchBuilder(math.MaxInt64) // TODO: configure max batch size
	if err != nil {
		return nil, fmt.Errorf("failed to initialize batch builder: %w", err)
	}
	l2Client := eth.NewLazilyDialedEthClient(cfg.L2().GetEndpoint())
	return disseminator.NewBatchDisseminator(cfg.Sequencer(), batchBuilder, l1TxMgr, l2Client), nil
}

func createValidator(
	ctx context.Context,
	cfg *services.SystemConfig,
) (*validator.Validator, error) {
	l1TxMgr, err := createTxManager(ctx, "validator", cfg.L1(), cfg.Validator())
	if err != nil {
		return nil, fmt.Errorf("failed to initialize l1 tx manager: %w", err)
	}
	l1Client, err := eth.DialWithRetry(ctx, cfg.L1().GetEndpoint())
	if err != nil {
		return nil, fmt.Errorf("failed to initialize l1 client: %w", err)
	}
	l1BridgeClient, err := bridge.NewBridgeClient(l1Client, cfg.L1())
	if err != nil {
		return nil, fmt.Errorf("failed to initialize l1 bridge client: %w", err)
	}
	l1State := eth.NewEthState()
	l1Syncer := eth.NewEthSyncer(l1State)
	l1Syncer.Start(ctx, l1Client)
	l2Client := eth.NewLazilyDialedEthClient(cfg.L2().GetEndpoint())
	return validator.NewValidator(cfg.Validator(), l1TxMgr, l1BridgeClient, l1State, l2Client), nil
}

func createTxManager(
	ctx context.Context,
	name string,
	l1Cfg services.L1Config,
	serCfg serviceCfg,
) (*bridge.TxManager, error) {
	transactor, err := createTransactor(
		serCfg.GetAccountAddr(), serCfg.GetClefEndpoint(), serCfg.GetSecretKey(), l1Cfg.GetChainID(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize transactor: %w", err)
	}
	l1Client, err := eth.DialWithRetry(ctx, l1Cfg.GetEndpoint())
	if err != nil {
		return nil, fmt.Errorf("failed to initialize l1 client: %w", err)
	}
	signer := func(ctx context.Context, address common.Address, tx *ethTypes.Transaction) (*ethTypes.Transaction, error) {
		return transactor.Signer(address, tx)
	}
	return bridge.NewTxManager(txmgr.NewTxManager(log.New("service", name), serCfg.GetTxMgrCfg(), l1Client, signer), l1Cfg)
}

// Creates a transactor for the given account address, either using a clef endpoint (preferred) or secret key.
func createTransactor(
	accountAddress common.Address,
	clefEndpoint string,
	secretKey *ecdsa.PrivateKey,
	chainID uint64,
) (*bind.TransactOpts, error) {
	if clefEndpoint != "" {
		clef, err := external.NewExternalSigner(clefEndpoint)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize external signer from clef endpoint: %w", err)
		}
		return bind.NewClefTransactor(clef, accounts.Account{Address: accountAddress}), nil
	}
	log.Warn("No external signer specified, using geth signer")
	return bind.NewKeyedTransactorWithChainID(secretKey, new(big.Int).SetUint64(chainID))
}
