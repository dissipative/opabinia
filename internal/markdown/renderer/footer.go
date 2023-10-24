package renderer

import (
	"io"

	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
)

// RenderFooter writes HTML document footer.
func (r *Renderer) RenderFooter(w io.Writer, _ ast.Node) {
	if r.state.documentMatter != ast.DocumentMatterNone {
		r.Outs(w, "</section>\n")
	}

	if r.h != nil && r.h.isUsed {
		err := r.h.renderCSS(w)
		if err != nil {
			r.logger.Warn("CSS render error", err)
		}
	}

	if r.Opts.Flags&html.CompletePage == 0 {
		return
	}
	io.WriteString(w, "\n</body>\n</html>\n")
}

// DocumentMatter just copies parent method using child's `state` property
func (r *Renderer) DocumentMatter(w io.Writer, node *ast.DocumentMatter, entering bool) {
	if !entering {
		return
	}
	if r.state.documentMatter != ast.DocumentMatterNone {
		r.Outs(w, "</section>\n")
	}
	switch node.Matter {
	case ast.DocumentMatterFront:
		r.Outs(w, `<section data-matter="front">`)
	case ast.DocumentMatterMain:
		r.Outs(w, `<section data-matter="main">`)
	case ast.DocumentMatterBack:
		r.Outs(w, `<section data-matter="back">`)
	}
	r.state.documentMatter = node.Matter
}
