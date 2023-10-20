package drivers

import "github.com/gofiber/fiber/v2"

func (c *Controller) GetPopularTags(fc *fiber.Ctx) error {
	tags, err := c.articleService.GetPopularTags(fc.Context())

	if err != nil {
		return err
	}

	return renderComponent(popularTags(tags), fc)
}
