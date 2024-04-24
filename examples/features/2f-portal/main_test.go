package main_test

import (
	"testing"

	"github.com/VJftw/vigor/examples/features"
	"github.com/chromedp/chromedp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test2fPortal(t *testing.T) {
	page := features.BuildServeAndGetVigorPage(t)

	var appElHtml string
	require.NoError(t, chromedp.Run(page,
		chromedp.OuterHTML("#app", &appElHtml),
	))

	assert.Equal(t,
		`<div id="app"><div><p>Just some text inside a div that has a restricted size.</p></div></div>`,
		appElHtml,
	)

	var rootPortalElHTML string
	require.NoError(t, chromedp.Run(page,
		chromedp.OuterHTML("#vigor-portal-root", &rootPortalElHTML),
	))
	assert.Equal(t,
		`<vigor-portal-root id="vigor-portal-root"><vigor-portal><div class="popup"><h1>Popup</h1><p>Some text you might need for something or other.</p></div></vigor-portal></vigor-portal-root>`,
		rootPortalElHTML,
	)
}
