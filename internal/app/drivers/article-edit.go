package drivers

import "github.com/gofiber/fiber/v2"

func (c *Controller) getEditArticle(fc *fiber.Ctx) error {
  slug := fc.Params("slug")
	props := editArticleProps{
		layoutProps: getLayoutProps(fc),
	}

	props.layoutProps.title = "Edit Article"
	props.isNew = slug == ""

	return renderComponent(editArticleComp(props), fc)
}
