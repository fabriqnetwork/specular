package config

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"strings"
)

const (
	defaultLogLevel = logrus.InfoLevel
	serviceName     = "sidecar"
)

type Config struct {
	ServiceName    string `mapstructure:"SERVICE_NAME"`
	ServiceVersion string `mapstructure:"SERVICE_VERSION"`
	LogLevel       string `mapstructure:"LOG_LEVEL"`
}

func (c *Config) GetLogLevel(defaultLevel logrus.Level) logrus.Level {
	level, err := logrus.ParseLevel(c.LogLevel)
	if err != nil {
		level = defaultLevel
	}

	return level
}

func LoadConfig(log *logrus.Logger, configObject *Config, fileNames ...string) (*viper.Viper, error) {
	mainConfig := viper.New()
	fileNames = append([]string{"default.yaml", "config/default.yaml"}, fileNames...)

	for _, fileName := range fileNames {
		viperConfig := viper.New()
		viperConfig.SetConfigFile(fileName)

		if strings.Contains(fileName, "default.") {
			viper.AutomaticEnv()
		}

		if err := viperConfig.MergeInConfig(); err != nil {
			log.WithError(err).Info("config not found; skipping")
			continue
		}

		if err := mainConfig.MergeConfigMap(viperConfig.AllSettings()); err != nil {
			return nil, err
		}
	}

	if err := mainConfig.Unmarshal(configObject); err != nil {
		return nil, errors.Wrap(err, "config parsing error")
	}

	return mainConfig, nil
}

func newConfig(configFiles []string) (*Config, error) {
	var (
		tmpLog = logrus.New()
		cfg    Config
	)
	tmpLog.SetOutput(os.Stdout)
	tmpLog.SetLevel(defaultLogLevel)

	_, err := LoadConfig(tmpLog, &cfg, configFiles...)
	if err != nil {
		return nil, err
	}

	//err = cfg.IsValid()
	//if err != nil {
	//	return nil, err
	//}

	return &cfg, nil
}

func NewConfig() (*Config, error) {
	configFiles := []string{
		"default.yaml",
		"config/default.yaml",
		"/config/config.yaml",
		"/vault/secrets/config.yaml",
		".env",
	}

	return newConfig(configFiles)
}
