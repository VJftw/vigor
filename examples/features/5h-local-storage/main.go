//go:build js && wasm

/*
web.LocalStorage is a signal kept in the browser's localStorage: it starts as
what's stored (or a default), and setting it stores the new value, so it's
still there on the next visit.
*/
package main

import (
	"context"

	"github.com/VJftw/vigor/html"
	"github.com/VJftw/vigor/web"
)

func App() html.Node {
	name, setName := web.LocalStorage("vigor-example-name", "anonymous")

	return html.El("div",
		html.El("input", html.Attr("id", "name"), html.BindValue(name, setName)),
		html.El("p", html.Attr("id", "hello"), "hello ", name),
	)
}

func main() {
	ctx := context.Background()
	web.Mount(ctx, "app", App)

	<-ctx.Done()
}
