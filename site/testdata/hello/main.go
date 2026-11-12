//go:build js && wasm

package main

import (
	"context"

	"github.com/VJftw/vigor/html"
	"github.com/VJftw/vigor/web"
)

func main() {
	ctx := context.Background()
	web.Mount(ctx, "app", func() html.Node { return html.El("h1", "hello") })

	<-ctx.Done()
}
