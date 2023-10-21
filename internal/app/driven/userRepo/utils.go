package userRepo

import (
	"database/sql"
	"strconv"
	"strings"

	"github.com/berkeleytrue/conduit/internal/core/domain"
	"github.com/berkeleytrue/conduit/internal/infra/data/krono"
	"github.com/berkeleytrue/conduit/internal/infra/data/password"
	"golang.org/x/exp/slog"
)

func formatToDomain(user User, followers *[]int64) *domain.User {
	createdAt, err := krono.FromString(user.CreatedAt)

	if err != nil {
		slog.Debug("error parsing createdAt", "error", err)
	}

	updatedAt, err := krono.FromNullString(user.UpdatedAt)

	if err != nil {
		slog.Debug("error parsing updatedAt", "error", err)
	}

	dUser := &domain.User{
		UserId:    int(user.ID),
		Username:  user.Username,
		Password:  password.HashedPassword(user.Password),
		Email:     user.Email,
		Bio:       user.Bio.String,
		Image:     user.Image.String,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	if followers != nil {
		dUser.Followers = make([]int, len(*followers))

		for idx, id := range *followers {
			dUser.Followers[idx] = int(id)
		}

	} else {
		dUser.Followers = []int{}
	}

	return dUser
}

func formatStringToFollowers(followers string) *[]int64 {
	followersString := strings.Split(followers, ",")
	followersInt := make([]int64, len(followersString))

	for idx, id := range followersString {

		if id != "" {
			id, err := strconv.ParseInt(id, 10, 64)

			if err != nil {
				slog.Debug("error parsing follower id", "error", err)
			} else {
				followersInt[idx] = id
			}
		}
	}

	return &followersInt
}

func fromArgsToUser(
	ID int64,
	Username string,
	Email string,
	Password string,
	Bio sql.NullString,
	Image sql.NullString,
	CreatedAt string,
	UpdatedAt sql.NullString,
) *User {
	return &User{
		ID:        ID,
		Username:  Username,
		Email:     Email,
		Password:  Password,
		Bio:       Bio,
		Image:     Image,
		CreatedAt: CreatedAt,
		UpdatedAt: UpdatedAt,
	}
}

func formatFromRowToDomain(byId *getByIdRow, byEmail *getByEmailRow, byUsername *getByUsernameRow) *domain.User {
	user := new(User)

	if byId != nil {
		user = fromArgsToUser(
			byId.ID,
			byId.Username,
			byId.Email,
			byId.Password,
			byId.Bio,
			byId.Image,
			byId.CreatedAt,
			byId.UpdatedAt,
		)
		return formatToDomain(*user, formatStringToFollowers(byId.Followers))
	}

	if byEmail != nil {
		user = fromArgsToUser(
			byEmail.ID,
			byEmail.Username,
			byEmail.Email,
			byEmail.Password,
			byEmail.Bio,
			byEmail.Image,
			byEmail.CreatedAt,
			byEmail.UpdatedAt,
		)

		return formatToDomain(*user, formatStringToFollowers(byEmail.Followers))
	}

	if byUsername != nil {
		user = fromArgsToUser(
			byUsername.ID,
			byUsername.Username,
			byUsername.Email,
			byUsername.Password,
			byUsername.Bio,
			byUsername.Image,
			byUsername.CreatedAt,
			byUsername.UpdatedAt,
		)

		return formatToDomain(*user, formatStringToFollowers(byUsername.Followers))
	}

	return nil
}
