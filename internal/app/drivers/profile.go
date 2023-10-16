package drivers

import (
	"github.com/gofiber/fiber/v2"
)

func (c *Controller) GetProfile(ctx *fiber.Ctx) error {
	authorname := ctx.Params("username")
	userId, ok := ctx.Locals("userId").(int)

	if !ok {
		userId = 0
	}

	username, ok := ctx.Locals("username").(string)

	if !ok {
		username = ""
	}

	_profile, err := c.userService.GetProfile(
		0,
		authorname,
		userId,
	)

	if err != nil {
		c.log.Debug("Error getting profile", "error", err)
		return ctx.Redirect("/", 303)
	}

	if _profile == nil {
		c.log.Debug("Profile not found", "username", authorname)
		return ctx.Redirect("/", 303)
	}

	_layoutProps := getLayoutProps(ctx)

	_layoutProps.title = "Profile"

	props := profileProps{
		layoutProps: _layoutProps,
		Profile:     *_profile,
		IsMyself:    _profile.Username == username,
	}

	return renderComponent(profile(props), ctx)
}
