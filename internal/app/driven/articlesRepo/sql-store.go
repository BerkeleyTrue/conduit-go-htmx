package articlesRepo

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"go.uber.org/fx"

	"github.com/berkeleytrue/conduit/internal/core/domain"
	"github.com/berkeleytrue/conduit/internal/infra/data/krono"
	"github.com/berkeleytrue/conduit/internal/infra/data/slug"
	"github.com/berkeleytrue/conduit/internal/infra/db"
)

type (
	ArticleStore struct {
		db.SqlStore
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

    CREATE TABLE IF NOT EXISTS favorites (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        user_id INTEGER NOT NULL,
        article_id INTEGER NOT NULL,
        UNIQUE(user_id, article_id),

        FOREIGN KEY(user_id) REFERENCES users(id),
        FOREIGN KEY(article_id) REFERENCES articles(id)
    )

  `
	_ domain.ArticleRepository = (*ArticleStore)(nil)

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
		SqlStore: db.SqlStore{Db: _db},
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

	err := s.CreateTx(func(tx *sqlx.Tx) error {
		query := `
      INSERT INTO articles (slug, title, description, body, author_id, created_at, updated_at)
      VALUES (:slug, :title, :description, :body, :author_id, :created_at, :updated_at);
    `

		res, err := tx.NamedExec(query, article)

		if err != nil {
			return fmt.Errorf("error creating article: %w", err)
		}

		articleId, err := res.LastInsertId()

		if err != nil {
			return fmt.Errorf("error getting article id: %w", err)
		}

		article.ArticleId = int(articleId)

		// insert tags into tags table and create article_tag records
		for _, tag := range input.Tags {
			query = `
        INSERT INTO tags (tag)
        VALUES ($1)
        ON CONFLICT (tag) DO UPDATE SET tag = $1
      `
			_, err = tx.Exec(query, tag)

			if err != nil {
				return fmt.Errorf("error creating article: %w", err)
			}

			query = `
        INSERT INTO article_tags (article_id, tag_id)
        VALUES ($1, (SELECT id FROM tags WHERE tag = $2))
      `

			_, err = tx.Exec(query, articleId, tag)

			if err != nil {
				return fmt.Errorf("error creating article: %w", err)
			}
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error in commit: %w", err)
	}

	err = s.Db.Get(&article, "SELECT * FROM articles WHERE slug = $1", slug)

	if err != nil {
		return nil, fmt.Errorf("error getting res article: %w", err)
	}

	fmt.Printf("%+v\n", article)

	return &article, err
}

func (s *ArticleStore) GetById(articleId string) (*domain.Article, error) {
	var article domain.Article
	err := s.Db.Get(&article, "SELECT * FROM articles WHERE id = $1", articleId)

	if err != nil {
		return nil, fmt.Errorf("error getting article: %w", err)
	}

	return &article, nil
}

func (s *ArticleStore) GetBySlug(mySlug string) (*domain.Article, error) {
	var article domain.Article
	err := s.Db.Get(&article, "SELECT * FROM articles WHERE slug = $1", mySlug)

	if err != nil {
		return nil, fmt.Errorf("error getting article: %w", err)
	}

	return &article, nil
}

func (s *ArticleStore) List(
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

	err := s.Db.Select(&articles, query, input.Limit, input.Offset)

	if err != nil {
		fmt.Printf("error getting articles: %v\n", err)
		return nil, err
	}

	return articles, nil
}

func (s *ArticleStore) Update(
	slug string,
	updater domain.Updater[domain.Article],
) (*domain.Article, error) {

	article, err := s.GetBySlug(slug)

	if err != nil {
		return nil, fmt.Errorf("sql-store: error getting article: %w", err)
	}

	updater(article)

	_, err = s.Db.NamedExec(`
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

func (s *ArticleStore) Favorite(slug string, userId int) (*domain.Article, error) {

	article, err := s.GetBySlug(slug)

	if err != nil {
		return nil, fmt.Errorf("sql-store: error getting article: %w", err)
	}

	_, err = s.Db.Exec(`
    INSERT INTO favorites (user_id, article_id)
    VALUES ($1, $2)
  `, userId, article.ArticleId)

	if err != nil {
		return nil, fmt.Errorf("sql-store: error favoriting article: %w", err)
	}

	return article, nil
}

func (s *ArticleStore) Unfavorite(slug string, userId int) (*domain.Article, error) {
	article, err := s.GetBySlug(slug)

	if err != nil {
		return nil, fmt.Errorf("sql-store: error getting article: %w", err)
	}

	_, err = s.Db.Exec(`
      DELETE FROM favorites
      WHERE user_id = $1 AND article_id = $2
    `, userId, article.ArticleId)

	if err != nil {
		return nil, fmt.Errorf("sql-store: error unfavoriting article: %w", err)
	}

	return article, nil
}

func (s *ArticleStore) Delete(slug string) error {
	_, err := s.Db.Exec("DELETE FROM articles WHERE slug = $1", slug)

	if err != nil {
		return fmt.Errorf("sql-store: error deleting article: %w", err)
	}

	return nil
}
