package drivers

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
	"golang.org/x/exp/slog"

	"github.com/berkeleytrue/conduit/config"
	"github.com/berkeleytrue/conduit/internal/core/services"
	"github.com/berkeleytrue/conduit/internal/infra/session"
)

type (
	Controller struct {
		userService    *services.UserService
		articleService *services.ArticleService
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
			NewController,
			fx.OnStart(func(c *Controller) error {
				c.onStart <- struct{}{}
				return nil
			}),
		)),
	)
)

func NewController(
	userService *services.UserService,
	articleService *services.ArticleService,
) *Controller {

	return &Controller{
		userService:    userService,
		articleService: articleService,
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

	app.Get("/profile/:username", c.GetProfile)
	app.Get("/article/:slug", c.getArticle)
	app.Get("/articles", c.GetArticles)
	app.Get("/tags", c.GetPopularTags)

	app.Use(authMiddleware)

	app.Get("/articles/feed", c.getFeed)
	app.Get("/editor", c.getEditArticle)
	app.Get("/editor/:slug", c.getEditArticle)
	app.Get("/settings", c.GetSettings)
	app.Post("/settings", c.UpdateSettings)
	app.Post("/logout", c.Logout)
}

func (c *Controller) Index(fc *fiber.Ctx) error {
	flashes, err := session.GetFlashes(fc)

	if err != nil {
		return err
	}

	_layoutProps := getLayoutProps(fc)

	_layoutProps.title = "Home"
	_layoutProps.flashes = flashes

	p := indexProps{
		layoutProps: _layoutProps,
	}

	return renderComponent(index(p), fc)
}
