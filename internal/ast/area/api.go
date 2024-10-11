package area

import (
	"fmt"
	"path/filepath"
)

type AreaInterface interface {
	Prefix() string
	Subareas() []AreaInterface
	Pages() []PageInterface
}

type PageInterface interface {
	Link() string
}

type areainterface struct {
	prefix   string
	subareas []AreaInterface
	pages    []PageInterface
}

func (a *areainterface) Prefix() string            { return a.prefix }
func (a *areainterface) Subareas() []AreaInterface { return a.subareas }
func (a *areainterface) Pages() []PageInterface    { return a.pages }

func (a *Area) Interface() (AreaInterface, error) {
	return a.geninterface("/")
}

func (a *Area) geninterface(parentdir string) (AreaInterface, error) {
	dir := a.augmentdir(parentdir)
	subareas, err := a.getsubareas(dir)
	if err != nil {
		return nil, fmt.Errorf("cannot generate subareas: %w", err)
	}
	pages, err := a.genpages(dir)
	if err != nil {
		return nil, fmt.Errorf("cannot generate pages: %w", err)
	}
	return &areainterface{a.prefix, subareas, pages}, nil
}

func (a *Area) augmentdir(parentdir string) string {
	if a.prefix != "." {
		return filepath.Join(parentdir, a.prefix)
	}
	return parentdir
}

func (a *Area) getsubareas(dir string) ([]AreaInterface, error) {
	subareas := make([]AreaInterface, len(a.subareas))
	for i := range a.subareas {
		sub, err := a.subareas[i].geninterface(dir)
		if err != nil {
			return nil, fmt.Errorf("cannot make subarea: %w", err)
		}
		subareas[i] = sub
	}
	return subareas, nil
}

type pageinterface string

func (p pageinterface) Link() string { return string(p) }

func (a *Area) genpages(dir string) ([]PageInterface, error) {
	var pages []PageInterface
	for name := range a.pages {
		pg := a.pages[name]
		path, err := hostpath(&pg, name, dir, "/")
		if err != nil {
			return nil, fmt.Errorf("cannot make hostpath: %w", err)
		}
		pages = append(pages, pageinterface(path))
	}
	return pages, nil
}
func genlink(name, dir string) string {
	if name == indexFile {
		return dir
	}
	return filepath.Join(dir, replaceext(name, ""))
}
