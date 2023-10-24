package renderer

import (
	"bytes"
	"io"

	"github.com/gomarkdown/markdown/ast"
)

// Link converts .md links to .html links and adds starting slash to local pages.
func (r *Renderer) Link(w io.Writer, link *ast.Link, entering bool) {
	suffixMD := []byte(".md")
	suffixRemote := []byte(":")

	// convert .md links to .html links
	if bytes.HasSuffix(link.Destination, suffixMD) {
		link.Destination = bytes.ReplaceAll(link.Destination, suffixMD, []byte(".html"))
	}

	// add starting slash to local pages
	if !bytes.HasSuffix(link.Destination, suffixRemote) {
		link.Destination = append([]byte("/"), link.Destination...)
	}

	// call the default renderer
	r.Renderer.Link(w, link, entering)
}
