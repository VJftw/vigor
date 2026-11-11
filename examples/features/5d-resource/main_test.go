package main_test

import (
	"testing"

	"github.com/VJftw/vigor/examples/features"
	"github.com/stretchr/testify/assert"
)

func TestResource(t *testing.T) {
	page := features.BuildServeAndGetVigorPage(t)
	ready := func() string { return features.Text(t, page, "#item .ready") }

	assert.Equal(t, "loading", features.Text(t, page, "#item .pending"))
	features.Eventually(t, page, "loaded", func() bool { return ready() == "item-1#1" })

	// The source changed: fetched again, and the ready view updated in place.
	features.Eval(t, page, `window.__p = document.querySelector('#item .ready'); true`, nil)
	features.Click(t, page, "#next")
	features.Eventually(t, page, "the next item", func() bool { return ready() == "item-2#2" })
	var same bool
	features.Eval(t, page, `document.querySelector('#item .ready') === window.__p`, &same)
	assert.True(t, same, "ready wasn't rebuilt")

	// Refetch.
	features.Click(t, page, "#refetch")
	features.Eventually(t, page, "refetched", func() bool { return ready() == "item-2#3" })

	// A failed fetch keeps the last value, and reports the error.
	features.Click(t, page, "#bad")
	features.Eventually(t, page, "the error", func() bool { return features.Text(t, page, "#error") == "no such item" })
	assert.Equal(t, "item-2#3", ready())

	// Failing before it ever loaded shows failed.
	features.Eventually(t, page, "failed", func() bool { return features.Text(t, page, "#broken .failed") == "down" })
}
