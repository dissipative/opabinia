package handlers

import (
	"mime"
	"net/http"
)

func init() {
	mime.AddExtensionType(".webmanifest", "application/manifest+json")
}

func AssetsHandler(w http.ResponseWriter, r *http.Request) {
	fileServer := http.FileServer(http.Dir("./assets"))
	http.StripPrefix(
		"/assets",
		fileServer,
	).ServeHTTP(w, r)
}
