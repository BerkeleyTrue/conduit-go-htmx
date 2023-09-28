package domain

import (
	"time"

	"github.com/berkeleytrue/conduit/internal/infra/data/password"
)

type (
	User struct {
		UserId    string                  `json:"userId"`
		Username  string                  `json:"username"`
		Email     string                  `json:"email"`
		Password  password.HashedPassword `json:"password"`  // hashed password
		Followers []string                `json:"following"` // []UserId
		Bio       string                  `json:"bio"`       // nullable
		Image     string                  `json:"image"`
		CreatedAt time.Time               `json:"createdAt"`
		UpdatedAt time.Time               `json:"updatedAt"`
	}

	Article struct {
		AuthorId  string `json:"authorId"`
		ArticleId string `json:"articleId"`

		Title string `json:"title"`
		Slug  string `json:"slug"`

		Description string   `json:"description"`
		Body        string   `json:"body"`
		Tags        []string `json:"Tags"`

		CreatedAt time.Time `json:"createdAt"`
		UpdatedAt time.Time `json:"updatedAt"`

		FavoritedBy []string `json:"favoritedBy"` // []UserId
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
