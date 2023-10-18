package domain

import (
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
		Create(UserCreateInput) (*User, error)
	  // Get a user by id
		GetByID(id int) (*User, error)
	  // Get a user by email
		GetByEmail(email string) (*User, error)
	  // Get a user by username
		GetByUsername(username string) (*User, error)
	  // Get authors a user is following
	  GetFollowing(userId int) ([]int, error)
	  // Update a user
		Update(userId int, updater Updater[User]) (*User, error)
	  // A user follows an author
		Follow(userId, authorId int) (*User, error)
	  // A user unfollows an author
		Unfollow(userId, authorId int) (*User, error)
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
		Create(input ArticleCreateInput) (*Article, error)
		GetById(articleId int) (*Article, error)
		GetBySlug(mySlug string) (*Article, error)
		List(input ArticleListInput) ([]*Article, error)
		GetPopularTags() ([]string, error)
		Update(slug string, updater Updater[Article]) (*Article, error)
		Favorite(slug string, userId int) (*Article, error)
		Unfavorite(slug string, userId int) (*Article, error)
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
