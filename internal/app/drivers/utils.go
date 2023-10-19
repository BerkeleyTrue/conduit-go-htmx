package drivers

import (
	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
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

func (lp layoutProps) addAlert(class, msg string) layoutProps {

	lp.alerts = append(lp.alerts, alertPackage{
		class: class,
		msg:   msg,
	})

	return lp
}

func (lp layoutProps) storeAlerts(ctx *fiber.Ctx) layoutProps {

	session, ok := ctx.Locals("session").(*session.Session)

	if ok {
		session.Set("alerts", lp.alerts)
	}
	return lp
}
