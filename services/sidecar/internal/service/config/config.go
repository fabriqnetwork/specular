package config

import (
	"os"
	"strings"

	"github.com/specularL2/specular/services/sidecar/utils/log"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

const (
	defaultLogLevel    = log.LvlInfo
	defaultServiceName = "sidecar"
	defaultVersion     = "unknown"
	defaultUsageDesc   = "launch a validator and/or disseminator"
)

type Config struct {
	ServiceName    string `mapstructure:"SERVICE_NAME"`
	ServiceVersion string `mapstructure:"SERVICE_VERSION"`

	UsageDesc string `mapstructure:"USAGE_DESC"`
	LogLevel  string `mapstructure:"LOG_LEVEL"`
}

func (c *Config) IsValid() error {
	if len(c.ServiceName) == 0 {
		return errors.New("invalid config: SERVICE_NAME cannot be empty")
	}

	if len(c.UsageDesc) == 0 {
		return errors.New("invalid config: USAGE_DESC cannot be empty")
	}

	return nil
}

func LoadConfig(log log.Logger, configObject *Config, fileNames ...string) (*viper.Viper, error) {
	mainConfig := viper.New()
	fileNames = append([]string{"default.yaml", "config/default.yaml"}, fileNames...)

	for _, fileName := range fileNames {
		log.Debug("Loading config", "fileName", fileName)
		viperConfig := viper.New()
		viperConfig.SetConfigFile(fileName)

		if strings.Contains(fileName, "default.") {
			viper.AutomaticEnv()
		}

		if err := viperConfig.MergeInConfig(); err != nil {
			log.Debug("config not found; skipping")
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

func (c *Config) LoadDefaults() {
	if len(c.ServiceName) == 0 {
		c.ServiceName = defaultServiceName
	}

	if len(c.UsageDesc) == 0 {
		c.UsageDesc = defaultUsageDesc
	}

	if len(c.ServiceVersion) == 0 {
		c.ServiceVersion = defaultVersion
	}
}

func newConfig(configFiles []string) (*Config, error) {
	var (
		tmpLog = log.New()
		cfg    Config
	)
	gLogger := log.NewGlogHandler(log.StreamHandler(os.Stdout, log.TerminalFormat(false)))
	gLogger.Verbosity(defaultLogLevel)
	tmpLog.SetHandler(gLogger)

	cfg.LoadDefaults()

	_, err := LoadConfig(tmpLog, &cfg, configFiles...)
	if err != nil {
		return nil, err
	}

	err = cfg.IsValid()
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

func NewConfig() (*Config, error) {
	configFiles := []string{
		"default.yaml",
		"config/default.yaml",
		// "/config/config.yaml",
		// ".env",
	}

	return newConfig(configFiles)
}
