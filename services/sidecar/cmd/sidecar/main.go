package main

import (
	"context"
	"os"

	"github.com/urfave/cli/v2"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/specularL2/specular/services/sidecar/rollup/services"
	"github.com/specularL2/specular/services/sidecar/rollup/services/disseminator"
	"github.com/specularL2/specular/services/sidecar/rollup/services/validator"
	"github.com/specularL2/specular/services/sidecar/utils/fmt"
	"github.com/specularL2/specular/services/sidecar/utils/log"
)

func initializeAccountManager() *accounts.AccountManager {
	acctMgr := accounts.NewManager(&accounts.Config{InsecureUnlockAllowed: false})

	// < copied from go-ethereum/cmd/config.go:setAccountManagerBackends >
	keydir := stack.KeyStoreDir()
	scryptN := keystore.StandardScryptN
	scryptP := keystore.StandardScryptP
	if conf.UseLightweightKDF {
		scryptN = keystore.LightScryptN
		scryptP = keystore.LightScryptP
	}

	// Assemble the supported backends
	if len(conf.ExternalSigner) > 0 {
		log.Info("Using external signer", "url", conf.ExternalSigner)
		if extapi, err := external.NewExternalBackend(conf.ExternalSigner); err == nil {
			am.AddBackend(extapi)
			return nil
		} else {
			return fmt.Errorf("error connecting to external signer: %v", err)
		}
	}

	// For now, we're using EITHER external signer OR local signers.
	// If/when we implement some form of lockfile for USB and keystore wallets,
	// we can have both, but it's very confusing for the user to see the same
	// accounts in both externally and locally, plus very racey.
	am.AddBackend(keystore.NewKeyStore(keydir, scryptN, scryptP))
	if conf.USB {
		// Start a USB hub for Ledger hardware wallets
		if ledgerhub, err := usbwallet.NewLedgerHub(); err != nil {
			log.Warn(fmt.Sprintf("Failed to start Ledger hub, disabling: %v", err))
		} else {
			am.AddBackend(ledgerhub)
		}
		// Start a USB hub for Trezor hardware wallets (HID version)
		if trezorhub, err := usbwallet.NewTrezorHubWithHID(); err != nil {
			log.Warn(fmt.Sprintf("Failed to start HID Trezor hub, disabling: %v", err))
		} else {
			am.AddBackend(trezorhub)
		}
		// Start a USB hub for Trezor hardware wallets (WebUSB version)
		if trezorhub, err := usbwallet.NewTrezorHubWithWebUSB(); err != nil {
			log.Warn(fmt.Sprintf("Failed to start WebUSB Trezor hub, disabling: %v", err))
		} else {
			am.AddBackend(trezorhub)
		}
	}
	if len(conf.SmartCardDaemonPath) > 0 {
		// Start a smart card hub
		if schub, err := scwallet.NewHub(conf.SmartCardDaemonPath, scwallet.Scheme, keydir); err != nil {
			log.Warn(fmt.Sprintf("Failed to start smart card hub, disabling: %v", err))
		} else {
			am.AddBackend(schub)
		}
	}

	// </ copied from go-ethereum/cmd/config.go:setAccountManagerBackends >


}

func startService(cliCtx *cli.Context) error {

	cfg := ParseSystemConfig(cliCtx)

	// TODO: verify if hard-coded false is a problem
	acctMgr := initializeAccountManager()

	var disseminator *disseminator.BatchDisseminator
	var validator *validator.Validator

	if cfg.Sequencer().GetIsEnabled() {
		disseminator, err := createDisseminator(context.Background(), cfg, acctMgr)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize disseminator: %w", err)
		}
	}
	if cfg.Validator().GetIsEnabled() {
		validator, err := createValidator(context.Background(), cfg, acctMgr)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize validator: %w", err)
		}
	}

	disseminatorErrorGroup, disseminatorCtx := errgroup.WithContext(context.Background()
	disseminator.Start(disseminatorCtx, disseminatorCtx)

	validatorErrorGroup, validatorCtx := errgroup.WithContext(context.Background()
	validator.Start(validatorCtx, validatorErrorGroup)


}

func main() {
    app := &cli.App{
        Name:  "sidecar",
        Usage: "launch validator+disseminator",
        Action: startServices,
    }

    if err := app.Run(os.Args); err != nil {
        log.Fatal(err)
    }
}
