package page

import (
	"io"

	"github.com/xr0-org/progstack-ssg/internal/ast/area/sitefile"
	"github.com/xr0-org/progstack-ssg/internal/theme"
)

type Page interface {
	Title() (string, error)
	Link(path string, pi PageInfo) (string, error)

	GenerateIndex(w io.Writer, posts []Post, pi PageInfo) error
	Generate(w io.Writer, pi PageInfo, index Page) error
	GenerateWithoutIndex(w io.Writer, pi PageInfo) error

	IsPost() bool
	AsPost(category, link string) *Post

	ToFile(path string, pi PageInfo) (sitefile.File, error)
}

type PageInfo interface {
	Theme() *theme.Theme
	Head() string
	Foot() string
	Root() string
	DynamicLinks() bool
}
