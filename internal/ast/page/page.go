package page

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/rhinoman/go-commonmark"
	"github.com/xr0-org/progstack-ssg/internal/assert"
	"github.com/xr0-org/progstack-ssg/internal/theme"
)

type Page struct {
	title, url string
	doc        *commonmark.CMarkNode
}

func ParsePage(path string) (*Page, error) {
	buf, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("cannot read file: %w", err)
	}
	components, err := separate(string(buf))
	if err != nil {
		return nil, fmt.Errorf("cannot separate: %w", err)
	}
	doc := commonmark.ParseDocument(
		components.content, commonmark.CMARK_OPT_DEFAULT,
	)
	m, err := parsemetadata(components.metadata)
	if err != nil {
		return nil, fmt.Errorf("cannot parse metadata: %w", err)
	}
	return &Page{title: gettitle(doc), url: m.URL, doc: doc}, nil
}

type components struct {
	metadata string
	content  string
}

func separate(s string) (*components, error) {
	s = strings.TrimSpace(s)
	if len(s) < 3 || s[:3] != "---" {
		return &components{"", s}, nil
	}
	s = s[3:]
	endmeta := strings.Index(s, "---")
	if endmeta == -1 {
		return nil, fmt.Errorf("unclosed metadata section")
	}
	return &components{
		strings.TrimSpace(s[:endmeta]),
		strings.TrimSpace(s[endmeta+3:]),
	}, nil
}

func gettitle(doc *commonmark.CMarkNode) string {
	iter := commonmark.NewCMarkIter(doc)
	defer iter.Free()
	for iter.Next() != commonmark.CMARK_EVENT_DONE {
		n := iter.GetNode()
		if n.GetNodeType() == commonmark.CMARK_NODE_HEADING &&
			n.GetHeaderLevel() == 1 {
			return strings.TrimSpace(n.FirstChild().GetLiteral())
		}
	}
	return ""
}

func (pg *Page) Link(path, rootdir string, dynamiclinks bool) (string, error) {
	if pg.url != "" {
		if dynamiclinks {
			return pg.url, nil
		}
		log.Printf("warning: %q has custom url in static mode\n", path)
	}
	url, err := filepath.Rel(rootdir, rightextpath(path, dynamiclinks))
	if err != nil {
		return "", fmt.Errorf("cannot get relative path: %w", err)
	}
	return "/" + url, nil
}

func rightextpath(path string, dynamiclinks bool) string {
	return replaceext(path, rightext(dynamiclinks))
}

func rightext(dynamiclinks bool) string {
	if dynamiclinks {
		return ""
	}
	return ".html"
}

func replaceext(path, newext string) string {
	ext := filepath.Ext(path)
	return path[:len(path)-len(ext)] + newext
}

func (pg *Page) customurl() (string, bool) {
	return pg.url, len(pg.url) > 0
}

func (pg *Page) GenerateIndex(
	w io.Writer, posts []Post, themedir string,
) error {
	thm, err := theme.ParseTheme(themedir)
	if err != nil {
		return fmt.Errorf("cannot parse theme: %w", err)
	}
	return thm.ExecuteIndex(w, &theme.IndexData{
		Title:   pg.title,
		Content: pg.doc.RenderHtml(commonmark.CMARK_OPT_DEFAULT),
		Posts:   tothemeposts(posts),
	})
}

type Post struct {
	title, category, link string
}

func CreatePost(page *Page, category, link string) *Post {
	return &Post{page.title, category, link}
}

func tothemeposts(posts []Post) []theme.Post {
	themeposts := make([]theme.Post, len(posts))
	for i, p := range posts {
		themeposts[i] = theme.Post{
			Title:    p.title,
			Category: p.category,
			Link:     p.link,
		}
	}
	return themeposts
}

func (pg *Page) GenerateWithoutIndex(w io.Writer, themedir string) error {
	thm, err := theme.ParseTheme(themedir)
	if err != nil {
		return fmt.Errorf("cannot parse theme: %w", err)
	}
	return thm.ExecuteDefault(w, &theme.DefaultData{
		Title:   pg.title,
		Content: pg.doc.RenderHtml(commonmark.CMARK_OPT_DEFAULT),
	})
}

func (pg *Page) Generate(w io.Writer, themedir string, index *Page) error {
	assert.Assert(index != nil)

	thm, err := theme.ParseTheme(themedir)
	if err != nil {
		return fmt.Errorf("cannot parse theme: %w", err)
	}
	return thm.ExecuteDefault(w, &theme.DefaultData{
		Title:     pg.title,
		Content:   pg.doc.RenderHtml(commonmark.CMARK_OPT_DEFAULT),
		SiteTitle: index.title,
	})
}
