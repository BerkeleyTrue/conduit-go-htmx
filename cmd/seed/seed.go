package main

import (
	"fmt"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/jmoiron/sqlx"
	"go.uber.org/fx"

	"github.com/berkeleytrue/conduit/config"
	"github.com/berkeleytrue/conduit/internal/app/driven/articlesRepo"
	"github.com/berkeleytrue/conduit/internal/app/driven/userRepo"
	"github.com/berkeleytrue/conduit/internal/core/domain"
	"github.com/berkeleytrue/conduit/internal/core/services"
	"github.com/berkeleytrue/conduit/internal/infra/data/password"
	"github.com/berkeleytrue/conduit/internal/infra/db"
)

type (
	Generate struct {
		userService    *services.UserService
		articleService *services.ArticleService
		shutdown       fx.Shutdowner
	}

	UserOutputPlusId struct {
		*services.UserOutput
		userId int
	}
)

func generateUser(
	userService *services.UserService,
) (*UserOutputPlusId, error) {
	pass := gofakeit.Password(true, true, true, true, false, 10)

	input := domain.UserCreateInput{
		Username:       gofakeit.Username(),
		Email:          gofakeit.Email(),
		HashedPassword: password.HashedPassword(pass),
	}

	userId, err := userService.Register(input)

	if err != nil {
		return nil, fmt.Errorf("error registering user: %w", err)
	}

	userOutput, err := userService.Update(
		services.UserIdOrUsername{UserId: userId},
		services.UpdateUserInput{
			Bio:   gofakeit.Sentence(10),
			Image: gofakeit.ImageURL(200, 200),
		},
	)

	return &UserOutputPlusId{
		UserOutput: userOutput,
		userId:     userId,
	}, nil
}

func seed(
	userService *services.UserService,
	articleService *services.ArticleService,
	shutdown fx.Shutdowner,
) {
	numOfUsers := 30
	numOfArticles := 20

	users := make([]*UserOutputPlusId, numOfUsers)

	for i := 0; i < numOfUsers; i++ {
		user, err := generateUser(userService)

		if err != nil {
			panic(err)
		}

		users[i] = user
	}

	for _, user := range users {
		for i := 0; i < numOfArticles; i++ {
			_, err := articleService.Create(
				user.userId,
				services.ArticleCreateInput{
					Title:       gofakeit.Sentence(5),
					Description: gofakeit.Sentence(10),
					Body:        gofakeit.Sentence(20),
					Tags:        []string{gofakeit.HipsterWord(), gofakeit.HipsterWord()},
				},
			)

			if err != nil {
				fmt.Println(err)
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
