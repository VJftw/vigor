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
	page.MustSetWindow(0, 0, 640, 480)

	appEl := page.MustElement("#app")
	appHTML := appEl.MustHTML()

	assert.Equal(t,
		`<div id="app"><div><p>Just some text inside a div that has a restricted size.</p></div></div>`,
		appHTML,
	)

	rootPortalEl := page.MustElement("#vigor-portal-root")
	rootPortalHTML := rootPortalEl.MustHTML()
	assert.Equal(t,
		`<vigor-portal-root id="vigor-portal-root"><vigor-portal><div class="popup"><h1>Popup</h1><p>Some text you might need for something or other.</p></div></vigor-portal></vigor-portal-root>`,
		rootPortalHTML,
	)
}
