package generator

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/rhinoman/go-commonmark"
	"github.com/xr0-org/progstack-ssg/internal/assert"
)

func Generate(srcdir, targetdir, themedir string) error {
	gen, err := newGeneration(srcdir, targetdir, themedir)
	if err != nil {
		return fmt.Errorf("cannot make generation: %w", err)
	}
	return filepath.WalkDir(
		srcdir,
		func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return fmt.Errorf("path error: %w", err)
			}
			isfile, err := checkisfile(path)
			if err != nil {
				return fmt.Errorf(
					"cannot check if file: %w", err,
				)
			}
			if isfile && filepath.Ext(path) == ".md" {
				return gen.render(path)
			}
			return nil
		},
	)
}

func checkisfile(path string) (bool, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return false, fmt.Errorf("cannot stat: %w", err)
	}
	return fi.Mode().IsRegular(), nil
}

// generation represents the operation of generating the target blog from the
// source directory
type generation struct {
	src, target string
	theme       *theme
}

const (
	themeDir  = "_theme"
	indexFile = "index.md"
)

func newGeneration(src, target, themedir string) (*generation, error) {
	assert.Printf(themedir == "", "themes NOT IMPLEMENTED\n")
	theme, err := newTheme(filepath.Join(src, themeDir))
	if err != nil {
		return nil, fmt.Errorf("cannot get theme: %w", err)
	}
	return &generation{src, target, theme}, nil
}

func (gen *generation) render(srcfile string) error {
	outputpath := gen.changetotarget(replaceext(srcfile, ".html"))
	if err := os.MkdirAll(filepath.Dir(outputpath), 0777); err != nil {
		return fmt.Errorf("cannot make dir: %w", err)
	}
	f, err := os.Create(outputpath)
	if err != nil {
		return fmt.Errorf("cannot create file: %w", err)
	}
	defer f.Close()
	switch filepath.Base(srcfile) {
	case indexFile:
		return gen.renderIndex(srcfile, f)
	default:
		return gen.renderDefault(srcfile, f)
	}
}

func (gen *generation) changetotarget(path string) string {
	rel, err := filepath.Rel(gen.src, path)
	assert.Assert(err == nil)
	return filepath.Join(gen.target, rel)
}

func replaceext(path, newext string) string {
	ext := filepath.Ext(path)
	return path[:len(path)-len(ext)] + newext
}

func (gen *generation) renderIndex(srcfile string, w io.Writer) error {
	data, err := getindexdata(srcfile, gen.src)
	if err != nil {
		return fmt.Errorf("cannot get page data: %w", err)
	}
	return gen.theme.executeIndex(w, data)
}

type indexdata struct {
	Title, Content string
	Posts          []post
}

func getindexdata(path, rootdir string) (*indexdata, error) {
	def, err := getdefaultdata(path)
	if err != nil {
		return nil, fmt.Errorf("cannot get default data: %w", err)
	}
	posts, err := getposts(path, rootdir)
	if err != nil {
		return nil, fmt.Errorf("cannot get posts: %w", err)
	}
	return &indexdata{def.Title, def.Body, posts}, nil
}

type post struct {
	Name, Link, Category string
}

func getposts(path, rootdir string) ([]post, error) {
	var posts []post
	return posts, filepath.WalkDir(
		filepath.Dir(path),
		func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return fmt.Errorf("path error: %w", err)
			}
			isfile, err := checkisfile(path)
			if err != nil {
				return fmt.Errorf(
					"cannot check if file: %w", err,
				)
			}
			if isfile && filepath.Ext(path) == ".md" {
				switch n := filepath.Base(path); n {
				case indexFile:
					return nil
				}
				def, err := getdefaultdata(path)
				if err != nil {
					return fmt.Errorf("cannot get data: %w", err)
				}
				posts = append(
					posts,
					post{
						def.Title,
						getlink(path, rootdir),
						getcategory(path, rootdir),
					},
				)
			}
			return nil
		},
	)
}

func getcategory(path, rootdir string) string {
	rel, err := filepath.Rel(rootdir, filepath.Dir(path))
	assert.Assert(err == nil)
	if rel == "." {
		return ""
	}
	return rel
}

func getlink(path, rootdir string) string {
	rel, err := filepath.Rel(rootdir, filepath.Dir(path))
	assert.Assert(err == nil)
	return replaceext(filepath.Join(rel, filepath.Base(path)), ".html")
}

func (gen *generation) renderDefault(srcfile string, w io.Writer) error {
	data, err := getdefaultdata(srcfile)
	if err != nil {
		return fmt.Errorf("cannot get page data: %w", err)
	}
	return gen.theme.executeDefault(w, data)
}

type defaultdata struct {
	Title, Body string
}

func getdefaultdata(path string) (*defaultdata, error) {
	doc, err := commonmark.ParseFile(path, commonmark.CMARK_OPT_DEFAULT)
	if err != nil {
		return nil, fmt.Errorf("cannot parse: %w", err)
	}
	title, err := gettitle(doc)
	if err != nil {
		return nil, fmt.Errorf("cannot get title: %w", err)
	}
	return &defaultdata{
		title,
		doc.RenderHtml(commonmark.CMARK_OPT_DEFAULT),
	}, nil
}

func gettitle(doc *commonmark.CMarkNode) (string, error) {
	iter := commonmark.NewCMarkIter(doc)
	defer iter.Free()
	for iter.Next() != commonmark.CMARK_EVENT_DONE {
		n := iter.GetNode()
		if n.GetNodeType() == commonmark.CMARK_NODE_HEADING &&
			n.GetHeaderLevel() == 1 {
			title := strings.TrimSpace(n.FirstChild().GetLiteral())
			if title == "" {
				return "", fmt.Errorf("empty")
			}
			return title, nil
		}
	}
	return "", fmt.Errorf("no heading node")
}
