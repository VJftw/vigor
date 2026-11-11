package main_test

import (
	"testing"

	"github.com/VJftw/vigor/examples/features"
	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {
	page := features.BuildServeAndGetVigorPage(t)
	assert.Equal(t, "*", features.Text(t, page, "#a"))
	assert.Equal(t, "2", features.Text(t, page, "#len"))
	assert.Equal(t, "1", features.Text(t, page, "#a-runs"))
	assert.Equal(t, "1", features.Text(t, page, "#keys-runs"))

	// Changing another key's value touches neither the "a" reader nor the
	// keys reader.
	features.Click(t, page, "#bump-b")
	assert.Equal(t, "1", features.Text(t, page, "#a-runs"))
	assert.Equal(t, "1", features.Text(t, page, "#keys-runs"))

	// Changing "a" re-runs only its reader.
	features.Click(t, page, "#bump-a")
	assert.Equal(t, "**", features.Text(t, page, "#a"))
	assert.Equal(t, "2", features.Text(t, page, "#a-runs"))
	assert.Equal(t, "1", features.Text(t, page, "#keys-runs"))

	// A new key re-runs the keys reader (and not "a"'s).
	features.Click(t, page, "#add-c")
	assert.Equal(t, "3", features.Text(t, page, "#len"))
	assert.Equal(t, "2", features.Text(t, page, "#keys-runs"))
	assert.Equal(t, "2", features.Text(t, page, "#a-runs"))

	// Replace makes the map exactly those entries: "a" goes, so its reader runs.
	features.Click(t, page, "#replace")
	assert.Equal(t, "2", features.Text(t, page, "#len"))
	assert.Equal(t, "", features.Text(t, page, "#a"))
	assert.Equal(t, "3", features.Text(t, page, "#a-runs"))
}
