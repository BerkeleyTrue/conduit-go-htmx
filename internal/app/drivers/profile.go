package drivers

import (
	"context"

	"github.com/gofiber/fiber/v2"
)

func (c *Controller) GetProfile(fc *fiber.Ctx) error {
	ctx := context.Background()
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
		ctx,
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

	return renderComponent(profileComp(props), fc)
}

func (c *Controller) follow(fc *fiber.Ctx) error {
	ctx := context.Background()
	authorname := fc.Params("username")
	userId, ok := fc.Locals("userId").(int)

	if !ok {
		return fiber.ErrUnauthorized
	}

	_, err := c.userService.Follow(
		ctx,
		userId,
		0,
		authorname,
	)

	if err != nil {
		fc.Response().Header.Add("HX-Push-Url", "false")
		fc.Response().Header.Add("HX-Reswap", "none")

		return renderComponent(listErrors(map[string]error{"follow": err}), fc)
	}

	return fc.SendStatus(200)
}

func (c *Controller) unfollow(fc *fiber.Ctx) error {
	ctx := context.Background()
	authorname := fc.Params("username")
	userId, ok := fc.Locals("userId").(int)

	if !ok {
		return fiber.ErrUnauthorized
	}

	_, err := c.userService.Unfollow(
		ctx,
		userId,
		0,
		authorname,
	)

	if err != nil {
		fc.Response().Header.Add("HX-Push-Url", "false")
		fc.Response().Header.Add("HX-Reswap", "none")

		return renderComponent(listErrors(map[string]error{"follow": err}), fc)
	}

	return fc.SendStatus(200)
}
