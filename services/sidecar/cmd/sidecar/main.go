package main

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/specularL2/specular/services/sidecar/internal/service/di"
)

func main() {
	application, _, err := di.SetupApplication()
	if err != nil {
		log.Fatalf("failed to setup application %s", err)
	}

	exitCode := 0
	defer func() { os.Exit(exitCode) }()

	go func() {
		<-application.GetContext().Done()

		application.GetLogger().Info("application context canceled, cleaning up")
		application.ShutdownAndCleanup()
	}()

	if err := application.Run(); err != nil {
		if !errors.Is(err, context.Canceled) {
			application.GetLogger().WithError(err).Log(logrus.FatalLevel, "application failed")
			exitCode = 1
		}
	}
}
