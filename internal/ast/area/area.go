package area

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/xr0-org/progstack-ssg/internal/assert"
	"github.com/xr0-org/progstack-ssg/internal/ast/area/areainfo"
	"github.com/xr0-org/progstack-ssg/internal/ast/area/ignorer"
	"github.com/xr0-org/progstack-ssg/internal/ast/page"
)

const (
	indexFile  = "index.md"
	ignoreFile = ".progstackignore"
)

type Area struct {
	prefix   string
	pages    map[string]page.Page
	subareas []Area
}

func ParseArea(dir string) (*Area, error) {
	return parse(dir, dir, ignorer.Empty())
}

func parse(dir, parent string, oldign *ignorer.Ignorer) (*Area, error) {
	prefix, err := getprefix(dir, parent)
	if err != nil {
		return nil, fmt.Errorf("cannot get prefix: %w", err)
	}
	ign, err := oldign.Augment(filepath.Join(dir, ignoreFile))
	if err != nil {
		return nil, fmt.Errorf("cannot get ignore: %w", err)
	}
	A := Area{prefix, map[string]page.Page{}, []Area{}}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("cannot read dir: %w", err)
	}
	for _, e := range entries {
		if ign.ShouldIgnore(e.Name()) {
			continue
		}
		if e.IsDir() {
			if e.Name() == ".git" {
				continue
			}
			path := filepath.Join(dir, e.Name())
			a, err := parse(path, dir, ign)
			if err != nil {
				return nil, fmt.Errorf(
					"cannot parse subarea %q: %w", path, err,
				)
			}
			A.subareas = append(A.subareas, *a)
		}
	}
	for _, e := range entries {
		if ign.ShouldIgnore(e.Name()) {
			continue
		}
		if !e.IsDir() {
			t := e.Type()
			assert.Printf(t.IsRegular(), "unknown file type %v", t)
			path := filepath.Join(dir, e.Name())
			if filepath.Ext(path) != ".md" {
				continue
			}
			page, err := page.ParsePage(path)
			if err != nil {
				return nil, fmt.Errorf(
					"cannot parse page %q: %w",
					path, err,
				)
			}
			A.pages[e.Name()] = *page
		}
	}
	return &A, nil
}

func getprefix(dir, parent string) (string, error) {
	prefix, err := filepath.Rel(parent, dir)
	if err != nil {
		return "", fmt.Errorf("cannot get relative path: %w", err)
	}
	if prefix == "." {
		return "", nil
	}
	return prefix, nil
}

func (A *Area) GenerateSite(
	target string, theme string, p areainfo.Purpose,
) error {
	return A.generate(target, areainfo.Create(theme, target, nil, p))
}

func (A *Area) generate(target string, ainfo *areainfo.AreaInfo) error {
	if index, ok := A.pages[indexFile]; ok {
		ainfo = ainfo.WithNewIndex(&index)
	}
	dir := filepath.Join(target, A.prefix)
	if err := os.MkdirAll(dir, 0777); err != nil {
		return fmt.Errorf("cannot make dir: %w", err)
	}
	for _, a := range A.subareas {
		if err := a.generate(dir, ainfo); err != nil {
			return fmt.Errorf(
				"cannot generate subarea %q: %w",
				filepath.Join(dir, a.prefix), err,
			)
		}
	}
	for name, page := range A.pages {
		if err := A.generatepage(name, dir, &page, ainfo); err != nil {
			return fmt.Errorf(
				"cannot generate page: %q: %w", name, err,
			)
		}
	}
	return nil
}

func (A *Area) generatepage(
	name, dir string, page *page.Page, ainfo *areainfo.AreaInfo,
) error {
	f, err := os.Create(genpagehtmlpath(name, dir))
	if err != nil {
		return fmt.Errorf("cannot create file: %w", err)
	}
	defer f.Close()
	if name == indexFile {
		return page.GenerateIndex(
			f, A.getposts(dir, ainfo), ainfo.Theme(),
		)
	}
	if index, ok := ainfo.GetIndex(); ok {
		return page.Generate(f, ainfo.Theme(), index)
	}
	return page.GenerateWithoutIndex(f, ainfo.Theme())
}

func genpagehtmlpath(name, dir string) string {
	return filepath.Join(dir, replaceext(name, ".html"))
}

func replaceext(path, newext string) string {
	ext := filepath.Ext(path)
	return path[:len(path)-len(ext)] + newext
}

func (A *Area) getposts(dir string, ainfo *areainfo.AreaInfo) []page.Post {
	var posts []page.Post
	for _, a := range A.subareas {
		posts = append(
			posts,
			a.getposts(filepath.Join(dir, a.prefix), ainfo)...,
		)
	}
	for name := range A.pages {
		if name == indexFile {
			continue
		}
		link, err := filepath.Rel(
			ainfo.Root(), pagelink(name, dir, ainfo),
		)
		assert.Assert(err == nil)
		p := A.pages[name]
		posts = append(posts, *page.CreatePost(&p, A.prefix, link))
	}
	return posts
}

func pagelink(name, dir string, ainfo *areainfo.AreaInfo) string {
	p := filepath.Join(dir, name)
	if ainfo.LinksWithHtmlExt() {
		return replaceext(p, ".html")
	}
	return replaceext(p, "")
}

func (A *Area) Handler(target, theme string) (http.Handler, error) {
	purpose := areainfo.PurposeDynamicServe
	if err := A.GenerateSite(target, theme, purpose); err != nil {
		return nil, fmt.Errorf("cannot generate site: %w", err)
	}
	r := mux.NewRouter()
	return r, A.registerhandlers(
		target, areainfo.Create(theme, target, nil, purpose), r,
	)
}

func (A *Area) registerhandlers(
	target string, ainfo *areainfo.AreaInfo, mux *mux.Router,
) error {
	if index, ok := A.pages[indexFile]; ok {
		ainfo = ainfo.WithNewIndex(&index)
	}
	dir := filepath.Join(target, A.prefix)
	for _, a := range A.subareas {
		if err := a.registerhandlers(dir, ainfo, mux); err != nil {
			return fmt.Errorf(
				"cannot register subarea %q: %w",
				filepath.Join(dir, a.prefix), err,
			)
		}
	}
	for name := range A.pages {
		path, err := hostpath(name, dir, ainfo.Root())
		if err != nil {
			return fmt.Errorf(
				"cannot make path for %q: %w", name, err,
			)
		}
		url := fmt.Sprintf("/%s", path)
		mux.HandleFunc(url, handler(name, dir))
	}
	return nil
}

func hostpath(name, dir, rootdir string) (string, error) {
	if name == indexFile {
		path, err := filepath.Rel(rootdir, dir)
		if err != nil {
			return "", err
		}
		if path == "." {
			return "", nil
		}
		return path, nil
	}
	return filepath.Rel(rootdir, replaceext(filepath.Join(dir, name), ""))
}

func handler(name, dir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := genpagehtmlpath(name, dir)
		log.Println(r.URL, "->", path)
		http.ServeFile(
			w, r, genpagehtmlpath(name, dir),
		)
	}
}

type LiveHandler struct {
	src, theme string
}

func CreateLiveHandler(src, theme string) *LiveHandler {
	return &LiveHandler{src, theme}
}

func (lh *LiveHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h, target, err := lh.genhandler()
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(target)
	h.ServeHTTP(w, r)
}

func (lh *LiveHandler) genhandler() (http.Handler, string, error) {
	blog, err := ParseArea(lh.src)
	if err != nil {
		return nil, "", fmt.Errorf("cannot parse: %w", err)
	}
	target, err := os.MkdirTemp("", "")
	if err != nil {
		return nil, "", fmt.Errorf("cannot make tempdir: %w", err)
	}
	h, err := blog.Handler(target, lh.theme)
	if err != nil {
		return nil, target, fmt.Errorf("cannot get handler: %w", err)
	}
	return h, target, nil
}
