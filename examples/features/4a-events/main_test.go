package main_test

import (
	"testing"

	"github.com/VJftw/vigor/examples/features"
	"github.com/chromedp/cdproto/input"
	"github.com/chromedp/chromedp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test4aEvents(t *testing.T) {
	page := features.BuildServeAndGetVigorPage(t)

	var appElHtml string
	require.NoError(t, chromedp.Run(page,
		chromedp.OuterHTML("#app", &appElHtml),
	))

	require.NoError(t,
		chromedp.Run(page,
			chromedp.MouseEvent(input.MouseMoved, 0, 0),
			chromedp.OuterHTML("#app", &appElHtml),
		),
	)
	assert.Equal(t,
		`<div id="app"><div>The mouse position is: 0 x 0</div></div>`,
		appElHtml,
	)

	require.NoError(t,
		chromedp.Run(page,
			chromedp.MouseEvent(input.MouseMoved, 500, 10),
			chromedp.OuterHTML("#app", &appElHtml),
		),
	)
	assert.Equal(t,
		`<div id="app"><div>The mouse position is: 500 x 10</div></div>`,
		appElHtml,
	)
}
