package ignorer

import (
	"fmt"
	"os"
	"strings"

	"github.com/xr0-org/progstack-ssg/internal/assert"
)

type Ignorer struct {
	m map[string]bool
}

func Empty() *Ignorer { return &Ignorer{map[string]bool{}} }

func (ign *Ignorer) Augment(path string) (*Ignorer, error) {
	newign, err := fromfile(path)
	if err != nil {
		return nil, fmt.Errorf("cannot get from file: %w", err)
	}
	for k, old := range ign.m {
		if _, ok := newign.m[k]; !ok {
			newign.m[k] = old
		}
	}
	return newign, nil
}

func fromfile(path string) (*Ignorer, error) {
	lines, err := getignorelines(path)
	if err != nil {
		return nil, fmt.Errorf("cannot get ignore lines: %w", err)
	}
	ignore := map[string]bool{}
	for _, line := range lines {
		inst := parseinstruction(line)
		ignore[inst.name] = inst.shouldignore
	}
	return &Ignorer{ignore}, nil
}

func getignorelines(path string) ([]string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, fmt.Errorf("cannot read file: %w", err)
	}
	trimmed := strings.TrimSpace(string(b))
	if len(trimmed) == 0 {
		return []string{}, nil
	}
	return strings.Split(trimmed, "\n"), nil
}

type instruction struct {
	name         string
	shouldignore bool
}

func parseinstruction(s string) *instruction {
	assert.Assert(len(s) > 0)
	if s[0] == '!' {
		return &instruction{s[1:], false}
	}
	return &instruction{s, true}
}

func (ign *Ignorer) ShouldIgnore(name string) bool {
	shouldignore, ok := ign.m[name]
	return ok && shouldignore
}
