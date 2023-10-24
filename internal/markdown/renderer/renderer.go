package renderer

import (
	"io"
	"sync"

	"github.com/dissipative/opabinia/internal/infra/config"

	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
)

type Logger interface {
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

type Config interface {
	CodeDefaultLang() string
	SyntaxTheme() string
	TmplConf() *config.Template
}

type Renderer struct {
	*html.Renderer

	h *highlighter
	m *sync.Mutex

	state  *renderState
	opts   Options
	logger Logger
}

type Options struct {
	html.RendererOptions

	SyntaxDefaultLang string
	SyntaxTheme       string
	TmplConf          *config.Template
}

type renderState struct {
	documentMatter ast.DocumentMatters
}

func NewRenderer(opts Options, logger Logger) (*Renderer, error) {
	h, err := newHighlighter(opts.SyntaxTheme)
	if err != nil {
		return nil, err
	}

	return &Renderer{
		Renderer: html.NewRenderer(opts.RendererOptions),
		h:        h,
		state:    &renderState{},
		opts:     opts,
		logger:   logger,
		m:        new(sync.Mutex),
	}, nil
}

func NewDefaultRenderer(conf Config, logger Logger) (*Renderer, error) {
	return NewRenderer(
		Options{
			RendererOptions:   html.RendererOptions{Flags: html.CommonFlags | html.HrefTargetBlank},
			SyntaxDefaultLang: conf.CodeDefaultLang(),
			SyntaxTheme:       conf.SyntaxTheme(),
			TmplConf:          conf.TmplConf(),
		},
		logger,
	)
}

// RenderNode process images, links, and code blocks differently than default renderer
func (r *Renderer) RenderNode(w io.Writer, node ast.Node, entering bool) ast.WalkStatus {
	switch n := node.(type) {
	case *ast.Link:
		r.Link(w, n, entering)
		return ast.GoToNext
	case *ast.Image:
		if r.opts.Flags&html.SkipImages != 0 {
			return ast.SkipChildren
		}
		r.Image(w, n, entering)
		return ast.GoToNext
	case *ast.CodeBlock:
		r.CodeBlock(w, n)
		return ast.GoToNext
	case *ast.Code:
		r.Code(w, n)
		return ast.GoToNext
	}
	return r.Renderer.RenderNode(w, node, entering)
}
