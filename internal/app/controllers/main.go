package controllers

import (
	"fmt"

	"github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/gofiber/fiber/v2"
)

type (
	Controller struct {
	}
	Link struct {
		Uri   string
		Title string
	}
	RegisterRequest struct {
		Username string `form:"username"`
		Email    string `form:"email"`
		Password string `form:"password"`
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
		"Links":      UnAuthedLinks,
		"Page":       ctx.Path(),
	}, "layouts/main")
}

func (c *Controller) GetRegister(ctx *fiber.Ctx) error {
	return ctx.Render("auth", fiber.Map{
		"IsRegister": true,
		"Links":      UnAuthedLinks,
		"Page":       ctx.Path(),
	}, "layouts/main")
}

func (r *RegisterRequest) Validate() error {
	return validation.ValidateStruct(
		r,
		validation.Field(
			&r.Username,
			validation.Required,
			validation.Length(4, 32),
		),
		validation.Field(&r.Email, validation.Required, is.Email),
		validation.Field(
			&r.Password,
			validation.Required,
			validation.Length(8, 32),
		),
	)
}

func (c *Controller) Register(ctx *fiber.Ctx) error {
	registerReq := RegisterRequest{}

	if err := ctx.BodyParser(&registerReq); err != nil {
		return err
	}

	if err := registerReq.Validate(); err != nil {

		ctx.Response().Header.Add("HX-Push-Url", "false")
		ctx.Response().Header.Add("HX-Reswap", "none")

		return ctx.Render("partials/auth-errors", fiber.Map{
			"Errors": err,
		})
	}

	fmt.Printf("register success: %+v\n", registerReq)

  ctx.Response().Header.Add("HX-Push-Url", "/")
	return ctx.Redirect("/", 303)
}
