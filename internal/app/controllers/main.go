package controllers

import (
	"github.com/gofiber/fiber/v2"
)

type (
	Controller struct {
	}
)

func NewController() *Controller {
	return &Controller{}
}

func RegisterRoutes(app *fiber.App, c *Controller) {
	app.Get("/", c.Index)
}

func (c *Controller) Index(ctx *fiber.Ctx) error {
	return ctx.Render("index", fiber.Map{}, "layouts/main")
}
