package app

import (
	"go.uber.org/fx"

	"github.com/berkeleytrue/conduit/internal/app/controllers"
	"github.com/berkeleytrue/conduit/internal/app/driven"
)

var Module = fx.Options(
	driven.Module,
	fx.Provide(controllers.NewController),
	fx.Invoke(controllers.RegisterRoutes),
)
