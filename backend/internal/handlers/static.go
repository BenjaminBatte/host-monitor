package handlers

import (
	"embed"
	"io/fs"
	"net/http"
	"path"
	"strings"
)

//go:embed dist/browser/*
var uiFS embed.FS

func SPA() (http.Handler, error) {
	sub, err := fs.Sub(uiFS, "dist/browser") // sub is type fs.FS ✅
	if err != nil {
		return nil, err
	}

	fileServer := http.FileServer(http.FS(sub)) // FileServer needs http.FileSystem ✅

	serveIndex := func(w http.ResponseWriter, r *http.Request) {
		// Serve SPA entry directly from fs.FS
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		http.ServeFileFS(w, r, sub, "index.html") // <-- pass sub (fs.FS), not http.FS(sub)
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p := strings.TrimPrefix(path.Clean(r.URL.Path), "/")
		if p == "" || p == "index.html" {
			serveIndex(w, r)
			return
		}
		// If asset missing, SPA fallback
		if _, err := fs.Stat(sub, p); err != nil {
			serveIndex(w, r)
			return
		}
		fileServer.ServeHTTP(w, r)
	}), nil
}
