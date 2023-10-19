package drivers

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"go.uber.org/fx"
	"golang.org/x/exp/slog"

	"github.com/berkeleytrue/conduit/config"
	"github.com/berkeleytrue/conduit/internal/core/services"
)

type (
	Controller struct {
		store          *session.Store
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
	store *session.Store,
	userService *services.UserService,
	articleService *services.ArticleService,
) *Controller {

	return &Controller{
		store:          store,
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
	app.Use(func(ctx *fiber.Ctx) error {
		userId, ok := ctx.Locals("userId").(int)

		if !ok {
			userId = 0
		}
		links := UnAuthedLinks

		if userId != 0 {
			links = AuthedLinks
		}

		user, ok := ctx.Locals("user").(*services.UserOutput)

		if !ok {
			user = &services.UserOutput{}
		}

		alerts := []alertPackage{}

		session, ok := ctx.Locals("session").(*session.Session)

		if ok {
			_alerts, ok := session.Get("alerts").([]alertPackage)

			if ok {
				alerts = _alerts
				session.Delete("alerts")
			}
		}

		ctx.Locals("layoutProps", layoutProps{
			title:  "Conduit",
			page:   ctx.Path(),
			uri:    ctx.OriginalURL(),
			userId: userId,
			user:   *user,
			links:  links,
			isDev:  config.Release == "development",
			alerts: alerts,
		})

		return ctx.Next()
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

func (c *Controller) Index(ctx *fiber.Ctx) error {
	_layoutProps := getLayoutProps(ctx)

	_layoutProps.title = "Home"

	p := indexProps{
		layoutProps: _layoutProps,
	}

	return renderComponent(index(p), ctx)
}
