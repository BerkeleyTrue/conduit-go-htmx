package userRepo

import (
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
