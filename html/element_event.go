package html

import (
	"context"
	"syscall/js"

	"github.com/VJftw/vigor"
)

type JSEventFunc func(this js.Value, args []js.Value)

// https://www.solidjs.com/tutorial/bindings_events
type ElEvent interface{}

type elementEvent struct {
	event string
	fn    JSEventFunc
	opts  EventOptions
}

// EventOptions are addEventListener's options. Passive defaults to the
// browser's: touch and wheel listeners are passive, so they can't
// preventDefault; set Passive to false to be able to (e.g. to drag a slider
// without scrolling the page).
type EventOptions struct {
	Passive *bool
	Capture bool
	Once    bool
}

// Passive is a helper for EventOptions.Passive.
func Passive(b bool) *bool { return &b }

// On listens for an event on the element.
func On(ev string, fn JSEventFunc) *elementEvent {
	return &elementEvent{
		event: ev,
		fn:    fn,
	}
}

// OnWith is On with options.
func OnWith(ev string, opts EventOptions, fn JSEventFunc) *elementEvent {
	return &elementEvent{event: ev, fn: fn, opts: opts}
}

// OnClick runs fn when the element is clicked (or tapped).
func OnClick(fn func()) *elementEvent {
	return On("click", func(js.Value, []js.Value) { fn() })
}

// OnClick2 is OnClick with the event, e.g. to stop it propagating or see what
// was clicked.
func OnClick2(fn func(ev js.Value)) *elementEvent {
	return On("click", func(_ js.Value, args []js.Value) { fn(args[0]) })
}

// OnInput calls fn with the element's value as it's typed in.
func OnInput(fn func(value string)) *elementEvent {
	return On("input", func(this js.Value, _ []js.Value) { fn(this.Get("value").String()) })
}

// OnChange calls fn with the element's value when it's committed (a select
// chosen, an input left).
func OnChange(fn func(value string)) *elementEvent {
	return On("change", func(this js.Value, _ []js.Value) { fn(this.Get("value").String()) })
}

// PreventDefault wraps a handler so the browser's own action (following a
// link, submitting a form) doesn't happen.
func PreventDefault(fn JSEventFunc) JSEventFunc {
	return func(this js.Value, args []js.Value) {
		if len(args) > 0 {
			args[0].Call("preventDefault")
		}
		fn(this, args)
	}
}

func (o EventOptions) js() js.Value {
	opts := map[string]any{"capture": o.Capture, "once": o.Once}
	if o.Passive != nil {
		opts["passive"] = *o.Passive
	}

	return js.ValueOf(opts)
}

// listen attaches fn to target, and removes it when the current scope is
// disposed (so handlers of elements that went away don't leak).
func listen(target js.Value, event string, opts EventOptions, fn JSEventFunc) {
	f := js.FuncOf(func(this js.Value, args []js.Value) any {
		fn(this, args)

		return nil
	})
	o := opts.js()
	target.Call("addEventListener", event, f, o)
	vigor.OnCleanup(func() {
		target.Call("removeEventListener", event, f, o)
		f.Release()
	})
}

// OnWindow listens to an event on the window while the surrounding
// component is alive (resize, online, popstate, keydown...). It's a Node
// that renders nothing, so it goes among an element's children.
func OnWindow(ev string, fn JSEventFunc, opts ...EventOptions) Node {
	return onGlobal("window", ev, fn, opts)
}

// OnDocument is OnWindow for the document (visibilitychange, pointerup...).
func OnDocument(ev string, fn JSEventFunc, opts ...EventOptions) Node {
	return onGlobal("document", ev, fn, opts)
}

func onGlobal(name, ev string, fn JSEventFunc, opts []EventOptions) Node {
	var o EventOptions
	if len(opts) > 0 {
		o = opts[0]
	}

	return NewNode(func(_ context.Context, _ js.Value) js.Value {
		target := js.Global()
		if name == "document" {
			target = target.Get("document")
		}
		listen(target, ev, o, fn)

		return js.Undefined()
	})
}
