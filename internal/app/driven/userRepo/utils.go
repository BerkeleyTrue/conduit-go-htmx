package userRepo

import (
	"github.com/berkeleytrue/conduit/internal/core/domain"
	"github.com/berkeleytrue/conduit/internal/infra/data/krono"
	"golang.org/x/exp/slog"
)

func formatToDomain(user User, followers *[]int64) *domain.User {
	createdAt, err := krono.FromString(user.CreatedAt)

	if err != nil {
		slog.Debug("error parsing createdAt: %s", err)
	}

	updatedAt, err := krono.FromNullString(user.UpdatedAt)

	if err != nil {
		slog.Debug("error parsing updatedAt: %s", err)
	}

	dUser := &domain.User{
		UserId:    int(user.ID),
		Username:  user.Username,
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
