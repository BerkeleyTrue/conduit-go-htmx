package drivers

import (
	"errors"
	"os"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
	"golang.org/x/exp/slog"

	"github.com/berkeleytrue/conduit/config"
	"github.com/berkeleytrue/conduit/internal/core/services"
)

type (
	Controller struct {
		userService    *services.UserService
		articleService *services.ArticleService
		commentService *services.CommentService
		log            *slog.Logger
		onStart        chan struct{}
	}
	Link struct {
		Uri   string
		Title string
	}
)

var (
	UnAuthedLinks = []Link{
		{
			Uri:   "/",
			Title: "Home",
		},
		{
			Uri:   "/login",
			Title: "Login",
		},
		{
			Uri:   "/register",
			Title: "Sign up",
		},
	}
	AuthedLinks = []Link{
		{
			Uri:   "/",
			Title: "Home",
		},
		{
			Uri:   "/editor",
			Title: "New Article",
		},
		{
			Uri:   "/settings",
			Title: "Settings",
		},
	}
	Module = fx.Options(
		fx.Provide(fx.Annotate(
			newController,
			fx.OnStart(func(c *Controller) error {
				c.onStart <- struct{}{}
				return nil
			}),
		)),
	)
)

func notImplementedHandler(fc *fiber.Ctx) error {
	return errors.New("not implemented")
}

func newController(
	userService *services.UserService,
	articleService *services.ArticleService,
	commentService *services.CommentService,
) *Controller {

	return &Controller{
		userService:    userService,
		articleService: articleService,
		commentService: commentService,
		onStart:        make(chan struct{}, 1),

		log: slog.New(slog.NewTextHandler(os.Stdin, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})).WithGroup("drivers").WithGroup("controller"),
	}
}

func RegisterRoutes(
	app *fiber.App,
	c *Controller,
	authMiddleware fiber.Handler,
	config *config.Config,
) {
	if config.Release == "development" {
		c.log.Debug("Registering hot reload route")
		app.Get("/__hotreload", c.getSSE)
	}
	app.Use(func(fc *fiber.Ctx) error {
		userId, ok := fc.Locals("userId").(int)

		if !ok {
			userId = 0
		}
		links := UnAuthedLinks

		if userId != 0 {
			links = AuthedLinks
		}

		user, ok := fc.Locals("user").(*services.UserOutput)

		if !ok {
			user = &services.UserOutput{}
		}

		fc.Locals("layoutProps", layoutProps{
			title:  "Conduit",
			page:   fc.Path(),
			uri:    fc.OriginalURL(),
			userId: userId,
			user:   *user,
			links:  links,
			isDev:  config.Release == "development",
		})

		return fc.Next()
	})

	app.Get("/", c.Index)
	app.Get("/login", c.GetLogin)
	app.Post("/login", c.Login)
	app.Get("/register", c.GetRegister)
	app.Post("/register", c.Register)

	app.Get("/profiles/:username", c.GetProfile)
	app.Get("/tags", c.GetPopularTags)

	// auth required
	app.Get("/articles/feed", authMiddleware, c.getFeed)
	app.Get("/articles", c.getArticles)

	app.Get("/articles/:slug", c.getArticle)
	app.Get("/articles/:slug/comments", c.getComments)

	app.Use(authMiddleware)

	app.Post("/articles/:slug/comments", c.createComment)
	app.Delete("/articles/:slug/comments/:id", c.deleteComment)

	app.Post("/articles/:slug/favorite", notImplementedHandler)
	app.Delete("/articles/:slug/favorite", notImplementedHandler)

	app.Post("/profiles/:username/follow", c.follow)
	app.Delete("/profiles/:username/follow", c.unfollow)

	app.Get("/editor", c.getEditArticle)
	app.Get("/editor/:slug", c.getEditArticle)

	app.Get("/settings", c.GetSettings)
	app.Post("/settings", c.UpdateSettings)
	app.Post("/logout", c.Logout)
}
