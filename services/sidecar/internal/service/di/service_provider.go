package di

import (
	"github.com/google/wire"

	"github.com/specularL2/specular/services/sidecar/internal/sidecar/infra/services"
)

var DisseminatorProvider = wire.NewSet( //nolint:gochecknoglobals
	services.NewDisseminator,
)

var ValidatorProvider = wire.NewSet( //nolint:gochecknoglobals
	services.NewValidator,
)
