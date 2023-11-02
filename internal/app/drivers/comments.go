package drivers

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/gofiber/fiber/v2"

	"github.com/berkeleytrue/conduit/internal/core/services"
	"github.com/berkeleytrue/conduit/internal/infra/session"
)

type createCommentInput struct {
	Body string `json:"body"`
}

func (i *createCommentInput) validate() error {
	return validation.ValidateStruct(
		i,
		validation.Field(
			&i.Body,
			validation.Required,
			validation.Length(1, 254),
		),
	)
}

func (c *Controller) createComment(fc *fiber.Ctx) error {
	ctx := context.Background()
	slug := fc.Params("slug")
	input := createCommentInput{}

	userId, ok := fc.Locals("userId").(int)

	if !ok {
		return fiber.ErrUnauthorized
	}

	if err := fc.BodyParser(&input); err != nil {
		return err
	}

	if err := input.validate(); err != nil {
		fc.Response().Header.Add("HX-Push-Url", "false")
		fc.Response().Header.Add("HX-Reswap", "none")

		return renderComponent(listErrors(err.(validation.Errors)), fc)
	}

	comment, err := c.commentService.Create(ctx, services.CommentCreateInput{
		ArticleSlug: slug,
		Body:        input.Body,
		AuthorId:    userId,
	})

	if err != nil {
		fc.Response().Header.Add("HX-Push-Url", "false")
		fc.Response().Header.Add("HX-Reswap", "none")

		return renderComponent(listErrors(map[string]error{"comment": err}), fc)
	}

	session.AddFlash(fc, session.Success, "Comment created successfully!")

	return renderComponent(commentComp(commentProps{
		CommentOutput: *comment,
		slug:          slug,
	}), fc)
}

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
		slug:     slug,
	}

	return renderComponent(commentsComp(props), fc)
}

func (c *Controller) deleteComment(fc *fiber.Ctx) error {
	ctx := context.Background()
	commentId, err := fc.ParamsInt("id")

	if err != nil {
		return err
	}

	userId, ok := fc.Locals("userId").(int)

	if !ok {
		return fiber.ErrUnauthorized
	}

	err = c.commentService.Delete(ctx, commentId, userId)

	if err != nil {
		return err
	}

	return fc.SendStatus(fiber.StatusNoContent)
}
