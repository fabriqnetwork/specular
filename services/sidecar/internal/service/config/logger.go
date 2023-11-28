package config

import (
	"os"

	"github.com/specularL2/specular/services/sidecar/rollup/services"
	"github.com/specularL2/specular/services/sidecar/utils/log"
)

func NewLogger(systemCfg *services.SystemConfig) log.Logger {
	baseLog := log.New()
	gLogger := log.NewGlogHandler(log.StreamHandler(os.Stderr, log.TerminalFormat(false)))
	gLogger.Verbosity(systemCfg.Verbosity)
	baseLog.SetHandler(gLogger)
	// TODO: refactor is required to the "old" code to pass down the logger instance from the provider,
	//		 for now let's keep the global setting of the handler.
	log.Root().SetHandler(gLogger)
	return baseLog
}
