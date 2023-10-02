package services

import (
	"github.com/berkeleytrue/conduit/internal/core/domain"
	"github.com/berkeleytrue/conduit/internal/infra/data/krono"
	"github.com/berkeleytrue/conduit/internal/infra/data/slug"
)

type (
	ArticleService struct {
		repo domain.ArticleRepository
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
		TagList        []string
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

func formatArticle(article *domain.Article) ArticleOutput {
	return ArticleOutput{
		Slug:        article.Slug,
		Title:       article.Title,
		Description: article.Description,
		Body:        article.Body,
		TagList:     article.Tags,
		CreatedAt:   article.CreatedAt,
		UpdatedAt:   article.UpdatedAt,
		// TODO: Favorited
		Favorited: false,
		// TODO: FavoritesCount
		FavoritesCount: 0,
	}
}

func NewArticleService(repo domain.ArticleRepository) *ArticleService {
	return &ArticleService{repo: repo}
}

func (s *ArticleService) Create(userId int, input ArticleCreateInput) (ArticleOutput, error) {
	article, err := s.repo.Create(domain.ArticleCreateInput{
		Title:       input.Title,
		Description: input.Description,
		Body:        input.Body,
		Tags:        input.Tags,
		AuthorId:    userId,
	})

	if err != nil {
		return ArticleOutput{}, err
	}

	return formatArticle(article), nil
}

func (s *ArticleService) List(username string, input domain.ArticleListInput) ([]ArticleOutput, error) {
	// TODO: following
	// username, ok := maybeUser.Get()

	// if !ok {
	// 	username = ""
	// }
	//
	articles, err := s.repo.List(input)

	if err != nil {
		return nil, err
	}

	outputs := make([]ArticleOutput, len(articles))

	for _, article := range articles {
		outputs = append(outputs, formatArticle(article))
	}

	return outputs, err
}

func (s *ArticleService) GetBySlug(slug string, userId int) (ArticleOutput, error) {
	article, err := s.repo.GetBySlug(slug)

	if err != nil {
		return ArticleOutput{}, err
	}

	return formatArticle(article), nil
}

func (s *ArticleService) GetIdFromSlug(slug string) (int, error) {
	article, err := s.repo.GetBySlug(slug)

	if err != nil {
		return 0, err
	}

	return article.ArticleId, nil
}

func (s *ArticleService) Update(slug string, username string, input ArticleUpdateInput) (ArticleOutput, error) {
	article, err := s.repo.Update(slug, func(a *domain.Article) *domain.Article {
		a.Title = input.Title
		a.Description = input.Description
		a.Body = input.Body
		return a
	})

	if err != nil {
		return ArticleOutput{}, err
	}

	return formatArticle(article), nil
}

func (s *ArticleService) Favorite(slug string, username string) (ArticleOutput, error) {
	article, err := s.repo.Update(slug, func(a *domain.Article) *domain.Article {
		return a
	})

	if err != nil {
		return ArticleOutput{}, err
	}

	return formatArticle(article), nil
}

func (s *ArticleService) Delete(slug string, username string) error {
	return s.repo.Delete(slug)
}

// newtype ArticleService = ArticleService
//   , favorite :: MySlug -> Username -> Om {} (ArticleServiceErrs ()) ArticleOutput
//   , unfavorite :: MySlug -> Username -> Om {} (ArticleServiceErrs ()) ArticleOutput
//   }
