//go:build !js

// Package site is the page that runs a vigor app in a browser, and a server
// for it: everything around the WebAssembly program that isn't the app.
//
// A vigor app is a Go program built with GOOS=js GOARCH=wasm. To run it a
// browser needs a page with an element to mount into (the "app" element),
// Go's wasm_exec.js loader, and the program itself. Write puts those in a
// directory, Handler serves it (with the fallback to the page that
// client-side routes need), and Build compiles the program.
//
// Every part is optional: an app with its own page, service worker or
// hosting can ignore this package, and use only what it likes of it.
package site

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// File names in a written site.
const (
	IndexFile    = "index.html"
	WasmExecFile = "wasm_exec.js"
	WasmFile     = "app.wasm"
	WasmGzFile   = "app.wasm.gz"
)

// Options say what the page is like. All are optional.
type Options struct {
	// Title is the page's title; "vigor app" by default.
	Title string
	// Description is the page's description.
	Description string
	// Lang is the page's language; "en" by default.
	Lang string
	// Styles are stylesheets' URLs, loaded before the app starts.
	Styles []string
	// Icon is the URL of the page's icon (e.g. "/icon.svg").
	Icon string
	// Head is more HTML for the page's <head>: a manifest, meta tags.
	Head template.HTML
	// Loading is HTML shown in the app element until the program is running,
	// and cleared when it starts; "" for nothing. It is trusted.
	Loading template.HTML
	// ID is the id of the element the app mounts into; "app" by default (it
	// is what web.Mount's id argument names).
	ID string
	// WasmURL and WasmExecURL are where the program and Go's loader are
	// served from, if not at "/app.wasm" and "/wasm_exec.js" (where Write
	// puts them).
	WasmURL, WasmExecURL string
}

var page = template.Must(template.New("index").Parse(`<!doctype html>
<html lang="{{.Lang}}">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>{{.Title}}</title>
{{if .Description}}<meta name="description" content="{{.Description}}">
{{end}}{{if .Icon}}<link rel="icon" href="{{.Icon}}">
{{end}}{{range .Styles}}<link rel="stylesheet" href="{{.}}">
{{end}}{{.Head}}
<script src="{{.WasmExecURL}}?v={{.Version}}"></script>
</head>
<body>
<div id="{{.ID}}">{{.Loading}}</div>
<script>
const go = new Go();
const root = document.getElementById({{.JSID}});
WebAssembly.instantiateStreaming(fetch({{.JSWasm}}), go.importObject)
  .then((result) => { root.textContent = ""; return go.run(result.instance); })
  .catch((err) => { root.textContent = {{.JSTitle}} + " couldn't start: " + err; });
</script>
</body>
</html>
`))

// Index is the page, for an app whose program has the given content hash.
func Index(opts Options, version string) ([]byte, error) {
	if opts.Title == "" {
		opts.Title = "vigor app"
	}
	if opts.Lang == "" {
		opts.Lang = "en"
	}
	if opts.ID == "" {
		opts.ID = "app"
	}
	if opts.WasmURL == "" {
		opts.WasmURL = "/" + WasmFile
	}
	if opts.WasmExecURL == "" {
		opts.WasmExecURL = "/" + WasmExecFile
	}
	// The script's strings: JSON, which is JavaScript, with <, > and & escaped
	// so that nothing in them can end the script.
	js := func(v string) template.JS {
		b, _ := json.Marshal(v)

		return template.JS(b)
	}
	var b bytes.Buffer
	err := page.Execute(&b, struct {
		Options
		Version               string
		JSID, JSWasm, JSTitle template.JS
	}{opts, version, js(opts.ID), js(opts.WasmURL + "?v=" + version), js(opts.Title)})

	return b.Bytes(), err
}

// WasmExec is Go's loader for WebAssembly programs, from the Go installation
// that builds the program (they have to match).
func WasmExec() ([]byte, error) {
	for _, dir := range []string{"lib", "misc"} { // moved in Go 1.24
		if b, err := os.ReadFile(filepath.Join(runtime.GOROOT(), dir, "wasm", WasmExecFile)); err == nil {
			return b, nil
		}
	}

	return nil, fmt.Errorf("site: no %s in %s", WasmExecFile, runtime.GOROOT())
}

// Build compiles the Go package in dir (default ".") to WebAssembly at out.
func Build(ctx context.Context, dir, out string) error {
	if dir == "" {
		dir = "."
	}
	cmd := exec.CommandContext(ctx, "go", "build", "-o", out, ".")
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "GOOS=js", "GOARCH=wasm")
	if b, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("go build: %w\n%s", err, b)
	}

	return nil
}

// Write writes the site into dir: the page, Go's loader, the program, and a
// gzipped copy of it for Handler to send. It returns the program's content
// hash, which names its version.
func Write(dir, wasmPath string, opts Options) (version string, err error) {
	wasm, err := os.ReadFile(wasmPath)
	if err != nil {
		return "", err
	}
	loader, err := WasmExec()
	if err != nil {
		return "", err
	}
	// Versioned by content, so builds are reproducible and browsers pick up
	// new ones.
	version = fmt.Sprintf("%x", sha256.Sum256(wasm))[:16]
	index, err := Index(opts, version)
	if err != nil {
		return "", err
	}
	var gz bytes.Buffer
	zw, _ := gzip.NewWriterLevel(&gz, gzip.BestCompression)
	zw.Write(wasm)
	if err := zw.Close(); err != nil {
		return "", err
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	for name, content := range map[string][]byte{
		IndexFile: index, WasmExecFile: loader, WasmFile: wasm, WasmGzFile: gz.Bytes(),
	} {
		if err := os.WriteFile(filepath.Join(dir, name), content, 0o644); err != nil {
			return "", err
		}
	}

	return version, nil
}

func exists(dir, name string) bool {
	_, err := os.Stat(filepath.Join(dir, filepath.FromSlash(name)))

	return err == nil
}
