package controllers

import (
	"fmt"

	"github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/gofiber/fiber/v2"

	"github.com/berkeleytrue/conduit/internal/core/domain"
	"github.com/berkeleytrue/conduit/internal/infra/data/password"
)

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
		return err
	}

	if err := registerInput.validate(); err != nil {

		ctx.Response().Header.Add("HX-Push-Url", "false")
		ctx.Response().Header.Add("HX-Reswap", "none")

		return ctx.Render("partials/auth-errors", fiber.Map{
			"Errors": err,
		})
	}

	userId, err := c.userService.Register(domain.UserCreateInput{
		Username: registerInput.Username,
		Email:    registerInput.Email,
		Password: password.Password(registerInput.Password),
	})

	if err != nil {
		return err
	}

	fmt.Printf("register success: %+v\n", registerInput)

  session, err := c.store.Get(ctx)

  if err != nil {
    return err
  }

  session.Set("userId", userId)
  err = session.Save()

  if err != nil {
    return err
  }

	ctx.Response().Header.Add("HX-Push-Url", "/")
	return ctx.Redirect("/", 303)
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
    return err
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
    return err
  }

  session, err := c.store.Get(ctx)

  if err != nil {
    return err
  }

  session.Set("userId", userId)
  err = session.Save()

  if err != nil {
    return err
  }

  ctx.Response().Header.Add("HX-Push-Url", "/")
  return ctx.Redirect("/", 303)
}
