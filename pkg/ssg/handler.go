package ssg

import (
	"fmt"
	"time"

	"github.com/xr0-org/progstack-ssg/internal/ast/area"
	"github.com/xr0-org/progstack-ssg/internal/ast/area/sitefile"
)

// A Site is the outcome of generating from a hylo directory.
type Site interface {
	// The Title of the Site.
	Title() string

	// The Files that constitute the Site.
	Bindings() map[string]File
}

type site struct {
	title    string
	bindings map[string]File
}

func (s *site) Title() string             { return s.title }
func (s *site) Bindings() map[string]File { return s.bindings }

func GenerateSiteWithBindings(
	src, target, theme, chromastyle string,
	head, foot string,
	custompages map[string]CustomPage,
) (Site, error) {
	a, err := area.ParseArea(src, chromastyle)
	if err != nil {
		return nil, fmt.Errorf("cannot parse area: %w", err)
	}
	if err := a.Inject(toinjectmap(custompages)); err != nil {
		return nil, fmt.Errorf("injection error: %w", err)
	}
	bindings, err := a.GenerateWithBindings(target, theme, head, foot)
	if err != nil {
		return nil, fmt.Errorf("cannot generate: %w", err)
	}
	return &site{gettitle(a), tofilemap(bindings)}, nil
}

// A File is any URL-accessible resource in a site.
type File interface {
	// Path is the path on disk to the generated File.
	Path() string

	// IsPost indicates whether or not the File is a post.
	IsPost() bool

	// PostTitle is the title of the post.
	PostTitle() string

	// PostHtml is the body of the post without the header/footer portions
	// that render ordinarily.
	PostHtml() string

	// PostTime is the timestamp associated with the post, if available.
	PostTime() (time.Time, bool)
}

func toinjectmap(m1 map[string]CustomPage) map[string]sitefile.CustomPage {
	m2 := map[string]sitefile.CustomPage{}
	for k, v := range m1 {
		m2[k] = v
	}
	return m2
}

func gettitle(a *area.Area) string {
	if s, err := a.Title(); err == nil {
		return s
	}
	return ""
}

func tofilemap(m1 map[string]sitefile.File) map[string]File {
	m2 := map[string]File{}
	for k, v := range m1 {
		m2[k] = v
	}
	return m2
}

// A CustomPage is a page generated without existing in the source directory.
type CustomPage interface {
	Template() string
	Data() map[string]string
}

type custompage struct {
	template string
	data     map[string]string
}

func (pg *custompage) Template() string        { return pg.template }
func (pg *custompage) Data() map[string]string { return pg.data }

func NewSubscriberPage(url string) CustomPage {
	return &custompage{
		"subscribe.html",
		map[string]string{"FormAction": url},
	}
}

func NewMessagePage(title, message string) CustomPage {
	return &custompage{
		"message.html",
		map[string]string{"Title": title, "Message": message},
	}
}
