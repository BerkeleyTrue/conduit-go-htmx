package domain

import (
	"github.com/berkeleytrue/conduit/internal/infra/data/krono"
	"github.com/berkeleytrue/conduit/internal/infra/data/password"
	"github.com/berkeleytrue/conduit/internal/infra/data/slug"
)

type (
	User struct {
		UserId    int                     `json:"userId" db:"id"`
		Username  string                  `json:"username"`
		Email     string                  `json:"email"`
		Password  password.HashedPassword `json:"password"`  // hashed password
		Followers []int                   `json:"following"` // []UserId
		Bio       string                  `json:"bio"`       // nullable
		Image     string                  `json:"image"`
		CreatedAt krono.Krono             `json:"createdAt" db:"created_at"`
		UpdatedAt krono.Krono             `json:"updatedAt" db:"updated_at"`
	}

	Article struct {
		ArticleId int `json:"articleId" db:"id"`
		AuthorId  int `json:"authorId" db:"author_id"`

		Title string    `json:"title" db:"title"`
		Slug  slug.Slug `json:"slug" db:"slug"`

		Description string   `json:"description" db:"description"`
		Body        string   `json:"body" db:"body"`
		Tags        []string `json:"Tags"`

		CreatedAt krono.Krono `json:"createdAt" db:"created_at"`
		UpdatedAt krono.Krono `json:"updatedAt" db:"updated_at"`
	}

	Comment struct {
		CommentId string `json:"commentId"`
		AuthorId  string `json:"authorId"`
		ArticleId string `json:"articleId"`
		Body      string `json:"body"`
		CreatedAt string `json:"createdAt"`
		UpdatedAt string `json:"updatedAt"`
	}
)
