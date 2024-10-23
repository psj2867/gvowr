package api

import (
	"embed"
	"net/http"
	"strings"

	"github.com/labstack/echo"
)

//go:embed static
var static embed.FS

func StaticServer() http.Handler {
	return http.FileServerFS(static)
}

func StaticEchoServer(prefix string) echo.HandlerFunc {
	h := addPrefix("static", StaticServer())
	h = http.StripPrefix(prefix, h)
	
	return echo.WrapHandler(h)
}

func addPrefix(prefix string, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_prefix := prefix
		if !strings.HasPrefix(r.URL.Path, "/") {
			_prefix = prefix + "/"
		}
		p := _prefix + r.URL.Path
		r.URL.Path = p
		h.ServeHTTP(w, r)
	})
}
