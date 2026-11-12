//go:build !js

// Command vigor builds and serves vigor apps.
//
//	vigor build [-o dir] [-title t] [-style url]... [package]
//	vigor serve [-addr host:port] [-title t] [-style url]... [package]
//
// package is the Go package of the app's main (the current directory by
// default). build writes a site (see package site) to dir. serve builds the
// app and serves it, building again whenever the page is loaded, so editing
// the app and reloading the browser shows the change.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/VJftw/vigor/site"
)

type list []string

func (l *list) String() string     { return strings.Join(*l, ",") }
func (l *list) Set(s string) error { *l = append(*l, s); return nil }

func main() {
	if len(os.Args) < 2 || (os.Args[1] != "build" && os.Args[1] != "serve") {
		fmt.Fprintln(os.Stderr, "usage: vigor build|serve [flags] [package]")
		os.Exit(2)
	}
	cmd := os.Args[1]
	fs := flag.NewFlagSet(cmd, flag.ExitOnError)
	title := fs.String("title", "", "the page's title")
	var styles list
	fs.Var(&styles, "style", "a stylesheet's URL (repeatable)")
	out := fs.String("o", "dist", "build: the directory to write the site to")
	addr := fs.String("addr", "127.0.0.1:8080", "serve: the address to serve on")
	fs.Parse(os.Args[2:])
	dir := "."
	if fs.NArg() > 0 {
		dir = fs.Arg(0)
	}
	opts := site.Options{Title: *title, Styles: styles}

	switch cmd {
	case "build":
		if _, err := assemble(context.Background(), dir, *out, opts); err != nil {
			log.Fatal(err)
		}
		log.Printf("wrote %s", *out)
	case "serve":
		serve(dir, *addr, opts)
	}
}

// assemble builds the app and writes the site to out.
func assemble(ctx context.Context, dir, out string, opts site.Options) (string, error) {
	tmp, err := os.MkdirTemp("", "vigor-wasm-")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(tmp)
	wasm := filepath.Join(tmp, "app.wasm")
	if err := site.Build(ctx, dir, wasm); err != nil {
		return "", err
	}

	return site.Write(out, wasm, opts)
}

func serve(dir, addr string, opts site.Options) {
	out, err := os.MkdirTemp("", "vigor-site-")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(out)

	var mu sync.Mutex
	rebuild := func() error {
		mu.Lock()
		defer mu.Unlock()
		_, err := assemble(context.Background(), dir, out, opts)

		return err
	}
	if err := rebuild(); err != nil {
		log.Fatal(err)
	}

	files := site.Handler(out)
	log.Printf("serving %s at http://%s", dir, addr)
	log.Fatal(http.ListenAndServe(addr, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Loading the page builds again (the Go build cache makes that quick
		// when nothing changed); a build that fails is shown instead.
		if r.Header.Get("Sec-Fetch-Dest") == "document" || (filepath.Ext(r.URL.Path) == "" && r.Method == http.MethodGet) {
			if err := rebuild(); err != nil {
				log.Print(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)

				return
			}
		}
		w.Header().Set("Cache-Control", "no-store")
		files.ServeHTTP(w, r)
	})))
}
