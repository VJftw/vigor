package html

import "syscall/js"

type JSEventFunc func(this js.Value, args []js.Value)

// https://www.solidjs.com/tutorial/bindings_events
type ElEvent interface{}

type elementEvent struct {
	event string
	fn    JSEventFunc
}

func On(ev string, fn JSEventFunc) *elementEvent {
	return &elementEvent{
		event: ev,
		fn:    fn,
	}
}
