package page

import (
	"fmt"
	"strings"

	katex "github.com/FurqanSoftware/goldmark-katex"
	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters/html"
	"github.com/yuin/goldmark"
	hl "github.com/yuin/goldmark-highlighting"
	"github.com/yuin/goldmark/ast"
	gm_ast "github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
	"go.abhg.dev/goldmark/anchor"
)

type mdpage struct {
	title, content string
}

func parsemdpage(content, style string) (*mdpage, error) {
	g := goldmark.New(
		goldmark.WithParserOptions(
			parser.WithAttribute(),
			parser.WithAutoHeadingID(),
		),
		goldmark.WithExtensions(
			extension.NewFootnote(),
			hl.NewHighlighting(
				getstyle(style),
				hl.WithFormatOptions(html.WithLineNumbers(false)),
			),
			&anchor.Extender{Texter: &texter{}},
			&katex.Extender{},
		),
		goldmark.WithExtensions(extension.NewTable()),
	)
	doc := g.Parser().Parse(text.NewReader([]byte(content)))
	s, err := render(g.Renderer(), doc, content)
	if err != nil {
		return nil, fmt.Errorf("cannot render content: %w", err)
	}
	return &mdpage{gettitle(doc, content), s}, nil
}

func getstyle(name string) hl.Option {
	style, ok := map[string]*chroma.Style{
		"based": chroma.MustNewStyle(
			"based",
			chroma.StyleEntries{
				chroma.Background:    "bg:#FFF",
				chroma.Text:          "#000",
				chroma.Comment:       "#777",
				chroma.Keyword:       "#000",
				chroma.Operator:      "#000",
				chroma.Punctuation:   "#000",
				chroma.Name:          "#000",
				chroma.LiteralString: "#000",
			},
		),
	}[name]
	if !ok {
		return hl.WithStyle(name)
	}
	return hl.WithCustomStyle(style)
}

// only render anchors for heading levels above h1
type texter struct{}

func (*texter) AnchorText(h *anchor.HeaderInfo) []byte {
	if h.Level == 1 {
		return nil
	}
	return []byte("§")
}

func render(r renderer.Renderer, doc gm_ast.Node, content string) (string, error) {
	var s strings.Builder
	if err := r.Render(&s, []byte(content), doc); err != nil {
		return "", err
	}
	return s.String(), nil
}

func gettitle(doc gm_ast.Node, content string) string {
	var title string
	ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if h, ok := n.(*ast.Heading); ok && h.Level == 1 {
			if entering {
				title = string(h.Text([]byte(content)))
				return ast.WalkStop, nil
			}
		}
		return ast.WalkContinue, nil
	})
	return title
}
