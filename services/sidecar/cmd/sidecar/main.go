package main

import (
	"context"
	"bytes"
	"math"
	"math/big"
	"reflect"
	"os"

	"github.com/urfave/cli/v2"
	"golang.org/x/sync/errgroup"

	"github.com/ethereum/go-ethereum/accounts"
	bind "github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/external"
	"github.com/ethereum/go-ethereum/accounts/keystore"
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

func initializeAccountManager(cfg services.KeyStoreConfig) *accounts.Manager {
	keyDir := cfg.GetKeyStoreDir()
	log.Info(keyDir)

	// TODO: verify if hard-coded true is a problem
	acctMgr := accounts.NewManager(&accounts.Config{InsecureUnlockAllowed: true})

	// TODO: Enable external signer
	//if (cfg.ExternalSigner) {
		//log.Info("Using external signer", "url", cfg.ExternalSigner)
		//if extapi, err := external.NewExternalBackend(cfg.ExternalSigner); err == nil {
			//am.AddBackend(extapi)
			//return nil
		//} else {
			//return fmt.Errorf("error connecting to external signer: %v", err)
		//}
	//}

	keystore := keystore.NewKeyStore(keyDir, keystore.StandardScryptN, keystore.StandardScryptP)

	acctMgr.AddBackend(keystore)

	return acctMgr
}

func createDisseminator(
	ctx context.Context,
	cfg *services.SystemConfig,
	accountMgr accountManager,
) (*disseminator.BatchDisseminator, error) {
	l1TxMgr, err := createTxManager(ctx, "disseminator", cfg.L1(), cfg.Sequencer(), accountMgr)
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
	accountMgr accountManager,
) (*validator.Validator, error) {
	l1TxMgr, err := createTxManager(ctx, "validator", cfg.L1(), cfg.Validator(), accountMgr)
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
	accountMgr accountManager,
) (*bridge.TxManager, error) {
	transactor, err := createTransactor(
		accountMgr, serCfg.GetAccountAddr(), serCfg.GetClefEndpoint(), serCfg.GetPassphrase(), l1Cfg.GetChainID(),
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
	mgr accountManager,
	accountAddress common.Address,
	clefEndpoint string,
	passphrase string,
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
	var ks *keystore.KeyStore
	if keystores := mgr.Backends(keystore.KeyStoreType); len(keystores) > 0 {
		ks = keystores[0].(*keystore.KeyStore)
	} else {
		return nil, fmt.Errorf("keystore not found")
	}
	json, err := ks.Export(accounts.Account{Address: accountAddress}, passphrase, "")
	if err != nil {
		return nil, fmt.Errorf("failed to export account for %s: %w", accountAddress, err)
	}
	transactor, err := bind.NewTransactorWithChainID(bytes.NewReader(json), passphrase, new(big.Int).SetUint64(chainID))
	if err != nil {
		return nil, fmt.Errorf("failed to create transactor: %w", err)
	}
	return transactor, nil
}


func startService(cliCtx *cli.Context) error {

	log.Info("Reading configuration")

	cfg, err := services.ParseSystemConfig(cliCtx)
	if err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}

	log.Info("Initializing Account Manager")

	acctMgr := initializeAccountManager(cfg.KeyStore())

	log.Info("Starting disseminator+validator")

	var disseminator *disseminator.BatchDisseminator
	var validator *validator.Validator

	if cfg.Sequencer().GetIsEnabled() {
		disseminator, err = createDisseminator(context.Background(), cfg, acctMgr)
		if err != nil {
			return fmt.Errorf("failed to initialize disseminator: %w", err)
		}
	}
	if cfg.Validator().GetIsEnabled() {
		validator, err = createValidator(context.Background(), cfg, acctMgr)
		if err != nil {
			return fmt.Errorf("failed to initialize validator: %w", err)
		}
	}

	disseminatorErrorGroup, disseminatorCtx := errgroup.WithContext(context.Background())
	disseminator.Start(disseminatorCtx, disseminatorErrorGroup)

	validatorErrorGroup, validatorCtx := errgroup.WithContext(context.Background())
	validator.Start(validatorCtx, validatorErrorGroup)

	return nil
}

func main() {
    app := &cli.App{
        Name:  "sidecar",
        Usage: "launch validator+disseminator",
        Action: startService,
    }

	app.Flags = services.CLIFlags()

    if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
    }
}
