package drivers

import (
	"errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/gofiber/fiber/v2"

	"github.com/berkeleytrue/conduit/internal/core/services"
)

type GetArticlesInput struct {
	Author    string `query:"author"`
	Favorited string `query:"favorited"`
	Tag       string `query:"tag"`
	Limit     int    `query:"limit"`
	Offset    int    `query:"offset"`
}

func (i *GetArticlesInput) validate() error {

	return validation.ValidateStruct(
		i,
		validation.Field(&i.Offset, validation.Min(0), validation.Max(10000)),
		validation.Field(&i.Limit, validation.Min(0), validation.Max(10000)),
	)
}

// get a list of articles,
// optionally filtered by query parameters
// author=authorname, favorited=authorname, tag=string, limit=int, offset=int
// is authenticated, check if articles is favorited by user
func (c *Controller) GetArticles(ctx *fiber.Ctx) error {
	input := new(GetArticlesInput)

	if err := ctx.QueryParser(input); err != nil {
		return err
	}

	if err := input.validate(); err != nil {
		return err
	}

	userId, ok := ctx.Locals("userId").(int)

	if !ok {
		userId = 0
	}

	if input.Limit == 0 {
		input.Limit = 20
	}

	articles, err := c.articleService.List(
		ctx.Context(),
		userId,
		services.ListArticlesInput{
			Tag:        input.Tag,
			Favorited:  input.Favorited,
			Limit:      input.Limit,
			Offset:     input.Offset,
			Authorname: input.Author,
		},
	)

	if err != nil {
		return err
	}

	return renderComponent(articleList(articlesProps{
		articles: articles,
		// TODO: get total articles count
		showPagination: len(articles) > 20,
		numOfPages:     len(articles) / 20,
		// TODO: get current page
		currentPage: 1,
	}), ctx)
}

type getFeedParams struct {
	Limit  int `query:"limit"`
	Offset int `query:"offset"`
}

func (c *Controller) getFeed(ctx *fiber.Ctx) error {
	input := new(getFeedParams)

	if err := ctx.QueryParser(input); err != nil {
		return err
	}

	userId, ok := ctx.Locals("userId").(int)

	if !ok {

		ctx.Response().Header.Add("HX-Push-Url", "false")
		ctx.Response().Header.Add("HX-Reswap", "none")

		return ctx.Redirect("/login", fiber.StatusSeeOther)
	}

	if input.Limit == 0 {
		input.Limit = 20
	}

	articles, err := c.articleService.List(
		ctx.Context(),
		userId,
		services.ListArticlesInput{
			Limit:  input.Limit,
			Offset: input.Offset,
			Feed:   true,
		},
	)

	if err != nil {
		if errors.Is(err, services.WarningNoFollowers) {
			return renderComponent(articleList(articlesProps{
				currentPage:    1,
				showPagination: false,
				articles:       articles,
				hasNoFollowing: true,
			}), ctx)
		}
		return err
	}

	return renderComponent(articleList(articlesProps{
		articles: articles,
		// TODO: get total articles count
		showPagination: len(articles) > 20,
		numOfPages:     len(articles) / 20,
		// TODO: get current page
		currentPage: 1,
	}), ctx)
}
