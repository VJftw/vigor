package main_test

import (
	"testing"

	"github.com/VJftw/vigor/examples/features"
	"github.com/chromedp/chromedp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test1bJSX(t *testing.T) {
	page := features.BuildServeAndGetVigorPage(t)

	var appElHtml string
	require.NoError(t, chromedp.Run(page,
		chromedp.OuterHTML("#app", &appElHtml),
	))

	assert.Equal(t, `<div id="app"><div><div>Hello Solid!</div><svg height="300" width="400"><defs><linearGradient id="gr1" x1="0%" y1="60%" x2="100%" y2="0%"><stop offset="5%" style="stop-color: rgb(255, 255, 3); stop-opacity: 1;"></stop><stop offset="100%" style="stop-color: rgb(255, 0, 0); stop-opacity: 1;"></stop></linearGradient></defs><ellipse cx="125" cy="150" rx="100" ry="60" fill="url(#gr1)"></ellipse>Sorry but this browser does not support inline SVG.</svg></div></div>`, appElHtml)

	namespaceUri := ""
	require.NoError(t, chromedp.Run(page,
		chromedp.Evaluate(
			`let el = document.getElementsByTagNameNS("http://www.w3.org/2000/svg", "svg"); el ? el[0].namespaceURI : ''`,
			&namespaceUri,
		),
	))

	assert.Equal(t, "http://www.w3.org/2000/svg", namespaceUri)
}
