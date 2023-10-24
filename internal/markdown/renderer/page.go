package renderer

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"os"

	"github.com/dissipative/opabinia/internal/infra/config"
	parser2 "github.com/dissipative/opabinia/internal/markdown/parser"

	"github.com/gomarkdown/markdown"
	mdParser "github.com/gomarkdown/markdown/parser"
)

var (
	ErrReadFile = errors.New("file reading error")
	ErrTpl      = errors.New("template error")
)

type Page struct {
	Title    string
	Filename string
	Body     template.HTML
	*config.Template
}

func (r *Renderer) RenderPage(mdFilename string) ([]byte, error) {
	r.m.Lock()
	defer r.m.Unlock()

	// reset highlighter after the page rendering
	defer r.h.endPage()

	var rendered bytes.Buffer

	fileData, err := os.ReadFile(mdFilename)
	if err != nil {
		return nil, fmt.Errorf("%s: %w: %w", mdFilename, ErrReadFile, err)
	}

	extensions := mdParser.CommonExtensions | mdParser.Footnotes | mdParser.SuperSubscript | mdParser.MathJax
	parser := mdParser.NewWithExtensions(extensions)

	htmlData := markdown.ToHTML(fileData, parser, r)

	t, err := template.ParseFiles(r.opts.TmplConf.File)
	if err != nil {
		return nil, r.wrapTemplateErr(err)
	}

	if err = t.Execute(&rendered, &Page{
		Title:    parser2.ExtractTitle(fileData),
		Filename: mdFilename,
		Body:     template.HTML(htmlData),
		Template: r.opts.TmplConf,
	}); err != nil {
		return nil, r.wrapTemplateErr(err)
	}

	return rendered.Bytes(), nil
}

func (r *Renderer) wrapTemplateErr(err error) error {
	return fmt.Errorf("%s: %w: %w", r.opts.TmplConf.File, ErrTpl, err)
}
