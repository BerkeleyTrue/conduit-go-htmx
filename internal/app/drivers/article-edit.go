package drivers

import "github.com/gofiber/fiber/v2"

func (c *Controller) getEditArticle(fc *fiber.Ctx) error {
	props := editArticleProps{
		layoutProps: getLayoutProps(fc),
	}

	props.layoutProps.title = "Edit Article"

	return renderComponent(editArticle(props), fc)
}
