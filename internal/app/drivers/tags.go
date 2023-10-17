package drivers

import "github.com/gofiber/fiber/v2"

func (c *Controller) GetPopularTags(ctx *fiber.Ctx) error {
	tags, err := c.articleService.GetPopularTags()

	if err != nil {
		return err
	}

	return renderComponent(popularTags(tags), ctx)
}
