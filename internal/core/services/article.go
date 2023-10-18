package services

import (
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

	ArticleCreateInput struct {
		Title       string
		Description string
		Body        string
		Tags        []string
	}

	ArticleOutput struct {
		Slug           slug.Slug
		Title          string
		Description    string
		Body           string
		Tags           []string
		CreatedAt      krono.Krono
		UpdatedAt      krono.Krono
		Favorited      bool
		FavoritesCount int
		Author         PublicProfile
	}

	ArticleUpdateInput struct {
		Title       string
		Description string
		Body        string
	}
)

func formatArticle(
	article *domain.Article,
	profile PublicProfile,
) ArticleOutput {
	return ArticleOutput{
		Slug:        article.Slug,
		Title:       article.Title,
		Description: article.Description,
		Body:        article.Body,
		Tags:        article.Tags,
		CreatedAt:   article.CreatedAt,
		UpdatedAt:   article.UpdatedAt,
		// TODO: Favorited
		Favorited: false,
		// TODO: FavoritesCount
		FavoritesCount: 0,
		Author:         profile,
	}
}

func NewArticleService(
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

func (s *ArticleService) Create(
	userId int,
	input ArticleCreateInput,
) (ArticleOutput, error) {
	article, err := s.repo.Create(domain.ArticleCreateInput{
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

	profile, err := s.userService.GetProfile(article.AuthorId, "", userId)

	if err != nil {
		return ArticleOutput{}, err
	}

	return formatArticle(article, *profile), nil
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
	userId int,
	input ListArticlesInput,
) ([]ArticleOutput, error) {
	params := domain.ArticleListInput{
		Limit:  input.Limit,
		Offset: input.Offset,
	}

	if input.Feed {
		followed, err := s.userService.GetFollowing(userId)

		if err != nil {
			return nil, err
		}

		if len(followed) == 0 {
			return nil, WarningNoFollowers
		}

		params.Authors = followed
	} else {
		params.Tag = input.Tag
		if input.Authorname != "" {
			authorId, err := s.userService.GetIdFromUsername(input.Authorname)

			if err != nil {
				return nil, err

			}

			params.AuthorId = authorId
		}

		if input.Favorited != "" {
			favoritedId, err := s.userService.GetIdFromUsername(input.Favorited)
			if err != nil {
				return nil, err
			}
			params.Favorited = favoritedId
		}
	}

	articles, err := s.repo.List(params)

	if err != nil {
		return nil, err
	}

	outputs := make([]ArticleOutput, len(articles))

	for idx, article := range articles {
		var _article ArticleOutput
		profile, err := s.userService.GetProfile(article.AuthorId, "", userId)

		if err != nil {
			s.log.Debug("error getting profile", "error", err)
		}

		_article = formatArticle(article, *profile)

		outputs[idx] = _article
	}

	return outputs, err
}

func (s *ArticleService) GetPopularTags() ([]string, error) {
	return s.repo.GetPopularTags()
}

func (s *ArticleService) GetBySlug(
	slug string,
	userId int,
) (ArticleOutput, error) {
	article, err := s.repo.GetBySlug(slug)

	if err != nil {
		return ArticleOutput{}, err
	}

	profile, err := s.userService.GetProfile(article.AuthorId, "", userId)

	if err != nil {
		return ArticleOutput{}, err
	}

	return formatArticle(article, *profile), nil
}

func (s *ArticleService) GetIdFromSlug(slug string, userId int) (int, error) {
	article, err := s.repo.GetBySlug(slug)

	if err != nil {
		return 0, err
	}

	return article.ArticleId, nil
}

func (s *ArticleService) Update(
	slug string,
	username string,
	input ArticleUpdateInput,
) (ArticleOutput, error) {
	article, err := s.repo.Update(slug, func(a domain.Article) domain.Article {
		if input.Title != "" {
			a.Title = input.Title
		}

		if input.Description != "" {
			a.Description = input.Description
		}

		if input.Body != "" {
			a.Body = input.Body
		}
		return a
	})

	if err != nil {
		return ArticleOutput{}, err
	}

	profile, err := s.userService.GetProfile(article.AuthorId, "", 0)

	if err != nil {
		return ArticleOutput{}, err
	}

	return formatArticle(article, *profile), nil
}

func (s *ArticleService) Favorite(
	slug string,
	username string,
) (ArticleOutput, error) {
	article, err := s.repo.Update(slug, func(a domain.Article) domain.Article {
		return a
	})

	if err != nil {
		return ArticleOutput{}, err
	}

	profile, err := s.userService.GetProfile(article.AuthorId, "", 0)

	if err != nil {
		return ArticleOutput{}, err
	}

	return formatArticle(article, *profile), nil
}

func (s *ArticleService) Delete(slug string, username string) error {
	return s.repo.Delete(slug)
}
