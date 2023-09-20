package domain

type (
	User struct {
		ID       string `json:"id"`
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	Article struct {
		ID      string `json:"id"`
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	Comment struct {
		ID      string `json:"id"`
		Content string `json:"content"`
		UserID  string `json:"userId"`
	}
)
