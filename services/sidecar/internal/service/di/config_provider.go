package di

import (
	"github.com/google/wire"
	"github.com/specularL2/specular/services/sidecar/internal/service/config"
	"github.com/specularL2/specular/services/sidecar/rollup/services"
)

var ConfigProvider = wire.NewSet( //nolint:gochecknoglobals
	config.NewConfig,
)

var SystemConfigProvider = wire.NewSet( //nolint:gochecknoglobals
	services.ParseSystemConfig,
)
