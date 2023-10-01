package infra

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"go.uber.org/fx"

	"github.com/berkeleytrue/conduit/config"
	"github.com/berkeleytrue/conduit/internal/app/controllers"
	"github.com/berkeleytrue/conduit/internal/app/driven/articlesRepo"
	"github.com/berkeleytrue/conduit/internal/app/driven/userRepo"
	"github.com/berkeleytrue/conduit/internal/core/domain"
	"github.com/berkeleytrue/conduit/internal/core/services"
	"github.com/berkeleytrue/conduit/internal/infra/session"
)

var (
	Module = fx.Options(
		fx.Provide(NewServer),
		fx.Provide(NewDB),
		fx.Provide(
			fx.Annotate(
				userRepo.NewSqlStore,
				fx.As(new(domain.UserRepository)),
			),
		),
		fx.Provide(
			fx.Annotate(
				articlesRepo.NewSqlStore,
				fx.As(new(domain.ArticleRepository)),
			),
		),
		fx.Provide(services.NewUserService),
		fx.Provide(services.NewArticleService),
		fx.Provide(session.NewSessionStore),
		fx.Provide(
			fx.Annotate(
				session.NewAuthMiddleware,
				fx.ResultTags(`name:"authMiddleware"`),
			),
		),
		fx.Provide(controllers.NewController),

		fx.Invoke(userRepo.RegisterUserSchema),
		fx.Invoke(articlesRepo.RegisterArticleSchema),
		fx.Invoke(AddMiddlewares),
		fx.Invoke(session.RegisterSessionMiddleware),
		fx.Invoke(
			fx.Annotate(
				controllers.RegisterRoutes,
				fx.ParamTags("", "", `name:"authMiddleware"`),
			),
		),
		fx.Invoke(RegisterServer),
	)
)

func NewServer(cfg *config.Config) *fiber.App {
	isDev := cfg.Release == "development"

	engine := html.New("./web/views", ".gohtml")
	engine.Reload(isDev)

	app := fiber.New(fiber.Config{
		Views:             engine,
		PassLocalsToViews: true,
	})

	return app
}

func StartServer(app *fiber.App, config *config.Config) error {
	cfg := config.HTTP
	port := cfg.Port
	go app.Listen(":" + port)

	return nil
}

func StopServer(app *fiber.App) error {
	return app.Shutdown()
}

func RegisterServer(lc fx.Lifecycle, app *fiber.App, config *config.Config) {
	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			StartServer(app, config)
			return nil
		},
		OnStop: func(_ context.Context) error {
			return StopServer(app)
		},
	})
}
