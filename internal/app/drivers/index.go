package drivers

import (
	"github.com/gofiber/fiber/v2"

	"github.com/berkeleytrue/conduit/internal/infra/session"
)

func (c *Controller) Index(fc *fiber.Ctx) error {
	flashes, err := session.GetFlashes(fc)

	if err != nil {
		return err
	}

	_layoutProps := getLayoutProps(fc)

	_layoutProps.title = "Home"
	_layoutProps.flashes = flashes

	p := indexProps{
		layoutProps: _layoutProps,
	}

	return renderComponent(index(p), fc)
}
