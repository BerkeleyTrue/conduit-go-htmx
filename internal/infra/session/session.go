package session

import (
	"fmt"

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

// SaveUser saves the user id in the session
func SaveUser(fc *fiber.Ctx, userId int) error {
	session, ok := fc.Locals("session").(*session.Session)

	if !ok {
		return fmt.Errorf("session not found")
	}

	session.Set("userId", userId)

	return nil
}

// Logout destroys the session
func Logout(fc *fiber.Ctx) error {
	session, ok := fc.Locals("session").(*session.Session)

	if !ok {
		return fmt.Errorf("session not found")
	}

	return session.Destroy()
}

func RegisterSessionMiddleware(
	app *fiber.App,
	store *session.Store,
	userService *services.UserService,
) {
	store.RegisterType(flashMap{})

	app.Use(func(fc *fiber.Ctx) error {
		session, err := store.Get(fc)

		if err != nil {
			return err
		}

		fc.Locals("session", session)

		userId, ok := session.Get("userId").(int)

		if !ok {
			fc.Locals("userId", 0)
		} else {
			user, err := userService.GetUser(fc.Context(), userId)

			if err != nil {
				session.Destroy()
				return fc.Status(fiber.StatusForbidden).Redirect("/login")
			}

			fc.Locals("userId", userId)
			fc.Locals("user", user)
			fc.Locals("username", user.Username)
		}

		err = fc.Next()

		if err != nil {
			return err
		}

		// Persist session data to the storage after request
		return session.Save()
	})
}

func NewAuthMiddleware(
	app *fiber.App,
	store *session.Store,
	userService *services.UserService,
) fiber.Handler {
	return func(fc *fiber.Ctx) error {
		session, err := store.Get(fc)

		if err != nil {
			return err
		}

		userId, ok := session.Get("userId").(int)

		if !ok || userId == 0 {
			return fc.Redirect("/login", fiber.StatusSeeOther)
		}

		user, err := userService.GetUser(fc.Context(), userId)

		if err != nil {
			session.Destroy()
			return fc.Redirect("/login", fiber.StatusSeeOther)
		}

		fc.Locals("user", user)
		fc.Locals("userId", userId)
		fc.Locals("username", user.Username)

		return fc.Next()
	}
}
