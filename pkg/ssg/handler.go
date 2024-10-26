package ssg

import (
	"fmt"

	"github.com/xr0-org/progstack-ssg/internal/ast/area"
)

func GenerateSiteWithBindings(
	src, target, theme, chromastyle string,
) (map[string]string, error) {
	a, err := area.ParseArea(src, chromastyle)
	if err != nil {
		return nil, fmt.Errorf("cannot parse area: %w", err)
	}
	return a.GenerateWithBindings(target, theme)
}
