package generator

import (
	"os/exec"
	"testing"
)

func TestParse(t *testing.T) {
	const (
		src      = "test/src"
		target   = "test/target"
		expected = "test/expected"
	)
	if err := Generate(src, target, ""); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command("diff", "-rq", target, expected)
	if b, err := cmd.Output(); err != nil {
		t.Error(string(b))
		t.Fatal(err)
	}
}
