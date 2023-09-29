package core

import (
	"github.com/berkeleytrue/conduit/internal/core/services"
	"go.uber.org/fx"
)

var Module = fx.Options(
	services.Module,
)
