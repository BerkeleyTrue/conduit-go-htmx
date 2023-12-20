package drivers

import (
	"context"
	"errors"

	"github.com/berkeleytrue/conduit/internal/infra/session"
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

func (c *Controller) deleteArticle(fc *fiber.Ctx) error {
	ctx := context.Background()
	slug := fc.Params("slug")
	err := c.articleService.Delete(ctx, slug)

	if err != nil {
		c.log.Debug("Error deleting article", "error", err)

		fc.Response().Header.Add("HX-Push-Url", "false")
		fc.Response().Header.Add("HX-Reswap", "none")

		return renderComponent(
			listErrors(
				map[string]error{
					"article": errors.New("Error deleting article"),
				},
			),
			fc,
		)
	}

	session.AddFlash(fc, session.Info, "Article removed successfully!")

	return fc.Redirect("/", 303)
}
