package main_test

import (
	"testing"

	"github.com/VJftw/vigor/examples/features"
	"github.com/chromedp/chromedp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test2eDynamic(t *testing.T) {
	page := features.BuildServeAndGetVigorPage(t)

	var appElHtml string
	require.NoError(t, chromedp.Run(page,
		chromedp.OuterHTML("#app", &appElHtml),
	))

	assert.Equal(t,
		`<div id="app"><div><select><option value="red">red</option><option value="green">green</option><option value="blue">blue</option></select><strong style="color: red;">Red Thing</strong></div></div>`,
		appElHtml,
	)

	require.NoError(t, chromedp.Run(page,
		chromedp.SetValue(`select`, "green", chromedp.ByQuery),
		chromedp.OuterHTML("#app", &appElHtml),
	))
	assert.Equal(t,
		`<div id="app"><div><select><option value="red">red</option><option value="green">green</option><option value="blue">blue</option></select><strong style="color: green;">Green Thing</strong></div></div>`,
		appElHtml,
	)

	require.NoError(t, chromedp.Run(page,
		chromedp.SetValue(`select`, "blue", chromedp.ByQuery),
		chromedp.OuterHTML("#app", &appElHtml),
	))
	assert.Equal(t,
		`<div id="app"><div><select><option value="red">red</option><option value="green">green</option><option value="blue">blue</option></select><strong style="color: blue;">Blue Thing</strong></div></div>`,
		appElHtml,
	)
}
