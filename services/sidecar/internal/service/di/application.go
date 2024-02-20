package di

import (
	"context"
	"os"
	"os/signal"

	"github.com/specularL2/specular/services/sidecar/utils/fmt"
	"github.com/specularL2/specular/services/sidecar/utils/log"

	"golang.org/x/sync/errgroup"

	"github.com/specularL2/specular/services/sidecar/internal/service/config"
	"github.com/specularL2/specular/services/sidecar/rollup/rpc/eth"
	"github.com/specularL2/specular/services/sidecar/rollup/services"
	"github.com/specularL2/specular/services/sidecar/rollup/services/disseminator"
	"github.com/specularL2/specular/services/sidecar/rollup/services/validator"
)

type WaitGroup interface {
	Add(int)
	Done()
	Wait()
}

type Application struct {
	ctx               context.Context
	log               log.Logger
	config            *config.Config
	systemConfig      *services.SystemConfig
	l1State           *eth.EthState
	batchDisseminator *disseminator.BatchDisseminator
	validator         *validator.Validator
}

func (app *Application) Run() error {
	var _, cancel = context.WithCancel(app.ctx)
	var err error

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	defer func() {
		signal.Stop(c)
		cancel()
	}()

	errGroup, _ := errgroup.WithContext(app.ctx)

	if app.systemConfig.Disseminator().GetIsEnabled() {
		app.log.Info("Starting disseminator...")
		err := app.batchDisseminator.Start(app.ctx, errGroup)
		if err != nil {
			return err
		}
	}

	if app.systemConfig.Validator().GetIsEnabled() {
		app.log.Info("Starting validator...")
		err := app.validator.Start(app.ctx, errGroup)
		if err != nil {
			return err
		}
	}

	if err := errGroup.Wait(); err != nil {
		return fmt.Errorf("service failed while running: %w", err)
	}
	app.log.Info("app stopped")

	return err
}

func (app *Application) ShutdownAndCleanup(exitCode int) {
	if exitCode == 0 {
		app.log.Info("app shutting down")
	} else {
		app.log.Crit("app shutting down due to error, exit code: %i", exitCode)
	}
	os.Exit(exitCode)
}

func (app *Application) GetLogger() log.Logger {
	return app.log
}

func (app *Application) GetContext() context.Context {
	return app.ctx
}

func (app *Application) GetConfig() *config.Config {
	return app.config
}

type TestApplication struct {
	*Application

	Ctx    context.Context
	log    log.Logger
	Config *config.Config
}
