package session

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/sqlite3"

	"github.com/berkeleytrue/conduit/config"
	"github.com/berkeleytrue/conduit/internal/core/services"
)

func NewSessionStore(cfg *config.Config) *session.Store {
	storage := sqlite3.New(sqlite3.Config{
		Database: cfg.DB,
	})

	return session.New(session.Config{
		Storage: storage,
	})
}

func RegisterSessionMiddleware(app *fiber.App, store *session.Store, userService *services.UserService) {
	app.Use(func(ctx *fiber.Ctx) error {
		session, err := store.Get(ctx)

		if err != nil {
			return err
		}

		ctx.Locals("session", session)

		userId, ok := session.Get("userId").(int)

		if !ok {
			ctx.Locals("userId", 0)
		} else {
			user, err := userService.GetUser(userId)

			if err != nil {
				session.Destroy()
				return ctx.Status(fiber.StatusForbidden).Redirect("/login")
			}

			ctx.Locals("userId", userId)
			ctx.Locals("user", user)
			ctx.Locals("username", user.Username)
		}

		return ctx.Next()
	})
}

func NewAuthMiddleware(app *fiber.App, store *session.Store, userService *services.UserService) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		session, err := store.Get(ctx)

		if err != nil {
			return err
		}

		userId, ok := session.Get("userId").(int)

		if !ok || userId == 0 {
			return ctx.Redirect("/login", fiber.StatusSeeOther)
		}

		user, err := userService.GetUser(userId)

		if err != nil {
			session.Destroy()
			return ctx.Redirect("/login", fiber.StatusSeeOther)
		}

		ctx.Locals("user", user)
		ctx.Locals("userId", userId)
		ctx.Locals("username", user.Username)

		return ctx.Next()
	}
}
