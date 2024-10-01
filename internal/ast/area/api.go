package area

import "path/filepath"

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

func (a *Area) Interface() AreaInterface {
	return a.geninterface("/")
}

func (a *Area) geninterface(parentdir string) AreaInterface {
	dir := a.augmentdir(parentdir)
	return &areainterface{a.prefix, a.getsubareas(dir), a.genpages(dir)}
}

func (a *Area) augmentdir(parentdir string) string {
	if a.prefix != "." {
		return filepath.Join(parentdir, a.prefix)
	}
	return parentdir
}

func (a *Area) getsubareas(dir string) []AreaInterface {
	subareas := make([]AreaInterface, len(a.subareas))
	for i := range a.subareas {
		subareas[i] = a.subareas[i].geninterface(dir)
	}
	return subareas
}

type pageinterface string

func (p pageinterface) Link() string { return string(p) }

func (a *Area) genpages(dir string) []PageInterface {
	var pages []PageInterface
	for name := range a.pages {
		pages = append(
			pages,
			pageinterface(filepath.Join(dir, replaceext(name, ""))),
		)
	}
	return pages
}
