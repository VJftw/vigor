package tests

import (
	"testing"

	"github.com/VJftw/vigor"
	"github.com/stretchr/testify/assert"
)

func TestSignal(t *testing.T) {
	t.Run("created signal returns initial value", func(t *testing.T) {
		a, setA := vigor.CreateSignal("a")
		_ = setA

		assert.Equal(t, "a", a())
	})

	t.Run("created signal returns newly set value", func(t *testing.T) {
		a, setA := vigor.CreateSignal("a")
		setA("b")

		assert.Equal(t, "b", a())
	})

	t.Run("signal triggers subscriptions upon setting a different value", func(t *testing.T) {
		a, setA := vigor.CreateSignal("a")

		listenerTriggerCount := 0

		a(vigor.NewFnSubscriber().SetFn(func() {
			listenerTriggerCount++
		}))

		assert.Equal(t, listenerTriggerCount, 0)

		setA("b")
		assert.Equal(t, listenerTriggerCount, 1)
	})

	t.Run("signal does not trigger subscriptions upon setting the same value", func(t *testing.T) {
		a, setA := vigor.CreateSignal("a")

		listenerTriggerCount := 0

		a(vigor.NewFnSubscriber().SetFn(func() {
			listenerTriggerCount++
		}))

		assert.Equal(t, listenerTriggerCount, 0)

		setA("a")
		assert.Equal(t, listenerTriggerCount, 0)
	})

	t.Run("signal does not duplicate subscriptions", func(t *testing.T) {
		a, setA := vigor.CreateSignal("a")

		listenerTriggerCount := 0
		subscriptionFn := vigor.NewFnSubscriber()
		subscriptionFn.SetFn(func() {
			listenerTriggerCount++
			a(subscriptionFn)
		}).Run()

		setA("b")
		assert.Equal(t, 2, listenerTriggerCount)
		setA("c")
		assert.Equal(t, 3, listenerTriggerCount)
		setA("d")
		assert.Equal(t, 4, listenerTriggerCount)
		setA("e")
		assert.Equal(t, 5, listenerTriggerCount)
		setA("f")
		assert.Equal(t, 6, listenerTriggerCount)
		setA("g")
		assert.Equal(t, 7, listenerTriggerCount)
		setA("h")
	})
}
