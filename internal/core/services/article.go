package services

import (
	"context"
	"errors"
	"os"

	"golang.org/x/exp/slog"

	"github.com/berkeleytrue/conduit/internal/core/domain"
	"github.com/berkeleytrue/conduit/internal/infra/data/krono"
	"github.com/berkeleytrue/conduit/internal/infra/data/slug"
)

type (
	ArticleService struct {
		repo        domain.ArticleRepository
		userService *UserService
		log         *slog.Logger
	}

	ArticleOutput struct {
		Slug           slug.Slug
		Title          string
		Description    string
		Body           string
		Tags           []string
		CreatedAt      krono.Krono
		UpdatedAt      krono.Krono
		IsFavorited    bool
		FavoritesCount int
		Author         PublicProfile
	}
)

func formatArticle(
	article *domain.Article,
	profile PublicProfile,
	numOfFavorites int,
	isFavorite bool,
) ArticleOutput {
	return ArticleOutput{
		Slug:           article.Slug,
		Title:          article.Title,
		Description:    article.Description,
		Body:           article.Body,
		Tags:           article.Tags,
		CreatedAt:      article.CreatedAt,
		UpdatedAt:      article.UpdatedAt,
		IsFavorited:    isFavorite,
		FavoritesCount: numOfFavorites,
		Author:         profile,
	}
}

func newArticleService(
	repo domain.ArticleRepository,
	userService *UserService,
) *ArticleService {
	return &ArticleService{
		repo: repo,
		log: slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})).WithGroup("services").WithGroup("article"),
		userService: userService,
	}
}

type ArticleCreateInput struct {
	Title       string
	Description string
	Body        string
	Tags        []string
}

func (s *ArticleService) Create(
	ctx context.Context,
	userId int,
	input ArticleCreateInput,
) (ArticleOutput, error) {
	article, err := s.repo.Create(ctx, domain.ArticleCreateInput{
		Title:       input.Title,
		Description: input.Description,
		Body:        input.Body,
		Tags:        input.Tags,
		AuthorId:    userId,
		CreatedAt:   krono.Now(),
	})

	s.log.Debug("created", "article", article)

	if err != nil {
		return ArticleOutput{}, err
	}

	profile, err := s.userService.GetProfile(ctx, article.AuthorId, "", userId)

	if err != nil {
		return ArticleOutput{}, err
	}

	return formatArticle(article, *profile, 0, false), nil
}

type ListArticlesInput struct {
	Limit      int
	Offset     int
	Tag        string
	Favorited  string
	Authorname string
	Feed       bool
}

var WarningNoFollowers = errors.New("user is not following anyone")

// Get all articles, filtered by authorId, tag, favorited, limit, offset, or by following
func (s *ArticleService) List(
	ctx context.Context,
	userId int,
	input ListArticlesInput,
) ([]ArticleOutput, error) {
	params := domain.ArticleListInput{
		Limit:  input.Limit,
		Offset: input.Offset,
	}

	if input.Feed {
		params.FollowedBy = userId
	} else {
		params.Tag = input.Tag
		if input.Authorname != "" {
			authorId, err := s.userService.GetIdFromUsername(ctx, input.Authorname)

			if err != nil {
				return nil, err

			}

			params.AuthorId = authorId
		}

		if input.Favorited != "" {
			favoritedId, err := s.userService.GetIdFromUsername(ctx, input.Favorited)
			if err != nil {
				return nil, err
			}
			params.Favorited = favoritedId
		}
	}

	articles, err := s.repo.List(ctx, params)

	if err != nil {
		return nil, err
	}

	outputs := make([]ArticleOutput, len(articles))

	for idx, article := range articles {
		profile, err := s.userService.GetProfile(
			ctx,
			article.AuthorId,
			"",
			userId,
		)

		if err != nil {
			s.log.Debug("error getting profile", "error", err)
			continue
		}

		numOfFavorites, err := s.repo.GetNumOfFavorites(ctx, article.ArticleId)

		if err != nil {
			s.log.Debug("error getting num of favorites", "error", err)
			continue
		}

		isFavorited, err := s.repo.IsFavoritedByUser(
			ctx,
			article.ArticleId,
			userId,
		)

		if err != nil {
			s.log.Debug("error getting is favorited", "error", err)
			continue
		}

		outputs[idx] = formatArticle(article, *profile, numOfFavorites, isFavorited)
	}

	return outputs, err
}

func (s *ArticleService) GetPopularTags(ctx context.Context) ([]string, error) {
	return s.repo.GetPopularTags(ctx)
}

func (s *ArticleService) GetBySlug(
	ctx context.Context,
	slug string,
	userId int,
) (ArticleOutput, error) {
	article, err := s.repo.GetBySlug(ctx, slug)

	if err != nil {
		return ArticleOutput{}, err
	}

	profile, err := s.userService.GetProfile(ctx, article.AuthorId, "", userId)

	if err != nil {
		return ArticleOutput{}, err
	}

	numOfFavorites, err := s.repo.GetNumOfFavorites(ctx, article.ArticleId)

	if err != nil {
		return ArticleOutput{}, err
	}

	isFavorited, err := s.repo.IsFavoritedByUser(
		ctx,
		article.ArticleId,
		userId,
	)

	if err != nil {
		s.log.Debug("error getting is favorited", "error", err)
	}

	return formatArticle(article, *profile, numOfFavorites, isFavorited), nil
}

func (s *ArticleService) GetIdFromSlug(
	ctx context.Context,
	slug string,
) (int, error) {
	article, err := s.repo.GetBySlug(ctx, slug)

	if err != nil {
		return 0, err
	}

	return article.ArticleId, nil
}

type ArticleUpdateInput struct {
	Title       string
	Description string
	Body        string
	Tags        []string
}

func (s *ArticleService) Update(
	ctx context.Context,
	slug string,
	username string,
	input ArticleUpdateInput,
) (ArticleOutput, error) {
	article, err := s.repo.Update(
		ctx,
		slug,
		func(a domain.Article) domain.Article {
			if input.Title != "" {
				a.Title = input.Title
			}

			if input.Description != "" {
				a.Description = input.Description
			}

			if input.Body != "" {
				a.Body = input.Body
			}

			if input.Tags != nil {
				a.Tags = input.Tags
			}
			return a
		},
	)

	if err != nil {
		return ArticleOutput{}, err
	}

	profile, err := s.userService.GetProfile(ctx, article.AuthorId, "", 0)

	if err != nil {
		return ArticleOutput{}, err
	}

	numOfFavorites, err := s.repo.GetNumOfFavorites(ctx, article.ArticleId)

	if err != nil {
		return ArticleOutput{}, err
	}

	return formatArticle(article, *profile, numOfFavorites, true), nil
}

func (s *ArticleService) Favorite(
	ctx context.Context,
	slug string,
	userId int,
) (ArticleOutput, error) {
	article, err := s.repo.Favorite(ctx, slug, userId)

	if err != nil {
		return ArticleOutput{}, err
	}

	profile, err := s.userService.GetProfile(ctx, article.AuthorId, "", userId)

	if err != nil {
		return ArticleOutput{}, err
	}

	numOfFavorites, err := s.repo.GetNumOfFavorites(ctx, article.ArticleId)

	if err != nil {
		return ArticleOutput{}, err
	}

	return formatArticle(article, *profile, numOfFavorites, true), nil
}

func (s *ArticleService) Unfavorite(
	ctx context.Context,
	slug string,
	userId int,
) (ArticleOutput, error) {
	article, err := s.repo.Unfavorite(ctx, slug, userId)

	if err != nil {
		return ArticleOutput{}, err
	}

	profile, err := s.userService.GetProfile(ctx, article.AuthorId, "", userId)

	if err != nil {
		return ArticleOutput{}, err
	}

	numOfFavorites, err := s.repo.GetNumOfFavorites(ctx, article.ArticleId)

	if err != nil {
		return ArticleOutput{}, err
	}

	return formatArticle(article, *profile, numOfFavorites, false), nil
}

func (s *ArticleService) Delete(
	ctx context.Context,
	slug string,
) error {
	return s.repo.Delete(ctx, slug)
}
