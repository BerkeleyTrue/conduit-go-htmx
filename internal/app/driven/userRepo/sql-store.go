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

func (q *Queries) Create(input domain.UserCreateInput) (*domain.User, error) {
	ctx := context.Background()
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

func (q *Queries) GetByID(id int) (*domain.User, error) {
	user, err := q.getById(context.Background(), int64(id))

	if err != nil {
		return nil, fmt.Errorf("sql-store: error getting user: %w", err)
	}

	return formatToDomain(user, nil), nil
}

func (q *Queries) GetByEmail(email string) (*domain.User, error) {
	user, err := q.getByEmail(context.Background(), email)
	if err != nil {
		return nil, fmt.Errorf("sql-store: error getting user: %w", err)
	}

	return formatToDomain(user, nil), nil
}

func (q *Queries) GetByUsername(username string) (*domain.User, error) {
	user, err := q.getByUsername(context.Background(), username)

	if err != nil {
		return nil, fmt.Errorf("sql-store: error getting user: %w", err)
	}

	return formatToDomain(user, nil), nil
}

func (q *Queries) Update(
	userId int,
	updater domain.Updater[domain.User],
) (*domain.User, error) {
	ctx := context.Background()
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

func (q *Queries) Follow(userId, authorId int) (*domain.User, error) {
	ctx := context.Background()
	_, err := q.follow(ctx, followParams{UserID: int64(userId), FollowerID: int64(authorId)})

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

func (q *Queries) Unfollow(userId, authorId int) (*domain.User, error) {
	ctx := context.Background()
	_, err := q.unfollow(ctx, unfollowParams{UserID: int64(userId), FollowerID: int64(authorId)})

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
