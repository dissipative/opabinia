package middleware

import (
	"net/http"
	"strings"
)

func FileSystem(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			http.Error(w, "Directory listings are forbidden", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
