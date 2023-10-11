package slug

import (
	"regexp"
	"strings"
)

// Slug is a custom data type for representing slugs.
type Slug string

// Remove any non-alphanumeric characters (except dashes)
var regex *regexp.Regexp = regexp.MustCompile("[^a-zA-Z0-9-]+")

// NewSlug creates a new Slug from a title string.
func NewSlug(title string) Slug {
	// Convert to lowercase
	slug := strings.ToLower(title)

	// Replace spaces with dashes
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = regex.ReplaceAllString(slug, "")

	return Slug(slug)
}

// String returns the string representation of a Slug.
func (s Slug) String() string {
	return string(s)
}
