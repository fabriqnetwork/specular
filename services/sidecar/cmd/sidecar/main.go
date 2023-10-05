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

func startService(cliCtx *cli.Context) error {

	cfg := ParseSystemConfig(cliCtx)
	acctMgr := _

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
