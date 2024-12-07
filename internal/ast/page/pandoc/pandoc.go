package pandoc

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
)

func ConvertPlaintext(markdown string, stdout io.Writer) error {
	cmd := exec.Command(
		"pandoc",
		"-t", "plain",
		"--columns=72",
	)
	cmd.Stdin = bytes.NewBufferString(markdown)
	cmd.Stdout = stdout
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("exec error: %w", err)
	}
	if stderr.Len() > 0 {
		return fmt.Errorf("stderr: %s", stderr.String())
	}
	return nil
}
