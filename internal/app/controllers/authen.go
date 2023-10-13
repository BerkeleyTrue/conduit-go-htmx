package controllers

import (
	"fmt"

	"github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/gofiber/fiber/v2"

	"github.com/berkeleytrue/conduit/internal/core/services"
	"github.com/berkeleytrue/conduit/internal/infra/data/password"
)

func (c *Controller) GetLogin(ctx *fiber.Ctx) error {
	return ctx.Render("auth", fiber.Map{
		"IsRegister": false,
	}, "layouts/main")
}

func (c *Controller) GetRegister(ctx *fiber.Ctx) error {
	return ctx.Render("auth", fiber.Map{
		"IsRegister": true,
	}, "layouts/main")
}

type RegisterInput struct {
	Username string `form:"username"`
	Email    string `form:"email"`
	Password string `form:"password"`
}

func (r *RegisterInput) validate() error {
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
	registerInput := RegisterInput{}

	if err := ctx.BodyParser(&registerInput); err != nil {
		return fmt.Errorf("error parsing register input: %w", err)
	}

	if err := registerInput.validate(); err != nil {

		ctx.Response().Header.Add("HX-Push-Url", "false")
		ctx.Response().Header.Add("HX-Reswap", "none")

		return ctx.Render("partials/auth-errors", fiber.Map{
			"Errors": err,
		})
	}

	pass, err := password.New(registerInput.Password)

	if err != nil {
		ctx.Response().Header.Add("HX-Push-Url", "false")
		ctx.Response().Header.Add("HX-Reswap", "none")

		return ctx.Render("partials/auth-errors", fiber.Map{
			"Errors": []error{err},
		})
	}

	userId, err := c.userService.Register(services.RegisterParams{
		Username: registerInput.Username,
		Email:    registerInput.Email,
		Password: pass,
	})

	if err != nil {
		return fmt.Errorf("error registering user: %w", err)
	}

	// fmt.Printf("register success: %+v\n", registerInput)

	session, err := c.store.Get(ctx)

	if err != nil {
		return fmt.Errorf("error getting session: %w", err)
	}

	session.Set("userId", userId)
	err = session.Save()

	if err != nil {
		return fmt.Errorf("error saving session: %w", err)
	}

	ctx.Response().Header.Add("HX-Push-Url", "/")
	return ctx.Redirect("/", 303)
}

type LoginInput struct {
	Email    string `form:"email"`
	Password string `form:"password"`
}

func (i *LoginInput) validate() error {
	return validation.ValidateStruct(
		i,
		validation.Field(&i.Email, validation.Required, is.Email),
		validation.Field(
			&i.Password,
			validation.Required,
			validation.Length(8, 32),
		),
	)
}

func (c *Controller) Login(ctx *fiber.Ctx) error {
	loginInput := LoginInput{}
	if err := ctx.BodyParser(&loginInput); err != nil {
		return fmt.Errorf("error parsing login input: %w", err)
	}

	if err := loginInput.validate(); err != nil {
		ctx.Response().Header.Add("HX-Push-Url", "false")
		ctx.Response().Header.Add("HX-Reswap", "none")

		return ctx.Render("partials/auth-errors", fiber.Map{
			"Errors": err,
		})
	}

	userId, err := c.userService.Login(loginInput.Email, loginInput.Password)

	if err != nil {
		fmt.Printf("error logging in: %+v\n", err)

		ctx.Response().Header.Add("HX-Push-Url", "false")
		ctx.Response().Header.Add("HX-Reswap", "none")

		return ctx.Render("partials/auth-errors", fiber.Map{
			"Errors": map[string]error{"login": err},
		})
	}

	session, err := c.store.Get(ctx)

	if err != nil {
		return fmt.Errorf("error getting session: %w", err)
	}

	session.Set("userId", userId)
	err = session.Save()

	if err != nil {
		return fmt.Errorf("error saving session: %w", err)
	}

	ctx.Response().Header.Add("HX-Push-Url", "/")
	return ctx.Redirect("/", 303)
}

func (c *Controller) Logout(ctx *fiber.Ctx) error {
	session, err := c.store.Get(ctx)

	if err != nil {
		return fmt.Errorf("error getting session: %w", err)
	}

	session.Destroy()

	ctx.Response().Header.Add("HX-Push-Url", "/")
	return ctx.Redirect("/", 303)
}
