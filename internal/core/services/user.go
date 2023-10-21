package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	"golang.org/x/exp/slog"

	"github.com/berkeleytrue/conduit/internal/core/domain"
	"github.com/berkeleytrue/conduit/internal/infra/data/krono"
	"github.com/berkeleytrue/conduit/internal/infra/data/password"
	pss "github.com/berkeleytrue/conduit/internal/infra/data/password"
	"github.com/berkeleytrue/conduit/internal/utils"
)

type (
	UserService struct {
		repo domain.UserRepository
		log  *slog.Logger
	}

	UserOutput struct {
		Email    string
		Username string
		Bio      string
		Image    string
	}

	// sent to any third party user
	PublicProfile struct {
		Username  string
		Bio       string
		Image     string
		Following bool
	}
)

var (
	ErrNoUser = errors.New(
		"No user found with that email and password",
	)
	ErrInvalidIdOrUsername = errors.New("Invalid userId or username")
)

func formatUser(user *domain.User) *UserOutput {
	return &UserOutput{
		Email:    user.Email,
		Username: user.Username,
		Bio:      user.Bio,
		Image:    user.Image,
	}
}

// is this user a follower of this author?
func isFollowing(author *domain.User, userId int) bool {
	return utils.Some(func(follower int) bool {
		return follower == userId
	}, author.Followers)
}

func formatToPublicProfile(author *domain.User, following bool) *PublicProfile {
	return &PublicProfile{
		Username:  author.Username,
		Bio:       author.Bio,
		Image:     author.Image,
		Following: following,
	}
}

func newUserService(repo domain.UserRepository) *UserService {
	logger := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	return &UserService{
		repo: repo,
		log:  slog.New(logger).WithGroup("services").WithGroup("users"),
	}
}

type RegisterParams struct {
	Username string
	Email    string
	Password password.Password
}

func (s *UserService) Register(ctx context.Context, input RegisterParams) (int, error) {
	hashedPassword, err := pss.HashPassword(input.Password)

	if err != nil {
		return 0, err
	}

	user, err := s.repo.Create(ctx, domain.UserCreateInput{
		Username:       input.Username,
		Email:          input.Email,
		HashedPassword: hashedPassword,
		CreatedAt:      krono.Krono{Time: time.Now()},
	})

	if err != nil {
		return 0, fmt.Errorf("error creating user: %w", err)
	}

	s.log.Debug("created", "user", user)

	return user.UserId, nil
}

func (s *UserService) Login(ctx context.Context, email, rawPass string) (int, error) {
	user, err := s.repo.GetByEmail(ctx, email)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrNoUser
		}

		return 0, fmt.Errorf("error getting user: %w", err)
	}

	password, err := pss.New(rawPass)

	if err != nil {
		fmt.Printf("error creating password: %v\n", err)
		return 0, err
	}

	if err := pss.CompareHashAndPassword(user.Password, password); err != nil {
		return 0, ErrNoUser
	}

	return user.UserId, nil
}

func (s *UserService) GetUser(ctx context.Context, userId int) (*UserOutput, error) {
	user, err := s.repo.GetByID(ctx, userId)

	if err != nil {
		return nil, fmt.Errorf("error getting user: %w", err)
	}

	return formatUser(user), nil
}

func (s *UserService) GetIdFromUsername(ctx context.Context, username string) (int, error) {
	user, err := s.repo.GetByUsername(ctx, username)

	if err != nil {
		return 0, fmt.Errorf("error getting user: %w", err)
	}

	return user.UserId, nil
}

func (s *UserService) GetProfile(
	ctx context.Context,
	authorId int,
	authorname string,
	userId int, // the user who is requesting the profile, if any
) (*PublicProfile, error) {
	var (
		author       *domain.User
		err          error
		_isFollowing bool = false
	)

	if authorId != 0 {
		author, err = s.repo.GetByID(ctx, authorId)
	} else if authorname != "" {
		author, err = s.repo.GetByUsername(ctx, authorname)
	} else {
		return nil, ErrInvalidIdOrUsername
	}

	if err != nil {
		return nil, fmt.Errorf("error getting author: %w", err)
	}

	if userId != 0 {
		_isFollowing = isFollowing(author, userId)
	}

	return formatToPublicProfile(author, _isFollowing), nil
}

// get all the authors that this user is following
func (s *UserService) GetFollowing(
	ctx context.Context,
	userId int,
) ([]int, error) {
	following, err := s.repo.GetFollowing(ctx, userId)

	if err != nil {
		return nil, err
	}

	return following, nil
}

type UpdateUserInput struct {
	Email    string
	Username string
	Image    string
	Bio      string
	Password pss.Password
}

func (s *UserService) Update(
	ctx context.Context,
	userId int,
	username string,
	input UpdateUserInput,
) (*UserOutput, error) {
	now := time.Now()
	var err error

	if userId == 0 {
		if username != "" {

			userId, err = s.GetIdFromUsername(ctx, username)

			if err != nil {
				return nil, err
			}

		} else {
			return nil, ErrInvalidIdOrUsername
		}
	}

	var updater domain.Updater[domain.User] = func(u domain.User) domain.User {
		if input.Email != "" {
			u.Email = input.Email
		}

		if input.Image != "" {
			u.Image = input.Image
		}

		if input.Bio != "" {
			u.Bio = input.Bio
		}

		u.UpdatedAt = krono.Krono{Time: now}
		return u
	}

	user, err := s.repo.Update(ctx, userId, updater)

	if err != nil {
		return nil, err
	}

	return formatUser(user), nil
}

func (s *UserService) Follow(
	ctx context.Context,
	userId int,
	authorId int,
	authorname string,
) (*PublicProfile, error) {
	var (
		err error
	)

	if authorId == 0 {
		if authorname != "" {
			authorId, err = s.GetIdFromUsername(ctx, authorname)

			if err != nil {
				return nil, err
			}
		} else {
			return nil, ErrInvalidIdOrUsername
		}
	}

	user, err := s.repo.Follow(ctx, userId, authorId)

	if err != nil {
		return nil, err
	}

	return formatToPublicProfile(user, true), nil
}

func (s *UserService) Unfollow(
	ctx context.Context,
	userId int,
	authorId int,
	authorname string,
) (*PublicProfile, error) {
	var (
		err error
	)

	if authorId == 0 {
		if authorname != "" {
			authorId, err = s.GetIdFromUsername(ctx, authorname)
		} else {
			return nil, ErrInvalidIdOrUsername
		}
	}

	user, err := s.repo.Unfollow(ctx, userId, authorId)

	if err != nil {
		return nil, err
	}

	return formatToPublicProfile(user, false), nil
}
