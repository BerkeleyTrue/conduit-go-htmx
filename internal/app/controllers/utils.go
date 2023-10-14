package controllers

import (
	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
)

func RenderComponent(comp templ.Component, ctx *fiber.Ctx) error {
	ctx.Type("html")

	return comp.Render(ctx.Context(), ctx)
}
