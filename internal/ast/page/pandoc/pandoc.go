package pandoc

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

func ConvertPlaintext(markdown string) (string, error) {
	cmd := exec.Command(
		"pandoc",
		"-t", "plain",
		"--columns", "72",
	)
	cmd.Stdin = bytes.NewBufferString(markdown)
	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("exec error: %w", err)
	}
	if stderr.Len() > 0 {
		return "", fmt.Errorf("stderr: %s", stderr.String())
	}
	return stdout.String(), nil
}
