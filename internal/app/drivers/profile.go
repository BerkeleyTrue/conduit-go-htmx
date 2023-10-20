package drivers

import (
	"github.com/gofiber/fiber/v2"
)

func (c *Controller) GetProfile(fc *fiber.Ctx) error {
	authorname := fc.Params("username")
	userId, ok := fc.Locals("userId").(int)

	if !ok {
		userId = 0
	}

	username, ok := fc.Locals("username").(string)

	if !ok {
		username = ""
	}

	_profile, err := c.userService.GetProfile(
		fc.Context(),
		0,
		authorname,
		userId,
	)

	if err != nil {
		c.log.Debug("Error getting profile", "error", err)
		return fc.Redirect("/", 303)
	}

	if _profile == nil {
		c.log.Debug("Profile not found", "username", authorname)
		return fc.Redirect("/", 303)
	}

	_layoutProps := getLayoutProps(fc)

	_layoutProps.title = "Profile"

	props := profileProps{
		layoutProps: _layoutProps,
		Profile:     *_profile,
		IsMyself:    _profile.Username == username,
	}

	return renderComponent(profile(props), fc)
}
