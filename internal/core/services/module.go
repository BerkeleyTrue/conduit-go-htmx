package services

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(newUserService),
	fx.Provide(newArticleService),
	fx.Provide(newCommentService),
)
