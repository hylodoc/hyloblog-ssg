package ssg

import (
	"fmt"
	"net/http"

	"github.com/xr0-org/progstack-ssg/internal/ast/area"
)

// A Handler is the outcome of compiling a progstack directory into a blog. It
// can be used to serve http requests or examined for information about the
// blog.
type Handler interface {
	http.Handler
	AreaInterface() (area.AreaInterface, error)
	Destroy() error
}

type handler struct {
	h    *area.Handler
	blog *area.Area
}

func NewHandler(src, theme string) (Handler, error) {
	blog, err := area.ParseArea(src)
	if err != nil {
		return nil, fmt.Errorf("cannot parse area: %w", err)
	}
	h, err := blog.Handler(theme)
	if err != nil {
		return nil, fmt.Errorf("cannot make http handler: %w", err)
	}
	return &handler{h, blog}, nil
}

func (h *handler) Destroy() error {
	return h.h.Destroy()
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.ServeHTTP(w, r)
}

func (h *handler) AreaInterface() (area.AreaInterface, error) {
	return h.blog.Interface()
}
