package services

import "github.com/berkeleytrue/conduit/internal/core/domain"

type UserService struct {
	repo *domain.UserRepository
}

func New(repo *domain.UserRepository) *UserService {
	return &UserService{repo: repo}
}
