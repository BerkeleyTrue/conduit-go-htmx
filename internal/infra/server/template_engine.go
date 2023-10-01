package server

import (
	"github.com/gofiber/template/html/v2"

	"github.com/berkeleytrue/conduit/config"
)

func NewEngine(cfg *config.Config) *html.Engine {
	isDev := cfg.Release == "development"

	engine := html.New("./web/views", ".gohtml")
	engine.Reload(isDev)
	engine.AddFunc("Iterate", func(count int) []int {
		cnts := make([]int, count)
		for i := 0; i < count; i++ {
			cnts[i] = i
		}
		return cnts
	})

	return engine
}
