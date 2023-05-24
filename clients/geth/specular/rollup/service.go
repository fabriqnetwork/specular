package rollup

import (
	"bytes"
	"context"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/external"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/node"
	"github.com/specularl2/specular/clients/geth/specular/proof"
	"github.com/specularl2/specular/clients/geth/specular/rollup/rpc/bridge"
	"github.com/specularl2/specular/clients/geth/specular/rollup/rpc/client"
	"github.com/specularl2/specular/clients/geth/specular/rollup/rpc/txmgr"
	"github.com/specularl2/specular/clients/geth/specular/rollup/services"
	"github.com/specularl2/specular/clients/geth/specular/rollup/services/derivation"
	"github.com/specularl2/specular/clients/geth/specular/rollup/services/derivation/stage"
	"github.com/specularl2/specular/clients/geth/specular/rollup/services/sequencer"
	"github.com/specularl2/specular/clients/geth/specular/rollup/services/validator"
	"github.com/specularl2/specular/clients/geth/specular/rollup/utils/fmt"
	"github.com/specularl2/specular/clients/geth/specular/rollup/utils/log"
)

// TODO: this is Geth-specific; generalize system initialization.
type Node interface {
	RegisterLifecycle(lifecycle node.Lifecycle)
	AccountManager() *accounts.Manager
}

// RegisterRollupService registers rollup service configured by ctx
// Either a sequncer service or a validator service will be registered
func RegisterRollupServices(
	stack Node,
	execBackend services.ExecutionBackend,
	proofBackend proof.Backend,
	cfg *services.SystemConfig,
) error {
	ctx := context.Background()
	rollupState, err := createRollupState(ctx, cfg)
	if err != nil {
		return fmt.Errorf("Failed to create rollup state: %w", err)
	}
	l1State := client.NewL1State()
	// Register driver
	driver, err := createDriver(ctx, cfg, execBackend, rollupState, l1State)
	if err != nil {
		return fmt.Errorf("Failed to create driver: %w", err)
	}
	stack.RegisterLifecycle(driver)
	// Register sequencer
	if (cfg.Sequencer().AccountAddr() != common.Address{}) {
		sequencer, err := createSequencer(ctx, cfg, stack.AccountManager(), execBackend)
		if err != nil {
			return fmt.Errorf("Failed to create sequencer: %w", err)
		}
		stack.RegisterLifecycle(sequencer)
	}
	// Register validator
	if (cfg.Validator().AccountAddr() != common.Address{}) {
		validator, err := createValidator(ctx, cfg, stack.AccountManager(), rollupState, proofBackend)
		if err != nil {
			return fmt.Errorf("Failed to create validator: %w", err)
		}
		stack.RegisterLifecycle(validator)
	}
	return nil
}

func createDriver(
	ctx context.Context,
	cfg *services.SystemConfig,
	execBackend services.ExecutionBackend,
	rollupState *derivation.RollupState,
	l1State stage.L1State,
) (*derivation.Driver, error) {
	l1Client, err := client.DialWithRetry(ctx, cfg.L1().Endpoint(), nil)
	if err != nil {
		return nil, fmt.Errorf("Failed to initialize l1 client: %w", err)
	}
	l2ClientCreatorFn := func(ctx context.Context) (stage.EthClient, error) {
		return client.DialWithRetry(ctx, cfg.L2().Endpoint(), nil)
	}
	terminalStage := stage.CreatePipeline(cfg.L1(), execBackend, rollupState, l2ClientCreatorFn, l1Client, l1State)
	driver := derivation.NewDriver(cfg.Driver(), terminalStage)
	return driver, nil
}

func createSequencer(
	ctx context.Context,
	cfg *services.SystemConfig,
	accountMgr *accounts.Manager,
	execBackend services.ExecutionBackend,
) (*sequencer.Sequencer, error) {
	l1TxMgr, err := createTxManager(
		ctx, cfg.L1(), accountMgr, cfg.Sequencer().AccountAddr(), cfg.L2().ClefEndpoint(), cfg.Sequencer().Passphrase(),
	)
	if err != nil {
		return nil, fmt.Errorf("Failed to initialize l1 tx manager: %w", err)
	}
	batchBuilder, err := sequencer.NewBatchBuilder(math.MaxInt64) // TODO: configure max batch size
	if err != nil {
		return nil, fmt.Errorf("Failed to initialize batch builder: %w", err)
	}
	l2ClientCreatorFn := func(ctx context.Context) (sequencer.L2Client, error) {
		return client.DialWithRetry(ctx, cfg.L2().Endpoint(), nil)
	}
	return sequencer.NewSequencer(cfg, execBackend, l2ClientCreatorFn, l1TxMgr, batchBuilder), nil
}

func createValidator(
	ctx context.Context,
	cfg *services.SystemConfig,
	accountMgr *accounts.Manager,
	rollupState *derivation.RollupState,
	proofBackend proof.Backend,
) (*validator.Validator, error) {
	l2ClientCreatorFn := func(ctx context.Context) (validator.L2Client, error) {
		return client.DialWithRetry(ctx, cfg.L2().Endpoint(), nil)
	}
	// Create tx manager, used to send transactions to L1.
	l1TxMgr, err := createTxManager(
		ctx, cfg.L1(), accountMgr, cfg.Validator().AccountAddr(), cfg.L2().ClefEndpoint(), cfg.Validator().Passphrase(),
	)
	if err != nil {
		return nil, fmt.Errorf("Failed to create tx manager: %w", err)
	}
	return validator.NewValidator(cfg, l2ClientCreatorFn, l1TxMgr, proofBackend, rollupState), nil
}

func createRollupState(ctx context.Context, cfg *services.SystemConfig) (*derivation.RollupState, error) {
	bridgeClient, err := bridge.DialWithRetry(ctx, cfg.L1())
	if err != nil {
		return nil, fmt.Errorf("Failed to initialize l1 bridge client: %w", err)
	}
	rollupState := derivation.NewRollupState(bridgeClient)
	if err != nil {
		return nil, fmt.Errorf("Failed to initialize assertion manager: %w", err)
	}
	return rollupState, nil
}

// func createSyncer(
// 	ctx context.Context,
// 	cfg *services.SystemConfig,
// 	accountMgr *accounts.Manager,
// 	execBackend services.ExecutionBackend,
// ) (*derivation.Syncer, error) {
// l1Client, err := client.DialWithRetry(ctx, cfg.L1Endpoint, nil)
// if err != nil {
// 	return nil, fmt.Errorf("Failed to connect to L1 node: %w", err)
// }
// l2Client, err := client.DialWithRetry(ctx, cfg.L2Endpoint, nil)
// if err != nil {
// 	return nil, fmt.Errorf("Failed to connect to L2 node: %w", err)
// }
// rollupState.StartSync(ctx, l1Client, l2Client)
////////////////////////////////////////////////////////////
// transactor, err := createTransactor(
// 	accountMgr, cfg.ValidatorAccountAddr, cfg.L2ClefEndpoint, cfg.ValidatorPassphrase, cfg.L1ChainID,
// )
// if err != nil {
// 	return nil, fmt.Errorf("Failed to initialize transactor: %w", err)
// }
// l1BridgeClient, err := client.NewEthBridgeClient(
// 	ctx, nil, cfg.L1Endpoint, cfg.L1RollupGenesisBlock, cfg.SequencerInboxAddr, cfg.RollupAddr, transactor, nil,
// )
// if err != nil {
// 	return nil, fmt.Errorf("Failed to initialize l1 bridge client: %w", err)
// }
// return derivation.NewSyncer(execBackend, l1BridgeClient, rollupState.L1Syncer), nil
// }

func createTxManager(
	ctx context.Context,
	cfg *services.L1Config,
	accountMgr *accounts.Manager,
	accountAddr common.Address,
	clefEndpoint string,
	passphrase string,
) (*bridge.TxManager, error) {
	transactor, err := createTransactor(accountMgr, accountAddr, clefEndpoint, passphrase, cfg.ChainID())
	if err != nil {
		return nil, fmt.Errorf("Failed to initialize transactor for tx manager: %w", err)
	}
	l1Client, err := client.DialWithRetry(ctx, cfg.Endpoint(), nil)
	if err != nil {
		return nil, fmt.Errorf("Failed to initialize l1 client for tx manager: %w", err)
	}
	signer := func(ctx context.Context, address common.Address, tx *types.Transaction) (*types.Transaction, error) {
		return transactor.Signer(address, tx)
	}
	// TODO: config
	return bridge.NewTxManager(
		txmgr.NewTxManager(txmgr.DefaultConfig(transactor.From), l1Client, signer),
		cfg,
	)
}

// createTransactor creates a transactor for the given account address,
// either using the clef endpoint or passphrase.
func createTransactor(
	mgr *accounts.Manager,
	accountAddress common.Address,
	clefEndpoint string,
	passphrase string,
	chainID uint64,
) (*bind.TransactOpts, error) {
	if clefEndpoint != "" {
		clef, err := external.NewExternalSigner(clefEndpoint)
		if err != nil {
			return nil, fmt.Errorf("Failed to create external signer from clef endpoint: %w", err)
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
