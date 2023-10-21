package main

import (
	"context"
	"crypto/ecdsa"
	"math"
	"math/big"
	"os"
	"reflect"

	"github.com/urfave/cli/v2"
	"golang.org/x/sync/errgroup"

	"github.com/ethereum/go-ethereum/accounts"
	bind "github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/external"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"

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

// TODO: this is the last Geth-specific interface here; remove.
type accountManager interface {
	Backends(reflect.Type) []accounts.Backend
}

type serviceCfg interface {
	GetAccountAddr() common.Address
	GetClefEndpoint() string
	GetPassphrase() string
	GetTxMgrCfg() txmgr.Config
}

func createDisseminator(
	ctx context.Context,
	cfg *services.SystemConfig,
) (*disseminator.BatchDisseminator, error) {

	hexKey := os.Getenv("SEQUENCER_PRIVATE_KEY")
	log.Info("getting private key from env")
	privateKey, err := crypto.HexToECDSA(hexKey[2:])
	if err != nil && cfg.Sequencer().ClefEndpoint == "" {
		return nil, fmt.Errorf("could not read private key: %w", err)
	}

	l1TxMgr, err := createTxManager(ctx, "disseminator", cfg.L1(), cfg.Sequencer(), privateKey)
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

	hexKey := os.Getenv("VALIDATOR_PRIVATE_KEY")
	privateKey, err := crypto.HexToECDSA(hexKey[2:])
	if err != nil && cfg.Validator().ClefEndpoint == "" {
		return nil, fmt.Errorf("could not read private key: %w", err)
	}

	l1TxMgr, err := createTxManager(ctx, "validator", cfg.L1(), cfg.Validator(), privateKey)
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
	privateKey *ecdsa.PrivateKey,
) (*bridge.TxManager, error) {
	transactor, err := createTransactor(
		serCfg.GetAccountAddr(), serCfg.GetClefEndpoint(), l1Cfg.GetChainID(), privateKey,
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

// Creates a transactor for the given account address, either using the clef endpoint or passphrase.
func createTransactor(
	accountAddress common.Address,
	clefEndpoint string,
	chainID uint64,
	privateKey *ecdsa.PrivateKey,
) (*bind.TransactOpts, error) {
	if clefEndpoint != "" {
		clef, err := external.NewExternalSigner(clefEndpoint)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize external signer from clef endpoint: %w", err)
		}
		return bind.NewClefTransactor(clef, accounts.Account{Address: accountAddress}), nil
	}

	return bind.NewKeyedTransactorWithChainID(privateKey, new(big.Int).SetUint64(chainID))
}


func startService(cliCtx *cli.Context) error {

	log.Info("Reading configuration")

	cfg, err := services.ParseSystemConfig(cliCtx)
	if err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}

	disseminatorErrorGroup, disseminatorCtx := errgroup.WithContext(context.Background())
	validatorErrorGroup, validatorCtx := errgroup.WithContext(context.Background())

	// start services
	if cfg.Sequencer().GetIsEnabled() {
		disseminator, err := createDisseminator(context.Background(), cfg)
		if err != nil {
			return fmt.Errorf("failed to initialize disseminator: %w", err)
		}
		disseminator.Start(disseminatorCtx, disseminatorErrorGroup)
	}
	if cfg.Validator().GetIsEnabled() {
		validator, err := createValidator(context.Background(), cfg)
		if err != nil {
			return fmt.Errorf("failed to initialize validator: %w", err)
		}
		validator.Start(validatorCtx, validatorErrorGroup)
	}

	// wait for services to finish
	if cfg.Sequencer().GetIsEnabled() {
		<-disseminatorCtx.Done()
	}
	if cfg.Validator().GetIsEnabled() {
		<-validatorCtx.Done()
	}

	return nil
}

func main() {
    app := &cli.App{
        Name:  "sidecar",
        Usage: "launch the specular validator and/or disseminator",
        Action: startService,
    }

	app.Flags = services.CLIFlags()

    if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
    }
}
