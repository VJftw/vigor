package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-rod/rod"
	"github.com/stretchr/testify/assert"
)

func TestApp(t *testing.T) {
	fs := http.FileServer(http.Dir("../build"))
	s := httptest.NewServer(fs)
	defer s.Close()

	page := rod.New().MustConnect().Timeout(10 * time.Second).Logger(rod.DefaultLogger).MustPage(s.URL)
	defer page.MustClose()
	page.Race().Element("#vigor-info").MustDo()

	appEl := page.MustElement("#app")

	appHTML := appEl.MustHTML()

	assert.Equal(t,
		`<div id="app"><vigor-switch><p>7 is between 5 and 10</p></vigor-switch></div>`,
		appHTML,
	)
}
