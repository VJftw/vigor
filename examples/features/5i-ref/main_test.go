package main_test

import (
	"testing"

	"github.com/VJftw/vigor/examples/features"
	"github.com/stretchr/testify/assert"
)

func TestRef(t *testing.T) {
	page := features.BuildServeAndGetVigorPage(t)
	var active string
	features.Eval(t, page, `document.activeElement.id`, &active)
	assert.NotEqual(t, "field", active)

	features.Click(t, page, "#focus")
	features.Eval(t, page, `document.activeElement.id`, &active)
	assert.Equal(t, "field", active)
}
