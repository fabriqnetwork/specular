package rollup

import (
	"bytes"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts"
	bind "github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/node"
	"github.com/specularl2/specular/clients/geth/specular/proof"
	"github.com/specularl2/specular/clients/geth/specular/rollup/services"
	"github.com/specularl2/specular/clients/geth/specular/rollup/services/sequencer"
	"github.com/specularl2/specular/clients/geth/specular/rollup/services/validator"
)

// RegisterRollupService registers rollup service configured by ctx
// Either a sequncer service or a validator service will be registered
func RegisterRollupService(stack *node.Node, eth services.Backend, proofBackend proof.Backend, cfg *services.Config) {
	// Unlock account for L1 transaction signer
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
	auth, err := bind.NewTransactorWithChainID(bytes.NewReader(json), cfg.Passphrase, chainID)
	if err != nil {
		log.Crit("Failed to register the Rollup service", "err", err)
	}

	// Register services
	if cfg.Node == services.NODE_SEQUENCER {
		sequencer.RegisterService(stack, eth, proofBackend, cfg, auth)
	} else if cfg.Node == services.NODE_VALIDATOR {
		validator.RegisterService(stack, eth, proofBackend, cfg, auth)
	} else if cfg.Node == services.NODE_INDEXER {
		validator.RegisterService(stack, eth, proofBackend, cfg, auth)
	} else {
		log.Crit("Failed to register the Rollup service: Node type unkown", "type", cfg.Node)
	}
}
