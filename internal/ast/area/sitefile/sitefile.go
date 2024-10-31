package sitefile

type File interface {
	Path() string
	IsPost() bool
}

type file struct {
	path   string
	ispost bool
}

func NewFile(path string, ispost bool) File {
	return &file{path, ispost}
}

func (f *file) Path() string { return f.path }
func (f *file) IsPost() bool { return f.ispost }
