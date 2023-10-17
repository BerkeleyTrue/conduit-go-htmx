package articlesRepo

import (
	"context"
	"fmt"
	"strings"

	"github.com/berkeleytrue/conduit/internal/core/domain"
	"github.com/berkeleytrue/conduit/internal/infra/data/krono"
	"github.com/berkeleytrue/conduit/internal/infra/data/slug"
)

func (s *ArticleStore) CreateTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := s.db.BeginTx(ctx, nil)

	if err != nil {
		return err
	}

	q := New(tx)

	err = fn(q)

	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %s, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

func formatToDomain(article Article) *domain.Article {
	createdAt, err := krono.FromString(article.CreatedAt)

	if err != nil {
		fmt.Printf("error parsing createdAt: %s", err)
	}

	updatedAt, err := krono.FromNullString(article.UpdatedAt)

	if err != nil {
		fmt.Printf("error parsing createdAt: %s", err)
	}

	return &domain.Article{
		ArticleId:   int(article.ID),
		Slug:        slug.NewSlug(article.Slug),
		Title:       article.Title,
		Description: article.Description,
		Body:        article.Body,
		AuthorId:    int(article.AuthorID),
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
}

func formatRowToDomain(row listRow) *domain.Article {
	createdAt, err := krono.FromString(row.CreatedAt)

	if err != nil {
		fmt.Printf("error parsing createdAt: %s", err)
	}

	updatedAt, err := krono.FromNullString(row.UpdatedAt)

	if err != nil {
		fmt.Printf("error parsing createdAt: %s", err)
	}

	tags := []string{}

	if row.Tags != "" {
		tags = strings.Split(row.Tags, ",")
	}

	return &domain.Article{
		ArticleId:   int(row.ID),
		Slug:        slug.NewSlug(row.Slug),
		Title:       row.Title,
		Description: row.Description,
		Body:        row.Body,
		AuthorId:    int(row.AuthorID),
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
		Tags:        tags,
	}
}
