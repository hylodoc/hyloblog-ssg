package area

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/xr0-org/progstack-ssg/internal/assert"
	"github.com/xr0-org/progstack-ssg/internal/ast/area/areainfo"
	"github.com/xr0-org/progstack-ssg/internal/ast/area/readdir"
	"github.com/xr0-org/progstack-ssg/internal/ast/area/sitefile"
	"github.com/xr0-org/progstack-ssg/internal/ast/page"
)

const (
	indexFile  = "index.md"
	ignoreFile = ".progstackignore"
)

type Area struct {
	prefix     string
	subareas   []Area
	pages      map[string]page.Page
	otherfiles map[string]readdir.File
}

func newarea(prefix string) *Area {
	return &Area{
		prefix,
		[]Area{},
		map[string]page.Page{},
		map[string]readdir.File{},
	}
}

func ParseArea(dir, chromastyle string) (*Area, error) {
	return parse(dir, dir, areainfo.NewParseInfo(chromastyle))
}

func parse(dir, parent string, info *areainfo.ParseInfo) (*Area, error) {
	info, err := info.Descend(dir, ignoreFile)
	if err != nil {
		return nil, fmt.Errorf("cannot descend info: %w", err)
	}
	prefix, err := getprefix(dir, parent)
	if err != nil {
		return nil, fmt.Errorf("cannot get prefix: %w", err)
	}
	A := newarea(prefix)
	areadir, err := readdir.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("cannot read dir: %w", err)
	}
	for _, d := range areadir.Directories() {
		path := d.Path()
		base := filepath.Base(path)
		if info.ShouldIgnore(base) {
			continue
		}
		if base == ".git" {
			continue
		}
		a, err := parse(path, dir, info)
		if err != nil {
			return nil, fmt.Errorf(
				"cannot parse subarea %q: %w", path, err,
			)
		}
		A.subareas = append(A.subareas, *a)
	}
	for _, f := range areadir.Files() {
		path := f.Path()
		base := filepath.Base(path)
		if info.ShouldIgnore(base) {
			continue
		}
		if filepath.Ext(base) != ".md" {
			if includefile(base) {
				A.otherfiles[base] = readdir.NewFile(path)
			}
			continue
		}
		page, err := parsepage(path, info)
		if err != nil {
			return nil, fmt.Errorf(
				"cannot parse page %q: %w",
				path, err,
			)
		}
		A.pages[base] = *page
	}
	return A, nil
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

func parsepage(path string, info *areainfo.ParseInfo) (*page.Page, error) {
	if gitdir, ok := info.GitDir(); ok {
		return page.ParsePageGit(path, gitdir, info.ChromaStyle())
	}
	return page.ParsePage(path, info.ChromaStyle())
}

func includefile(name string) bool {
	switch filepath.Ext(name) {
	case ".png", ".jpg", ".jpeg", ".gif":
		return true
	default:
		return false
	}
}

func (A *Area) GenerateSite(
	target string, theme string, p areainfo.Purpose,
) error {
	return A.generate(target, areainfo.NewGenInfo(theme, target, nil, p))
}

func (A *Area) generate(target string, g *areainfo.GenInfo) error {
	if index, ok := A.pages[indexFile]; ok {
		g = g.WithNewIndex(&index)
	}
	dir := filepath.Join(target, A.prefix)
	if err := os.MkdirAll(dir, 0777); err != nil {
		return fmt.Errorf("cannot make dir: %w", err)
	}
	for _, a := range A.subareas {
		if err := a.generate(dir, g); err != nil {
			return fmt.Errorf(
				"cannot generate subarea %q: %w",
				filepath.Join(dir, a.prefix), err,
			)
		}
	}
	for name, page := range A.pages {
		if err := A.generatepage(name, dir, &page, g); err != nil {
			return fmt.Errorf(
				"cannot generate page: %q: %w", name, err,
			)
		}
	}
	for name, f := range A.otherfiles {
		if err := fcopy(f.Path(), filepath.Join(dir, name)); err != nil {
			return fmt.Errorf("cannot copy %q: %w", name, err)
		}
	}
	return nil
}

func (A *Area) generatepage(
	name, dir string, page *page.Page, g *areainfo.GenInfo,
) error {
	f, err := os.Create(genpagehtmlpath(name, dir))
	if err != nil {
		return fmt.Errorf("cannot create file: %w", err)
	}
	defer f.Close()
	if name == indexFile {
		return page.GenerateIndex(
			f, A.getposts(dir, g), g.Theme(),
		)
	}
	if index, ok := g.GetIndex(); ok {
		return page.Generate(f, g.Theme(), index)
	}
	return page.GenerateWithoutIndex(f, g.Theme())
}

func genpagehtmlpath(name, dir string) string {
	return filepath.Join(dir, replaceext(name, ".html"))
}

func replaceext(path, newext string) string {
	ext := filepath.Ext(path)
	return path[:len(path)-len(ext)] + newext
}

func (A *Area) getposts(dir string, g *areainfo.GenInfo) []page.Post {
	var posts []page.Post
	for _, a := range A.subareas {
		posts = append(
			posts,
			a.getposts(filepath.Join(dir, a.prefix), g)...,
		)
	}
	for name := range A.pages {
		if name == indexFile {
			continue
		}
		pg := A.pages[name]
		link, err := pg.Link(
			filepath.Join(dir, name),
			g.Root(),
			g.DynamicLinks(),
		)
		assert.Assert(err == nil)
		posts = append(posts, *page.CreatePost(&pg, A.prefix, link))
	}
	return posts
}

func fcopy(srcpath, dstpath string) error {
	src, err := os.Open(srcpath)
	if err != nil {
		return fmt.Errorf("cannot open source: %w", err)
	}
	defer src.Close()
	dst, err := os.Create(dstpath)
	if err != nil {
		return fmt.Errorf("cannot create destination: %w", err)
	}
	defer dst.Close()
	if _, err := io.Copy(dst, src); err != nil {
		return fmt.Errorf("io copy error: %w", err)
	}
	if err := dst.Sync(); err != nil {
		return fmt.Errorf("flush error: %w", err)
	}
	return nil
}

type Handler struct {
	h         http.Handler
	targetdir string
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.h.ServeHTTP(w, r)
}

func (h *Handler) Destroy() error {
	return os.RemoveAll(h.targetdir)
}

func (A *Area) Handler(theme string) (*Handler, error) {
	target, err := os.MkdirTemp("", "")
	if err != nil {
		return nil, fmt.Errorf("cannot make tempdir: %w", err)
	}
	purpose := areainfo.PurposeDynamicServe
	if err := A.GenerateSite(target, theme, purpose); err != nil {
		return nil, fmt.Errorf("cannot generate site: %w", err)
	}
	r := mux.NewRouter()
	r.StrictSlash(true)
	return &Handler{r, target}, A.registerhandlers(
		target, areainfo.NewGenInfo(theme, target, nil, purpose), r,
	)
}

func (A *Area) registerhandlers(
	target string, g *areainfo.GenInfo, mux *mux.Router,
) error {
	if index, ok := A.pages[indexFile]; ok {
		g = g.WithNewIndex(&index)
	}
	dir := filepath.Join(target, A.prefix)
	for _, a := range A.subareas {
		if err := a.registerhandlers(dir, g, mux); err != nil {
			return fmt.Errorf(
				"cannot register subarea %q: %w",
				filepath.Join(dir, a.prefix), err,
			)
		}
	}
	for name := range A.pages {
		pg := A.pages[name]
		path, err := pagehostpath(&pg, name, dir, g.Root())
		if err != nil {
			return fmt.Errorf(
				"cannot make path for %q: %w", name, err,
			)
		}
		mux.HandleFunc(path, filehandler(genpagehtmlpath(name, dir)))
	}
	for name := range A.otherfiles {
		path, err := filehostpath(name, dir, g.Root())
		if err != nil {
			return fmt.Errorf(
				"cannot make path for %q: %w", name, err,
			)
		}
		mux.HandleFunc(path, filehandler(filepath.Join(dir, name)))
	}
	return nil
}

func pagehostpath(pg *page.Page, name, dir, rootdir string) (string, error) {
	if name == indexFile {
		path, err := filepath.Rel(rootdir, dir)
		if err != nil {
			return "", err
		}
		if path == "." {
			return "/", nil
		}
		return fmt.Sprintf("/%s", path), nil
	}
	return pg.Link(filepath.Join(dir, name), rootdir, true)
}

func filehostpath(name, dir, rootdir string) (string, error) {
	rel, err := filepath.Rel(rootdir, filepath.Join(dir, name))
	if err != nil {
		return "", fmt.Errorf("cannot get relative path: %w", err)
	}
	assert.Assert(rel != ".")
	return filepath.Join("/", rel), nil
}

func filehandler(path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL, "->", path)
		http.ServeFile(w, r, path)
	}
}

func (A *Area) GenerateWithBindings(
	target, theme string,
) (map[string]sitefile.File, error) {
	g := areainfo.NewGenInfo(
		theme, target, nil, areainfo.PurposeDynamicServe,
	)
	if err := A.generate(target, g); err != nil {
		return nil, fmt.Errorf("cannot generate: %w", err)
	}
	bindings := map[string]sitefile.File{}
	if err := A.handlebindings(target, g, bindings); err != nil {
		return nil, fmt.Errorf("cannot get bindings: %w", err)
	}
	return bindings, nil
}

func (A *Area) handlebindings(
	target string, g *areainfo.GenInfo, m map[string]sitefile.File,
) error {
	if index, ok := A.pages[indexFile]; ok {
		g = g.WithNewIndex(&index)
	}
	dir := filepath.Join(target, A.prefix)
	for _, a := range A.subareas {
		if err := a.handlebindings(dir, g, m); err != nil {
			return fmt.Errorf(
				"cannot handle subarea %q: %w",
				filepath.Join(dir, a.prefix), err,
			)
		}
	}
	for name := range A.pages {
		pg := A.pages[name]
		path, err := pagehostpath(&pg, name, dir, g.Root())
		if err != nil {
			return fmt.Errorf(
				"cannot make path for %q: %w", name, err,
			)
		}
		m[path] = sitefile.NewFile(genpagehtmlpath(name, dir), true)
	}
	for name := range A.otherfiles {
		path, err := filehostpath(name, dir, g.Root())
		if err != nil {
			return fmt.Errorf(
				"cannot make path for %q: %w", name, err,
			)
		}
		m[path] = sitefile.NewFile(filepath.Join(dir, name), false)
	}
	return nil
}

type LiveHandler struct {
	src, theme, chromastyle string
}

func CreateLiveHandler(src, theme, chromastyle string) *LiveHandler {
	return &LiveHandler{src, theme, chromastyle}
}

func (lh *LiveHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h, err := lh.genhandler()
	if err != nil {
		log.Fatal(err)
	}
	defer h.Destroy()
	h.ServeHTTP(w, r)
}

func (lh *LiveHandler) genhandler() (*Handler, error) {
	blog, err := ParseArea(lh.src, lh.chromastyle)
	if err != nil {
		return nil, fmt.Errorf("cannot parse: %w", err)
	}
	h, err := blog.Handler(lh.theme)
	if err != nil {
		return nil, fmt.Errorf("cannot get handler: %w", err)
	}
	return h, nil
}
