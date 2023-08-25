package rollup

import (
	"bytes"
	"context"
	"math/big"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/ethereum/go-ethereum/accounts"
	bind "github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/external"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/node"
	"github.com/specularl2/specular/clients/geth/specular/proof"
	"github.com/specularl2/specular/clients/geth/specular/rollup/client"
	"github.com/specularl2/specular/clients/geth/specular/rollup/services"
	"github.com/specularl2/specular/clients/geth/specular/rollup/services/indexer"
	"github.com/specularl2/specular/clients/geth/specular/rollup/services/sequencer"
	"github.com/specularl2/specular/clients/geth/specular/rollup/utils/fmt"
)

// Tries to connect to Clef as a signer and if not reverts to using Geth as the signer
func GetAuth(stack *node.Node, cfg *services.Config) *bind.TransactOpts {
	var ks *keystore.KeyStore
	if keystores := stack.AccountManager().Backends(keystore.KeyStoreType); len(keystores) > 0 {
		ks = keystores[0].(*keystore.KeyStore)
	}
	if ks == nil {
		log.Crit("Failed to register the Rollup service: keystore not found")
	}

	chainID := big.NewInt(int64(cfg.L1ChainID))
	json, err := ks.Export(accounts.Account{Address: cfg.Coinbase}, cfg.Passphrase, "")
	if err != nil {
		log.Crit("Failed to register the Rollup service", "err", err)
	}

	var auth *bind.TransactOpts
	if cfg.ClefEndpoint != "" {
		clef, err := external.NewExternalSigner(cfg.ClefEndpoint)
		if err != nil {
			log.Crit("Failed to create external signer from clef endpoint", "err", err)
		}
		auth = bind.NewClefTransactor(clef, accounts.Account{Address: cfg.Coinbase})
	} else {
		log.Warn("no external signer specified, using geth signer")
		auth, err = bind.NewTransactorWithChainID(bytes.NewReader(json), cfg.Passphrase, chainID)
		if err != nil {
			log.Crit("Failed to register the Rollup service", "err", err)
		}
	}

	return auth
}

// Constructs an EthBridgeClient for communicating with L1
func GetEthBridgeClient(stack *node.Node, cfg *services.Config) *client.EthBridgeClient {

	auth := GetAuth(stack, cfg)

	ctx := context.Background()

	retryOpts := []retry.Option{
		retry.Context(ctx),
		retry.Attempts(3),
		retry.Delay(5 * time.Second),
		retry.LastErrorOnly(true),
		retry.RetryIf(func(err error) bool {
			return true
		}),
		retry.OnRetry(func(n uint, err error) {
			log.Error("Failed attempt", "attempt", n, "err", err)
		}),
	}

	l1Client, err := client.NewEthBridgeClient(
		ctx,
		cfg.L1Endpoint,
		cfg.L1RollupGenesisBlock,
		cfg.SequencerInboxAddr,
		cfg.RollupAddr,
		auth,
		retryOpts,
	)
	if err != nil {
		log.Crit("Failed to register the Rollup service: cannot create l1 client", "err", err)
	}

	return l1Client

}

// RegisterRollupService registers rollup service configured by ctx
func RegisterRollupService(stack *node.Node, eth services.Backend, proofBackend proof.Backend, cfg *services.Config) {

	l1Client := GetEthBridgeClient(stack, cfg)

	var service node.Lifecycle
	var err error
	
  switch cfg.Node {
	case services.NODE_SEQUENCER:
		service, err = sequencer.New(eth, proofBackend, l1Client, cfg)
	case services.NODE_INDEXER:
		service, err = indexer.New(eth, proofBackend, l1Client, cfg)
	default:
		err = fmt.Errorf("Node type unkown: %v", cfg.Node)
	}
	if err != nil {
		log.Crit("Failed to register rollup service", "err", err)
	}
	stack.RegisterLifecycle(service)
}
