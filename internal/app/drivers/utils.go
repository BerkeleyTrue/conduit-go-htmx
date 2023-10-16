package drivers

import (
	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
)

func renderComponent(comp templ.Component, ctx *fiber.Ctx) error {
	ctx.Type("html")

	return comp.Render(ctx.Context(), ctx)
}

func getLayoutProps(ctx *fiber.Ctx) layoutProps {
	_layoutProps, ok := ctx.Locals("layoutProps").(layoutProps)

	if !ok {
		_layoutProps = layoutProps{}
	}

	return _layoutProps
}
