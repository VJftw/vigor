package main_test

import (
	"testing"

	"github.com/VJftw/vigor/examples/features"
	"github.com/chromedp/chromedp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test2gErrorBoundary(t *testing.T) {
	page := features.BuildServeAndGetVigorPage(t)
	var appElHtml string
	require.NoError(t, chromedp.Run(page,
		chromedp.OuterHTML("#app", &appElHtml),
	))

	assert.Equal(t,
		`<div id="app"><div><div>Before</div><div>Error: Oh No</div><div>After</div></div></div>`,
		appElHtml,
	)
}
