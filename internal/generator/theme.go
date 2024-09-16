package generator

import (
	"fmt"
	"io"
	"path/filepath"
	"text/template"
)

type theme struct {
	index, def *template.Template
}

const (
	themeIndex   = "index.html"
	themeDefault = "_default.html"
)

func newTheme(dir string) (*theme, error) {
	index, err := template.ParseFiles(filepath.Join(dir, themeIndex))
	if err != nil {
		return nil, fmt.Errorf("cannot get index: %w", err)
	}
	def, err := template.ParseFiles(filepath.Join(dir, themeDefault))
	if err != nil {
		return nil, fmt.Errorf("cannot get default: %w", err)
	}
	return &theme{index, def}, nil
}

func (thm *theme) executeIndex(w io.Writer, data any) error {
	return thm.index.Execute(w, data)
}

func (thm *theme) executeDefault(w io.Writer, data any) error {
	return thm.def.Execute(w, data)
}
