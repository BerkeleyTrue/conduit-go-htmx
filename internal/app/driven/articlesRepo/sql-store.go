package articlesRepo

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"go.uber.org/fx"

	"github.com/berkeleytrue/conduit/internal/core/domain"
	"github.com/berkeleytrue/conduit/internal/infra/data/krono"
	"github.com/berkeleytrue/conduit/internal/infra/data/slug"
)

type (
	SqlStore struct {
		db *sqlx.DB
	}
)

var (
	//sqlite datatypes
	schema = `
    CREATE TABLE IF NOT EXISTS articles (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        slug TEXT NOT NULL UNIQUE,
        title TEXT NOT NULL,
        description TEXT NOT NULL,
        body TEXT NOT NULL,
        author_id INTEGER NOT NULL,
        created_at TEXT NOT NULL,
        updated_at TEXT NOT NULL,
        FOREIGN KEY(author_id) REFERENCES users(id)
    );

    CREATE TABLE IF NOT EXISTS tags (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        tag TEXT NOT NULL UNIQUE
    );

    CREATE TABLE IF NOT EXISTS article_tags (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        article_id INTEGER NOT NULL,
        tag_id INTEGER NOT NULL,
        UNIQUE(article_id, tag_id),
        FOREIGN KEY(article_id) REFERENCES articles(id),
        FOREIGN KEY(tag_id) REFERENCES tags(id)
    );

    CREATE TABLE IF NOT EXISTS favoires (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        user_id INTEGER NOT NULL,
        article_id INTEGER NOT NULL,
        UNIQUE(user_id, article_id),

        FOREIGN KEY(user_id) REFERENCES users(id),
        FOREIGN KEY(article_id) REFERENCES articles(id)
    )

  `
	_ domain.ArticleRepository = (*SqlStore)(nil)
)

func NewSqlStore(db *sqlx.DB) *SqlStore {
	return &SqlStore{
		db: db,
	}
}

func RegisterArticleSchema(lc fx.Lifecycle, db *sqlx.DB) {
	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			_, err := db.Exec(schema)

			if err != nil {
				return err
			}

			return nil
		},
	})
}

func (s *SqlStore) Create(
	input domain.ArticleCreateInput,
) (*domain.Article, error) {
	slug := slug.NewSlug(input.Title)
	now := krono.Now()
	article := domain.Article{
		Slug:        slug,
		Title:       input.Title,
		Description: input.Description,
		Body:        input.Body,
		Tags:        input.Tags,
		AuthorId:    input.AuthorId,
		CreatedAt:   now,
	}

	// TODO: add tag creation and article_tag creation
	query := `
    INSERT INTO articles (slug, title, description, body, created_at, updated_at)
    VALUES (:slug, :title, :description, :body, :created_at, :updated_at)

  `

	_, err := s.db.NamedExec(query, article)

	if err != nil {
		return nil, fmt.Errorf("error creating article: %w", err)
	}

	err = s.db.Get(&article, "SELECT * FROM articles WHERE slug = $1", slug)

	if err != nil {
		return nil, fmt.Errorf("error getting article: %w", err)
	}

	return &article, err
}

func (s *SqlStore) GetById(articleId string) (*domain.Article, error) {
	var article domain.Article
	err := s.db.Get(&article, "SELECT * FROM articles WHERE id = $1", articleId)

	if err != nil {
		return nil, fmt.Errorf("error getting article: %w", err)
	}

	return &article, nil
}

func (s *SqlStore) GetBySlug(mySlug string) (*domain.Article, error) {
	var article domain.Article
	err := s.db.Get(&article, "SELECT * FROM articles WHERE slug = $1", mySlug)

	if err != nil {
		return nil, fmt.Errorf("error getting article: %w", err)
	}

	return &article, nil
}

func (s *SqlStore) List(
	input domain.ArticleListInput,
) ([]*domain.Article, error) {
	var articles []*domain.Article
	// TODO: add tags, author, favorited
	query := `
    SELECT * FROM articles
    ORDER BY created_at DESC
    LIMIT $1
    OFFSET $2
  `

	err := s.db.Select(&articles, query, input.Limit, input.Offset)

	if err != nil {
		fmt.Printf("error getting articles: %v\n", err)
		return nil, err
	}

	return articles, nil
}

func (s *SqlStore) Update(
	slug string,
	updater domain.Updater[domain.Article],
) (*domain.Article, error) {

	article, err := s.GetBySlug(slug)

	if err != nil {
		return nil, fmt.Errorf("sql-store: error getting article: %w", err)
	}

	updater(article)

	_, err = s.db.NamedExec(`
    UPDATE articles SET
      title = :title,
      description = :description,
      body = :body,
      updated_at = :updated_at
    WHERE slug = :slug
  `, article)

	if err != nil {
		return nil, fmt.Errorf("sql-store: error updating article: %w", err)
	}

	return article, nil
}

func (s *SqlStore) Favorite(slug string, userId int) (*domain.Article, error) {

	article, err := s.GetBySlug(slug)

	if err != nil {
		return nil, fmt.Errorf("sql-store: error getting article: %w", err)
	}

	_, err = s.db.Exec(`
    INSERT INTO favorites (user_id, article_id)
    VALUES ($1, $2)
  `, userId, article.ArticleId)

	if err != nil {
		return nil, fmt.Errorf("sql-store: error favoriting article: %w", err)
	}

	return article, nil
}

func (s *SqlStore) Delete(slug string) error {
	_, err := s.db.Exec("DELETE FROM articles WHERE slug = $1", slug)

	if err != nil {
		return fmt.Errorf("sql-store: error deleting article: %w", err)
	}

	return nil
}
