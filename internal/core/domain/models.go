package domain

import (
	"github.com/berkeleytrue/conduit/internal/infra/data/krono"
	"github.com/berkeleytrue/conduit/internal/infra/data/password"
	"github.com/berkeleytrue/conduit/internal/infra/data/slug"
)

type (
	User struct {
		UserId    int
		Username  string
		Email     string
		Password  password.HashedPassword
		Followers []int
		Bio       string
		Image     string
		CreatedAt krono.Krono
		UpdatedAt krono.Krono
	}

	Article struct {
		ArticleId   int
		AuthorId    int
		Title       string
		Slug        slug.Slug
		Description string
		Body        string
		Tags        []string
		CreatedAt   krono.Krono
		UpdatedAt   krono.Krono
	}

	Comment struct {
		CommentId int
		AuthorId  int
		ArticleId int
		Body      string
		CreatedAt krono.Krono
	}
)
