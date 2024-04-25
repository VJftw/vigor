package main_test

import (
	"testing"

	"github.com/VJftw/vigor/examples/features"
	"github.com/chromedp/chromedp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_2A_Show(t *testing.T) {
	page := features.BuildServeAndGetVigorPage(t)

	var appElHtml string
	require.NoError(t, chromedp.Run(page,
		chromedp.OuterHTML("#app", &appElHtml),
	))

	assert.Equal(t,
		`<div id="app"><div><button>Log in</button></div></div>`,
		appElHtml,
	)

	require.NoError(t, chromedp.Run(page,
		chromedp.Click("button"),
		chromedp.OuterHTML("#app", &appElHtml),
	))
	assert.Equal(t,
		`<div id="app"><div><button>Log out</button></div></div>`,
		appElHtml,
	)

	require.NoError(t, chromedp.Run(page,
		chromedp.Click("button"),
		chromedp.OuterHTML("#app", &appElHtml),
	))
	assert.Equal(t,
		`<div id="app"><div><button>Log in</button></div></div>`,
		appElHtml,
	)
}
