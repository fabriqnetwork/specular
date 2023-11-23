package di

import (
	"context"
	"github.com/cockroachdb/errors"
	"github.com/specularL2/specular/services/sidecar/internal/service/config"
	"github.com/specularL2/specular/services/sidecar/rollup/services"
	"github.com/urfave/cli/v2"
	"golang.org/x/sync/errgroup"
	"os"
	"os/signal"

	"github.com/sirupsen/logrus"
)

type WaitGroup interface {
	Add(int)
	Done()
	Wait()
}

type Application struct {
	cli    *cli.App
	ctx    context.Context
	log    *logrus.Logger
	config *config.Config

	// TODO: extract providers for
	// systemConfig *services.SystemConfig
	// disseminator *disseminator.BatchDisseminator
	// validator    *validator.Validator
}

func (app *Application) Run(ctx *cli.Context) error {
	var _, cancel = context.WithCancel(app.ctx)
	var err error

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	defer func() {
		signal.Stop(c)
		cancel()
	}()

	errGroup, _ := errgroup.WithContext(app.ctx)

	systemConfig, err := services.ParseSystemConfig(ctx)
	// TODO: use further on
	_ = systemConfig
	if err != nil {
		return errors.Newf("failed to parse config: %w", err)
	}

	app.log.Info("Starting L1 sync...")
	// TODO: refactor main_old.createL1State() function into a infra/service with a provider

	if systemConfig.Disseminator().GetIsEnabled() {
		app.log.Info("Starting disseminator...")
		// TODO: refactor main_old.createDisseminator() and start() functions into a infra/service with a provider
		//		 incl. dependencies
	}

	if systemConfig.Validator().GetIsEnabled() {
		app.log.Info("Starting validator...")
		// TODO: refactor main_old.createValidator() and start() functions into a infra/service with a provider
		//		 incl. dependencies
	}

	if err := errGroup.Wait(); err != nil {
		return errors.Newf("service failed while running: %w", err)
	}
	app.log.Info("app stopped")

	return err
}

func (app *Application) ShutdownAndCleanup() {
	app.log.Info("app shutting down")
}

func (app *Application) GetLogger() *logrus.Logger {
	return app.log
}

func (app *Application) GetContext() context.Context {
	return app.ctx
}

func (app *Application) GetConfig() *config.Config {
	return app.config
}

func (app *Application) GetCli() *cli.App {
	return app.cli
}

type TestApplication struct {
	*Application

	Ctx    context.Context
	Log    *logrus.Logger
	Config *config.Config
}
