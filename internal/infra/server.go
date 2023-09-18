package infra

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/template/html/v2"
	"go.uber.org/fx"

	"github.com/berkeleytrue/conduit/config"
	"github.com/berkeleytrue/conduit/internal/app"
)

var (
	Module = fx.Options(
		fx.Provide(NewServer),
		fx.Provide(NewSessionStore),
		fx.Provide(NewValidator),
		fx.Invoke(AddMiddlewares),
		app.Module,
		fx.Invoke(RegisterServer),
	)
)

func NewServer(cfg *config.Config) *fiber.App {
	isDev := cfg.Release == "development"

	engine := html.New("./web/views", ".gohtml")
	engine.Reload(isDev)

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	return app
}

func NewSessionStore(cfg *config.Config) *session.Store {
	return session.New()
}

func NewValidator() *validator.Validate {
  return validator.New(validator.WithRequiredStructEnabled())
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
