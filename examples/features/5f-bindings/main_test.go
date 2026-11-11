package main_test

import (
	"testing"

	"github.com/VJftw/vigor/examples/features"
	"github.com/stretchr/testify/assert"
)

func TestBindings(t *testing.T) {
	page := features.BuildServeAndGetVigorPage(t)
	str := func(js string) string {
		var s string
		features.Eval(t, page, js, &s)

		return s
	}
	boolean := func(js string) bool {
		var b bool
		features.Eval(t, page, js, &b)

		return b
	}

	// Two-way: the input shows the signal, and typing sets it.
	assert.Equal(t, "hi", str(`document.querySelector('#text').value`))
	features.Eval(t, page, `(() => { const i = document.querySelector('#text'); i.value = 'hello'; i.dispatchEvent(new Event('input')) })()`, nil)
	assert.Equal(t, "you typed: hello", features.Text(t, page, "#echo"))
	assert.Equal(t, "5px", str(`document.querySelector('#box').style.getPropertyValue('--width')`))

	// A bool attribute is there or not; nil removes an attribute; a class toggles.
	assert.True(t, boolean(`document.querySelector('#go').disabled`))
	assert.False(t, boolean(`document.querySelector('#box').hasAttribute('title')`))
	assert.Equal(t, "", str(`document.querySelector('#box').className`))
	features.Click(t, page, "#armed")
	assert.False(t, boolean(`document.querySelector('#go').disabled`))
	assert.Equal(t, "armed", str(`document.querySelector('#box').getAttribute('title')`))
	assert.Equal(t, "on", str(`document.querySelector('#box').className`))
	features.Click(t, page, "#armed")
	assert.True(t, boolean(`document.querySelector('#go').disabled`))
	assert.False(t, boolean(`document.querySelector('#box').hasAttribute('title')`))
}
