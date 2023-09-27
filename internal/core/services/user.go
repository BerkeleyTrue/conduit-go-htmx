package services

import (
	"errors"

	"github.com/berkeleytrue/conduit/internal/core/domain"
	pss "github.com/berkeleytrue/conduit/internal/infra/data/password"
)

type (
	UserService struct {
		repo domain.UserRepository
	}

	UserOutput struct {
		email    string
		username string
		token    string
		bio      string
		image    string
	}

	// sent to any third party user
	PublicProfile struct {
		username  string
		bio       string
		image     string
		following bool
	}

	UpdateUserInput struct {
		email    string
		username string
		image    string
		bio      string
		password string // unhashed password
	}
)

func formatUser(user *domain.User) *UserOutput {
	return &UserOutput{
		email:    user.Email,
		username: user.Username,
		token:    "",
		bio:      user.Bio,
		image:    user.Image,
	}
}

func New(repo domain.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Register(input domain.UserCreateInput) (*UserOutput, error) {
	hashedPassword, err := pss.HashPassword(input.Password)

	if err != nil {
		return nil, err
	}

	user, err := s.repo.Create(domain.UserCreateInput{
		Username:       input.Username,
		Email:          input.Email,
		HashedPassword: hashedPassword,
	})

	if err != nil {
		return nil, err
	}

	return formatUser(user), nil
}

func (s *UserService) Login(email, rawPass string) (*UserOutput, error) {
  user, err := s.repo.GetByEmail(email)

  if err != nil {
    return nil, err
  }

  password, err := pss.New(rawPass)

  if err != nil {
    return nil, err
  }

  if err := pss.CompareHashAndPassword(user.Password, password); err != nil {
    return nil, errors.New("Invalid password")
  }

  return formatUser(user), nil
}
// func (s *UserService) GetUser(userId string) (*UserOutput, error) {
// }
// func (s *UserService) GetIdFromUsername(username string) (string, error) {
// }
// func (s *UserService) GetProfile(authorIdOrAuthorname string, username string) (*PublicProfile, error) {
// }
// func (s *UserService) Update(userIdOrUsername string, input UpdateUserInput) (*UserOutput, error) {
// }
// func (s *UserService) Follow(userId, authorIdOrAuthorname string) (*PublicProfile, error) {
// }
// func (s *UserService) Unfollow(userId, authorIdOrAuthorname string) (*PublicProfile, error) {
// }
