package controllers

import (
	"github.com/gofiber/fiber/v2"
)

func (c *Controller) GetSettings(ctx *fiber.Ctx) error {
	return ctx.Render("settings", fiber.Map{}, "layouts/main")
}
