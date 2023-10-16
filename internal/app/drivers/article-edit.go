package drivers

import "github.com/gofiber/fiber/v2"

func (c *Controller) getEditArticle(ctx *fiber.Ctx) error {
	props := editArticleProps{
		layoutProps: getLayoutProps(ctx),
	}

	props.layoutProps.title = "Edit Article"

	return renderComponent(editArticle(props), ctx)
}
