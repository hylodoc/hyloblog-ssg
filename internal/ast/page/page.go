package page

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/xr0-org/progstack-ssg/internal/assert"
	"github.com/xr0-org/progstack-ssg/internal/theme"
)

type Page struct {
	title, url string
	timing     *timing
	doc        string
	authors    []authordef          // authors of this page
	authordefs map[string]authordef // definitions on this page
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
	mdpage, err := parsemdpage(components.content)
	if err != nil {
		return nil, fmt.Errorf("cannot parse content: %w", err)
	}
	m, err := parsemetadata(components.metadata)
	if err != nil {
		return nil, fmt.Errorf("cannot parse metadata: %w", err)
	}
	return &Page{
		title:      mdpage.title,
		url:        m.URL,
		timing:     m.timing(),
		doc:        mdpage.content,
		authors:    m.definedauthors(),
		authordefs: m.AuthorDefs,
	}, nil
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

func ParsePageGit(path, gitdir string) (*Page, error) {
	pg, err := ParsePage(path)
	if err != nil {
		return nil, err
	}
	info, err := getgitinfo(path, gitdir)
	if err != nil {
		return nil, fmt.Errorf(
			"cannot get timing from git: %w", err,
		)
	}
	if pg.timing == nil {
		pg.timing = &info.timing
	}
	if len(pg.authors) == 0 && info.author != "" {
		pg.authors = []authordef{authordef{Name: info.author}}
	}
	return pg, nil
}

type gitinfo struct {
	timing timing
	author string
}

type timing struct {
	published, updated time.Time
}

func getgitinfo(path, gitdir string) (*gitinfo, error) {
	repo, err := git.PlainOpen(gitdir)
	if err != nil {
		return nil, fmt.Errorf("cannot open repo: %w", err)
	}
	rel, err := filepath.Rel(filepath.Join(gitdir, ".."), path)
	if err != nil {
		return nil, fmt.Errorf("cannot get relative path: %w", err)
	}
	log, err := repo.Log(&git.LogOptions{FileName: &rel})
	if err != nil {
		return nil, fmt.Errorf("cannot get log: %w", err)
	}
	defer log.Close()
	commits := getcommits(log)
	published, err := getcreated(commits)
	if err != nil {
		return nil, fmt.Errorf("cannot get created: %w", err)
	}
	author, err := getauthor(commits)
	if err != nil {
		return nil, fmt.Errorf("cannot get author: %w", err)
	}
	updated, err := getupdated(commits)
	if err != nil {
		return nil, fmt.Errorf("cannot get updated: %w", err)
	}
	return &gitinfo{timing{*published, *updated}, author}, nil
}

func getcommits(iter object.CommitIter) []object.Commit {
	var commits []object.Commit
	err := iter.ForEach(func(c *object.Commit) error {
		commits = append(commits, *c)
		return nil
	})
	assert.Assert(err == nil)
	return commits
}

func getcreated(commits []object.Commit) (*time.Time, error) {
	var t time.Time
	for _, c := range commits {
		when := c.Author.When
		if t.IsZero() || when.Before(t) {
			t = when
		}
	}
	return &t, nil
}

func getauthor(commits []object.Commit) (string, error) {
	if len(commits) == 0 {
		return "", nil
	}
	return commits[0].Author.Name, nil
}

func getupdated(commits []object.Commit) (*time.Time, error) {
	var t time.Time
	for _, c := range commits {
		when := c.Author.When
		if t.IsZero() || when.After(t) {
			t = when
		}
	}
	return &t, nil
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
		Content: pg.doc,
		Posts:   tothemeposts(posts, pg),
	})
}

type Post struct {
	title, category, link string
	timing                *timing
	authors               []authordef
}

func CreatePost(pg *Page, category, link string) *Post {
	return &Post{pg.title, category, link, pg.timing, pg.authors}
}

func tothemeposts(posts []Post, index *Page) []theme.Post {
	sort.Slice(posts, func(i, j int) bool {
		t0, t1 := posts[i].timing, posts[j].timing
		return t0 != nil && t1 != nil &&
			t0.published.After(t1.published)
	})
	themeposts := make([]theme.Post, len(posts))
	for i, p := range posts {
		themeposts[i] = theme.Post{
			Title:    p.title,
			Category: p.category,
			Link:     p.link,
			Date:     getdate(p.timing),
			Authors:  getauthors(p.authors, index.authordefs),
		}
	}
	return themeposts
}

func getauthors(undef []authordef, defs map[string]authordef) []theme.Author {
	var authors []theme.Author
	for _, author := range defineauthors(tostrings(undef), defs) {
		authors = append(authors, theme.Author(author))
	}
	return authors
}

func (pg *Page) GenerateWithoutIndex(w io.Writer, themedir string) error {
	thm, err := theme.ParseTheme(themedir)
	if err != nil {
		return fmt.Errorf("cannot parse theme: %w", err)
	}
	return thm.ExecuteDefault(w, &theme.DefaultData{
		Title:   pg.title,
		Content: pg.doc,
		Date:    getdate(pg.timing),
		Authors: getauthorsnoindex(pg.authors),
	})
}

func getdate(t *timing) string {
	if t == nil || t.published.IsZero() {
		return ""
	}
	return t.published.Format("Jan 02, 2006")
}

func (pg *Page) Generate(w io.Writer, themedir string, index *Page) error {
	assert.Assert(index != nil)

	thm, err := theme.ParseTheme(themedir)
	if err != nil {
		return fmt.Errorf("cannot parse theme: %w", err)
	}
	return thm.ExecuteDefault(w, &theme.DefaultData{
		Title:     pg.title,
		Content:   pg.doc,
		SiteTitle: index.title,
		Date:      getdate(pg.timing),
		Authors:   getauthors(pg.authors, index.authordefs),
	})
}

func getauthorsnoindex(authors []authordef) []theme.Author {
	return getauthors(authors, map[string]authordef{})
}
