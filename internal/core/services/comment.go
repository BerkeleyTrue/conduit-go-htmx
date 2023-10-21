package services

import (
	"context"
	"errors"
	"os"

	"golang.org/x/exp/slog"

	"github.com/berkeleytrue/conduit/internal/core/domain"
	"github.com/berkeleytrue/conduit/internal/infra/data/krono"
)

type (
	CommentService struct {
		repo           domain.CommentRepository
		userService    *UserService
		articleService *ArticleService
		log            *slog.Logger
	}
	CommentOutput struct {
		Id        int
		Body      string
		CreatedAt krono.Krono
		Author    PublicProfile
		IsAuthor  bool
	}
)

func newCommentService(
	repo domain.CommentRepository,
	userService *UserService,
	articleService *ArticleService,
) *CommentService {
	return &CommentService{
		repo:           repo,
		userService:    userService,
		articleService: articleService,
		log: slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})).WithGroup("services").WithGroup("comment"),
	}
}

func formatToCommentOutput(
	comment *domain.Comment,
	profile *PublicProfile,
	userId int,
) *CommentOutput {
	return &CommentOutput{
		Id:        comment.CommentId,
		Body:      comment.Body,
		CreatedAt: comment.CreatedAt,
		Author:    *profile,
		IsAuthor:  comment.AuthorId == userId,
	}
}

type CommentCreateInput struct {
	AuthorId    int
	ArticleSlug string
	Body        string
}

var ErrorNoArticleFound = errors.New("no article found for slug")

func (s *CommentService) Create(
	ctx context.Context,
	input CommentCreateInput,
) (*CommentOutput, error) {
	articleId, err := s.articleService.GetIdFromSlug(ctx, input.ArticleSlug)

	if err != nil {
		return nil, ErrorNoArticleFound
	}

	comment, err := s.repo.Create(ctx, domain.CommentCreateInput{
		ArticleId: articleId,
		Body:      input.Body,
		AuthorId:  input.AuthorId,
	})

	if err != nil {
		return nil, err
	}

	profile, err := s.userService.GetProfile(ctx, input.AuthorId, "", 0)

	if err != nil {
		return nil, errors.New("no user found with that id")
	}

	return formatToCommentOutput(comment, profile, comment.AuthorId), nil
}

func (s *CommentService) GetBySlug(
	ctx context.Context,
	slug string,
	userId int,
) ([]*CommentOutput, error) {
	articleId, err := s.articleService.GetIdFromSlug(ctx, slug)

	if err != nil {
		return nil, err
	}

	comments, err := s.repo.GetByArticleId(ctx, articleId)

	if err != nil {
		return nil, err
	}

	commentsOut := []*CommentOutput{}

	for _, comment := range comments {
		profile, err := s.userService.GetProfile(
			ctx,
			comment.AuthorId,
			"",
			userId,
		)

		if err != nil {
			s.log.Debug("error getting profile", "error", err)
			continue
		}

		commentsOut = append(
			commentsOut,
			formatToCommentOutput(comment, profile, userId),
		)
	}

	return commentsOut, nil

}

func (s *CommentService) Delete(
	ctx context.Context,
	commentId int,
	userId int,
) error {
	return s.repo.Delete(ctx, commentId, userId)
}
