package features

import (
	"context"
	"log"
	"net/http/httptest"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/VJftw/vigor/site"
	"github.com/chromedp/chromedp"
	"github.com/stretchr/testify/require"
)

func BuildServeAndGetVigorPage(t *testing.T, chromedpOpts ...func(*chromedp.Context)) context.Context {
	tempDir := t.TempDir()

	// A cold build cache, with other examples building too, takes a while.
	ctx, cancel := context.WithTimeout(t.Context(), 3*time.Minute)
	defer cancel()
	wasm := filepath.Join(t.TempDir(), "main.wasm")
	require.NoError(t, site.Build(ctx, "./", wasm))
	_, err := site.Write(tempDir, wasm, site.Options{})
	require.NoError(t, err)

	chromeCtx, cancelChromeCtx := chromedp.NewContext(t.Context(), chromedpOpts...)
	t.Cleanup(func() {
		cancelChromeCtx()
	})

	s := httptest.NewServer(site.Handler(tempDir))
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
