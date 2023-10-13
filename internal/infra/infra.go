package infra

import (
	"go.uber.org/fx"

	"github.com/berkeleytrue/conduit/internal/app/controllers"
	"github.com/berkeleytrue/conduit/internal/app/driven/articlesRepo"
	"github.com/berkeleytrue/conduit/internal/app/driven/userRepo"
	"github.com/berkeleytrue/conduit/internal/core/services"
	"github.com/berkeleytrue/conduit/internal/infra/db"
	"github.com/berkeleytrue/conduit/internal/infra/server"
	"github.com/berkeleytrue/conduit/internal/infra/session"
)

var (
	Module = fx.Options(
		fx.Provide(server.NewEngine),
		fx.Provide(server.NewServer),
		db.Module,
		userRepo.Module,
		articlesRepo.Module,
		services.Module,
		fx.Provide(session.NewSessionStore),
		fx.Provide(
			fx.Annotate(
				session.NewAuthMiddleware,
				fx.ResultTags(`name:"authMiddleware"`),
			),
		),
		fx.Provide(controllers.NewController),

		fx.Invoke(AddMiddlewares),
		fx.Invoke(session.RegisterSessionMiddleware),
		fx.Invoke(
			fx.Annotate(
				controllers.RegisterRoutes,
				fx.ParamTags("", "", `name:"authMiddleware"`),
			),
		),
		fx.Invoke(server.RegisterServer),
	)
)
