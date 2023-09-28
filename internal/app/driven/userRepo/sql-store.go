package userRepo

import (
	"github.com/jmoiron/sqlx"

	"github.com/berkeleytrue/conduit/internal/core/domain"
)

type (
	SqlStore struct {
		db *sqlx.DB
	}
)

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
