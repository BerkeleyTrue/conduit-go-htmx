package controllers

import (
	"fmt"

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
	app.Get("/login", c.GetLogin)
	app.Get("/register", c.GetRegister)
	app.Post("/register", c.Register)
}

func (c *Controller) Index(ctx *fiber.Ctx) error {
	return ctx.Render("index", fiber.Map{
		"Links": UnAuthedLinks,
		"Page":  ctx.Path(),
	}, "layouts/main")
}

func (c *Controller) GetLogin(ctx *fiber.Ctx) error {
	return ctx.Render("auth", fiber.Map{
		"IsRegister": false,
		"Links":   UnAuthedLinks,
		"Page":    ctx.Path(),
	}, "layouts/main")
}

func (c *Controller) GetRegister(ctx *fiber.Ctx) error {
	return ctx.Render("auth", fiber.Map{
		"IsRegister": true,
		"Links":   UnAuthedLinks,
		"Page":    ctx.Path(),
	}, "layouts/main")
}

func (c *Controller) Register(ctx *fiber.Ctx) error {
  req := struct {
    Username string `form:"username"`
    Email    string `form:"email"`
    Password string `form:"password"`
  }{}

  if err := ctx.BodyParser(&req); err != nil {
    return err
  }

  fmt.Printf("regsiter: %+v\n", req)

  return ctx.Redirect("/", 303)
}
