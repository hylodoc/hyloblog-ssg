package ssg

import (
	"fmt"
	"os"
	"testing"
)

func TestHandler(t *testing.T) {
	if err := testHandler(); err != nil {
		t.Fatal(err)
	}
}

func testHandler() error {
	target, err := os.MkdirTemp("", "")
	if err != nil {
		return fmt.Errorf("cannot make tempdir: %w", err)
	}
	bindings, err := GenerateSiteWithBindings(
		"test", target, "../../theme/lit", "algol_nu",
	)
	if err != nil {
		return fmt.Errorf("cannot generate: %w", err)
	}
	for _, url := range []string{
		"/",
		"/abc/def",
		"/nest/post",
		"/nest-no-ignore/README",
		"/nest-no-ignore/post",
	} {
		if file, ok := bindings[url]; !ok {
			return fmt.Errorf("%q not found", url)
		} else {
			fmt.Println(url, file)
		}
	}
	return nil
}
