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
		`<div id="app"><ul><vigor-index><li><a target="_blank" href="https://www.youtube.com/watch?v=J---aiyznGQ">1: Keyboard Cat</a></li><li><a target="_blank" href="https://www.youtube.com/watch?v=z_AbfPXTKms">2: Maru</a></li><li><a target="_blank" href="https://www.youtube.com/watch?v=OUtn3pvWmpg">3: Henri The Existential Cat</a></li></vigor-index></ul></div>`,
		appHTML,
	)
}
