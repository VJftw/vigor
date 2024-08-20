package vigor_test

import (
	"testing"

	"github.com/VJftw/vigor"
	"github.com/stretchr/testify/assert"
)

func TestMemo(t *testing.T) {
	t.Run("memo is returns initial signal value", func(t *testing.T) {
		a, setA := vigor.CreateSignal("a")
		_ = setA

		memo := vigor.CreateMemo(func() any {
			return "Hello " + a().(string)
		}, a)

		assert.Equal(t, "Hello a", memo())
	})

	t.Run("memo is returns updated signal value", func(t *testing.T) {
		a, setA := vigor.CreateSignal("a")

		memo := vigor.CreateMemo(func() any {
			return "Hello " + a().(string)
		}, a)

		setA("b")

		assert.Equal(t, "Hello b", memo())
	})

	t.Run("memo behaves as a cache", func(t *testing.T) {
		a, setA := vigor.CreateSignal("a")

		expensiveCallCount := 0
		memo := vigor.CreateMemo(func() any {
			expensiveCallCount++
			return "Hello " + a().(string)
		}, a)
		assert.Equal(t, 1, expensiveCallCount)

		setA("b")
		setA("b")
		setA("b")
		setA("b")
		setA("b")
		setA("b")
		setA("b")
		setA("b")

		assert.Equal(t, "Hello b", memo())
		assert.Equal(t, 2, expensiveCallCount)
	})

	t.Run("memo supports multiple signals", func(t *testing.T) {
		firstName, setFirstName := vigor.CreateSignal("Alice")
		lastName, setLastName := vigor.CreateSignal("Smith")

		fullName := vigor.CreateMemo(func() any {
			return firstName().(string) + " " + lastName().(string)
		}, firstName, lastName)
		assert.Equal(t, "Alice Smith", fullName())

		setFirstName("Eve")
		assert.Equal(t, "Eve Smith", fullName())

		setLastName("Rogers")
		assert.Equal(t, "Eve Rogers", fullName())
	})
}
