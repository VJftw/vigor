package tests

import (
	"testing"

	"github.com/VJftw/vigor"
	"github.com/stretchr/testify/assert"
)

func TestEffect(t *testing.T) {
	t.Run("effect is triggered upon signal change", func(t *testing.T) {
		a, setA := vigor.CreateSignal("a")

		effectTriggerCount := 0

		vigor.CreateEffect(func() {
			effectTriggerCount++
		}, a)

		assert.Equal(t, effectTriggerCount, 0)

		setA("b")
		assert.Equal(t, effectTriggerCount, 1)
	})

	t.Run("effect is not triggered if signal does not change", func(t *testing.T) {
		a, setA := vigor.CreateSignal("a")

		effectTriggerCount := 0

		vigor.CreateEffect(func() {
			effectTriggerCount++
		}, a)

		assert.Equal(t, effectTriggerCount, 0)

		setA("a")
		assert.Equal(t, effectTriggerCount, 0)
	})

	t.Run("panic if no dependencies passed", func(t *testing.T) {
		a, setA := vigor.CreateSignal("a")
		_ = a
		_ = setA
		assert.Panics(t, func() {
			vigor.CreateEffect(func() {})
		})
	})
}
