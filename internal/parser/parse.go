package parse

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/rhinoman/go-commonmark"
	"github.com/xr0-org/progstack-ssg/internal/assert"
)

const markdownExt = ".md"

func Generate(src, target string) error {
	title, err := gettitle(filepath.Join(src, "index.md"))
	if err != nil {
		return fmt.Errorf("cannot get title: %w", err)
	}
	gen := generation{src, target, title}
	log.Printf("title %q\n", title)
	return filepath.WalkDir(
		src,
		func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return fmt.Errorf("path error: %w", err)
			}
			isfile, err := checkisfile(path)
			if err != nil {
				return fmt.Errorf(
					"cannot check if file: %w", err,
				)
			}
			if isfile && filepath.Ext(path) == markdownExt {
				return gen.render(path)
			}
			return nil
		},
	)
}

func gettitle(index string) (string, error) {
	doc, err := commonmark.ParseFile(index, commonmark.CMARK_OPT_DEFAULT)
	if err != nil {
		return "", fmt.Errorf("cannot parse: %w", err)
	}
	iter := commonmark.NewCMarkIter(doc)
	defer iter.Free()
	for iter.Next() != commonmark.CMARK_EVENT_DONE {
		n := iter.GetNode()
		if n.GetNodeType() == commonmark.CMARK_NODE_HEADING &&
			n.GetHeaderLevel() == 1 {
			title := strings.TrimSpace(n.FirstChild().GetLiteral())
			if title == "" {
				return "", fmt.Errorf("empty")
			}
			return title, nil
		}
	}
	return "", fmt.Errorf("no heading node")
}

func checkisfile(path string) (bool, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return false, fmt.Errorf("cannot stat: %w", err)
	}
	return fi.Mode().IsRegular(), nil
}

// generation represents the operation of generating the target blog from the
// source directory
type generation struct {
	src, target string
	title       string
}

func (gen *generation) String() string {
	return fmt.Sprintf(
		"{%s -> %s}", filepath.Clean(gen.src), filepath.Clean(gen.target),
	)
}

func (gen *generation) render(srcfile string) error {
	b, err := os.ReadFile(srcfile)
	if err != nil {
		return fmt.Errorf("cannot read file: %w", err)
	}
	outputpath := gen.changetotarget(replaceext(srcfile, ".html"))
	if err := os.MkdirAll(filepath.Dir(outputpath), 0777); err != nil {
		return fmt.Errorf("cannot make dir: %w", err)
	}
	return os.WriteFile(
		outputpath,
		[]byte(commonmark.Md2Html(string(b), commonmark.CMARK_OPT_DEFAULT)),
		0666,
	)
}

func (gen *generation) changetotarget(path string) string {
	rel, err := filepath.Rel(gen.src, path)
	assert.Assert(err == nil)
	return filepath.Join(gen.target, rel)
}

func replaceext(path, newext string) string {
	ext := filepath.Ext(path)
	return path[:len(path)-len(ext)] + newext
}
