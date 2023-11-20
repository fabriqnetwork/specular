package config

import (
	"context"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

type CancelChannel chan struct{}

func NewCancelChannel() CancelChannel {
	return make(chan struct{}, 1)
}

func NewContext(log *logrus.Logger, termination CancelChannel) context.Context {
	ctx, cancel := context.WithCancel(context.Background())

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for {
			select {
			case sig := <-quit:
				log.WithField("signal", sig).Info("os signal - shutting down")
				cancel()
				return
			case <-termination:
				log.Info("term signal - shutting down")
				cancel()
				return
			}
		}
	}()

	return ctx
}
