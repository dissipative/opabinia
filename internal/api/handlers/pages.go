package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/dissipative/opabinia/internal/api/response"
	"github.com/dissipative/opabinia/internal/markdown/renderer"
)

type Renderer interface {
	RenderPage(mdFilename string) ([]byte, error)
}

type Cache interface {
	Get(key any) (any, error)
	Set(key, val any)
}

type ErrLogger interface {
	Error(msg string, args ...any)
}

type PagesHandler struct {
	renderer Renderer
	cache    Cache
	logger   ErrLogger
}

func NewPagesHandler(renderer Renderer, cache Cache, logger ErrLogger) *PagesHandler {
	return &PagesHandler{renderer: renderer, cache: cache, logger: logger}
}

func parseHTMLFilename(filename string) (string, error) {
	if strings.Contains(filename, ".html") {
		return strings.Replace(filename, ".html", ".md", 1), nil
	}
	return "", fmt.Errorf("%s is not an html file", filename)
}

func (p *PagesHandler) Test() http.HandlerFunc {
	return p.Pages
}

func (p *PagesHandler) Pages(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Path[len("/"):]
	cacheKey, err := parseHTMLFilename(filename)
	if err != nil {
		response.WriteInternalError(w, fmt.Errorf("error parsing URL path: %w", err))
		return
	}

	// get from cache
	cached, err := p.cache.Get(cacheKey)
	if err == nil {
		if _, err = w.Write(cached.([]byte)); err != nil {
			response.WriteInternalError(w, fmt.Errorf("error on writing cached template: %w", err))
		}
		return
	}

	// render page and then set cache
	content, err := p.renderer.RenderPage(cacheKey)
	if err != nil {
		if errors.Is(renderer.ErrReadFile, err) {
			http.Error(w, "Page reading error", http.StatusNotFound)
			p.logger.Error("template reading error", err)
			return
		}

		response.WriteInternalError(w, err)
		return
	}

	p.cache.Set(cacheKey, content)

	if _, err = w.Write(content); err != nil {
		response.WriteInternalError(w, err)
	}
}
