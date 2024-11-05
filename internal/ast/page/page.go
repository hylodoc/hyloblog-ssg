package page

import (
	"io"

	"github.com/xr0-org/progstack-ssg/internal/ast/area/sitefile"
)

type Page interface {
	Link(path string, pi PageInfo) (string, error)

	GenerateIndex(w io.Writer, posts []Post, pi PageInfo) error
	Generate(w io.Writer, pi PageInfo, index Page) error
	GenerateWithoutIndex(w io.Writer, pi PageInfo) error

	IsPost() bool
	AsPost(category, link string) *Post

	ToFile(path string) sitefile.File
}

type PageInfo interface {
	Theme() string
	Head() string
	Foot() string
	Root() string
	DynamicLinks() bool
}
