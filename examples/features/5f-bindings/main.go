//go:build js && wasm

/*
Elements take plain getters (func() T) wherever they take values, and update
just that part when the signals inside change. Form helpers bind inputs to
signals both ways (BindValue, BindChecked); a bool attribute is present or
absent; nil removes an attribute; ClassIf toggles a class; Style sets CSS,
including custom properties.
*/
package main

import (
	"context"
	"fmt"

	"github.com/VJftw/vigor"
	"github.com/VJftw/vigor/html"
	"github.com/VJftw/vigor/web"
)

func App() html.Node {
	text, setText := vigor.NewSignal("hi")
	armed, setArmed := vigor.NewSignal(false)

	return html.El("div",
		html.El("input", html.Attr("id", "text"), html.BindValue(text, setText)),
		html.El("p", html.Attr("id", "echo"), "you typed: ", text),
		html.El("input", html.Attr("id", "armed"), html.Attr("type", "checkbox"), html.BindChecked(armed, setArmed)),
		html.El("button", html.Attr("id", "go"), html.Attr("disabled", func() bool { return !armed() }), "go"),
		html.El("p", html.Attr("id", "box"),
			html.ClassIf("on", armed),
			html.Attr("title", func() any {
				if armed() {
					return "armed"
				}

				return nil
			}),
			html.Style("--width", func() string { return fmt.Sprintf("%dpx", len(text())) }),
			"box",
		),
	)
}

func main() {
	ctx := context.Background()
	web.Mount(ctx, "app", App)

	<-ctx.Done()
}
