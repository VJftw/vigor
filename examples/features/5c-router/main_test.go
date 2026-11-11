package main_test

import (
	"testing"

	"github.com/VJftw/vigor/examples/features"
	"github.com/stretchr/testify/assert"
)

func TestRouter(t *testing.T) {
	page := features.BuildServeAndGetVigorPage(t)
	assert.Equal(t, "/", features.Text(t, page, "#path"))
	assert.Equal(t, "home", features.Text(t, page, "#view"))

	// A link navigates without reloading the page.
	features.Eval(t, page, `window.__loaded = true`, nil)
	features.Click(t, page, "#to-1")
	assert.Equal(t, "/items/1", features.Text(t, page, "#path"))
	assert.Equal(t, "item 1 (built 1)", features.Text(t, page, "#view"))
	var pathname string
	var loaded bool
	features.Eval(t, page, `location.pathname`, &pathname)
	features.Eval(t, page, `window.__loaded === true`, &loaded)
	assert.Equal(t, "/items/1", pathname)
	assert.True(t, loaded, "no page reload")

	// Only the parameter changed: the view isn't rebuilt.
	features.Click(t, page, "#to-2")
	assert.Equal(t, "item 2 (built 1)", features.Text(t, page, "#view"))

	// Back, and forward, are followed.
	features.Eval(t, page, `history.back(); true`, nil)
	features.Eventually(t, page, "back", func() bool { return features.Text(t, page, "#view") == "item 1 (built 1)" })

	// A path nothing matches.
	features.Eval(t, page, `history.pushState(null, '', '/nope'); dispatchEvent(new PopStateEvent('popstate')); true`, nil)
	features.Eventually(t, page, "not found", func() bool { return features.Text(t, page, "#view") == "not found" })

	// Query parameters are signals, changed without adding history.
	var before, after float64
	features.Eval(t, page, `history.length`, &before)
	features.Click(t, page, "#sort-asc")
	assert.Equal(t, "asc", features.Text(t, page, "#sort"))
	var search string
	features.Eval(t, page, `location.search`, &search)
	features.Eval(t, page, `history.length`, &after)
	assert.Equal(t, "?sort=asc", search)
	assert.Equal(t, before, after)
}
