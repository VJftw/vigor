package main_test

import (
	"testing"

	"github.com/VJftw/vigor/examples/features"
	"github.com/chromedp/chromedp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test2cIndex(t *testing.T) {
	page := features.BuildServeAndGetVigorPage(t)

	var appElHtml string
	require.NoError(t, chromedp.Run(page,
		chromedp.OuterHTML("#app", &appElHtml),
	))

	assert.Equal(t,
		`<div id="app"><ul><li><a target="_blank" href="https://www.youtube.com/watch?v=J---aiyznGQ">1: Keyboard Cat</a></li><li><a target="_blank" href="https://www.youtube.com/watch?v=z_AbfPXTKms">2: Maru</a></li><li><a target="_blank" href="https://www.youtube.com/watch?v=OUtn3pvWmpg">3: Henri The Existential Cat</a></li></ul></div>`,
		appElHtml,
	)
}
