package drivers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func (c *Controller) getArticle(ctx *fiber.Ctx) error {
	slug := ctx.Params("slug")
	userId, ok := ctx.Locals("userId").(int)

	if !ok {
		userId = 0
	}

	_article, err := c.articleService.GetBySlug(ctx.Context(), slug, userId)

	fmt.Printf("%v\n", _article)

	if err != nil {
		c.log.Debug("Error getting article", "error", err)
		return ctx.Redirect("/", 303)
	}

	_layoutProps := getLayoutProps(ctx)
	_layoutProps.title = _article.Title

	c.log.Debug("layoutProps", "layoutProps", _layoutProps)

	props := articleProps{
		ArticleOutput: _article,
		layoutProps:   _layoutProps,
		isMyArticle:   _article.Author.Username == _layoutProps.user.Username,
	}

	return renderComponent(article(props), ctx)
}
