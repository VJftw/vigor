//go:build js && wasm

/*
A JavaScript error originating in the UI shouldn’t break the whole app. Error boundaries are components that catch JavaScript errors anywhere in their child component tree, log those errors, and display a fallback UI instead of the component tree that crashed.

A component has crashed our example. Let's wrap it in an Error Boundary that displays the error.
*/
package main

import (
	"context"
	"syscall/js"

	"github.com/VJftw/vigor/html"
	"github.com/VJftw/vigor/web"
)

func Broken() html.Node {
	return html.NewNode(func(v js.Value) js.Value {
		panic("Oh No")
	})
}

func App() html.Node {
	return html.El("div",
		html.El("div", "Before"),
		html.ErrorBoundary(
			func(err any) html.Node {
				return html.El("div", "Error: ", err)
			},
			Broken(),
		),
		html.El("div", "After"),
	)
}

func main() {
	if err := web.RenderToElementID(
		context.Background(),
		App(), "app",
	); err != nil {
		panic(err)
	}
}
