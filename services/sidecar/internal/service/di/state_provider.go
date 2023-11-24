package di

import (
	"github.com/google/wire"

	"github.com/specularL2/specular/services/sidecar/internal/sidecar/infra/services"
)

var L1StateProvider = wire.NewSet( //nolint:gochecknoglobals
	services.NewL1State,
)
