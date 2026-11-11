//go:build js && wasm

/*
html.Each renders a list keyed by something stable. When the list changes, rows
whose key is still there keep their DOM nodes (they're moved, not rebuilt),
new keys are rendered, and removed keys are disposed: their effects stop and
their OnCleanup callbacks run. A fallback shows when the list is empty.
*/
package main

import (
	"context"

	"github.com/VJftw/vigor"
	"github.com/VJftw/vigor/html"
	"github.com/VJftw/vigor/web"
)

func App() html.Node {
	items, setItems := vigor.NewSignal([]string{"a", "b", "c"})
	cleaned, setCleaned := vigor.NewSignal(0)

	return html.El("div",
		html.El("ul", html.Attr("id", "list"),
			html.Each(items, func(s string) string { return s },
				func(item func() string, index func() int) html.Node {
					vigor.OnCleanup(func() { setCleaned(vigor.Untrack(cleaned) + 1) })

					return html.El("li", html.Attr("data-key", item), index, ": ", item)
				},
				func() html.Node { return html.El("li", html.Class("empty"), "nothing") },
			),
		),
		html.El("p", html.Attr("id", "cleaned"), cleaned),
		html.El("button", html.Attr("id", "reverse"), "reverse", html.OnClick(func() {
			cur := vigor.Untrack(items)
			next := make([]string, len(cur))
			for i, s := range cur {
				next[len(cur)-1-i] = s
			}
			setItems(next)
		})),
		html.El("button", html.Attr("id", "drop-first"), "drop first", html.OnClick(func() { setItems(vigor.Untrack(items)[1:]) })),
		html.El("button", html.Attr("id", "clear"), "clear", html.OnClick(func() { setItems(nil) })),
	)
}

func main() {
	ctx := context.Background()
	web.Mount(ctx, "app", App)

	<-ctx.Done()
}
