package textutil

import (
	"regexp"
	"strings"
)

// nonAlnum matches any run of non-alphanumeric characters.
var nonAlnum = regexp.MustCompile(`[^a-zA-Z0-9]+`)

// Slugify converts a string into a URL-friendly slug.
func Slugify(s string) string {
	s = strings.ToLower(s)
	s = nonAlnum.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	return s
}
