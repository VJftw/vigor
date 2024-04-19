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
	page.Race().Element("#vigor-info").MustDo()

	appEl := page.MustElement("#app")

	appHTML, err := appEl.HTML()
	require.NoError(t, err)

	assert.Equal(t,
		`<div id="app"><div>Count: 0</div></div>`,
		appHTML,
	)

	assert.Eventually(t, func() bool {
		appHTML, err := appEl.HTML()
		require.NoError(t, err)

		return appHTML == `<div id="app"><div>Count: 2</div></div>`
	}, 3*time.Second, 100*time.Millisecond)

	assert.Eventually(t, func() bool {
		appHTML, err := appEl.HTML()
		require.NoError(t, err)
		return appHTML == `<div id="app"><div>Count: 4</div></div>`
	}, 3*time.Second, 100*time.Millisecond)
}
