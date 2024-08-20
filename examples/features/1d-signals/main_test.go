package main_test

import (
	"testing"
	"time"

	"github.com/VJftw/vigor/examples/features"
	"github.com/chromedp/chromedp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test1dSignals(t *testing.T) {
	page := features.BuildServeAndGetVigorPage(t)

	var appElHtml string
	require.NoError(t, chromedp.Run(page,
		chromedp.InnerHTML("#vigor-main", &appElHtml),
	))

	assert.Equal(t,
		`<div>Count: 0</div>`,
		appElHtml,
	)

	assert.Eventually(t, func() bool {
		require.NoError(t, chromedp.Run(page,
			chromedp.InnerHTML("#vigor-main", &appElHtml),
		))

		return appElHtml == `<div>Count: 1</div>`
	}, 5*time.Second, 100*time.Millisecond)

	assert.Eventually(t, func() bool {
		require.NoError(t, chromedp.Run(page,
			chromedp.InnerHTML("#vigor-main", &appElHtml),
		))

		return appElHtml == `<div>Count: 2</div>`
	}, 5*time.Second, 100*time.Millisecond)
}
