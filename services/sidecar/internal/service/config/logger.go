package config

import (
	"os"

	"github.com/sirupsen/logrus"
)

func NewLogger(cfg *Config) *logrus.Logger {
	level := cfg.GetLogLevel(defaultLogLevel)
	log := logrus.New()
	log.SetLevel(level)
	log.SetOutput(os.Stdout)
	log.WithField("name", cfg.ServiceName).Info("service is starting")

	return log
}
