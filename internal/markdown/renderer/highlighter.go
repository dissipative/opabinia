package renderer

import (
	"errors"
	"fmt"
	"io"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/gomarkdown/markdown/ast"
	baseRenderer "github.com/gomarkdown/markdown/html"
)

type highlighter struct {
	*html.Formatter
	*chroma.Style
	isUsed bool // whether highlighter is used on page
}

func newHighlighter(syntaxTheme string) (*highlighter, error) {
	htmlFormatter := html.New(html.WithClasses(true), html.TabWidth(2), html.WithLineNumbers(true))
	if htmlFormatter == nil {
		return nil, errors.New("couldn't create html formatter for code block")
	}
	highlightStyle := styles.Get(syntaxTheme)
	if highlightStyle == nil {
		return nil, errors.New("style for code block not found")
	}

	return &highlighter{
		Formatter: htmlFormatter,
		Style:     highlightStyle,
	}, nil
}

// based on https://github.com/alecthomas/chroma/blob/master/quick/quick.go
func (h *highlighter) highlight(w io.Writer, source, lang, defaultLang string) error {
	if lang == "" {
		lang = defaultLang
	}
	l := lexers.Get(lang)
	if l == nil {
		l = lexers.Analyse(source)
	}
	if l == nil {
		l = lexers.Fallback
	}
	l = chroma.Coalesce(l)

	it, err := l.Tokenise(nil, source)
	if err != nil {
		return err
	}
	return h.Format(w, h.Style, it)
}

func (h *highlighter) setIsUsed() {
	h.isUsed = true
}

func (h *highlighter) endPage() {
	h.isUsed = false
}

func (h *highlighter) renderCSS(w io.Writer) error {
	_, err := fmt.Fprint(w, "<style>\n")
	if err != nil {
		return fmt.Errorf("writing highlight css error: %w", err)
	}
	err = h.WriteCSS(w, h.Style)
	if err != nil {
		return fmt.Errorf("writing highlight css error: %w", err)
	}
	_, err = fmt.Fprint(w, "</style>\n")
	if err != nil {
		return fmt.Errorf("writing highlight css error: %w", err)
	}

	return nil
}

// Code renders code tag with syntax highlighting.
func (r *Renderer) Code(w io.Writer, code *ast.Code) {
	defer r.h.setIsUsed()

	r.Outs(w, "<code class=\"chroma op-code\"><span class=\"cl\">")
	baseRenderer.EscapeHTML(w, code.Literal)
	r.Outs(w, "</span></code>")
}

// CodeBlock renders code block with syntax highlighting.
func (r *Renderer) CodeBlock(w io.Writer, codeBlock *ast.CodeBlock) {
	defer r.h.setIsUsed()

	err := r.h.highlight(w, string(codeBlock.Literal), string(codeBlock.Info), r.opts.SyntaxDefaultLang)
	if err != nil {
		r.logger.Error("highlight error", err)
	}
}
