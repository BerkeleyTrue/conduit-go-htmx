package main

import (
	"fmt"

	"github.com/brianvoe/gofakeit/v6"
	"go.uber.org/fx"

	"github.com/berkeleytrue/conduit/config"
	"github.com/berkeleytrue/conduit/internal/app/driven/articlesRepo"
	"github.com/berkeleytrue/conduit/internal/app/driven/userRepo"
	"github.com/berkeleytrue/conduit/internal/core/domain"
	"github.com/berkeleytrue/conduit/internal/core/services"
	"github.com/berkeleytrue/conduit/internal/infra/data/password"
	"github.com/berkeleytrue/conduit/internal/infra/db"
)

func generateUser(userService *services.UserService) (*services.UserOutput, error) {
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

	userOutput, err := userService.Update(services.UserIdOrUsername{UserId: userId}, services.UpdateUserInput{
		Bio:   gofakeit.Sentence(10),
		Image: gofakeit.ImageURL(200, 200),
	})

	return userOutput, nil
}

func generate(userService *services.UserService, articleService *services.ArticleService) {
	numOfUsers := 30
	// numOfArticles := 20

	users := make([]*services.UserOutput, numOfUsers)

	for i := 0; i < numOfUsers; i++ {
		user, err := generateUser(userService)

		if err != nil {
			fmt.Println(err)
		}

		users = append(users, user)
	}

	// for user := range users {
	//   for i := 0; i < numOfArticles; i++ {
	//     _, err := articleService.Create(services.UserIdOrUsername{UserId: user.UserId}, services.CreateArticleInput{
	//       Title:       gofakeit.Sentence(5),
	//       Description: gofakeit.Sentence(10),
	//       Body:        gofakeit.Sentence(20),
	//       Tags:        []string{gofakeit.HipsterWord()},
	//     })
	//
	//     if err != nil {
	//       fmt.Println(err)
	//     }
	//   }
	// }
}

func main() {

	app := fx.New(
		config.Module,

		fx.Provide(db.NewDB),
		fx.Provide(
			fx.Annotate(
				userRepo.NewSqlStore,
				fx.As(new(domain.UserRepository)),
			),
		),
		fx.Provide(
			fx.Annotate(
				articlesRepo.NewSqlStore,
				fx.As(new(domain.ArticleRepository)),
			),
		),
		fx.Provide(services.NewUserService),
		fx.Provide(services.NewArticleService),
		fx.Invoke(generate),
	)

	app.Run()
}
