package userRepo

import (
	"github.com/berkeleytrue/conduit/internal/core/domain"
	"github.com/berkeleytrue/conduit/internal/infra/data/krono"
)

func formatToDomain(user User, followers *[]int64) *domain.User {
	createdAt := krono.Krono{}
	createdAt.Scan(user.CreatedAt)
	updatedAt := krono.Krono{}
	updatedAt.Scan(user.UpdatedAt)

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
