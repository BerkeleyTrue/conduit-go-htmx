package drivers

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func (c *Controller) getArticle(fc *fiber.Ctx) error {
	ctx := context.Background()
	slug := fc.Params("slug")
	userId, ok := fc.Locals("userId").(int)

	if !ok {
		userId = 0
	}

	_article, err := c.articleService.GetBySlug(ctx, slug, userId)

	fmt.Printf("%v\n", _article)

	if err != nil {
		c.log.Debug("Error getting article", "error", err)
		return fc.Redirect("/", 303)
	}

	_layoutProps := getLayoutProps(fc)
	_layoutProps.title = _article.Title

	c.log.Debug("layoutProps", "layoutProps", _layoutProps)

	props := articleProps{
		ArticleOutput: _article,
		layoutProps:   _layoutProps,
		isMyArticle:   _article.Author.Username == _layoutProps.user.Username,
	}

	return renderComponent(article(props), fc)
}
