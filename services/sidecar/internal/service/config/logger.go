package config

import (
	"github.com/sirupsen/logrus"
	"os"
)

func NewLogger(cfg *Config) *logrus.Logger {
	level := cfg.GetLogLevel(defaultLogLevel)
	log := logrus.New()
	log.SetLevel(level)
	log.SetOutput(os.Stdout)
	log.WithField("name", cfg.ServiceName).Info("service is starting")

	return log
}
