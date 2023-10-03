package userRepo

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"go.uber.org/fx"

	"github.com/berkeleytrue/conduit/internal/core/domain"
	"github.com/berkeleytrue/conduit/internal/infra/data/krono"
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
        username TEXT NOT NULL UNIQUE,
        email TEXT NOT NULL UNIQUE,
        password TEXT NOT NULL,
        bio TEXT,
        image TEXT,
        created_at TEXT NOT NULL,
        updated_at TEXT
    );

    CREATE TABLE IF NOT EXISTS followers (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        user_id INTEGER NOT NULL,
        follower_id INTEGER NOT NULL,
        created_at TEXT NOT NULL,
        UNIQUE(user_id, follower_id)
    );
  `
	// compile time check to make sure SqlStore implements domain.UserRepository
	_ domain.UserRepository = (*SqlStore)(nil)

	Module = fx.Options(
		fx.Provide(fx.Annotate(
			newSqlStore,
			fx.As(new(domain.UserRepository))),
		),
		fx.Invoke(registerUserSchema),
	)
)

func newSqlStore(db *sqlx.DB) *SqlStore {
	return &SqlStore{
		db: db,
	}
}

func registerUserSchema(db *sqlx.DB) error {
	_, err := db.Exec(schema)

	if err != nil {
		return err
	}

	return nil
}

// get followers for a user
func (s *SqlStore) getFollowers(userId int) ([]int, error) {
	var followers []int
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
		CreatedAt: krono.Krono{Time: now},
	}

	query := `
    INSERT INTO users (username, email, password, bio, image, created_at, updated_at)
    VALUES (:username, :email, :password, :bio, :image, :created_at, :updated_at)
  `
	_, err := s.db.NamedExec(query, user)

	if err != nil {
		return nil, fmt.Errorf("sql-store: error creating user: %w", err)
	}

	err = s.db.Get(&user, "SELECT * FROM users WHERE email = $1 LIMIT 1", input.Email)

	if err != nil {
		return nil, fmt.Errorf("sql-store: error getting new user: %w", err)
	}

	return &user, err
}

func (s *SqlStore) GetByID(id int) (*domain.User, error) {
	var user domain.User
	err := s.db.Get(&user, "SELECT * FROM users WHERE id = $1 LIMIT 1", id)

	if err != nil {
		return nil, fmt.Errorf("sql-store: error getting user: %w", err)
	}

	return &user, nil
}

func (s *SqlStore) GetByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := s.db.Get(&user, "SELECT * FROM users WHERE email = $1 LIMIT 1", email)
	if err != nil {
		return nil, fmt.Errorf("sql-store: error getting user: %w", err)
	}

	return &user, nil
}

func (s *SqlStore) GetByUsername(username string) (*domain.User, error) {
	var user domain.User
	err := s.db.Get(&user, "SELECT * FROM users WHERE username = $1 LIMIT 1", username)

	if err != nil {
		return nil, fmt.Errorf("sql-store: error getting user: %w", err)
	}

	return &user, nil
}

func (s *SqlStore) Update(
	userId int,
	updater domain.Updater[domain.User],
) (*domain.User, error) {
	var user domain.User
	err := s.db.Get(&user, "SELECT * FROM users WHERE id = $1 LIMIT 1", userId)
	if err != nil {
		return nil, fmt.Errorf("sql-store: error getting user: %w", err)
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
		return nil, fmt.Errorf("sql-store: error updating user: %w", err)
	}

	return updatedUser, nil
}

func (s *SqlStore) Follow(userId, authorId int) (*domain.User, error) {
	var author domain.User
	err := s.db.Get(&author, "SELECT * FROM users WHERE id = $1 LIMIT 1", authorId)

	if err != nil {
		return nil, fmt.Errorf("sql-store: error getting user: %w", err)
	}

	_, err = s.db.Exec(`
    INSERT INTO followers (user_id, follower_id)
    VALUES ($1, $2)
    WHERE id = $2
  `, userId, authorId)

	if err != nil {
		return nil, fmt.Errorf("sql-store: error following user: %w", err)
	}

	followers, err := s.getFollowers(authorId)

	if err != nil {
		return nil, fmt.Errorf("sql-store: error getting followers: %w", err)
	}
	author.Followers = followers

	return &author, nil
}

func (s *SqlStore) Unfollow(userId, authorId int) (*domain.User, error) {
	var author domain.User
	err := s.db.Get(&author, "SELECT * FROM users WHERE id = $1 LIMIT 1", authorId)

	if err != nil {
		return nil, fmt.Errorf("sql-store: error getting user: %w", err)
	}

	_, err = s.db.Exec(`
    DELETE FROM followers
    WHERE user_id = $1 AND follower_id = $2
  `, userId, authorId)

	if err != nil {
		return nil, fmt.Errorf("sql-store: error unfollowing user: %w", err)
	}

	followers, err := s.getFollowers(authorId)

	if err != nil {
		return nil, fmt.Errorf("sql-store: error getting followers: %w", err)
	}

	author.Followers = followers

	return &author, nil
}
