package drivers

import (
	"context"
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/gofiber/fiber/v2"

	"github.com/berkeleytrue/conduit/internal/core/services"
	"github.com/berkeleytrue/conduit/internal/infra/data/password"
	"github.com/berkeleytrue/conduit/internal/infra/session"
)

func (c *Controller) GetLogin(fc *fiber.Ctx) error {
	_layoutProps := getLayoutProps(fc)
	_layoutProps.title = "Login"

	props := authProps{
		isRegister:  false,
		layoutProps: _layoutProps,
	}

	return renderComponent(auth(props), fc)
}

func (c *Controller) GetRegister(fc *fiber.Ctx) error {
	_layoutProps := getLayoutProps(fc)
	_layoutProps.title = "Register"

	props := authProps{
		isRegister:  true,
		layoutProps: _layoutProps,
	}

	return renderComponent(auth(props), fc)
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

func (c *Controller) Register(fc *fiber.Ctx) error {
	ctx := context.Background()
	registerInput := RegisterInput{}

	if err := fc.BodyParser(&registerInput); err != nil {
		return fmt.Errorf("error parsing register input: %w", err)
	}

	if err := registerInput.validate(); err != nil {

		fc.Response().Header.Add("HX-Push-Url", "false")
		fc.Response().Header.Add("HX-Reswap", "none")

		return renderComponent(listErrors(err.(validation.Errors)), fc)
	}

	pass, err := password.New(registerInput.Password)

	if err != nil {
		fc.Response().Header.Add("HX-Push-Url", "false")
		fc.Response().Header.Add("HX-Reswap", "none")

		return renderComponent(
			listErrors(map[string]error{"password": err}),
			fc,
		)
	}

	userId, err := c.userService.Register(
		ctx,
		services.RegisterParams{
			Username: registerInput.Username,
			Email:    registerInput.Email,
			Password: pass,
		},
	)

	if err != nil {
		return fmt.Errorf("error registering user: %w", err)
	}

	err = session.SaveUser(fc, userId)

	if err != nil {
		return fmt.Errorf("error saving user: %w", err)
	}

	err = session.AddFlash(fc, "success", "Welcome to Conduit!")

	if err != nil {
		c.log.Debug("error adding flash", "error", err)
	}

	return fc.Redirect("/", fiber.StatusSeeOther)
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

func (c *Controller) Login(fc *fiber.Ctx) error {
	ctx := context.Background()
	loginInput := LoginInput{}
	if err := fc.BodyParser(&loginInput); err != nil {
		return fmt.Errorf("error parsing login input: %w", err)
	}

	if err := loginInput.validate(); err != nil {
		fc.Response().Header.Add("HX-Push-Url", "false")
		fc.Response().Header.Add("HX-Reswap", "none")

		return renderComponent(listErrors(err.(validation.Errors)), fc)
	}

	userId, err := c.userService.Login(
		ctx,
		loginInput.Email,
		loginInput.Password,
	)

	if err != nil {
		fmt.Printf("error logging in: %+v\n", err)

		fc.Response().Header.Add("HX-Push-Url", "false")
		fc.Response().Header.Add("HX-Reswap", "none")

		return renderComponent(listErrors(map[string]error{"login": err}), fc)
	}

	err = session.SaveUser(fc, userId)

	if err != nil {
		return fmt.Errorf("error saving session: %w", err)
	}

	err = session.AddFlash(fc, session.Success, "Logged in successfully!")

	if err != nil {
		c.log.Debug("error adding flash", "error", err)
	}

	return fc.Redirect("/", fiber.StatusSeeOther)
}

func (c *Controller) Logout(fc *fiber.Ctx) error {
	err := session.Logout(fc)

	if err != nil {
		return fmt.Errorf("error logging out: %w", err)
	}

	fc.Response().Header.Add("HX-Push-Url", "/")
	return fc.Redirect("/", fiber.StatusSeeOther)
}
