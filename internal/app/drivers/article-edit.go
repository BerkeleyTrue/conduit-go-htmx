package drivers

import "github.com/gofiber/fiber/v2"

func (c *Controller) getEditArticle(fc *fiber.Ctx) error {
	slug := fc.Params("slug")
	props := editArticleProps{
		layoutProps: getLayoutProps(fc),
	}

	if slug == "" {
		props.layoutProps.title = "New Article"
		props.isNew = true
		return renderComponent(editArticleComp(props), fc)
	}

	props.layoutProps.title = "Edit Article"

	article, err := c.articleService.GetBySlug(fc.Context(), slug, 0)

	if err != nil {
		c.log.Debug("Error getting article", "error", err)
		return fc.Redirect("/", fiber.StatusNotFound)
	}

	props.article = &article

	return renderComponent(editArticleComp(props), fc)
}
