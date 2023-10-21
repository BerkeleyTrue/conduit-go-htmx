package userRepo

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"

	"github.com/jmoiron/sqlx"
	"go.uber.org/fx"

	"github.com/berkeleytrue/conduit/internal/core/domain"
	"github.com/berkeleytrue/conduit/internal/infra/data/krono"
	"github.com/berkeleytrue/conduit/internal/infra/db"
	"github.com/berkeleytrue/conduit/internal/utils"
)

type (
	UserStore struct {
		db.SqlStore
	}
)

var (
	//go:embed sql/schema.sql
	schema string
	// compile time check to make sure SqlStore implements domain.UserRepository
	_ domain.UserRepository = (*Queries)(nil)

	Module = fx.Options(
		fx.Provide(fx.Annotate(
			newSqlStore,
			fx.As(new(domain.UserRepository))),
		),
		fx.Invoke(registerUserSchema),
	)
)

func newSqlStore(_db *sqlx.DB) *Queries {
	return &Queries{db: _db}
}

func registerUserSchema(db *sqlx.DB) error {
	_, err := db.Exec(schema)

	if err != nil {
		return err
	}

	return nil
}

func (q *Queries) Create(ctx context.Context, input domain.UserCreateInput) (*domain.User, error) {
	params := createParams{
		Username:  input.Username,
		Email:     input.Email,
		Password:  string(input.HashedPassword),
		Bio:       sql.NullString{},
		Image:     sql.NullString{},
		CreatedAt: input.CreatedAt.ToString(),
		UpdatedAt: sql.NullString{},
	}

	user, err := q.create(ctx, params)

	if err != nil {
		return nil, fmt.Errorf("sql-store: error creating user: %w", err)
	}

	return formatToDomain(user, nil), err
}

func (q *Queries) GetByID(ctx context.Context, id int) (*domain.User, error) {
	user, err := q.getById(ctx, int64(id))

	if err != nil {
		return nil, fmt.Errorf("sql-store: error getting user: %w", err)
	}

	return formatToDomain(user, nil), nil
}

func (q *Queries) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	user, err := q.getByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("sql-store: error getting user: %w", err)
	}

	return formatToDomain(user, nil), nil
}

func (q *Queries) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	user, err := q.getByUsername(ctx, username)

	if err != nil {
		return nil, fmt.Errorf("sql-store: error getting user: %w", err)
	}

	return formatToDomain(user, nil), nil
}

func (q *Queries) GetFollowing(ctx context.Context, userId int) ([]int, error) {
	following, err := q.getFollowing(ctx, int64(userId))

	if err != nil {
		return nil, fmt.Errorf("sql-store: error getting following: %w", err)
	}

	return utils.Map(func(userId int64) int {
		return int(userId)
	}, following), nil
}

func (q *Queries) Update(
	ctx context.Context,
	userId int,
	updater domain.Updater[domain.User],
) (*domain.User, error) {
	user, err := q.getById(ctx, int64(userId))

	if err != nil {
		return nil, fmt.Errorf("sql-store: error getting user: %w", err)
	}

	updates := updater(*(formatToDomain(user, nil)))
	params := updateParams{ID: user.ID}

	if updates.Username != "" {
		params.Username = updates.Username
	} else {
		params.Username = user.Username
	}

	if updates.Email != "" {
		params.Email = updates.Email
	} else {
		params.Email = user.Email
	}

	if updates.Password != "" {
		params.Password = string(updates.Password)
	} else {
		params.Password = user.Password
	}

	if updates.Bio != "" {
		params.Bio = sql.NullString{}
		params.Bio.Scan(updates.Bio)
	} else {
		params.Bio = user.Bio
	}

	if updates.Image != "" {
		params.Image = sql.NullString{}
		params.Image.Scan(updates.Image)
	} else {
		params.Image = user.Image
	}

	// for seeding, can not be done directly through the user service
	if updates.CreatedAt.IsZero() {
		params.CreatedAt = user.CreatedAt
	} else {
		params.CreatedAt = updates.CreatedAt.ToString()
	}

	params.UpdatedAt = sql.NullString{}
	params.UpdatedAt.Scan(krono.Now())

	user, err = q.update(ctx, params)

	if err != nil {
		return nil, fmt.Errorf("sql-store: error updating user: %w", err)
	}

	return formatToDomain(user, nil), nil
}

func (q *Queries) Follow(ctx context.Context, userId, authorId int) (*domain.User, error) {
	_, err := q.follow(
		ctx,
		followParams{
			AuthorID:   int64(authorId),
			FollowerID: int64(userId),
			CreatedAt:  krono.Now().ToString(),
		},
	)

	if err != nil {
		return nil, fmt.Errorf("sql-store: error following author: %w", err)
	}

	author, err := q.getById(ctx, int64(authorId))

	if err != nil {
		return nil, fmt.Errorf("sql-store: error getting author: %w", err)
	}

	followers, err := q.getFollowers(ctx, int64(authorId))

	if err != nil {
		return nil, fmt.Errorf("sql-store: error getting followers: %w", err)
	}

	return formatToDomain(author, &followers), nil
}

func (q *Queries) Unfollow(ctx context.Context, userId, authorId int) (*domain.User, error) {
	_, err := q.unfollow(ctx, unfollowParams{
		FollowerID: int64(userId),
		AuthorID:   int64(authorId),
	})

	if err != nil {
		return nil, fmt.Errorf("sql-store: error unfollowing author: %w", err)
	}

	author, err := q.getById(ctx, int64(authorId))

	if err != nil {
		return nil, fmt.Errorf("sql-store: error getting author: %w", err)
	}

	followers, err := q.getFollowers(ctx, int64(authorId))

	if err != nil {
		return nil, fmt.Errorf("sql-store: error getting followers: %w", err)
	}

	return formatToDomain(author, &followers), nil
}
