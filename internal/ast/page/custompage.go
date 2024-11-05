package page

import (
	"fmt"
	"io"
	"path/filepath"

	"github.com/xr0-org/progstack-ssg/internal/assert"
	"github.com/xr0-org/progstack-ssg/internal/ast/area/sitefile"
	"github.com/xr0-org/progstack-ssg/internal/theme"
)

type custompage struct {
	title, content string
}

func CustomPage(title, content string) *custompage {
	return &custompage{title, content}
}

func (pg *custompage) Link(path string, pi PageInfo) (string, error) {
	url, err := filepath.Rel(
		pi.Root(),
		rightextpath(path, pi.DynamicLinks()),
	)
	if err != nil {
		return "", fmt.Errorf("cannot get relative path: %w", err)
	}
	return "/" + url, nil
}

func (pg *custompage) GenerateIndex(
	w io.Writer, posts []Post, pi PageInfo,
) error {
	return fmt.Errorf("custom page cannot generate index")
}

func (pg *custompage) GenerateWithoutIndex(w io.Writer, pi PageInfo) error {
	thm, err := theme.ParseTheme(pi.Theme())
	if err != nil {
		return fmt.Errorf("cannot parse theme: %w", err)
	}
	return thm.ExecuteCustom(w, &theme.CustomData{
		Title:   pg.title,
		Content: pg.content,
	})
}

func (pg *custompage) Generate(w io.Writer, pi PageInfo, index Page) error {
	return pg.GenerateWithoutIndex(w, pi)
}

func (pg *custompage) IsPost() bool { return false }

func (pg *custompage) AsPost(_, _ string) *Post {
	assert.Assert(false)
	return nil
}

func (pg *custompage) ToFile(path string) sitefile.File {
	return sitefile.PostFile(path, pg.title)
}
