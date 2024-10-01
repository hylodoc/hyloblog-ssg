package ssg

import (
	"fmt"
	"net/http"
	"os"

	"github.com/xr0-org/progstack-ssg/internal/ast/area"
)

// A Handler is the outcome of compiling a progstack directory into a blog. It
// can be used to serve http requests or examined for information about the
// blog.
type Handler interface {
	http.Handler
	AreaInterface() area.AreaInterface
	Destroy() error
}

type handler struct {
	h      http.Handler
	blog   *area.Area
	target string
}

func NewHandler(src, theme string) (Handler, error) {
	blog, err := area.ParseArea(src)
	if err != nil {
		return nil, fmt.Errorf("cannot parse area: %w", err)
	}
	target, err := os.MkdirTemp("", "")
	if err != nil {
		return nil, fmt.Errorf("cannot make tempdir: %w", err)
	}
	h, err := blog.Handler(target, theme)
	if err != nil {
		return nil, fmt.Errorf("cannot make http handler: %w", err)
	}
	return &handler{h, blog, target}, nil
}

func (h *handler) Destroy() error {
	return os.RemoveAll(h.target)
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.ServeHTTP(w, r)
}

func (h *handler) AreaInterface() area.AreaInterface {
	return h.blog.Interface()
}
