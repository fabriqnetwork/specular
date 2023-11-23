package api

import (
	"github.com/specularL2/specular/services/sidecar/internal/service/config"
	"github.com/specularL2/specular/services/sidecar/rollup/services"
	"github.com/urfave/cli/v2"
)

func NewCli(cfg *config.Config) (*cli.App, error) {
	cli := &cli.App{
		Name:  "sidecar",
		Usage: "launch a validator and/or disseminator",
	}
	cli.Flags = services.CLIFlags()
	return cli, nil
}
