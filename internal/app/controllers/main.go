package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"

	"github.com/berkeleytrue/conduit/internal/core/services"
)

type (
	Controller struct {
		store          *session.Store
		userService    *services.UserService
		articleService *services.ArticleService
	}
	Link struct {
		Uri   string
		Title string
	}
)

var (
	UnAuthedLinks = []Link{
		{
			Uri:   "/",
			Title: "Home",
		},
		{
			Uri:   "/login",
			Title: "Login",
		},
		{
			Uri:   "/register",
			Title: "Sign up",
		},
	}
	AuthedLinks = []Link{
		{
			Uri:   "/",
			Title: "Home",
		},
		{
			Uri:   "/editor",
			Title: "New Article",
		},
		{
			Uri:   "/settings",
			Title: "Settings",
		},
	}
)

func NewController(
	store *session.Store,
	userService *services.UserService,
	articleService *services.ArticleService,
) *Controller {
	return &Controller{
		store:          store,
		userService:    userService,
		articleService: articleService,
	}
}

func RegisterRoutes(app *fiber.App, c *Controller, authMiddleware fiber.Handler) {
	app.Use(func(ctx *fiber.Ctx) error {
		userId := ctx.Locals("userId")
		links := UnAuthedLinks

		if userId != 0 {
			links = AuthedLinks
		}

		ctx.Bind(fiber.Map{
			"Links": links,
			"Page":  ctx.Path(),
		})

		return ctx.Next()
	})

	app.Get("/", c.Index)
	app.Get("/login", c.GetLogin)
	app.Post("/login", c.Login)
	app.Get("/register", c.GetRegister)
	app.Post("/register", c.Register)

	app.Get("/profile/:username", c.GetProfile)
	app.Get("/articles", c.GetArticles)
	app.Get("/tags", c.GetPopularTags)

	app.Use(authMiddleware)

	app.Get("/settings", c.GetSettings)
	app.Post("/settings", c.UpdateSettings)
	app.Post("/logout", c.Logout)
}

func (c *Controller) Index(ctx *fiber.Ctx) error {

	return ctx.Render("index", fiber.Map{}, "layouts/main")
}
