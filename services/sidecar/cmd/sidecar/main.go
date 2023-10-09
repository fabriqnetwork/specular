package main

import (
	"context"
	"os"

	"github.com/urfave/cli/v2"
	"golang.org/x/sync/errgroup"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/specularL2/specular/services/sidecar/rollup"
	"github.com/specularL2/specular/services/sidecar/rollup/services"
	"github.com/specularL2/specular/services/sidecar/rollup/services/disseminator"
	"github.com/specularL2/specular/services/sidecar/rollup/services/validator"
	"github.com/specularL2/specular/services/sidecar/utils/fmt"
	"github.com/specularL2/specular/services/sidecar/utils/log"
)

func initializeAccountManager(cfg services.KeyStoreConfig) *accounts.Manager {
	keyDir := cfg.GetKeyStoreDir()
	acctMgr := accounts.NewManager(&accounts.Config{InsecureUnlockAllowed: false})

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

	acctMgr.AddBackend(keystore.NewKeyStore(keyDir, keystore.StandardScryptN, keystore.StandardScryptP))

	return acctMgr
}

func startService(cliCtx *cli.Context) error {

	log.Info("Reading configuration")

	cfg, err := services.ParseSystemConfig(cliCtx)
	if err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}

	log.Info("Initializing Account Manager")

	// TODO: verify if hard-coded false is a problem
	acctMgr := initializeAccountManager(cfg.KeyStore())

	log.Info("Starting disseminator+validator")

	var disseminator *disseminator.BatchDisseminator
	var validator *validator.Validator

	if cfg.Sequencer().GetIsEnabled() {
		disseminator, err = rollup.CreateDisseminator(context.Background(), cfg, acctMgr)
		if err != nil {
			return fmt.Errorf("failed to initialize disseminator: %w", err)
		}
	}
	if cfg.Validator().GetIsEnabled() {
		validator, err = rollup.CreateValidator(context.Background(), cfg, acctMgr)
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

    if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
    }
}
