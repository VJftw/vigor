package features

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/stretchr/testify/require"
)

func BuildServeAndGetVigorPage(t *testing.T, chromedpOpts ...func(*chromedp.Context)) context.Context {
	wd, err := os.Getwd()
	require.NoError(t, err)
	wasmExecJs, err := os.ReadFile(filepath.Join(runtime.GOROOT(), "misc", "wasm", "wasm_exec.js"))
	require.NoError(t, err)

	indexHTML, err := os.ReadFile(filepath.Join(wd, "..", "..", "_common", "index.html"))
	require.NoError(t, err)

	tempDir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(tempDir, "wasm_exec.js"), wasmExecJs, 0o600))
	require.NoError(t, os.WriteFile(filepath.Join(tempDir, "index.html"), indexHTML, 0o600))

	cmd := exec.Command("go", "build", "-o", filepath.Join(tempDir, "main.wasm"), "./")
	cmd.Env = append(os.Environ(),
		"GOARCH=wasm",
		"GOOS=js",
	)
	out, err := cmd.CombinedOutput()
	require.NoErrorf(t, err, "output: %s", out)

	chromeCtx, cancelChromeCtx := chromedp.NewContext(context.TODO(), chromedpOpts...)
	t.Cleanup(func() {
		cancelChromeCtx()
	})

	fs := http.FileServer(http.Dir(tempDir))
	s := httptest.NewServer(fs)

	if err := chromedp.Run(chromeCtx,
		chromedp.Navigate(s.URL),
	); err != nil {
		t.Logf("could not navigate to '%s': %s", s.URL, err)
		require.NoError(t, err)
	}

	chromeCtx, cancel := context.WithTimeout(chromeCtx, 30*time.Second)
	t.Cleanup(func() { cancel() })
	if err := chromedp.Run(chromeCtx,
		chromedp.WaitReady("#vigor-info"),
	); err != nil {
		require.NoError(t, err)
	}

	t.Cleanup(func() {
		s.Close()
	})

	return chromeCtx
}
