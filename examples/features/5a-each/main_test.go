package main_test

import (
	"testing"

	"github.com/VJftw/vigor/examples/features"
	"github.com/stretchr/testify/assert"
)

func TestEach(t *testing.T) {
	page := features.BuildServeAndGetVigorPage(t)
	rows := func() string {
		var s string
		features.Eval(t, page, `[...document.querySelectorAll('#list li')].map(li => li.textContent).join(',')`, &s)

		return s
	}
	assert.Equal(t, "0: a,1: b,2: c", rows())

	// Reordering moves the same nodes, and their index updates.
	features.Eval(t, page, `window.__a = document.querySelector('[data-key=a]'); true`, nil)
	features.Click(t, page, "#reverse")
	assert.Equal(t, "0: c,1: b,2: a", rows())
	var same bool
	features.Eval(t, page, `document.querySelector('[data-key=a]') === window.__a`, &same)
	assert.True(t, same, "the row kept its DOM node")
	assert.Equal(t, "0", features.Text(t, page, "#cleaned"))

	// A removed key is disposed.
	features.Click(t, page, "#drop-first")
	assert.Equal(t, "0: b,1: a", rows())
	assert.Equal(t, "1", features.Text(t, page, "#cleaned"))

	// Empty shows the fallback; the rest are disposed.
	features.Click(t, page, "#clear")
	assert.Equal(t, "nothing", rows())
	assert.Equal(t, "3", features.Text(t, page, "#cleaned"))
}
