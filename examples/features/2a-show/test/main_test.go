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
		`<div id="app"><div><vigor-show><button>Log in</button></vigor-show></div></div>`,
		appHTML,
	)

	buttonEl := page.MustElement("button")

	buttonEl.MustClick()

	appHTML = appEl.MustHTML()
	assert.Equal(t,
		`<div id="app"><div><vigor-show><button>Log out</button></vigor-show></div></div>`,
		appHTML,
	)
}
