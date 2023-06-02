package services

import (
	"bytes"
	"context"
	"math"
	"math/big"
	"reflect"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/external"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/specularl2/specular/clients/geth/specular/proof"
	"github.com/specularl2/specular/clients/geth/specular/rollup/rpc/bridge"
	"github.com/specularl2/specular/clients/geth/specular/rollup/rpc/eth"
	"github.com/specularl2/specular/clients/geth/specular/rollup/rpc/eth/txmgr"
	"github.com/specularl2/specular/clients/geth/specular/rollup/services/api"
	"github.com/specularl2/specular/clients/geth/specular/rollup/services/derivation"
	"github.com/specularl2/specular/clients/geth/specular/rollup/services/derivation/stage"
	"github.com/specularl2/specular/clients/geth/specular/rollup/services/sequencer"
	"github.com/specularl2/specular/clients/geth/specular/rollup/services/validator"
	"github.com/specularl2/specular/clients/geth/specular/rollup/types/da"
	"github.com/specularl2/specular/clients/geth/specular/utils/fmt"
	"github.com/specularl2/specular/clients/geth/specular/utils/log"
)

// TODO: this is the last Geth-specific interface here; remove.
type accountManager interface {
	Backends(reflect.Type) []accounts.Backend
}

// Creates services configured by cfg:
// - Driver (always)
// - Sequencer (if configured)
// - Validator (if configured)
func CreateRollupServices(
	accMgr accountManager,
	execBackend api.ExecutionBackend,
	proofBackend proof.Backend,
	cfg *systemConfig,
) ([]api.Service, error) {
	var (
		ctx      = context.Background()
		services []api.Service
	)
	rollupState, err := createRollupState(ctx, cfg.L1())
	if err != nil {
		return nil, fmt.Errorf("Failed to initialize rollup state: %w", err)
	}

	// Create driver
	driver, err := createDriver(ctx, cfg, execBackend, rollupState)
	if err != nil {
		return nil, fmt.Errorf("Failed to initialize driver: %w", err)
	}
	services = append(services, driver)

	// Create sequencer
	if (cfg.Sequencer().GetAccountAddr() != common.Address{}) {
		sequencer, err := createSequencer(ctx, cfg, accMgr, execBackend)
		if err != nil {
			return nil, fmt.Errorf("Failed to initialize sequencer: %w", err)
		}
		services = append(services, sequencer)
	}
	// Create validator
	if (cfg.Validator().GetAccountAddr() != common.Address{}) {
		validator, err := createValidator(ctx, cfg, accMgr, rollupState, proofBackend)
		if err != nil {
			return nil, fmt.Errorf("Failed to create validator: %w", err)
		}
		services = append(services, validator)
	}
	return services, nil
}

// Creates driver.
// Two L1 clients are created; one for L1 state syncing and one for fetching L1 blocks.
// An L2 client factory fn is also created (lazily create l2 client since the blockchain hasn't started yet).
func createDriver(
	ctx context.Context,
	cfg *systemConfig,
	execBackend api.ExecutionBackend,
	rollupState *derivation.RollupState,
) (*derivation.Driver, error) {
	if err := bridge.EnsureUtilInit(); err != nil {
		return nil, fmt.Errorf("Failed to initialize bridge util: %w", err)
	}
	l1Client, err := eth.DialWithRetry(ctx, cfg.L1().Endpoint(), nil)
	if err != nil {
		return nil, fmt.Errorf("Failed to initialize l1 client: %w", err)
	}
	l1State, err := createSyncingL1State(ctx, cfg.L1()) // TODO: move into a service for proper cleanup on stop.
	if err != nil {
		return nil, fmt.Errorf("Failed to start l1 state sync: %w", err)
	}
	l2ClientCreatorFn := func(ctx context.Context) (stage.EthClient, error) {
		return eth.DialWithRetry(ctx, cfg.L2().GetEndpoint(), nil)
	}
	terminalStage := stage.CreatePipeline(cfg.L1(), execBackend, rollupState, l2ClientCreatorFn, l1Client, l1State)
	driver := derivation.NewDriver(cfg.Driver(), terminalStage)
	return driver, nil
}

func createSequencer(
	ctx context.Context,
	cfg *systemConfig,
	accountMgr accountManager,
	execBackend api.ExecutionBackend,
) (*sequencer.Sequencer, error) {
	l1TxMgr, err := createTxManager(
		ctx, cfg, accountMgr, cfg.Sequencer().GetAccountAddr(), cfg.L2().GetClefEndpoint(), cfg.Sequencer().GetPassphrase(),
	)
	if err != nil {
		return nil, fmt.Errorf("Failed to initialize l1 tx manager: %w", err)
	}
	batchBuilder, err := da.NewBatchBuilder(math.MaxInt64) // TODO: configure max batch size
	if err != nil {
		return nil, fmt.Errorf("Failed to initialize batch builder: %w", err)
	}
	l2ClientCreatorFn := func(ctx context.Context) (sequencer.L2Client, error) {
		return eth.DialWithRetry(ctx, cfg.L2().GetEndpoint(), nil)
	}
	return sequencer.NewSequencer(cfg.Sequencer(), execBackend, l2ClientCreatorFn, l1TxMgr, batchBuilder), nil
}

func createValidator(
	ctx context.Context,
	cfg *systemConfig,
	accountMgr accountManager,
	rollupState *derivation.RollupState,
	proofBackend proof.Backend,
) (*validator.Validator, error) {
	l2ClientCreatorFn := func(ctx context.Context) (validator.L2Client, error) {
		return eth.DialWithRetry(ctx, cfg.L2().GetEndpoint(), nil)
	}
	// Create tx manager, used to send transactions to L1.
	l1TxMgr, err := createTxManager(
		ctx, cfg, accountMgr, cfg.Validator().GetAccountAddr(), cfg.L2().GetClefEndpoint(), cfg.Validator().GetPassphrase(),
	)
	if err != nil {
		return nil, fmt.Errorf("Failed to initialize tx manager: %w", err)
	}
	return validator.NewValidator(cfg.Validator(), l2ClientCreatorFn, l1TxMgr, proofBackend, rollupState), nil
}

func createRollupState(ctx context.Context, cfg L1Config) (*derivation.RollupState, error) {
	bridgeClient, err := bridge.DialWithRetry(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("Failed to initialize l1 bridge client: %w", err)
	}
	rollupState := derivation.NewRollupState(bridgeClient)
	if err != nil {
		return nil, fmt.Errorf("Failed to initialize assertion manager: %w", err)
	}
	return rollupState, nil
}

func createSyncingL1State(ctx context.Context, cfg L1Config) (*eth.EthState, error) {
	l1State := eth.NewEthState()
	l1Client, err := eth.DialWithRetry(ctx, cfg.Endpoint(), nil)
	if err != nil {
		return nil, fmt.Errorf("Failed to initialize l1 client: %w", err)
	}
	syncer := eth.NewEthSyncer(l1State)
	syncer.Start(ctx, l1Client)
	return l1State, nil
}

func createTxManager(
	ctx context.Context,
	cfg *systemConfig,
	accountMgr accountManager,
	accountAddr common.Address,
	clefEndpoint string,
	passphrase string,
) (*bridge.TxManager, error) {
	transactor, err := createTransactor(accountMgr, accountAddr, clefEndpoint, passphrase, cfg.ChainID())
	if err != nil {
		return nil, fmt.Errorf("Failed to initialize transactor: %w", err)
	}
	l1Client, err := eth.DialWithRetry(ctx, cfg.L1().Endpoint(), nil)
	if err != nil {
		return nil, fmt.Errorf("Failed to initialize l1 client: %w", err)
	}
	signer := func(ctx context.Context, address common.Address, tx *types.Transaction) (*types.Transaction, error) {
		return transactor.Signer(address, tx)
	}
	// TODO: config
	return bridge.NewTxManager(txmgr.NewTxManager(cfg.Sequencer().GetTxMgrCfg(), l1Client, signer), cfg.L1())
}

// createTransactor creates a transactor for the given account address,
// either using the clef endpoint or passphrase.
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
			return nil, fmt.Errorf("Failed to initialize external signer from clef endpoint: %w", err)
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
		return nil, fmt.Errorf("Failed to export account: %w", err)
	}
	transactor, err := bind.NewTransactorWithChainID(bytes.NewReader(json), passphrase, new(big.Int).SetUint64(chainID))
	if err != nil {
		return nil, fmt.Errorf("Failed to create transactor: %w", err)
	}
	return transactor, nil
}
