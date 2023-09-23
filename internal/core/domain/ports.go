package domain

type (
	Updater[T any] func(*T) *T
	UserRepository interface {
		Create(password string) (*User, error)
		GetByID(id string) (*User, error)
		GetByEmail(email string) (*User, error)
		GetByUsername(username string) (*User, error)
		Update(userId string, updater Updater[*User]) (*User, error)
		Follow(userId, authorId string) (*User, error)
		Unfollow(userId, authorId string) (*User, error)
	}

	ArticleCreateInput struct {
		title       string
		description string
		body        string
		tags        []string
		authorId    string
	}

	ArticleListInput struct {
		tag       string
		author    string // authorId
		favorited string // authorId
		limit     int8
		offset    int8
	}

	ArticleRepository interface {
		Create(input ArticleCreateInput) (*Article, error)
		GetById(articleId string) (*Article, error)
		GetBySlug(mySlug string) (*Article, error)
		List(input ArticleListInput) []*Article
		Update(slug string, updater Updater[*Article]) (*Article, error)
		Delete(slug string) error
	}

	CommentCreateInput struct {
		body      string
		authorId  string // UserId
		articleId string // ArticleId
	}

	CommentRepository interface {
		Create(input CommentCreateInput) (*Comment, error)
		GetById(commentId string) (*Comment, error)
		GetByArticleId(articleId string) ([]*Comment, error)
		Update(commentId string, updater Updater[*Comment]) (*Comment, error)
		Delete(commentId string) error
	}
)
