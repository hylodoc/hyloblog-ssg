package ast

import (
	"fmt"
	"testing"

	"github.com/hylodoc/hylodoc-ssg/internal/ast/area"
)

func TestParse(t *testing.T) {
	if err := testParse(); err != nil {
		t.Fatal(err)
	}
}

func testParse() error {
	const (
		themes = "../../theme/lit"
		src    = "test/src"
		target = "test/target"
	)
	blog, err := area.ParseArea(src, src)
	if err != nil {
		return fmt.Errorf("cannot parse: %w", err)
	}
	if err := blog.GenerateSite(target, themes); err != nil {
		return fmt.Errorf("cannot generate: %w", err)
	}
	return nil
}
