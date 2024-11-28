package sitefile

import (
	"time"

	"github.com/xr0-org/progstack-ssg/internal/assert"
)

type File interface {
	Path() string
	IsPost() bool

	// post details
	PostTitle() string
	PostHtml() string
	PostPlaintext() string
	PostTime() (time.Time, bool)
}

type file struct {
	path string

	ispost          bool
	title           string
	html, plaintext string
	time            time.Time
}

func PostFile(path, title, html, plaintext string) File {
	return &file{path, true, title, html, plaintext, time.Time{}}
}

func TimedPostFile(path, title, html, plaintext string, time time.Time) File {
	return &file{path, true, title, html, plaintext, time}
}

func NonPostFile(path string) File {
	return &file{path: path}
}

func (f *file) Path() string { return f.path }
func (f *file) IsPost() bool { return f.ispost }

func (f *file) PostTitle() string {
	assert.Assert(f.ispost)
	return f.title
}

func (f *file) PostHtml() string {
	assert.Assert(f.ispost)
	return f.html
}

func (f *file) PostPlaintext() string {
	assert.Assert(f.ispost)
	return f.plaintext
}

func (f *file) PostTime() (time.Time, bool) {
	assert.Assert(f.ispost)
	return f.time, !f.time.IsZero()
}
