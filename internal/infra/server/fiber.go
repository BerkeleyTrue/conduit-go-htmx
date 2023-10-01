package server

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"go.uber.org/fx"

	"github.com/berkeleytrue/conduit/config"
)

func NewServer(engine *html.Engine) *fiber.App {
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
