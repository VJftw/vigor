//go:build js && wasm

/*
Scopes: what a component creates (effects, intervals, OnCleanup callbacks) is
disposed when the part of the page it's in goes away. html.When builds a branch
when it's shown and disposes it when it's hidden, and only when its condition
flips, not whenever what it reads changes. vigor.CreateContext passes values
down to components without threading them through every call.
*/
package main

import (
	"context"
	"time"

	"github.com/VJftw/vigor"
	"github.com/VJftw/vigor/html"
	"github.com/VJftw/vigor/web"
)

var theme = vigor.CreateContext("light")

func Badge() html.Node { return html.El("span", html.Class("badge"), theme.Use()) }

func App() html.Node {
	count, setCount := vigor.NewSignal(0)
	builds, setBuilds := vigor.NewSignal(0)
	cleanups, setCleanups := vigor.NewSignal(0)
	ticks, setTicks := vigor.NewSignal(0)

	var inside html.Node
	theme.Provide("dark", func() { inside = Badge() })

	return html.El("div",
		html.El("p", html.Attr("id", "outside"), Badge()),
		html.El("p", html.Attr("id", "inside"), inside),
		html.When(func() bool { return count() > 0 }, func() html.Node {
			setBuilds(vigor.Untrack(builds) + 1)
			vigor.OnCleanup(func() { setCleanups(vigor.Untrack(cleanups) + 1) })
			vigor.Interval(20*time.Millisecond, func() { setTicks(vigor.Untrack(ticks) + 1) })

			return html.El("p", html.Attr("id", "shown"), "count is ", count)
		}),
		html.El("p", html.Attr("id", "builds"), builds),
		html.El("p", html.Attr("id", "cleanups"), cleanups),
		html.El("p", html.Attr("id", "ticks"), ticks),
		html.El("button", html.Attr("id", "inc"), html.OnClick(func() { setCount(vigor.Untrack(count) + 1) })),
		html.El("button", html.Attr("id", "reset"), html.OnClick(func() { setCount(0) })),
	)
}

func main() {
	ctx := context.Background()
	web.Mount(ctx, "app", App)

	<-ctx.Done()
}
