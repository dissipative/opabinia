package renderer

import (
	"io"

	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
)

// Image renders img tag without trailing slash to match w3c recommendation
func (r *Renderer) Image(w io.Writer, node *ast.Image, entering bool) {
	if entering {
		r.Renderer.Image(w, node, entering)
	} else {
		r.imageExit(w, node)
	}
}

func (r *Renderer) imageExit(w io.Writer, image *ast.Image) {
	r.DisableTags--
	if r.DisableTags == 0 {
		if image.Title != nil {
			r.Outs(w, `" title="`)
			html.EscapeHTML(w, image.Title)
		}
		r.Outs(w, `">`)
	}
}
