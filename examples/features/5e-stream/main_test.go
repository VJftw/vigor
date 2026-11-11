package main_test

import (
	"testing"

	"github.com/VJftw/vigor/examples/features"
	"github.com/stretchr/testify/assert"
)

func TestStream(t *testing.T) {
	page := features.BuildServeAndGetVigorPage(t)

	// It dropped twice, reconnected each time, and stayed up the third.
	features.Eventually(t, page, "the third connection", func() bool { return features.Text(t, page, "#latest") == "30" })
	assert.Equal(t, "3", features.Text(t, page, "#connects"))
	assert.Equal(t, "true", features.Text(t, page, "#connected"))
	assert.Equal(t, "none", features.Text(t, page, "#error"))
}
