package controllers

import (
	"github.com/gofiber/fiber/v2"
)

type (
	Controller struct {
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

func NewController() *Controller {
	return &Controller{}
}

func RegisterRoutes(app *fiber.App, c *Controller) {
	app.Get("/", c.Index)
	app.Get("/login", c.Login)
	app.Get("/register", c.Register)
}

func (c *Controller) Index(ctx *fiber.Ctx) error {
	return ctx.Render("index", fiber.Map{
		"Links": UnAuthedLinks,
		"Page":  ctx.Path(),
	}, "layouts/main")
}

func (c *Controller) Login(ctx *fiber.Ctx) error {
	return ctx.Render("auth", fiber.Map{
		"IsLogin": true,
		"Links":   UnAuthedLinks,
		"Page":    ctx.Path(),
	}, "layouts/main")
}

func (c *Controller) Register(ctx *fiber.Ctx) error {
	return ctx.Render("auth", fiber.Map{
		"IsLogin": false,
		"Links":   UnAuthedLinks,
		"Page":    ctx.Path(),
	}, "layouts/main")
}
