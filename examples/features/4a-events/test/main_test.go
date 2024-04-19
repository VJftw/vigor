package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApp(t *testing.T) {
	fs := http.FileServer(http.Dir("../build"))
	s := httptest.NewServer(fs)
	defer s.Close()

	page := rod.New().MustConnect().
		Timeout(10 * time.Second).
		Logger(rod.DefaultLogger).
		MustPage(s.URL)
	defer page.MustClose()
	page.MustWaitLoad()
	page.MustSetWindow(0, 0, 1920, 1080)

	appEl := page.MustElement("#app")

	appHTML, err := appEl.HTML()
	require.NoError(t, err)

	require.NoError(t,
		page.Mouse.MoveTo(proto.Point{X: 0, Y: 0}),
	)

	assert.Equal(t,
		`<div id="app"><div>The mouse position is: 0 x 0</div></div>`,
		appHTML,
	)

	require.NoError(t,
		page.Mouse.MoveTo(proto.Point{X: 500, Y: 10}),
	)

	page.WaitStable(3 * time.Second)

	appHTML, err = appEl.HTML()
	require.NoError(t, err)
	assert.Equal(t,
		`<div id="app"><div>The mouse position is: 500 x 10</div></div>`,
		appHTML,
	)
}
