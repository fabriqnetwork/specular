package main

import (
	"github.com/sirupsen/logrus"
	"github.com/specularL2/specular/services/sidecar/internal/service/di"
	"github.com/specularL2/specular/services/sidecar/rollup/services"
	"os"
)

func main() {
	app, _, err := di.SetupApplication()
	if err != nil {
		logrus.Fatalf("failed to setup application: %s", err)
		os.Exit(1)
	}

	app.GetCli().Flags = services.CLIFlags()
	app.GetCli().Action = app.Run

	if err := app.GetCli().Run(os.Args); err != nil {
		app.GetLogger().WithError(err).Log(logrus.FatalLevel, "application failed")
		os.Exit(1)
	}
}
