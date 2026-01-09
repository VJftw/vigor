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
		chromedp.InnerHTML("#vigor-main", &appElHtml),
	))

	require.NoError(t,
		chromedp.Run(page,
			chromedp.MouseEvent(input.MouseMoved, 0, 0),
			chromedp.InnerHTML("#vigor-main", &appElHtml),
		),
	)
	assert.Equal(t,
		`<div>The mouse position is: 0 x 0</div>`,
		appElHtml,
	)

	require.NoError(t,
		chromedp.Run(page,
			chromedp.MouseEvent(input.MouseMoved, 500, 10),
			chromedp.InnerHTML("#vigor-main", &appElHtml),
		),
	)
	assert.Equal(t,
		`<div>The mouse position is: 500 x 10</div>`,
		appElHtml,
	)
}
