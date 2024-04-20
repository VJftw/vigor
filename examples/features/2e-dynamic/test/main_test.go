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

	assert.Equal(t,
		`<div id="app"><div><select><option value="red">red</option><option value="green">green</option><option value="blue">blue</option></select><strong style="color: red;">Red Thing</strong></div></div>`,
		appEl.MustHTML(),
	)

	page.MustElement("select").MustSelect("green")
	assert.Equal(t,
		`<div id="app"><div><select><option value="red">red</option><option value="green">green</option><option value="blue">blue</option></select><strong style="color: green;">Green Thing</strong></div></div>`,
		appEl.MustHTML(),
	)

	page.MustElement("select").MustSelect("blue")
	assert.Equal(t,
		`<div id="app"><div><select><option value="red">red</option><option value="green">green</option><option value="blue">blue</option></select><strong style="color: blue;">Blue Thing</strong></div></div>`,
		appEl.MustHTML(),
	)
}
