package services

import "fmt"

func (p PublicProfile) String() string {
	return fmt.Sprintf(
		"{Username: %v, Bio: %v, Image: %v, Following: %v}",
		p.Username,
		p.Bio,
		p.Image,
		p.IsFollowing,
	)
}

func (a ArticleOutput) String() string {
	return fmt.Sprintf(
		"{Slug: %v, Title: %v, Description: %v, Body: %v, Tags: %v, CreatedAt: %v, UpdatedAt: %v, Favorited: %v, FavoritesCount: %v, Author: %v}",
		a.Slug,
		a.Title,
		a.Description,
		a.Body,
		a.Tags,
		a.CreatedAt,
		a.UpdatedAt,
		a.IsFavorited,
		a.FavoritesCount,
		a.Author,
	)
}
