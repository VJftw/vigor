package vigor

import "reflect"

// GetterFn represents functions which return the current value of a Signal.
// Passing `Subscriber` implementations to them will call those `Subscriber`s
// whenever the value of the Signal changes.
type GetterFn func(...Subscriber) any

// SetterFn represents functions which set the current value of a Signal.
// Calling this triggers any `Subscriber`s for the Signal.
type SetterFn func(any)

// Subscriber abstracts the implementation of callbacks which may subscribe to
// Signal changes.
type Subscriber interface {
	// Run executes the callback.
	Run()
}

// CreateSignal returns a "getter" and "setter" functions with the given initial
// value. This implementation is a translation/interpretation of the JavaScript
// example code from https://www.solidjs.com/guides/reactivity#how-it-works.
func CreateSignal(value any) (GetterFn, SetterFn) {
	subscribers := map[Subscriber]struct{}{}

	getterFn := func(listeners ...Subscriber) any {
		for _, l := range listeners {
			if _, ok := subscribers[l]; !ok {
				subscribers[l] = struct{}{}
			}
		}

		return value
	}

	setterFn := func(nextValue any) {
		if reflect.DeepEqual(value, nextValue) {
			return
		}

		// convert all slices to []any
		if reflect.TypeOf(nextValue).Kind() == reflect.Slice {
			sliceVal := reflect.ValueOf(nextValue)
			newNextValue := make([]any, sliceVal.Len())
			for i := 0; i < sliceVal.Len(); i++ {
				newNextValue[i] = sliceVal.Index(i).Interface()
			}
			nextValue = newNextValue
		}

		value = nextValue

		for s := range subscribers {
			s.Run()
		}
	}

	return getterFn, setterFn
}

// NewFnSubscriber returns a new `Subscriber` implementation that calls the
// callback passed to `SetFn`. This is useful to easily wrap a Go anonymous
// function as a `Subscriber`.
func NewFnSubscriber() *fnSubscriber {
	return &fnSubscriber{
		fn: func() {},
	}
}

type fnSubscriber struct {
	fn func()
}

// SetFn sets the given anonymous function as the callback for this `Subscriber`.
func (l *fnSubscriber) SetFn(fn func()) *fnSubscriber {
	l.fn = fn

	return l
}

// Run implements `Subscriber.Run` by calling the configured function on this
// `Subscriber`.
func (l *fnSubscriber) Run() {
	l.fn()
}

// NewDerivedSignal creates a *Derived Signal* which allows you wrap a signal
// with an expression that is updated when the upstream is. This acts like a
// Signal but does not store any state.
func NewDerivedSignal(upstream GetterFn, fn func(v any) any) GetterFn {
	return func(s ...Subscriber) any {
		return fn(upstream(s...))
	}
}
