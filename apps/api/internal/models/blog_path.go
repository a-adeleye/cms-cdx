package models

import (
	"errors"
	"regexp"
	"strings"
)

var blogPathSegmentPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9_-]*$`)

// CanonicalBlogPath returns the site-relative path used for article pages.
func CanonicalBlogPath(value string) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "/articles", nil
	}
	if !strings.HasPrefix(value, "/") {
		value = "/" + value
	}
	value = strings.TrimRight(value, "/")
	if value == "" {
		return "", errors.New("blog path must not be the site root")
	}

	segments := strings.Split(strings.TrimPrefix(value, "/"), "/")
	for _, segment := range segments {
		if !blogPathSegmentPattern.MatchString(segment) {
			return "", errors.New("blog path must use slash-separated letters, numbers, hyphens, or underscores")
		}
	}
	return "/" + strings.Join(segments, "/"), nil
}
