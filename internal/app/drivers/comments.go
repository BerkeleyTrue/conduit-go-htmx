package drivers

import (
	"context"

	"github.com/gofiber/fiber/v2"
)

func (c *Controller) getComments(fc *fiber.Ctx) error {
	ctx := context.Background()
	slug := fc.Params("slug")
	userId, ok := fc.Locals("userId").(int)

	if !ok {
		userId = 0
	}

	comments, err := c.commentService.GetBySlug(ctx, slug, userId)

	if err != nil {
		return err
	}

	props := commentsProps{
		comments: comments,
	}

	return renderComponent(commentsComp(props), fc)
}
