package pkg

import (
	"regexp"
	"strings"
	"unicode"
)

func Slug(name string) string {

	// Convert to lowercase
	lowerName := strings.ToLower(name)

	// Remove any character that is not a letter, number, or space
	var cleanedName strings.Builder
	for _, r := range lowerName {
		if unicode.IsLetter(r) || unicode.IsNumber(r) || unicode.IsSpace(r) {
			cleanedName.WriteRune(r)
		}
	}

	// Collapse multiple spaces into a single space
	spaceRegexp := regexp.MustCompile(`\s+`)
	singleSpacedName := spaceRegexp.ReplaceAllString(cleanedName.String(), " ")

	// Replace spaces with hyphens
	slug := strings.ReplaceAll(singleSpacedName, " ", "-")

	return slug
}
