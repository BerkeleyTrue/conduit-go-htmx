package controllers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"

	"github.com/berkeleytrue/conduit/internal/core/services"
)

func (c *Controller) GetProfile(ctx *fiber.Ctx) error {
	authorname := ctx.Params("username")
	userId, ok := ctx.Locals("userId").(int)

	if !ok {
		fmt.Println("user not found")
		userId = 0
	}

	profile, err := c.userService.GetProfile(services.UserIdOrUsername{
		Username: authorname,
	}, userId)

	if err != nil {
		fmt.Println(err)
		return ctx.Status(404).SendString("profile not found")
	}

	return ctx.Render("profile", fiber.Map{
		"Profile": profile,
	}, "layouts/main")
}
