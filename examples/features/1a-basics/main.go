//go:build js && wasm

/*
Just like in SolidJS, a Vigor App is composed of functions that we call
components. The entry point for Vigor App is the `web.RenderTo()` and
`web.RenderToElementID()` functions.
*/
package main

import (
	"context"

	"github.com/VJftw/vigor/html"
	"github.com/VJftw/vigor/web"
)

func HelloWorld() html.Node {
	return html.El("div",
		html.Text("Hello Solid World!"),
	)
}

func main() {
	ctx := context.Background()

	web.RenderToElementID(
		ctx,
		HelloWorld(),
		"app",
	)

	<-ctx.Done()
}
