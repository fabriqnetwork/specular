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
	"github.com/specularl2/specular/services/cl_clients/ripcord/proof"
	"github.com/specularl2/specular/services/cl_clients/ripcord/rollup/client"
	"github.com/specularl2/specular/services/cl_clients/ripcord/rollup/derivation"
	"github.com/specularl2/specular/services/cl_clients/ripcord/rollup/rpc/bridge"
	"github.com/specularl2/specular/services/cl_clients/ripcord/rollup/rpc/eth"
	"github.com/specularl2/specular/services/cl_clients/ripcord/rollup/rpc/eth/txmgr"
	"github.com/specularl2/specular/services/cl_clients/ripcord/rollup/services"
	"github.com/specularl2/specular/services/cl_clients/ripcord/rollup/services/api"
	"github.com/specularl2/specular/services/cl_clients/ripcord/rollup/services/disseminator"
	"github.com/specularl2/specular/services/cl_clients/ripcord/rollup/services/indexer"
	"github.com/specularl2/specular/services/cl_clients/ripcord/rollup/services/sequencer"
	"github.com/specularl2/specular/services/cl_clients/ripcord/rollup/services/validator"
	"github.com/specularl2/specular/services/cl_clients/ripcord/utils/fmt"
	"github.com/specularl2/specular/services/cl_clients/ripcord/utils/log"
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
	if cfg.Sequencer().GetIsEnabled() {
		disseminator, err := createDisseminator(context.Background(), cfg, accMgr)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize sequencer: %w", err)
		}
		services = append(services, disseminator)
	}
	if cfg.Validator().GetIsEnabled() {
		validator, err := createValidator(context.Background(), cfg, accMgr)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize validator: %w", err)
		}
		services = append(services, validator)
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
	l1Client, err := client.NewEthBridgeClient(
		context.Background(),
		cfg.L1().Endpoint,
		cfg.SequencerInboxAddr,
		eth.DefaultRetryOpts,
	)
	if err != nil {
		return nil, err
	}
	// Register service
	var (
		service    api.Service
		serviceCfg = struct {
			services.SequencerConfig
			services.L1Config
		}{cfg.Sequencer(), cfg.L1()}
	)
	if cfg.Sequencer().GetIsEnabled() {
		service, err = sequencer.New(execBackend, proofBackend, l1Client, serviceCfg)
	} else {
		service, err = indexer.New(execBackend, proofBackend, l1Client, serviceCfg)
	}
	return service, err
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
