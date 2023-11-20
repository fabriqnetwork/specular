package di

import (
	"context"
	"github.com/specularL2/specular/services/sidecar/internal/service/config"
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
	ctx    context.Context
	log    *logrus.Logger
	config *config.Config
}

func (app *Application) Run() error {
	var _, cancel = context.WithCancel(app.ctx)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	defer func() {
		signal.Stop(c)
		cancel()
	}()

	errGroup, _ := errgroup.WithContext(app.ctx)

	err := errGroup.Wait()
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

type TestApplication struct {
	*Application

	Ctx    context.Context
	Log    *logrus.Logger
	Config *config.Config
}
