package controllers

import (
	"github.com/go-ozzo/ozzo-validation/v4"
	"github.com/gofiber/fiber/v2"

	"github.com/berkeleytrue/conduit/internal/core/domain"
)

type (
	GetArticlesInput struct {
		author    string `query:"author"`
		favorited string `query:"favorited"`
		tag       string `query:"tag"`
		limit     int    `query:"limit"`
		offset    int    `query:"offset"`
	}
)

func (i *GetArticlesInput) validate() error {

	return validation.ValidateStruct(
		i,
		validation.Field(&i.offset, validation.Min(0), validation.Max(10000)),
		validation.Field(&i.limit, validation.Min(0), validation.Max(10000)),
	)
}

// get a list of articles,
// optionally filtered by query parameters
// author=authorname, favorited=authorname, tag=string, limit=int, offset=int
// is authenticated, check if articles is favorited by user
func (c *Controller) GetArticles(ctx *fiber.Ctx) error {
	input := GetArticlesInput{}

	session, err := c.store.Get(ctx)

	if err != nil {
		return err
	}

	if err := ctx.QueryParser(&input); err != nil {
		return err
	}

	if err := input.validate(); err != nil {
		return err
	}

	username, ok := session.Get("username").(string)

	if !ok {
		username = ""
	}

	articles, err := c.articleService.List(
		username,
		domain.ArticleListInput{
			Tag:       input.tag,
			Favorited: input.favorited,
			Limit:     input.limit,
			Offset:    input.offset,
		},
	)

	if err != nil {
		return err
	}

	return ctx.Render("partials/articles", fiber.Map{
		"Articles": articles,
	})
}
