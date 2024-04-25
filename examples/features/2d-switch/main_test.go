package main_test

import (
	"testing"

	"github.com/VJftw/vigor/examples/features"
	"github.com/chromedp/chromedp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test2dSwitch(t *testing.T) {
	page := features.BuildServeAndGetVigorPage(t)

	var appElHtml string
	require.NoError(t, chromedp.Run(page,
		chromedp.OuterHTML("#app", &appElHtml),
	))

	assert.Equal(t,
		`<div id="app"><p>7 is between 5 and 10</p></div>`,
		appElHtml,
	)
}
