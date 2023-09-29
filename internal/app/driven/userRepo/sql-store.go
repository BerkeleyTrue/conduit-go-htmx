package userRepo

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"go.uber.org/fx"

	"github.com/berkeleytrue/conduit/internal/core/domain"
)

type (
	SqlStore struct {
		db *sqlx.DB
	}
)

var (
	//sqlite datatypes
	schema = `
    CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        username TEXT NOT NULL,
        email TEXT NOT NULL,
        password TEXT NOT NULL,
        bio TEXT,
        image TEXT,
        created_at TEXT NOT NULL,
        updated_at TEXT NOT NULL
    );

    CREATE TABLE IF NOT EXISTS followers (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        user_id INTEGER NOT NULL,
        follower_id INTEGER NOT NULL,
        created_at TEXT NOT NULL,
        UNIQUE(user_id, follower_id)
    );
  `
)

func RegisterUserSchema(lc fx.Lifecycle, db *sqlx.DB) {
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

func NewSqlStore(db *sqlx.DB) *SqlStore {
	return &SqlStore{
		db: db,
	}
}

// get followers for a user
func (s *SqlStore) getFollowers(userId int8) ([]int8, error) {
  var followers []int8
  err := s.db.Select(&followers, "SELECT follower_id FROM followers WHERE user_id = $1", userId)

  if err != nil {
    return nil, err
  }

  return followers, nil
}

func (s *SqlStore) Create(input domain.UserCreateInput) (*domain.User, error) {
	now := time.Now()
	user := domain.User{
		Username:  input.Username,
		Email:     input.Email,
		Password:  input.HashedPassword,
		Bio:       "",
		Image:     "",
		CreatedAt: now,
	}

	query := `
    INSERT INTO users (username, email, password, bio, image, created_at, updated_at)
    VALUES (:username, :email, :password, :bio, :image, :created_at, :updated_at)
  `
	_, err := s.db.NamedExec(query, user)

	return &user, err
}

func (s *SqlStore) GetByID(id string) (*domain.User, error) {
	var user domain.User
	err := s.db.Get(&user, "SELECT * FROM users WHERE id = $1 LIMIT 1", id)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *SqlStore) GetByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := s.db.Get(&user, "SELECT * FROM users WHERE email = $1 LIMIT 1", email)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *SqlStore) GetByUsername(username string) (*domain.User, error) {
	var user domain.User
	err := s.db.Get(&user, "SELECT * FROM users WHERE username = $1 LIMIT 1", username)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *SqlStore) Update(
	userId string,
	updater domain.Updater[domain.User],
) (*domain.User, error) {
	var user domain.User
	err := s.db.Get(&user, "SELECT * FROM users WHERE id = $1 LIMIT 1", userId)
	if err != nil {
		return nil, err
	}
	updatedUser := updater(&user)

	_, err = s.db.NamedExec(`
    UPDATE users
    SET username = :username,
        email = :email,
        password = :password,
        bio = :bio,
        image = :image,
        updated_at = :updated_at
    WHERE id = :id
  `, updatedUser)

	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}

func (s *SqlStore) Follow(userId, authorId int8) (*domain.User, error) {
  var author domain.User
  err := s.db.Get(&author, "SELECT * FROM users WHERE id = $1 LIMIT 1", authorId)

  if err != nil {
    return nil, err
  }

  _, err = s.db.Exec(`
    INSERT INTO followers (user_id, follower_id)
    VALUES ($1, $2)
    WHERE id = $2
  `, userId, authorId)

  if err != nil {
    return nil, err
  }

  followers, err := s.getFollowers(authorId)

  if err != nil {
    return nil, err
  }
  author.Followers = followers

  return &author, nil
}

func (s *SqlStore) Unfollow(userId, authorId int8) (*domain.User, error) {
  var author domain.User
  err := s.db.Get(&author, "SELECT * FROM users WHERE id = $1 LIMIT 1", authorId)

  if err != nil {
    return nil, err
  }

  _, err = s.db.Exec(`
    DELETE FROM followers
    WHERE user_id = $1 AND follower_id = $2
  `, userId, authorId)

  if err != nil {
    return nil, err
  }

  followers, err := s.getFollowers(authorId)

  if err != nil {
    return nil, err
  }

  author.Followers = followers

  return &author, nil
}
