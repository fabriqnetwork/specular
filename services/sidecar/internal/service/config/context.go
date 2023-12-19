package config

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/specularL2/specular/services/sidecar/utils/log"
)

type CancelChannel chan struct{}

func NewCancelChannel() CancelChannel {
	return make(chan struct{}, 1)
}

func NewContext(log log.Logger, termination CancelChannel) context.Context {
	ctx, cancel := context.WithCancel(context.Background())

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for {
			select {
			case sig := <-quit:
				log.Info("os signal - shutting down", "signal", sig)
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
