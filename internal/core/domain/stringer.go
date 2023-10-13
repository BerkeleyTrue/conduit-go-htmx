package domain

import "fmt"

func (u *User) String() string {
	return fmt.Sprintf(
		"{ UserId: %d, Username: %s, Email: %s, Password: ***, Bio: %s, Image: %s, CreatedAt: %s, UpdatedAt: %s }",
		u.UserId,
		u.Username,
		u.Email,
		u.Bio,
		u.Image,
		u.CreatedAt,
		u.UpdatedAt,
	)
}

func (a *Article) String() string {
	return fmt.Sprintf("{ ArticleId: %d, AuthorId: %d, Title: %s, Slug: %s, Description: %s, Body: %s, Tags: %v, CreatedAt: %s, UpdatedAt: %s, }", a.ArticleId,
		a.AuthorId,
		a.Title,
		a.Slug,
		a.Description,
		a.Body,
		a.Tags,
		a.CreatedAt,
		a.UpdatedAt,
	)
}
