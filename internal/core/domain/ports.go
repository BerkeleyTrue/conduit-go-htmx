package domain

import (
	"github.com/berkeleytrue/conduit/internal/infra/data/password"
)

type (
	Updater[T any] func(u *T) *T

	UserCreateInput struct {
		Username       string
		Email          string
		Password       password.Password
		HashedPassword password.HashedPassword
	}

	UserRepository interface {
		Create(UserCreateInput) (*User, error)
		GetByID(id int) (*User, error)
		GetByEmail(email string) (*User, error)
		GetByUsername(username string) (*User, error)
		Update(userId int, updater Updater[User]) (*User, error)
		Follow(userId, authorId int) (*User, error)
		Unfollow(userId, authorId int) (*User, error)
	}

	ArticleCreateInput struct {
		Title       string
		Description string
		Body        string
		Tags        []string
		AuthorId    int
	}

	ArticleListInput struct {
		Tag       string
		Author    string // authorId
		Favorited string // authorId
		Limit     int
		Offset    int
	}

	ArticleRepository interface {
		Create(input ArticleCreateInput) (*Article, error)
		GetById(articleId string) (*Article, error)
		GetBySlug(mySlug string) (*Article, error)
		List(input ArticleListInput) ([]*Article, error)
		Update(slug string, updater Updater[Article]) (*Article, error)
		Favorite(slug string, userId int) (*Article, error)
		Delete(slug string) error
	}

	CommentCreateInput struct {
		Body      string
		AuthorId  int // UserId
		ArticleId int // ArticleId
	}

	CommentRepository interface {
		Create(input CommentCreateInput) (*Comment, error)
		GetById(commentId string) (*Comment, error)
		GetByArticleId(articleId string) ([]*Comment, error)
		Update(commentId string, updater Updater[Comment]) (*Comment, error)
		Delete(commentId string) error
	}
)
