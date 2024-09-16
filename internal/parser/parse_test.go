package parse

import (
	"testing"
)

func TestParse(t *testing.T) {
	const (
		src    = "../../test/src"
		target = "../../test/target"
	)
	if err := Parse(src, target); err != nil {
		t.Fatal(err)
	}
}
