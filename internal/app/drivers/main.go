package drivers

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"golang.org/x/exp/slog"

	"github.com/berkeleytrue/conduit/internal/core/services"
)

type (
	Controller struct {
		store          *session.Store
		userService    *services.UserService
		articleService *services.ArticleService
		log            *slog.Logger
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
		log: slog.New(slog.NewTextHandler(os.Stdin, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})).WithGroup("drivers").WithGroup("controller"),
	}
}

func RegisterRoutes(
	app *fiber.App,
	c *Controller,
	authMiddleware fiber.Handler,
) {
	app.Use(func(ctx *fiber.Ctx) error {
		userId, ok := ctx.Locals("userId").(int)

		if !ok {
			userId = 0
		}
		links := UnAuthedLinks

		if userId != 0 {
			links = AuthedLinks
		}

		user, ok := ctx.Locals("user").(services.UserOutput)

		if !ok {
			user = services.UserOutput{}
		}

		ctx.Locals("layoutProps", layoutProps{
			title:  "Conduit",
			page:   ctx.Path(),
			uri:    ctx.OriginalURL(),
			userId: userId,
			user:   user,
			links:  links,
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
