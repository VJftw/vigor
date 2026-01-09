//go:build js && wasm

/*
Components are just functions like the above `HelloWorld()`. We can also nest
them like so.
*/
package main

import (
	"context"

	"github.com/VJftw/vigor/html"
	"github.com/VJftw/vigor/web"
)

func Nested() html.Node {
	return html.El("p",
		html.Text("This is a Paragraph"),
	)
}

func main() {
	ctx := context.Background()
	web.RenderToElementID(ctx,
		html.El("div",
			html.El("h1", html.Text("This is a Header")),
			Nested(),
		),
		"app",
	)

	<-ctx.Done()
}
