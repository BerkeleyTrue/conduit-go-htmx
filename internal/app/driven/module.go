package driven

import (
	"go.uber.org/fx"

	"github.com/berkeleytrue/conduit/internal/app/driven/userRepo"
)

var Module = fx.Options(
	userRepo.Module,
)
