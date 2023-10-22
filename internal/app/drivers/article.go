package drivers

import (
	"context"

	"github.com/gofiber/fiber/v2"
)

func (c *Controller) getArticle(fc *fiber.Ctx) error {
	ctx := context.Background()
	slug := fc.Params("slug")
	isOob := fc.Query("oob") == "true"
	userId, ok := fc.Locals("userId").(int)

	if !ok {
		userId = 0
	}

	_article, err := c.articleService.GetBySlug(ctx, slug, userId)

	if err != nil {
		c.log.Debug("Error getting article", "error", err)
		return fc.Redirect("/", 303)
	}

	_layoutProps := getLayoutProps(fc)
	_layoutProps.title = _article.Title

	props := articleProps{
		ArticleOutput: _article,
		layoutProps:   _layoutProps,
		isMyArticle:   _article.Author.Username == _layoutProps.user.Username,
	}

	if isOob {
		return renderComponent(articleOOBComp(props), fc)
	}

	return renderComponent(articleComp(props), fc)
}
