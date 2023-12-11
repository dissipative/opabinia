package renderer

import (
	"bytes"
	"io"

	"github.com/gomarkdown/markdown/ast"
)

// Link converts .md links to .html links and adds starting slash to local pages.
func (r *Renderer) Link(w io.Writer, link *ast.Link, entering bool) {
	link.Destination = processLink(link.Destination)

	// call the default renderer
	r.Renderer.Link(w, link, entering)
}

func processLink(link []byte) []byte {
	// Do nothing with empty and external links
	if len(link) == 0 || bytes.ContainsRune(link, ':') {
		return link
	}

	// Prepending '/' might create double slashes if link already starts with '/'
	if link[0] != '/' {
		link = append([]byte("/"), link...)
	}

	// convert .md links to .html links
	suffixMD := []byte(".md")
	if bytes.HasSuffix(link, suffixMD) {
		link = bytes.TrimSuffix(link, suffixMD)
		link = append(link, []byte(".html")...)
	}

	return link
}
