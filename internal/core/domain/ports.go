package domain

import (
	"context"

	"github.com/berkeleytrue/conduit/internal/infra/data/krono"
	"github.com/berkeleytrue/conduit/internal/infra/data/password"
)

type (
	Updater[T any] func(u T) T

	UserCreateInput struct {
		Username       string
		Email          string
		HashedPassword password.HashedPassword
		CreatedAt      krono.Krono
	}

	UserRepository interface {
		// Create a new user
		Create(ctx context.Context, input UserCreateInput) (*User, error)
		// Get a user by id
		GetByID(ctx context.Context, id int) (*User, error)
		// Get a user by email
		GetByEmail(ctx context.Context, email string) (*User, error)
		// Get a user by username
		GetByUsername(ctx context.Context, username string) (*User, error)
		// Get authors a user is following
		GetFollowing(ctx context.Context, userId int) ([]int, error)
		// Update a user
		Update(ctx context.Context, userId int, updater Updater[User]) (*User, error)
		// A user follows an author
		Follow(ctx context.Context, userId, authorId int) (*User, error)
		// A user unfollows an author
		Unfollow(ctx context.Context, userId, authorId int) (*User, error)
	}

	ArticleCreateInput struct {
		Title       string
		Description string
		Body        string
		Tags        []string
		AuthorId    int
		CreatedAt   krono.Krono
	}

	ArticleListInput struct {
		AuthorId  int   // authorId
		Favorited int   // authorId
		Authors   []int // for user following authors
		Tag       string
		Limit     int
		Offset    int
	}

	ArticleRepository interface {
		Create(ctx context.Context, input ArticleCreateInput) (*Article, error)
		GetById(ctx context.Context, articleId int) (*Article, error)
		GetBySlug(ctx context.Context, mySlug string) (*Article, error)
		List(ctx context.Context, input ArticleListInput) ([]*Article, error)
		GetPopularTags(ctx context.Context) ([]string, error)
		Update(ctx context.Context, slug string, updater Updater[Article]) (*Article, error)
		Favorite(ctx context.Context, slug string, userId int) (*Article, error)
		Unfavorite(ctx context.Context, slug string, userId int) (*Article, error)
		Delete(ctx context.Context, slug string) error
	}

	CommentCreateInput struct {
		Body      string
		AuthorId  int // UserId
		ArticleId int // ArticleId
		CreatedAt krono.Krono
	}

	CommentRepository interface {
		Create(ctx context.Context, input CommentCreateInput) (*Comment, error)
		GetById(ctx context.Context, commentId int) (*Comment, error)
		GetByArticleId(ctx context.Context, articleId int) ([]*Comment, error)
		GetByAuthorId(ctx context.Context, authorId int) ([]*Comment, error)
		Delete(ctx context.Context, commentId int) error
	}
)
