package ssg

import (
	"fmt"
	"testing"

	"github.com/xr0-org/progstack-ssg/internal/ast/area"
)

func TestHandler(t *testing.T) {
	if err := testHandler(); err != nil {
		t.Fatal(err)
	}
}

func testHandler() error {
	h, err := NewHandler("test", "../../theme/lit")
	if err != nil {
		return fmt.Errorf("cannot make handler: %w", err)
	}
	expected := []string{
		"/",
		"/post",
		"/nest/post",
		"/nest-no-ignore/README",
		"/nest-no-ignore/post",
	}
	return confirmsetequal(expected, getallpages(h.AreaInterface()))
}

func getallpages(ifc area.AreaInterface) []string {
	var pages []string
	for _, pg := range ifc.Pages() {
		pages = append(pages, pg.Link())
	}
	for _, a := range ifc.Subareas() {
		pages = append(pages, getallpages(a)...)
	}
	return pages
}

func confirmsetequal(arr, arr0 []string) error {
	if len(arr) != len(arr0) {
		return fmt.Errorf("length mismatch")
	}
	for _, s := range arr {
		if !contains(arr0, s) {
			return fmt.Errorf("%q not found", s)
		}
	}
	return nil
}
func contains(arr []string, item string) bool {
	for _, s := range arr {
		if s == item {
			return true
		}
	}
	return false
}
