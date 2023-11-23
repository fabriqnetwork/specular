package di

import (
	"github.com/google/wire"
	"github.com/specularL2/specular/services/sidecar/internal/sidecar/api"
)

var CliProvider = wire.NewSet( //nolint:gochecknoglobals
	api.NewCli,
)
