package api

import (
	"github.com/specularL2/specular/services/sidecar/internal/service/config"
	"github.com/specularL2/specular/services/sidecar/rollup/services"
	"github.com/urfave/cli/v2"
)

func NewCli(cfg *config.Config) (*cli.App, error) {
	cliApp := &cli.App{
		Name:  cfg.ServiceName,
		Usage: cfg.UsageDesc,
	}
	cliApp.Flags = services.CLIFlags()
	return cliApp, nil
}
