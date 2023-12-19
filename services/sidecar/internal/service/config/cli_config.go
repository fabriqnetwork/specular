package config

import (
	"os"

	"github.com/urfave/cli/v2"

	"github.com/specularL2/specular/services/sidecar/rollup/services"
)

type CLIExtractor struct {
	systemConfig *services.SystemConfig
}

func (t *CLIExtractor) ExtractFromCLIContext(cliCtx *cli.Context) error {
	var err error
	t.systemConfig, err = services.ParseSystemConfig(cliCtx)
	if err != nil {
		return err
	}
	return nil
}

// Nasty trick to extract parsed SystemConfig from the urfave/cli package wrapper and serve properly from a provider
func NewSystemConfig(cfg *Config) (*services.SystemConfig, error) {
	cliExtractor := &CLIExtractor{}

	cliApp := &cli.App{
		Name:   cfg.ServiceName,
		Usage:  cfg.UsageDesc,
		Action: cliExtractor.ExtractFromCLIContext,
	}
	cliApp.Flags = services.CLIFlags()

	if err := cliApp.Run(os.Args); err != nil {
		return nil, err
	}

	return cliExtractor.systemConfig, nil
}
