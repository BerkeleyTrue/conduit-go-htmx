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
		userId   int8
		username string
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
func isFollowing(author *domain.User, userId int8) bool {
	return utils.Some(func(follower int8) bool {
		return follower == userId
	}, author.Followers)
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

func (s *UserService) GetUser(userId int8) (*UserOutput, error) {
	user, err := s.repo.GetByID(userId)

	if err != nil {
		return nil, err
	}

	return formatUser(user), nil
}

func (s *UserService) GetIdFromUsername(username string) (int8, error) {
	user, err := s.repo.GetByUsername(username)

	if err != nil {
		return 0, err
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
	if authorIdOrAuthorname.userId != 0 {
		author, err = s.repo.GetByID(authorIdOrAuthorname.userId)
	} else if authorIdOrAuthorname.username != "" {
		author, err = s.repo.GetByUsername(authorIdOrAuthorname.username)
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
	var userId int8
	var err error

	if userIdOrUsername.userId != 0 {
		userId = userIdOrUsername.userId
	} else if userIdOrUsername.username != "" {

		userId, err = s.GetIdFromUsername(userIdOrUsername.username)

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

func (s *UserService) Follow(userId int8, authorIdOrAuthorname UserIdOrUsername) (*PublicProfile, error) {
	var (
		authorId int8
		err      error
	)

	if authorIdOrAuthorname.userId != 0 {
		authorId = authorIdOrAuthorname.userId
	} else if authorIdOrAuthorname.username != "" {
		authorId, err = s.GetIdFromUsername(authorIdOrAuthorname.username)

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

func (s *UserService) Unfollow(userId int8, authorIdOrAuthorname UserIdOrUsername) (*PublicProfile, error) {
	var (
		authorId int8
		err      error
	)

	if authorIdOrAuthorname.userId != 0 {
		authorId = authorIdOrAuthorname.userId
	} else {
		return nil, errors.New("Invalid authorId or authorname")
	}

	user, err := s.repo.Unfollow(userId, authorId)

	if err != nil {
		return nil, err
	}

	return formatToPublicProfile(user, false), nil
}
