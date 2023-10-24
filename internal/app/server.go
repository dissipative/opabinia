package app

import (
	"fmt"
	"net/http"

	"github.com/dissipative/opabinia/internal/api/middleware"

	"github.com/dissipative/opabinia/internal/api/handlers"
	"github.com/dissipative/opabinia/internal/markdown/renderer"
	"github.com/dissipative/opabinia/pkg/cache"

	"github.com/go-chi/chi/v5"
)

func (a *App) DoServe() error {
	mdRenderer, err := renderer.NewDefaultRenderer(a.config, a.logger)
	if err != nil {
		return err
	}

	pagesHandler := handlers.NewPagesHandler(mdRenderer, cache.NewCache(a.config.CacheSize), a.logger)

	router := chi.NewRouter()
	router.Use(middleware.NewLoggerMiddleware(a.logger).DebugLogger)

	router.Route("/assets", func(r chi.Router) {
		r.Use(middleware.FileSystem)
		r.Get("/*", handlers.AssetsHandler)
	})

	router.Get("/index.html", pagesHandler.Pages)
	router.Get("/pages/*", pagesHandler.Pages)
	router.Get("/", handlers.Redirect)

	a.logger.InfoContext(a.ctx, fmt.Sprintf("%s listening on %s", Name, a.config.ServerAddr()))

	return http.ListenAndServe(a.config.ServerAddr(), router)
}
