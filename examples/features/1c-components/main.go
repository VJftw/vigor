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
	if err := web.RenderToElementID(
		context.Background(),
		html.El("div",
			html.El("h1", html.Text("This is a Header")),
			Nested(),
		),
		"app",
	); err != nil {
		panic(err)
	}
}
