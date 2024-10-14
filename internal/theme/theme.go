package theme

import (
	"fmt"
	"io"
	"path/filepath"
	"text/template"
)

type Theme struct {
	index, def *template.Template
}

const (
	themeIndex   = "index.html"
	themeDefault = "_default.html"
)

func ParseTheme(dir string) (*Theme, error) {
	index, err := template.ParseFiles(filepath.Join(dir, themeIndex))
	if err != nil {
		return nil, fmt.Errorf("cannot get index: %w", err)
	}
	def, err := template.ParseFiles(filepath.Join(dir, themeDefault))
	if err != nil {
		return nil, fmt.Errorf("cannot get default: %w", err)
	}
	return &Theme{index, def}, nil
}

type IndexData struct {
	Title, Content string
	Posts          []Post
}

func (thm *Theme) ExecuteIndex(w io.Writer, data *IndexData) error {
	return thm.index.Execute(w, data)
}

type Post struct {
	Title, Link, Category, Date string
	Authors                     []Author
}

type Author struct {
	Name, Page string
}

type DefaultData struct {
	Title, Content string
	SiteTitle      string
	Date           string
	Authors        []Author
}

func (thm *Theme) ExecuteDefault(w io.Writer, data *DefaultData) error {
	return thm.def.Execute(w, data)
}
