package drivers

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/gofiber/fiber/v2"

	"github.com/berkeleytrue/conduit/internal/core/services"
	"github.com/berkeleytrue/conduit/internal/infra/data/password"
)

type (
	SettingsInput struct {
		Username string `form:"username"`
		Email    string `form:"email"`
		Password string `form:"password"`
		Bio      string `form:"bio"`
		Image    string `form:"image"`
	}
)

func (c *Controller) GetSettings(fc *fiber.Ctx) error {
	props := settingsProps{
		user: *fc.Locals("user").(*services.UserOutput),
	}
	return renderComponent(settings(props), fc)
}

func (r *SettingsInput) validate() error {
	return validation.ValidateStruct(
		r,
		validation.Field(
			&r.Username,
			validation.Length(4, 32),
		),
		validation.Field(&r.Email, validation.Required, is.Email),
		validation.Field(
			&r.Bio,
			validation.Length(0, 255),
		),
		validation.Field(
			&r.Image,
			validation.Length(0, 255),
			is.URL,
		),
	)
}

func (c *Controller) UpdateSettings(fc *fiber.Ctx) error {
	settingsInput := SettingsInput{}

	if err := fc.BodyParser(&settingsInput); err != nil {
		return fmt.Errorf("error parsing settings input: %w", err)
	}

	if err := settingsInput.validate(); err != nil {

		fc.Response().Header.Add("HX-Push-Url", "false")
		fc.Response().Header.Add("HX-Reswap", "none")

		return renderComponent(listErrors(err.(validation.Errors)), fc)
	}
	updates := services.UpdateUserInput{
		Username: settingsInput.Username,
		Email:    settingsInput.Email,
		Bio:      settingsInput.Bio,
		Image:    settingsInput.Image,
	}

	if settingsInput.Password != "" {
		fmt.Printf("settings input %+v\n", settingsInput)

		if pass, err := password.New(settingsInput.Password); err == nil {
			updates.Password = pass
		} else {

			fc.Response().Header.Add("HX-Push-Url", "false")
			fc.Response().Header.Add("HX-Reswap", "none")

			return renderComponent(listErrors(map[string]error{"password": err}), fc)
		}
	}

	user, err := c.userService.Update(
		fc.Context(),
		fc.Locals("userId").(int),
		"",
		updates,
	)

	if err != nil {
		fc.Response().Header.Add("HX-Push-Url", "false")
		fc.Response().Header.Add("HX-Reswap", "none")
		return renderComponent(listErrors(map[string]error{"user": err}), fc)
	}

	fc.Locals("user", user)

	_layoutProps := getLayoutProps(fc)

	_layoutProps.title = "Settings"

	props := settingsProps{
		layoutProps: _layoutProps,
		user:        *user,
	}
	return renderComponent(settings(props), fc)
}
