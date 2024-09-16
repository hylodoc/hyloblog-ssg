package assert

import "testing"

func TestAssert(t *testing.T) {
	Printf(false, "hello %q\n", "world")
}
