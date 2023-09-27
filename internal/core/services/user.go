package services

import (
	"errors"
	"time"

	"github.com/berkeleytrue/conduit/internal/core/domain"
	pss "github.com/berkeleytrue/conduit/internal/infra/data/password"
	"github.com/berkeleytrue/conduit/internal/utils"
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

	UserIdOrUsername struct {
		authorId   string
		authorname string
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

// is this user a follower of this author?
func isFollowing(author *domain.User, userId string) bool {
	return utils.Some(func(follower string) bool {
		return follower == userId
	}, author.Following)
}

func formatToPublicProfile(author *domain.User, following bool) *PublicProfile {
	return &PublicProfile{
		username:  author.Username,
		bio:       author.Bio,
		image:     author.Image,
		following: following,
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

func (s *UserService) GetUser(userId string) (*UserOutput, error) {
	user, err := s.repo.GetByID(userId)

	if err != nil {
		return nil, err
	}

	return formatUser(user), nil
}

func (s *UserService) GetIdFromUsername(username string) (string, error) {
	user, err := s.repo.GetByUsername(username)

	if err != nil {
		return "", err
	}

	return user.UserId, nil
}

func (s *UserService) GetProfile(authorIdOrAuthorname UserIdOrUsername, username string) (*PublicProfile, error) {

	user, err := s.repo.GetByUsername(username)

	if err != nil {
		return nil, err
	}
	err = nil

	var author *domain.User
	if authorIdOrAuthorname.authorId != "" {
		author, err = s.repo.GetByID(authorIdOrAuthorname.authorId)
	} else if authorIdOrAuthorname.authorname != "" {
		author, err = s.repo.GetByUsername(authorIdOrAuthorname.authorname)
	} else {
		return nil, errors.New("Invalid authorId or authorname")
	}

	if err != nil {
		return nil, err
	}

	return formatToPublicProfile(author, isFollowing(author, user.UserId)), nil
}

func (s *UserService) Update(userIdOrUsername UserIdOrUsername, input UpdateUserInput) (*UserOutput, error) {
	now := time.Now()
	var userId string
	var err error

	if userIdOrUsername.authorId != "" {
		userId = userIdOrUsername.authorId
	} else if userIdOrUsername.authorname != "" {
		userId, err = s.GetIdFromUsername(userIdOrUsername.authorname)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("Invalid authorId or authorname")
	}

	var updater domain.Updater[domain.User] = func(u *domain.User) *domain.User {
		u.Email = input.email
		u.Username = input.username
		u.Image = input.image
		u.Bio = input.bio
		u.UpdatedAt = now
		return u
	}

	user, err := s.repo.Update(userId, updater)

	return formatUser(user), nil
}

func (s *UserService) Follow(userId string, authorIdOrAuthorname UserIdOrUsername) (*PublicProfile, error) {
	var (
		authorId string
		err      error
	)

	if authorIdOrAuthorname.authorId != "" {
		authorId = authorIdOrAuthorname.authorId
	} else if authorIdOrAuthorname.authorname != "" {
		authorId, err = s.GetIdFromUsername(authorIdOrAuthorname.authorname)

		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("Invalid authorId or authorname")
	}

	user, err := s.repo.Follow(userId, authorId)

	if err != nil {
		return nil, err
	}

	return formatToPublicProfile(user, true), nil
}

func (s *UserService) Unfollow(userId, authorIdOrAuthorname string) (*PublicProfile, error) {
  var (
    authorId string
    err      error
  )

  if authorIdOrAuthorname != "" {
    authorId = authorIdOrAuthorname
  } else {
    return nil, errors.New("Invalid authorId or authorname")
  }

  user, err := s.repo.Unfollow(userId, authorId)

  if err != nil {
    return nil, err
  }

  return formatToPublicProfile(user, false), nil
}
