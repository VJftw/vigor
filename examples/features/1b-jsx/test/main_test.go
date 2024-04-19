package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-rod/rod"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApp(t *testing.T) {
	fs := http.FileServer(http.Dir("../build"))
	s := httptest.NewServer(fs)
	defer s.Close()

	page := rod.New().MustConnect().Timeout(10 * time.Second).Logger(rod.DefaultLogger).MustPage(s.URL)
	defer page.MustClose()
	page.MustWaitStable()

	appEl := page.MustElement("#app")

	appHTML, err := appEl.HTML()
	require.NoError(t, err)

	assert.Equal(t, `<div id="app"><div><div>Hello Solid!</div><svg height="300" width="400"><defs><linearGradient id="gr1" x1="0%" y1="60%" x2="100%" y2="0%"><stop offset="5%" style="stop-color: rgb(255, 255, 3); stop-opacity: 1;"></stop><stop offset="100%" style="stop-color: rgb(255, 0, 0); stop-opacity: 1;"></stop></linearGradient></defs><ellipse cx="125" cy="150" rx="100" ry="60" fill="url(#gr1)"></ellipse>Sorry but this browser does not support inline SVG.</svg></div></div>`, appHTML)

	namespaceUri, err := page.MustElement("svg").Property("namespaceURI")
	require.NoError(t, err)
	assert.Equal(t, "http://www.w3.org/2000/svg", namespaceUri.Str())
}
