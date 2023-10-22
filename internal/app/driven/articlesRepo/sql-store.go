package articlesRepo

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"

	"github.com/jmoiron/sqlx"
	"go.uber.org/fx"

	"github.com/berkeleytrue/conduit/internal/core/domain"
	"github.com/berkeleytrue/conduit/internal/infra/data/krono"
	"github.com/berkeleytrue/conduit/internal/infra/data/slug"
)

type (
	ArticleStore struct {
		*Queries
		db *sqlx.DB
	}
)

var (
	//go:embed sql/schema.sql
	schema string
	_      domain.ArticleRepository = (*ArticleStore)(nil)

	Module = fx.Options(
		fx.Provide(fx.Annotate(
			newSqlStore,
			fx.As(new(domain.ArticleRepository)),
		)),

		fx.Invoke(registerArticleSchema),
	)
)

func newSqlStore(_db *sqlx.DB) *ArticleStore {
	return &ArticleStore{
		Queries: New(_db),
		db:      _db,
	}
}

func registerArticleSchema(db *sqlx.DB) error {
	_, err := db.Exec(schema)

	if err != nil {
		return err
	}

	return nil
}

func (s *ArticleStore) Create(
	ctx context.Context,
	input domain.ArticleCreateInput,
) (*domain.Article, error) {
	slug := slug.NewSlug(input.Title)

	err := s.CreateTx(ctx, func(q *Queries) error {
		params := createParams{
			Slug:        slug.String(),
			Title:       input.Title,
			Description: input.Description,
			Body:        input.Body,
			AuthorID:    int64(input.AuthorId),
			CreatedAt:   input.CreatedAt.ToString(),
			UpdatedAt:   sql.NullString{},
		}

		article, err := q.create(ctx, params)

		if err != nil {
			return fmt.Errorf("error creating article: %w", err)
		}

		// insert tags into tags table and create article_tag records
		for _, tag := range input.Tags {
			_, err = q.createTag(ctx, tag)

			if err != nil {
				return fmt.Errorf("error creating tag: %w", err)
			}

			_, err = q.createArticleTag(ctx, createArticleTagParams{
				ArticleID: article.ID,
				Tag:       tag,
			})

			if err != nil {
				return fmt.Errorf("error creating article tag: %w", err)
			}
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error in commit: %w", err)
	}

	article, err := s.getBySlug(ctx, slug.String())

	if err != nil {
		return nil, fmt.Errorf("error getting res article: %w", err)
	}

	return formatToDomain(article), nil
}

func (s *ArticleStore) GetById(
	ctx context.Context,
	articleId int,
) (*domain.Article, error) {
	article, err := s.getById(ctx, int64(articleId))

	if err != nil {
		return nil, fmt.Errorf("error getting article: %w", err)
	}

	return formatToDomain(article), nil
}

func (s *ArticleStore) GetBySlug(
	ctx context.Context,
	mySlug string) (*domain.Article, error) {
	article, err := s.getBySlug(ctx, mySlug)

	if err != nil {
		return nil, fmt.Errorf("error getting article: %w", err)
	}

	return formatToDomain(article), nil
}

// TODO: get number of favorites per article
func (s *ArticleStore) List(
	ctx context.Context,
	input domain.ArticleListInput,
) ([]*domain.Article, error) {
	params := listParams{
		Limit:  int64(input.Limit),
		Offset: int64(input.Offset),
		AuthorID: sql.NullInt64{
			Valid: input.AuthorId != 0,
			Int64: int64(input.AuthorId),
		},
		Tag: sql.NullString{
			Valid:  input.Tag != "",
			String: input.Tag,
		},
		Favorited: sql.NullInt64{
			Valid: input.Favorited != 0,
			Int64: int64(input.Favorited),
		},
	}

	rows, err := s.list(ctx, params)

	if err != nil {
		fmt.Printf("error getting articles: %v\n", err)
		return nil, err
	}

	articles := make([]*domain.Article, len(rows))

	for idx, row := range rows {
		articles[idx] = formatRowToDomain(row)
	}

	return articles, nil
}

func (s *ArticleStore) GetPopularTags(ctx context.Context) ([]string, error) {
	var tags []string

	tags, err := s.getPopularTags(ctx)

	if err != nil {
		return nil, fmt.Errorf("error getting tags: %w", err)
	}

	return tags, nil
}

func (s *ArticleStore) GetNumOfFavorites(
	ctx context.Context,
	articleId int,
) (int, error) {
	count, err := s.getNumOfFavorites(ctx, int64(articleId))

	if err != nil {
		return 0, fmt.Errorf("error getting number of favorites: %w", err)
	}

	return int(count), nil
}

// TODO: update tags
func (s *ArticleStore) Update(
	ctx context.Context,
	slug string,
	updater domain.Updater[domain.Article],
) (*domain.Article, error) {
	article, err := s.getBySlug(ctx, slug)

	if err != nil {
		return nil, fmt.Errorf("sql-store: error getting article: %w", err)
	}

	updates := updater(*formatToDomain(article))

	params := updateParams{
		UpdatedAt: krono.Now().ToNullString(),
		Slug:      article.Slug,
	}

	if updates.Title != "" {
		params.Title = updates.Title
	} else {
		params.Title = article.Title
	}

	if updates.Description != "" {
		params.Description = updates.Description
	} else {
		params.Description = article.Description
	}

	if updates.Body != "" {
		params.Body = updates.Body
	} else {
		params.Body = article.Body
	}

	updatedArticle, err := s.update(ctx, params)

	if err != nil {
		return nil, fmt.Errorf("sql-store: error updating article: %w", err)
	}

	return formatToDomain(updatedArticle), nil
}

func (s *ArticleStore) Favorite(
	ctx context.Context,
	slug string,
	userId int,
) (*domain.Article, error) {

	article, err := s.getBySlug(ctx, slug)

	if err != nil {
		return nil, fmt.Errorf("sql-store: error getting article: %w", err)
	}

	_, err = s.favorite(ctx, favoriteParams{
		UserID:    int64(userId),
		ArticleID: article.ID,
	})

	if err != nil {
		return nil, fmt.Errorf("sql-store: error favoriting article: %w", err)
	}

	return formatToDomain(article), nil
}

func (s *ArticleStore) Unfavorite(
	ctx context.Context,
	slug string,
	userId int,
) (*domain.Article, error) {
	article, err := s.getBySlug(ctx, slug)

	if err != nil {
		return nil, fmt.Errorf("sql-store: error getting article: %w", err)
	}

	_, err = s.unfavorite(ctx, unfavoriteParams{
		UserID:    int64(userId),
		ArticleID: article.ID,
	})

	if err != nil {
		return nil, fmt.Errorf("sql-store: error unfavoriting article: %w", err)
	}

	return formatToDomain(article), nil
}

func (s *ArticleStore) Delete(ctx context.Context, slug string) error {
	_, err := s.delete(ctx, slug)

	if err != nil {
		return fmt.Errorf("sql-store: error deleting article: %w", err)
	}

	return nil
}
