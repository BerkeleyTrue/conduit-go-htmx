package main

import (
	"go.uber.org/fx"

	"github.com/berkeleytrue/conduit/config"
	"github.com/berkeleytrue/conduit/internal/infra"
)

func main() {
	app := fx.New(
		config.Module,
		infra.Module,
	)

	app.Run()
}
