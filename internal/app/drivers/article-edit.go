package drivers

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/gofiber/fiber/v2"

	"github.com/berkeleytrue/conduit/internal/core/services"
	"github.com/berkeleytrue/conduit/internal/infra/session"
)

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

type updateArticleInput struct {
	title       string
	description string
	body        string
	tags        string
}

func (i *updateArticleInput) validate() error {
	return validation.ValidateStruct(
		i,
		validation.Field(
			&i.title,
			validation.Required,
			validation.Length(1, 254),
		),
		validation.Field(
			&i.description,
			validation.Required,
			validation.Length(1, 254),
		),
		validation.Field(
			&i.body,
			validation.Required,
		),
		validation.Field(
			&i.tags,
		),
	)
}

func (c *Controller) updateArticle(fc *fiber.Ctx) error {
	ctx := context.Background()
	slug := fc.Params("slug")
	input := updateArticleInput{}

	username, ok := fc.Locals("username").(string)

	if !ok {
		return fiber.ErrUnauthorized
	}

	input.title = fc.FormValue("title")
	input.description = fc.FormValue("description")
	input.body = fc.FormValue("body")
	input.tags = fc.FormValue("tags")

	if err := input.validate(); err != nil {
		fc.Response().Header.Add("HX-Push-Url", "false")
		fc.Response().Header.Add("HX-Reswap", "none")

		return renderComponent(listErrors(err.(validation.Errors)), fc)
	}

	article, err := c.articleService.Update(
		ctx,
		slug,
		username,
		services.ArticleUpdateInput{
			Title:       input.title,
			Description: input.description,
			Body:        input.body,
		},
	)

	if err != nil {
		fc.Response().Header.Add("HX-Reswap", "none")

		return renderComponent(listErrors(map[string]error{"article": err}), fc)
	}

	session.AddFlash(fc, session.Success, "Article updated successfully!")

	return fc.Redirect("/editor/"+article.Slug.String(), fiber.StatusSeeOther)
}
