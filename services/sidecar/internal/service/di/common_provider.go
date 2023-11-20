package di

import (
	"github.com/google/wire"

	"github.com/specularL2/specular/services/sidecar/internal/service/config"
)

var CommonProvider = wire.NewSet( //nolint:gochecknoglobals
	config.NewLogger,
	config.NewCancelChannel,
	config.NewContext,
)
