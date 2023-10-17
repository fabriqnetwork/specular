package rollup

import (
	"bytes"
	"context"
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
	"github.com/specularl2/specular/services/cl_clients/ripcord/rollup/rpc/bridge"
	"github.com/specularl2/specular/services/cl_clients/ripcord/rollup/rpc/eth"
	"github.com/specularl2/specular/services/cl_clients/ripcord/rollup/rpc/eth/txmgr"
	"github.com/specularl2/specular/services/cl_clients/ripcord/rollup/services"
	"github.com/specularl2/specular/services/cl_clients/ripcord/rollup/services/api"
	"github.com/specularl2/specular/services/cl_clients/ripcord/rollup/services/indexer"
	"github.com/specularl2/specular/services/cl_clients/ripcord/rollup/services/sequencer"
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

func CreateSequencer(
	accMgr accountManager,
	execBackend api.ExecutionBackend,
	proofBackend proof.Backend,
	cfg *services.SystemConfig,
) (api.Service, error) {
	var services []api.Service
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
	if err != nil {
		return nil, err
	}
	return service, nil
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
