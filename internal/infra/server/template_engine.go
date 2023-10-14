package server

import (
	"github.com/gofiber/template/html/v2"

	"github.com/berkeleytrue/conduit/config"
	"github.com/berkeleytrue/conduit/internal/utils"
)

func NewEngine(cfg *config.Config) *html.Engine {
	isDev := cfg.Release == "development"

	engine := html.New("./web/views", ".gohtml")
	engine.Reload(isDev)
	engine.AddFunc("Iterate", utils.Iterate)

	return engine
}
