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
		chromedp.InnerHTML("#vigor-main", &appElHtml),
	))

	assert.Equal(t,
		`<div><select><!--vigor_for-start--><option value="red">red</option><option value="green">green</option><option value="blue">blue</option><!--vigor_for-end--></select><strong style="color: red;">Red Thing</strong></div>`,
		appElHtml,
	)

	require.NoError(t, chromedp.Run(page,
		chromedp.SetValue(`select`, "green", chromedp.ByQuery),
		chromedp.InnerHTML("#vigor-main", &appElHtml),
	))
	assert.Equal(t,
		`<div><select><!--vigor_for-start--><option value="red">red</option><option value="green">green</option><option value="blue">blue</option><!--vigor_for-end--></select><strong style="color: green;">Green Thing</strong></div>`,
		appElHtml,
	)

	require.NoError(t, chromedp.Run(page,
		chromedp.SetValue(`select`, "blue", chromedp.ByQuery),
		chromedp.InnerHTML("#vigor-main", &appElHtml),
	))
	assert.Equal(t,
		`<div><select><!--vigor_for-start--><option value="red">red</option><option value="green">green</option><option value="blue">blue</option><!--vigor_for-end--></select><strong style="color: blue;">Blue Thing</strong></div>`,
		appElHtml,
	)
}
