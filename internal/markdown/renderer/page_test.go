package renderer

import (
	"errors"
	"io"
	"log/slog"
	"testing"

	"github.com/dissipative/opabinia/internal/infra/config"
)

func TestRenderer_RenderPage(t *testing.T) {
	tests := []struct {
		name     string
		renderer *Renderer
		mdFile   string
		wantErr  bool
		err      error
	}{
		{
			name:     "Valid markdown file and template",
			renderer: makeTestRenderer("test_data/valid.tmpl"),
			mdFile:   "test_data/valid.md",
			wantErr:  false,
		},
		{
			name:     "Non-existent markdown file",
			renderer: makeTestRenderer("test_data/valid.tmpl"),
			mdFile:   "non_existent.md",
			wantErr:  true,
			err:      ErrReadFile,
		},
		{
			name:     "Invalid template",
			renderer: makeTestRenderer("test_data/invalid.tmpl"),
			mdFile:   "test_data/valid.md",
			wantErr:  true,
			err:      ErrTpl,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.renderer.RenderPage(tt.mdFile)
			if (err != nil) != tt.wantErr && !errors.Is(err, tt.err) {
				t.Errorf("RenderPage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func makeTestRenderer(template string) *Renderer {
	c := config.NewConfig()
	c.Template.File = template
	c.Markdown.SyntaxTheme = ""

	l := slog.New(slog.NewJSONHandler(io.Discard, nil))
	r, _ := NewDefaultRenderer(c, l)

	return r
}
