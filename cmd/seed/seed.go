package main

import (
	"fmt"
	"os"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/jmoiron/sqlx"
	"go.uber.org/fx"
	"golang.org/x/exp/slog"

	"github.com/berkeleytrue/conduit/config"
	"github.com/berkeleytrue/conduit/internal/app/driven/articlesRepo"
	"github.com/berkeleytrue/conduit/internal/app/driven/userRepo"
	"github.com/berkeleytrue/conduit/internal/core/domain"
	"github.com/berkeleytrue/conduit/internal/core/services"
	"github.com/berkeleytrue/conduit/internal/infra/data/krono"
	"github.com/berkeleytrue/conduit/internal/infra/data/password"
	"github.com/berkeleytrue/conduit/internal/infra/db"
)

type (
	UserOutputPlusId struct {
		*services.UserOutput
		userId    int
		createdAt krono.Krono
	}
)

var log *slog.Logger = slog.New(slog.NewTextHandler(os.Stdout, nil)).WithGroup("seed")

func genImage() string {
  return fmt.Sprintf("https://picsum.photos/id/%d/200/200", gofakeit.Number(1, 1000))
}

func generateUser(
	userService *services.UserService,
	userRepo domain.UserRepository,
) (*UserOutputPlusId, error) {
	pass := gofakeit.Password(true, true, true, true, false, 10)
	createdAt := gofakeit.DateRange(
		time.Now().AddDate(-1, 0, 0),
		time.Now(),
	)

	input := services.RegisterParams{
		Username: gofakeit.Username(),
		Email:    gofakeit.Email(),
		Password: password.Password(pass),
	}

	userId, err := userService.Register(input)

	if err != nil {
		return nil, fmt.Errorf("error registering user: %w", err)
	}

	user, err := userRepo.Update(
		userId,
		func(user domain.User) domain.User {
			user.Bio = gofakeit.Sentence(10)
			user.Image = genImage()
			user.CreatedAt = krono.Krono{Time: createdAt}
			return user
		},
	)

	if err != nil {
		return nil, fmt.Errorf("error updating user: %w", err)
	}

	return &UserOutputPlusId{
		UserOutput: &services.UserOutput{
			Email:    user.Email,
			Username: user.Username,
			Bio:      user.Bio,
			Image:    user.Image,
		},
		userId:    userId,
		createdAt: user.CreatedAt,
	}, nil
}

func seed(
	userRepo domain.UserRepository,
	userService *services.UserService,
	articleRepo domain.ArticleRepository,
	shutdown fx.Shutdowner,
) {
	numOfUsers := 30
	numOfArticles := 20

	users := make([]*UserOutputPlusId, numOfUsers)

	for i := 0; i < numOfUsers; i++ {
		user, err := generateUser(userService, userRepo)

		if err != nil {
			panic(err)
		}

		users[i] = user
	}

	for _, user := range users {
		for i := 0; i < numOfArticles; i++ {
			createdAt := krono.Krono{Time: gofakeit.DateRange(
				user.createdAt.Time,
				time.Now(),
			)}

			article, err := articleRepo.Create(
				domain.ArticleCreateInput{
					Title:       gofakeit.Sentence(5),
					Description: gofakeit.Sentence(10),
					Body:        gofakeit.Sentence(20),
					Tags: []string{
						gofakeit.HipsterWord(),
						gofakeit.HipsterWord(),
					},
					AuthorId:  user.userId,
					CreatedAt: createdAt,
				},
			)

			fmt.Printf("created article: %v\n", article)

			if err != nil {
				log.Error("error", err)
			}
		}
	}

	shutdown.Shutdown()
}

func clearDb(db *sqlx.DB) {
	fmt.Println("clearing db")

	_, err := db.Exec(`
    DELETE FROM users;
    DELETE FROM articles;
    -- DELETE FROM comments;
    DELETE FROM tags;
    DELETE FROM article_tags;
  `)

	if err != nil {
		panic(err)
	}
}

func main() {

	app := fx.New(
		config.Module,

		db.Module,

		userRepo.Module,
		articlesRepo.Module,

		fx.Invoke(clearDb),

		fx.Provide(services.NewUserService),
		fx.Provide(services.NewArticleService),

		fx.Invoke(seed),
	)

	app.Run()
}
