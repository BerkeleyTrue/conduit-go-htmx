package infra

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/session"
)

func AddMiddlewares(app *fiber.App, store *session.Store) {
	app.Use(logger.New())
	app.Use(cors.New())
	app.Use("/", func(ctx *fiber.Ctx) error {
    _, err := store.Get(ctx)

    if err != nil {
        return err
    }


    return ctx.Next()
  })
}
