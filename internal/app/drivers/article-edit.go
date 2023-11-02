package drivers

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/gofiber/fiber/v2"
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
  title string
  description string
  body string
  tags string
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
}
