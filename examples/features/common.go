package features

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"testing"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/stretchr/testify/require"
)

const indexHTML = `
<html>
<head>
	<meta charset="utf-8"/>
	<script src="wasm_exec.js"></script>
</head>
<body>
    <div id="app"></div>
	<script>
		document.addEventListener("DOMContentLoaded", (event) => {
		const go = new Go();
		WebAssembly.instantiateStreaming(fetch("main.wasm"), go.importObject).then((result) => {
			go.run(result.instance);
		});
});

	</script>
</body>
</html>
`

func BuildServeAndGetVigorPage(t *testing.T, chromedpOpts ...func(*chromedp.Context)) context.Context {
	wasmExecJs, err := os.ReadFile(filepath.Join(runtime.GOROOT(), "lib", "wasm", "wasm_exec.js"))
	require.NoError(t, err)

	tempDir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(tempDir, "wasm_exec.js"), wasmExecJs, 0o600))
	require.NoError(t, os.WriteFile(filepath.Join(tempDir, "index.html"), []byte(indexHTML), 0o600))

	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "go", "build", "-o", filepath.Join(tempDir, "main.wasm"), "./")
	cmd.Env = append(os.Environ(),
		"GOARCH=wasm",
		"GOOS=js",
	)
	out, err := cmd.CombinedOutput()
	require.NoErrorf(t, err, "output: %s", out)

	chromeCtx, cancelChromeCtx := chromedp.NewContext(t.Context(), chromedpOpts...)
	t.Cleanup(func() {
		cancelChromeCtx()
	})

	fs := http.FileServer(http.Dir(tempDir))
	s := httptest.NewServer(fs)
	log.Printf("Test Page at %s", s.URL)
	if err := chromedp.Run(chromeCtx,
		chromedp.Navigate(s.URL),
	); err != nil {
		t.Logf("could not navigate to '%s': %s", s.URL, err)
		require.NoError(t, err)
	}

	chromeCtx, cancel = context.WithTimeout(chromeCtx, 30*time.Second)
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

// Eval evaluates JavaScript in the page, into out (a pointer).
func Eval(t *testing.T, page context.Context, js string, out any) {
	t.Helper()
	require.NoError(t, chromedp.Run(page, chromedp.Evaluate(js, out)))
}

// Text is the text content of the first element matching a selector, or
// "<none>" if there isn't one.
func Text(t *testing.T, page context.Context, selector string) string {
	t.Helper()
	var s string
	Eval(t, page, "document.querySelector("+strconv.Quote(selector)+")?.textContent ?? '<none>'", &s)

	return s
}

// Click clicks the first element matching a selector, in the page.
func Click(t *testing.T, page context.Context, selector string) {
	t.Helper()
	Eval(t, page, "document.querySelector("+strconv.Quote(selector)+").click()", nil)
}

// Eventually waits for cond to hold, printing the page if it doesn't.
func Eventually(t *testing.T, page context.Context, what string, cond func() bool) {
	t.Helper()
	for end := time.Now().Add(5 * time.Second); time.Now().Before(end); time.Sleep(20 * time.Millisecond) {
		if cond() {
			return
		}
	}
	var html string
	Eval(t, page, "document.body.innerHTML", &html)
	t.Fatalf("timed out waiting for %s; page:\n%s", what, html)
}
