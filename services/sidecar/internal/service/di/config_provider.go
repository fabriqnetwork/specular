package di

import (
	"github.com/google/wire"

	"github.com/specularL2/specular/services/sidecar/internal/service/config"
)

var ConfigProvider = wire.NewSet( //nolint:gochecknoglobals
	config.NewConfig,
)

var SystemConfigProvider = wire.NewSet( //nolint:gochecknoglobals
	config.NewSystemConfig,
)
