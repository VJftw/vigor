package main_test

import (
	"strconv"
	"testing"
	"time"

	"github.com/VJftw/vigor/examples/features"
	"github.com/stretchr/testify/assert"
)

func TestContextAndCleanup(t *testing.T) {
	page := features.BuildServeAndGetVigorPage(t)

	// A context's value is what the nearest provider gave, else the default.
	assert.Equal(t, "light", features.Text(t, page, "#outside"))
	assert.Equal(t, "dark", features.Text(t, page, "#inside"))

	// Hidden until the condition holds.
	assert.Equal(t, "<none>", features.Text(t, page, "#shown"))
	features.Click(t, page, "#inc")
	assert.Equal(t, "count is 1", features.Text(t, page, "#shown"))
	assert.Equal(t, "1", features.Text(t, page, "#builds"))

	// The condition didn't flip, so the branch isn't rebuilt; what reads
	// the count updates.
	features.Click(t, page, "#inc")
	assert.Equal(t, "count is 2", features.Text(t, page, "#shown"))
	assert.Equal(t, "1", features.Text(t, page, "#builds"))
	assert.Equal(t, "0", features.Text(t, page, "#cleanups"))

	// Its interval runs while it's shown.
	ticks := func() int {
		n, _ := strconv.Atoi(features.Text(t, page, "#ticks"))

		return n
	}
	features.Eventually(t, page, "ticks", func() bool { return ticks() >= 2 })

	// Hidden: cleaned up, and the interval stopped.
	features.Click(t, page, "#reset")
	assert.Equal(t, "<none>", features.Text(t, page, "#shown"))
	assert.Equal(t, "1", features.Text(t, page, "#cleanups"))
	time.Sleep(60 * time.Millisecond)
	stopped := ticks()
	time.Sleep(100 * time.Millisecond)
	assert.Equal(t, stopped, ticks(), "the interval kept running")
}
