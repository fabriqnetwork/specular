package main

import (
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/specularL2/specular/services/sidecar/internal/service/di"
	"log"
	"os"
)

func main() {
	application, _, err := di.SetupApplication()
	if err != nil {
		log.Fatalf("failed to setup application #{err}")
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
