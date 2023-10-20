package drivers

import (
	"context"

	"github.com/gofiber/fiber/v2"
)

func (c *Controller) GetPopularTags(fc *fiber.Ctx) error {
	tags, err := c.articleService.GetPopularTags(context.Background())

	if err != nil {
		return err
	}

	return renderComponent(popularTags(tags), fc)
}
