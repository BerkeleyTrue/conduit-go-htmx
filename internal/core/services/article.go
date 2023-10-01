package services

import (
	"github.com/berkeleytrue/conduit/internal/core/domain"
	"github.com/berkeleytrue/conduit/internal/infra/data/krono"
	"github.com/berkeleytrue/conduit/internal/infra/data/maybe"
	"github.com/berkeleytrue/conduit/internal/infra/data/slug"
)

type (
	ArticleService struct {
		repo domain.ArticleRepository
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
		title       string
		description string
		body        string
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

func (s *ArticleService) List(maybeUser maybe.Maybe[string], input domain.ArticleListInput) ([]ArticleOutput, error) {
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

// newtype ArticleService = ArticleService
//   { list :: { username :: (Maybe Username), input :: ArticleListInput } -> Om {} (ArticleServiceErrs ()) (Array ArticleOutput)
//   , getBySlug :: MySlug -> (Maybe Username) -> Om {} (ArticleServiceErrs ()) ArticleOutput
//   , getIdFromSlug :: MySlug -> Om {} (ArticleServiceErrs ()) ArticleId
//   , update :: MySlug -> Username -> ArticleUpdateInput -> Om {} (ArticleServiceErrs ()) ArticleOutput
//   , delete :: MySlug -> Om {} (ArticleServiceErrs ()) Unit
//   , favorite :: MySlug -> Username -> Om {} (ArticleServiceErrs ()) ArticleOutput
//   , unfavorite :: MySlug -> Username -> Om {} (ArticleServiceErrs ()) ArticleOutput
//   }
