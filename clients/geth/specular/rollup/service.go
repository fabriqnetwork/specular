package rollup

import (
	"bytes"
	"context"
	"math"
	"math/big"
	"reflect"

	"github.com/ethereum/go-ethereum/accounts"
	bind "github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/external"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"github.com/specularl2/specular/clients/geth/specular/proof"
	"github.com/specularl2/specular/clients/geth/specular/rollup/client"
	"github.com/specularl2/specular/clients/geth/specular/rollup/derivation"
	"github.com/specularl2/specular/clients/geth/specular/rollup/rpc/bridge"
	"github.com/specularl2/specular/clients/geth/specular/rollup/rpc/eth"
	"github.com/specularl2/specular/clients/geth/specular/rollup/rpc/eth/txmgr"
	"github.com/specularl2/specular/clients/geth/specular/rollup/services"
	"github.com/specularl2/specular/clients/geth/specular/rollup/services/api"
	"github.com/specularl2/specular/clients/geth/specular/rollup/services/disseminator"
	"github.com/specularl2/specular/clients/geth/specular/rollup/services/indexer"
	"github.com/specularl2/specular/clients/geth/specular/rollup/services/sequencer"
	"github.com/specularl2/specular/clients/geth/specular/utils/fmt"
)

// TODO: this is the last Geth-specific interface here; remove.
type accountManager interface {
	Backends(reflect.Type) []accounts.Backend
}

// Creates services configured by cfg:
// - Sequencer (if sequencer configured)
// - L2 block data disseminator (if sequencer configured)
// - Validator (if validator configured)
func CreateRollupServices(
	accMgr accountManager,
	execBackend api.ExecutionBackend,
	proofBackend proof.Backend,
	cfg *services.SystemConfig,
) ([]api.Service, error) {
	var services []api.Service
	legacyService, err := createLegacyService(accMgr, execBackend, proofBackend, cfg)
	if err != nil {
		return nil, err
	}
	services = append(services, legacyService)
	if (cfg.Sequencer().GetAccountAddr() != common.Address{}) {
		disseminator, err := createDisseminator(context.Background(), cfg, accMgr)
		if err != nil {
			return nil, fmt.Errorf("Failed to initialize sequencer: %w", err)
		}
		services = append(services, disseminator)
	}
	return services, nil
}

// TODO: delete.
func createLegacyService(
	accMgr accountManager,
	execBackend api.ExecutionBackend,
	proofBackend proof.Backend,
	cfg *services.SystemConfig,
) (api.Service, error) {
	transactor, err := createTransactor(
		accMgr, cfg.GetSequencerInboxAddr(), cfg.GetClefEndpoint(), cfg.GetPassphrase(), cfg.L1().GetChainID(),
	)
	if err != nil {
		return nil, err
	}
	l1Client, err := client.NewEthBridgeClient(
		context.Background(),
		cfg.L1().Endpoint,
		cfg.L1().RollupGenesisBlock,
		cfg.SequencerInboxAddr,
		cfg.RollupAddr,
		transactor,
		eth.DefaultRetryOpts,
	)
	if err != nil {
		return nil, err
	}
	// Register service
	var service api.Service
	if cfg.Sequencer().GetAccountAddr() != (common.Address{}) {
		service, err = sequencer.New(execBackend, proofBackend, l1Client, cfg)
	} else {
		service, err = indexer.New(execBackend, proofBackend, l1Client, cfg)
	}
	return service, err
}

func createDisseminator(
	ctx context.Context,
	cfg *services.SystemConfig,
	accountMgr accountManager,
) (*disseminator.BatchDisseminator, error) {
	var (
		seqCfg       = cfg.Sequencer()
		l1TxMgr, err = createTxManager(
			ctx, cfg, accountMgr, seqCfg.GetAccountAddr(), seqCfg.GetClefEndpoint(), seqCfg.GetPassphrase(),
		)
	)
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

func createTxManager(
	ctx context.Context,
	cfg *services.SystemConfig,
	accountMgr accountManager,
	accountAddr common.Address,
	clefEndpoint string,
	passphrase string,
) (*bridge.TxManager, error) {
	transactor, err := createTransactor(accountMgr, accountAddr, clefEndpoint, passphrase, cfg.L1().GetChainID())
	if err != nil {
		return nil, fmt.Errorf("failed to initialize transactor: %w", err)
	}
	l1Client, err := eth.DialWithRetry(ctx, cfg.L1().GetEndpoint())
	if err != nil {
		return nil, fmt.Errorf("failed to initialize l1 client: %w", err)
	}
	signer := func(ctx context.Context, address common.Address, tx *ethTypes.Transaction) (*ethTypes.Transaction, error) {
		return transactor.Signer(address, tx)
	}
	// TODO: config
	return bridge.NewTxManager(txmgr.NewTxManager(cfg.Sequencer().GetTxMgrCfg(), l1Client, signer), cfg.L1())
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
		return nil, fmt.Errorf("failed to export account: %w", err)
	}
	transactor, err := bind.NewTransactorWithChainID(bytes.NewReader(json), passphrase, new(big.Int).SetUint64(chainID))
	if err != nil {
		return nil, fmt.Errorf("failed to create transactor: %w", err)
	}
	return transactor, nil
}
