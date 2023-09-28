package userRepo

import (
	"context"

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
  `
	Module = fx.Options(
		fx.Provide(NewSqlStore),
		fx.Invoke(RegisterUserSchema),
	)
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

func (s *SqlStore) Create(input domain.UserCreateInput) (*domain.User, error) {
	panic("implement me")
}

func (s *SqlStore) GetByID(id string) (*domain.User, error) {
	panic("implement me")
}

func (s *SqlStore) GetByEmail(email string) (*domain.User, error) {
	panic("implement me")
}

func (s *SqlStore) GetByUsername(username string) (*domain.User, error) {
	panic("implement me")
}

func (s *SqlStore) Update(userId string, updater domain.Updater[domain.User]) (*domain.User, error) {
	panic("implement me")
}

func (s *SqlStore) Follow(userId, authorId string) (*domain.User, error) {
	panic("implement me")
}

func (s *SqlStore) Unfollow(userId, authorId string) (*domain.User, error) {
	panic("implement me")
}
