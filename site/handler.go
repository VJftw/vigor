//go:build !js

package site

import (
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

// Handler serves a written site. Paths that aren't files (and have no
// extension) get the page, for the app's client-side routes (reload on
// /items/3 and the app routes it); paths that look like files but aren't
// there are 404s. The program is sent gzipped to browsers that accept it.
func Handler(dir string) http.Handler {
	files := http.FileServer(http.Dir(dir))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p := path.Clean("/" + r.URL.Path)
		switch {
		case p == "/"+WasmFile:
			serveWasm(w, r, dir, files)

			return
		case p == "/" || p == "/"+IndexFile || (path.Ext(p) == "" && !exists(dir, p)):
			// The page, which is also what a route of the app gets.
			serveIndex(w, r, dir)

			return
		}
		files.ServeHTTP(w, r)
	})
}

// serveIndex serves the page itself: not through the file server, which
// redirects requests for index.html to the directory.
func serveIndex(w http.ResponseWriter, r *http.Request, dir string) {
	f, err := os.Open(filepath.Join(dir, IndexFile))
	if err != nil {
		http.NotFound(w, r)

		return
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	http.ServeContent(w, r, IndexFile, info.ModTime(), f)
}

func serveWasm(w http.ResponseWriter, r *http.Request, dir string, files http.Handler) {
	w.Header().Set("Content-Type", "application/wasm")
	w.Header().Set("Vary", "Accept-Encoding")
	if accepts(r.Header.Get("Accept-Encoding"), "gzip") && exists(dir, "/"+WasmGzFile) {
		r2 := r.Clone(r.Context())
		r2.URL.Path = "/" + WasmGzFile
		w.Header().Set("Content-Encoding", "gzip")
		files.ServeHTTP(w, r2)

		return
	}
	files.ServeHTTP(w, r)
}

// accepts reports whether an Accept-Encoding header allows an encoding.
func accepts(header, encoding string) bool {
	for _, part := range strings.Split(header, ",") {
		name, params, _ := strings.Cut(strings.TrimSpace(part), ";")
		if !strings.EqualFold(strings.TrimSpace(name), encoding) {
			continue
		}
		for _, p := range strings.Split(params, ";") {
			if k, v, _ := strings.Cut(strings.TrimSpace(p), "="); strings.EqualFold(strings.TrimSpace(k), "q") {
				q, err := strconv.ParseFloat(strings.TrimSpace(v), 64)

				return err == nil && q > 0 // q=0 is "not acceptable"
			}
		}

		return true
	}

	return false
}
