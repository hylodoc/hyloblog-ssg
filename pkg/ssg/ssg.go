package ssg

import "github.com/xr0-org/progstack-ssg/internal/generator"

// Generate produces a blog at target using the files in src. The themedir
// provides the search path for any themes referenced in src.
func Generate(src, target, themedir string) error {
	return generator.Generate(src, target, themedir)
}
