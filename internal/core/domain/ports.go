package domain

type (
	UserRepository interface {
	  GetByID(id string) (*User, error)
	  GetByEmail(email string) (*User, error)
	}
	ArticleRepository interface {
	}
	CommentRepository interface {
	}

	UserService interface {
	}
	ArticleService interface {
	}
	CommentService interface {
	}
)
