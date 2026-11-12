//go:build !js

package site_test

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/VJftw/vigor/site"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// fakeWasm is enough of a program for the site's files.
var fakeWasm = []byte("\x00asm\x01\x00\x00\x00 pretend")

func written(t *testing.T, opts site.Options) (dir, version string) {
	t.Helper()
	wasm := filepath.Join(t.TempDir(), "in.wasm")
	require.NoError(t, os.WriteFile(wasm, fakeWasm, 0o600))
	dir = t.TempDir()
	version, err := site.Write(dir, wasm, opts)
	require.NoError(t, err)

	return dir, version
}

func TestWrite(t *testing.T) {
	dir, version := written(t, site.Options{
		Title:       "My app",
		Description: "Does things.",
		Lang:        "fr",
		Styles:      []string{"/app.css"},
		Icon:        "/icon.svg",
		Loading:     "Loading…",
		Head:        `<meta name="theme-color" content="#000">`,
	})
	assert.Len(t, version, 16)

	for _, f := range []string{site.IndexFile, site.WasmExecFile, site.WasmFile, site.WasmGzFile} {
		info, err := os.Stat(filepath.Join(dir, f))
		require.NoError(t, err, f)
		assert.NotZero(t, info.Size(), f)
	}
	index, err := os.ReadFile(filepath.Join(dir, site.IndexFile))
	require.NoError(t, err)
	for _, want := range []string{
		`<html lang="fr">`, `<title>My app</title>`, `content="Does things."`, `href="/app.css"`, `href="/icon.svg"`,
		`<meta name="theme-color" content="#000">`, `<div id="app">Loading…</div>`,
		`src="/wasm_exec.js?v=` + version, `fetch("/app.wasm?v=` + version,
	} {
		assert.Contains(t, string(index), want)
	}

	// The loader is the installed Go release's.
	loader, err := site.WasmExec()
	require.NoError(t, err)
	got, _ := os.ReadFile(filepath.Join(dir, site.WasmExecFile))
	assert.Equal(t, loader, got)
}

func TestWriteDefaultsAndEscaping(t *testing.T) {
	dir, _ := written(t, site.Options{Title: `<script>alert(1)</script>`})
	index, _ := os.ReadFile(filepath.Join(dir, site.IndexFile))
	assert.NotContains(t, string(index), "<script>alert(1)")
	assert.Contains(t, string(index), `<html lang="en">`)

	dir, _ = written(t, site.Options{})
	index, _ = os.ReadFile(filepath.Join(dir, site.IndexFile))
	assert.Contains(t, string(index), "<title>vigor app</title>")
}

func TestVersionFollowsTheProgram(t *testing.T) {
	_, a := written(t, site.Options{})
	_, b := written(t, site.Options{})
	assert.Equal(t, a, b, "same program, same version")

	other := filepath.Join(t.TempDir(), "other.wasm")
	require.NoError(t, os.WriteFile(other, []byte("\x00asm different"), 0o600))
	c, err := site.Write(t.TempDir(), other, site.Options{})
	require.NoError(t, err)
	assert.NotEqual(t, a, c)
}

func get(h http.Handler, path string, header ...string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodGet, path, nil)
	for i := 0; i+1 < len(header); i += 2 {
		req.Header.Set(header[i], header[i+1])
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	return rec
}

func TestHandler(t *testing.T) {
	dir, _ := written(t, site.Options{Title: "My app"})
	h := site.Handler(dir)

	// The page at the root, and for the app's routes (also deep, and with a
	// query).
	for _, path := range []string{"/", "/items/3", "/a/b/c?x=1", "/index.html"} {
		rec := get(h, path)
		assert.Equal(t, http.StatusOK, rec.Code, path)
		assert.Contains(t, rec.Body.String(), "<title>My app</title>", path)
	}

	// Files are themselves; missing files aren't the page.
	rec := get(h, "/wasm_exec.js")
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "Go")
	assert.Equal(t, http.StatusNotFound, get(h, "/missing.js").Code)
	assert.Equal(t, http.StatusNotFound, get(h, "/img/logo.png").Code)
}

func TestHandlerWasm(t *testing.T) {
	dir, _ := written(t, site.Options{})
	h := site.Handler(dir)

	plain := get(h, "/app.wasm")
	assert.Equal(t, http.StatusOK, plain.Code)
	assert.Equal(t, "application/wasm", plain.Header().Get("Content-Type"))
	assert.Empty(t, plain.Header().Get("Content-Encoding"))
	assert.Equal(t, fakeWasm, plain.Body.Bytes())

	for accept, wantGzip := range map[string]bool{
		"gzip, deflate, br": true,
		"gzip;q=0.5":        true,
		"gzip;q=0":          false,
		"deflate":           false,
		"":                  false,
	} {
		rec := get(h, "/app.wasm", "Accept-Encoding", accept)
		assert.Equal(t, "application/wasm", rec.Header().Get("Content-Type"), accept)
		assert.Contains(t, rec.Header().Get("Vary"), "Accept-Encoding")
		if !wantGzip {
			assert.Empty(t, rec.Header().Get("Content-Encoding"), accept)

			continue
		}
		assert.Equal(t, "gzip", rec.Header().Get("Content-Encoding"), accept)
		zr, err := gzip.NewReader(bytes.NewReader(rec.Body.Bytes()))
		require.NoError(t, err, accept)
		body, _ := io.ReadAll(zr)
		assert.Equal(t, fakeWasm, body, accept)
	}
}

func TestBuild(t *testing.T) {
	if testing.Short() {
		t.Skip("builds a program")
	}
	out := filepath.Join(t.TempDir(), "app.wasm")
	require.NoError(t, site.Build(t.Context(), "testdata/hello", out))
	b, err := os.ReadFile(out)
	require.NoError(t, err)
	assert.True(t, bytes.HasPrefix(b, []byte("\x00asm")), "a WebAssembly program")

	// A program that doesn't build says why.
	err = site.Build(t.Context(), t.TempDir(), filepath.Join(t.TempDir(), "x.wasm"))
	require.Error(t, err)
	assert.True(t, strings.Contains(err.Error(), "go build"), err.Error())
}
