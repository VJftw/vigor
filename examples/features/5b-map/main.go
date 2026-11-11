//go:build js && wasm

/*
vigor.Map is a reactive map for collections that change one entry at a time.
Reading a key subscribes only to that key; Keys, Len and Values subscribe to
entries being added or removed (Values to every value too). So an update
re-runs only what shows that entry.
*/
package main

import (
	"context"
	"strings"

	"github.com/VJftw/vigor"
	"github.com/VJftw/vigor/html"
	"github.com/VJftw/vigor/web"
)

func App() html.Node {
	m := vigor.NewMap[string, int]()
	m.Set("a", 1)
	m.Set("b", 1)
	var aRuns, keysRuns int

	aRunsSig, setARuns := vigor.NewSignal(0)
	keysRunsSig, setKeysRuns := vigor.NewSignal(0)
	vigor.Effect(func() { m.Get("a"); aRuns++; setARuns(aRuns) })
	vigor.Effect(func() { m.Keys(); keysRuns++; setKeysRuns(keysRuns) })

	return html.El("div",
		html.El("p", html.Attr("id", "a"), func() string { v, _ := m.Get("a"); return strings.Repeat("*", v) }),
		html.El("p", html.Attr("id", "len"), m.Len),
		html.El("p", html.Attr("id", "a-runs"), aRunsSig),
		html.El("p", html.Attr("id", "keys-runs"), keysRunsSig),
		html.El("button", html.Attr("id", "bump-a"), html.OnClick(func() { v, _ := m.Get("a"); m.Set("a", v+1) })),
		html.El("button", html.Attr("id", "bump-b"), html.OnClick(func() { v, _ := m.Get("b"); m.Set("b", v+1) })),
		html.El("button", html.Attr("id", "add-c"), html.OnClick(func() { m.Set("c", 1) })),
		html.El("button", html.Attr("id", "replace"), html.OnClick(func() { m.Replace(map[string]int{"c": 1, "d": 1}) })),
	)
}

func main() {
	ctx := context.Background()
	web.Mount(ctx, "app", App)

	<-ctx.Done()
}
