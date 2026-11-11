//go:build js && wasm

/*
html.Ref calls a function with an element once it's built, for what signals
can't express: focusing it, measuring it, calling a method on it.
*/
package main

import (
	"context"
	"syscall/js"

	"github.com/VJftw/vigor/html"
	"github.com/VJftw/vigor/web"
)

func App() html.Node {
	var input js.Value

	return html.El("div",
		html.El("input", html.Attr("id", "field"), html.Ref(func(el js.Value) { input = el })),
		html.El("button", html.Attr("id", "focus"), "focus", html.OnClick(func() { input.Call("focus") })),
	)
}

func main() {
	ctx := context.Background()
	web.Mount(ctx, "app", App)

	<-ctx.Done()
}
