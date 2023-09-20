package infra

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/sqlite3"

	"github.com/berkeleytrue/conduit/config"
)

func NewSessionStore(cfg *config.Config) *session.Store {
	storage := sqlite3.New(sqlite3.Config{
		Database: cfg.DB,
	})

	return session.New(session.Config{
		Storage: storage,
	})
}

func AddSessionMiddleware(app *fiber.App, store *session.Store) {
	app.Use(func(ctx *fiber.Ctx) error {
		session, err := store.Get(ctx)

		if err != nil {
			return err
		}

		var sid string
		// always set sid if session is fresh
		if session.Fresh() {

			sid = session.ID()
			session.Set("sid", sid)

			fmt.Printf("Fresh session: %v\n", sid)

			// not safe to use session after save
			if err := session.Save(); err != nil {
				return err
			}

		} else {
			sid = session.Get("sid").(string)
			fmt.Printf("Session: %v\n", sid)
		}

		ctx.Locals("sid", sid)

		return ctx.Next()
	})
}
