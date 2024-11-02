package sitefile

import (
	"time"
)

type File interface {
	Path() string
	IsPost() bool
	Time() (time.Time, bool)
}

type file struct {
	path   string
	ispost bool
	time   time.Time
}

func PostFile(path string) File {
	return &file{path: path, ispost: true}
}

func TimedPostFile(path string, time time.Time) File {
	return &file{path, true, time}
}

func NonPostFile(path string) File {
	return &file{path: path}
}

func (f *file) Path() string { return f.path }
func (f *file) IsPost() bool { return f.ispost }

func (f *file) Time() (time.Time, bool) {
	return f.time, !f.time.IsZero()
}
