package drivers

import (
	"github.com/a-h/templ"
	"github.com/berkeleytrue/conduit/internal/infra/session"
	"github.com/gofiber/fiber/v2"
)

func renderComponent(comp templ.Component, fc *fiber.Ctx) error {
	fc.Type("html")

	return comp.Render(fc.Context(), fc)
}

func getLayoutProps(fc *fiber.Ctx) layoutProps {
	_layoutProps, ok := fc.Locals("layoutProps").(layoutProps)

	if !ok {
		_layoutProps = layoutProps{}
	}

	return _layoutProps
}

func (lp layoutProps) addFlashes(fc *fiber.Ctx) layoutProps {
	flashes, err := session.GetFlashes(fc)

	if err == nil {
		lp.flashes = flashes
	}
	return lp
}
