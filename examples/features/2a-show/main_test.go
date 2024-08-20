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
		chromedp.InnerHTML("#vigor-main", &appElHtml),
	))

	assert.Equal(t,
		`<div><button>Log in</button></div>`,
		appElHtml,
	)

	require.NoError(t, chromedp.Run(page,
		chromedp.Click("button"),
		chromedp.InnerHTML("#vigor-main", &appElHtml),
	))
	assert.Equal(t,
		`<div><button>Log out</button></div>`,
		appElHtml,
	)

	require.NoError(t, chromedp.Run(page,
		chromedp.Click("button"),
		chromedp.InnerHTML("#vigor-main", &appElHtml),
	))
	assert.Equal(t,
		`<div><button>Log in</button></div>`,
		appElHtml,
	)
}
