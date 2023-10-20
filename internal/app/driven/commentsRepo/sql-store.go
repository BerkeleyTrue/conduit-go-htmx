package commentsRepo

import (
	_ "embed"

	"github.com/berkeleytrue/conduit/internal/core/domain"
	"github.com/jmoiron/sqlx"
	"go.uber.org/fx"
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
			fx.As(new(domain.ArticleRepository)),
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
