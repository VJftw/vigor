package web

import (
	"syscall/js"

	"github.com/VJftw/vigor"
)

// LocalStorage is a signal kept in the browser's localStorage under key: it
// starts as what's stored (or def), and setting it stores the new value, so
// it's still there on the next visit. Use it for what a device remembers
// (display settings, an API token). It falls back to memory only when
// storage isn't available (private windows may refuse it).
func LocalStorage(key, def string) (get func() string, set func(string)) {
	get, setSignal := vigor.NewSignal(readStorage(key, def))

	return get, func(v string) {
		writeStorage(key, v)
		setSignal(v)
	}
}

func storage() js.Value {
	defer func() { _ = recover() }() // the getter itself can throw
	s := js.Global().Get("localStorage")
	if s.IsUndefined() || s.IsNull() {
		return js.Undefined()
	}

	return s
}

func readStorage(key, def string) string {
	s := storage()
	if s.IsUndefined() {
		return def
	}
	var out = def
	func() {
		defer func() { _ = recover() }()
		if v := s.Call("getItem", key); !v.IsNull() {
			out = v.String()
		}
	}()

	return out
}

func writeStorage(key, v string) {
	s := storage()
	if s.IsUndefined() {
		return
	}
	defer func() { _ = recover() }() // quota or privacy errors
	s.Call("setItem", key, v)
}

// RemoveStorage forgets a stored key.
func RemoveStorage(key string) {
	if s := storage(); !s.IsUndefined() {
		func() {
			defer func() { _ = recover() }()
			s.Call("removeItem", key)
		}()
	}
}

// SetTitle sets the document's title, and restores the previous one when
// the surrounding component is disposed. title is a string or a getter of
// one (call it in a component to track it).
func SetTitle(title func() string) {
	doc := js.Global().Get("document")
	prev := doc.Get("title").String()
	vigor.Effect(func() { doc.Set("title", title()) })
	vigor.OnCleanup(func() { doc.Set("title", prev) })
}
