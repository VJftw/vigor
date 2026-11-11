package main_test

import (
	"testing"

	"github.com/VJftw/vigor/examples/features"
	"github.com/chromedp/chromedp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLocalStorage(t *testing.T) {
	page := features.BuildServeAndGetVigorPage(t)
	assert.Equal(t, "hello anonymous", features.Text(t, page, "#hello"))

	features.Eval(t, page, `(() => { const i = document.querySelector('#name'); i.value = 'vj'; i.dispatchEvent(new Event('input')) })()`, nil)
	assert.Equal(t, "hello vj", features.Text(t, page, "#hello"))
	var stored string
	features.Eval(t, page, `localStorage.getItem('vigor-example-name')`, &stored)
	assert.Equal(t, "vj", stored)

	// A new visit starts from what was stored.
	require.NoError(t, chromedp.Run(page, chromedp.Reload(), chromedp.WaitReady("#vigor-info")))
	features.Eventually(t, page, "the stored name", func() bool { return features.Text(t, page, "#hello") == "hello vj" })
}
