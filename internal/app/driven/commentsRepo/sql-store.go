package commentsRepo

import (
	"context"
	_ "embed"

	"github.com/berkeleytrue/conduit/internal/core/domain"
	"github.com/berkeleytrue/conduit/internal/infra/data/krono"
	"github.com/berkeleytrue/conduit/internal/utils"
	"github.com/jmoiron/sqlx"
	"go.uber.org/fx"
	"golang.org/x/exp/slog"
)

type CommentStore struct {
	*Queries
	db *sqlx.DB
}

var (
	//go:embed sql/schema.sql
	schema string
	// compile time check to make sure ArticleStore implements domain.CommentRepository
	_ domain.CommentRepository = (*CommentStore)(nil)

	Module = fx.Options(
		fx.Provide(fx.Annotate(
			newCommentStore,
			fx.As(new(domain.CommentRepository)),
		)),

		fx.Invoke(registerCommentSchema),
	)
)

func newCommentStore(_db *sqlx.DB) *CommentStore {
	return &CommentStore{
		Queries: New(_db),
		db:      _db,
	}
}

func registerCommentSchema(db *sqlx.DB) error {
	_, err := db.Exec(schema)

	if err != nil {
		return err
	}

	return nil
}

func formatToDomain(comment Comment) *domain.Comment {
	createdAt, err := krono.FromString(comment.CreatedAt)

	if err != nil {
		slog.Debug("error parsing createdAt", "error", err)
	}

	return &domain.Comment{
		CommentId: int(comment.ID),
		Body:      comment.Body,
		AuthorId:  int(comment.AuthorID),
		ArticleId: int(comment.ArticleID),
		CreatedAt: createdAt,
	}
}

func (r *CommentStore) Create(
	ctx context.Context,
	input domain.CommentCreateInput,
) (*domain.Comment, error) {
	now := krono.Now().ToString()

	if !input.CreatedAt.IsZero() {
		now = input.CreatedAt.ToString()
	}

	comment, err := r.create(ctx, createParams{
		Body:      input.Body,
		ArticleID: int64(input.ArticleId),
		AuthorID:  int64(input.AuthorId),
		CreatedAt: now,
	})

	if err != nil {
		return nil, err
	}

	return formatToDomain(comment), nil
}

func (r *CommentStore) GetById(
	ctx context.Context,
	commentId int,
) (*domain.Comment, error) {
	comment, err := r.getById(ctx, int64(commentId))

	if err != nil {
		return nil, err
	}

	return formatToDomain(comment), nil
}

func (r *CommentStore) GetByArticleId(
	ctx context.Context,
	articleId int,
) ([]*domain.Comment, error) {
	comments, err := r.getByArticleId(ctx, int64(articleId))

	if err != nil {
		return nil, err
	}

	return utils.Map(func(c Comment) *domain.Comment {
		return formatToDomain(c)
	}, comments), nil
}

func (r *CommentStore) GetByAuthorId(ctx context.Context, authorId int) ([]*domain.Comment, error) {
	comments, err := r.getByAuthorId(ctx, int64(authorId))

	if err != nil {
		return nil, err
	}

	return utils.Map(func(c Comment) *domain.Comment {
		return formatToDomain(c)
	}, comments), nil
}

func (r *CommentStore) Delete(ctx context.Context, commentId, userId int) error {
	return r.delete(ctx, deleteParams{
		ID:       int64(commentId),
		AuthorID: int64(userId),
	})
}
