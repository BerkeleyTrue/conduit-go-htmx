package controllers

import (
	"fmt"

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

	profile, err := c.userService.GetProfile(
		0,
		authorname,
		userId,
	)

	if err != nil {
		fmt.Println(err)
		return ctx.Status(404).SendString("profile not found")
	}

	return ctx.Render("profile", fiber.Map{
		"Profile":  profile,
		"IsMyself": profile.Username == username,
	}, "layouts/main")
}
