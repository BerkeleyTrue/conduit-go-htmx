package drivers

import (
	"fmt"

	"github.com/go-ozzo/ozzo-validation/v4"
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

func (c *Controller) GetSettings(ctx *fiber.Ctx) error {
	return ctx.Render("settings", fiber.Map{}, "layouts/main")
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

func (c *Controller) UpdateSettings(ctx *fiber.Ctx) error {
	settingsInput := SettingsInput{}

	if err := ctx.BodyParser(&settingsInput); err != nil {
		return fmt.Errorf("error parsing settings input: %w", err)
	}

	if err := settingsInput.validate(); err != nil {

		ctx.Response().Header.Add("HX-Push-Url", "false")
		ctx.Response().Header.Add("HX-Reswap", "none")

		return ctx.Render("partials/auth-errors", fiber.Map{
			"Errors": err,
		})
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

			ctx.Response().Header.Add("HX-Push-Url", "false")
			ctx.Response().Header.Add("HX-Reswap", "none")

			return ctx.Render("partials/auth-errors", fiber.Map{
				"Errors": map[string]error{"password": err},
			})

		}
	}

	// fmt.Printf("settings input %+v\n", settingsInput)

	user, err := c.userService.Update(
		ctx.Locals("userId").(int),
		"",
		updates,
	)

	if err != nil {
		ctx.Response().Header.Add("HX-Push-Url", "false")
		ctx.Response().Header.Add("HX-Reswap", "none")
		return ctx.Render("partials/auth-errors", fiber.Map{
			"Errors": []error{err},
		})
	}

	ctx.Locals("user", user)

	return ctx.Render("settings", fiber.Map{})
}
