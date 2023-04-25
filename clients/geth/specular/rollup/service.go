package rollup

import (
	"bytes"
	"context"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts"
	bind "github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/external"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/node"
	"github.com/specularl2/specular/clients/geth/specular/proof"
	"github.com/specularl2/specular/clients/geth/specular/rollup/comms/client"
	"github.com/specularl2/specular/clients/geth/specular/rollup/comms/txmgr"
	"github.com/specularl2/specular/clients/geth/specular/rollup/services"
	"github.com/specularl2/specular/clients/geth/specular/rollup/services/indexer"
	"github.com/specularl2/specular/clients/geth/specular/rollup/services/sequencer"
	"github.com/specularl2/specular/clients/geth/specular/rollup/services/state"
	"github.com/specularl2/specular/clients/geth/specular/rollup/services/validator"
	"github.com/specularl2/specular/clients/geth/specular/rollup/types/assertion"
	"github.com/specularl2/specular/clients/geth/specular/rollup/types/data"
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
	syncer, err := createSyncer(ctx, cfg, stack.AccountManager(), execBackend)
	if err != nil {
		return fmt.Errorf("Failed to create syncer: %w", err)
	}
	// Register services
	if (cfg.SequencerAccountAddr != common.Address{}) {
		// TODO: fix.
		_, err := syncer.SyncL2ChainToL1Head(ctx, cfg.L1RollupGenesisBlock)
		if err != nil {
			return fmt.Errorf("Failed to sync l2 chain to l1 head: %w", err)
		}
		sequencer, err := createSequencer(ctx, cfg, stack.AccountManager(), execBackend)
		if err != nil {
			return fmt.Errorf("Failed to create sequencer: %w", err)
		}
		stack.RegisterLifecycle(sequencer)
	}
	if (cfg.ValidatorAccountAddr != common.Address{}) {
		// TODO: fix.
		_, err := syncer.SyncL2ChainToL1Head(ctx, cfg.L1RollupGenesisBlock)
		if err != nil {
			return fmt.Errorf("Failed to sync l2 chain to l1 head: %w", err)
		}
		validator, err := createValidator(ctx, cfg, stack.AccountManager(), proofBackend)
		if err != nil {
			return fmt.Errorf("Failed to create validator: %w", err)
		}
		stack.RegisterLifecycle(validator)
	}
	if (cfg.IndexerAccountAddr != common.Address{}) {
		stack.RegisterLifecycle(indexer.NewIndexer(cfg, syncer))
	}
	return nil
}

func createSequencer(
	ctx context.Context,
	cfg *services.SystemConfig,
	accountMgr *accounts.Manager,
	execBackend services.ExecutionBackend,
) (*sequencer.Sequencer, error) {
	l1TxMgr, err := createTxManager(
		ctx, accountMgr, cfg.L1ChainID, cfg.L1Endpoint, cfg.SequencerAccountAddr, cfg.L2ClefEndpoint, cfg.SequencerPassphrase,
	)
	if err != nil {
		return nil, fmt.Errorf("Failed to initialize l1 tx manager: %w", err)
	}
	batchBuilder, err := data.NewBatchBuilder(math.MaxInt64) // TODO: configure max batch size
	if err != nil {
		return nil, fmt.Errorf("Failed to initialize batch builder: %w", err)
	}
	l2Client, err := client.DialWithRetry(ctx, "localhost", client.DefaultRetryOpts)
	if err != nil {
		return nil, fmt.Errorf("Failed to initialize l2 client: %w", err)
	}
	return sequencer.NewSequencer(cfg, execBackend, l2Client, l1TxMgr, batchBuilder), nil
}

func createValidator(
	ctx context.Context,
	cfg *services.SystemConfig,
	accountMgr *accounts.Manager,
	proofBackend proof.Backend,
) (*validator.Validator, error) {
	transactor, err := createTransactor(accountMgr, cfg.ValidatorAccountAddr, cfg.L2ClefEndpoint, cfg.ValidatorPassphrase, cfg.L1ChainID)
	if err != nil {
		return nil, fmt.Errorf("Failed to initialize transactor: %w", err)
	}
	l1BridgeClient, err := client.NewEthBridgeClient(
		ctx, nil, cfg.L1Endpoint, cfg.L1RollupGenesisBlock, cfg.SequencerInboxAddr, cfg.RollupAddr, transactor, client.DefaultRetryOpts,
	)
	if err != nil {
		return nil, fmt.Errorf("Failed to initialize l1 bridge client: %w", err)
	}
	l1TxMgr, err := createTxManager(
		ctx, accountMgr, cfg.L1ChainID, cfg.L1Endpoint, cfg.ValidatorAccountAddr, cfg.L2ClefEndpoint, cfg.ValidatorPassphrase,
	)
	if err != nil {
		return nil, fmt.Errorf("Failed to create tx manager: %w", err)
	}
	assertionMgr, err := assertion.NewAssertionManager(l1BridgeClient)
	if err != nil {
		return nil, fmt.Errorf("Failed to initialize assertion manager: %w", err)
	}
	l2Client, err := client.DialWithRetry(ctx, "localhost", client.DefaultRetryOpts)
	if err != nil {
		return nil, fmt.Errorf("Failed to initialize l2 client: %w", err)
	}
	return validator.NewValidator(cfg, l2Client, l1TxMgr, l1BridgeClient, proofBackend, assertionMgr), nil
}

func createSyncer(
	ctx context.Context,
	cfg *services.SystemConfig,
	accountMgr *accounts.Manager,
	execBackend services.ExecutionBackend,
) (*services.Syncer, error) {
	rollupState := state.NewRollupState()
	l1Client, err := client.DialWithRetry(ctx, cfg.L1Endpoint, client.DefaultRetryOpts)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to L1 node: %w", err)
	}
	l2Client, err := client.DialWithRetry(ctx, "localhost", client.DefaultRetryOpts)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to L2 node: %w", err)
	}
	rollupState.StartSync(ctx, l1Client, l2Client)
	transactor, err := createTransactor(
		accountMgr, cfg.ValidatorAccountAddr, cfg.L2ClefEndpoint, cfg.ValidatorPassphrase, cfg.L1ChainID,
	)
	if err != nil {
		return nil, fmt.Errorf("Failed to initialize transactor: %w", err)
	}
	l1BridgeClient, err := client.NewEthBridgeClient(
		ctx, nil, cfg.L1Endpoint, cfg.L1RollupGenesisBlock, cfg.SequencerInboxAddr, cfg.RollupAddr, transactor, client.DefaultRetryOpts,
	)
	if err != nil {
		return nil, fmt.Errorf("Failed to initialize l1 bridge client: %w", err)
	}
	return services.NewSyncer(execBackend, l1BridgeClient, rollupState.L1Syncer), nil
}

func createTxManager(
	ctx context.Context,
	accountMgr *accounts.Manager,
	l1ChainID uint64,
	l1Endpoint string,
	accountAddr common.Address,
	clefEndpoint string,
	passphrase string,
) (*txmgr.TxManager, error) {
	transactor, err := createTransactor(accountMgr, accountAddr, clefEndpoint, passphrase, l1ChainID)
	if err != nil {
		return nil, fmt.Errorf("Failed to initialize transactor: %w", err)
	}
	l1Client, err := client.DialWithRetry(ctx, l1Endpoint, client.DefaultRetryOpts)
	if err != nil {
		return nil, fmt.Errorf("Failed to initialize l1 client: %w", err)
	}
	return txmgr.NewTxManager(txmgr.DefaultConfig(transactor.From), l1Client, transactor.Signer), nil
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
