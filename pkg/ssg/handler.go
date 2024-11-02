package ssg

import (
	"fmt"
	"time"

	"github.com/xr0-org/progstack-ssg/internal/ast/area"
)

// A File is any URL-accessible resource in a site.
type File interface {
	// Path is the path on disk to the generated File.
	Path() string

	// IsPost indicates whether or not the File is a post.
	IsPost() bool

	// Title is the title of the post.
	PostTitle() string

	// Time is the timestamp associated with the post, if available.
	PostTime() (time.Time, bool)
}

func GenerateSiteWithBindings(
	src, target, theme, chromastyle string,
) (map[string]File, error) {
	a, err := area.ParseArea(src, chromastyle)
	if err != nil {
		return nil, fmt.Errorf("cannot parse area: %w", err)
	}
	m1, err := a.GenerateWithBindings(target, theme)
	if err != nil {
		return nil, fmt.Errorf("cannot generate: %w", err)
	}
	m2 := map[string]File{}
	for k, v := range m1 {
		m2[k] = v
	}
	return m2, nil
}
