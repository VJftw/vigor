//go:build !js

package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/VJftw/vigor/site"
)

// build: an app is compiled and written as a site.
func TestAssemble(t *testing.T) {
	if testing.Short() {
		t.Skip("builds a program")
	}
	out := t.TempDir()
	version, err := assemble(t.Context(), "../../site/testdata/hello", out, site.Options{Title: "Hello"})
	if err != nil {
		t.Fatal(err)
	}
	if len(version) != 16 {
		t.Errorf("version %q", version)
	}
	for _, f := range []string{site.IndexFile, site.WasmExecFile, site.WasmFile, site.WasmGzFile} {
		if _, err := os.Stat(filepath.Join(out, f)); err != nil {
			t.Error(err)
		}
	}
}
