package app

import (
	"go.uber.org/fx"

	"github.com/berkeleytrue/conduit/internal/app/controllers"
)

var Module = fx.Options(
	fx.Provide(controllers.NewController),
	fx.Invoke(controllers.RegisterRoutes),
)
