//go:build js && wasm

/*
vigor.Resource is a value loaded asynchronously: Value, Loading, Error and
Loaded are signals, and it fetches again when its source signal changes (or on
Refetch). html.Await shows its three states: pending until it first loads,
failed if it couldn't, and ready, which reads the value reactively so a
refetch updates it without rebuilding.
*/
package main

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/VJftw/vigor"
	"github.com/VJftw/vigor/html"
	"github.com/VJftw/vigor/web"
)

func App() html.Node {
	id, setID := vigor.NewSignal("1")
	var fetches atomic.Int32

	item := vigor.NewResourceFor(id, func(ctx context.Context, id string) (string, error) {
		n := fetches.Add(1)
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		case <-time.After(150 * time.Millisecond):
		}
		if id == "bad" {
			return "", errors.New("no such item")
		}

		return fmt.Sprintf("item-%s#%d", id, n), nil
	})
	broken := vigor.NewResource(func(context.Context) (string, error) { return "", errors.New("down") })

	return html.El("div",
		html.El("div", html.Attr("id", "item"), html.Await(item,
			func() html.Node { return html.El("p", html.Class("pending"), "loading") },
			func(err error) html.Node { return html.El("p", html.Class("failed"), err.Error()) },
			func(v func() string) html.Node { return html.El("p", html.Class("ready"), v) },
		)),
		html.El("p", html.Attr("id", "error"), func() string {
			if err := item.Error(); err != nil {
				return err.Error()
			}

			return "none"
		}),
		html.El("div", html.Attr("id", "broken"), html.Await(broken,
			nil,
			func(err error) html.Node { return html.El("p", html.Class("failed"), err.Error()) },
			func(v func() string) html.Node { return html.El("p", html.Class("ready"), v) },
		)),
		html.El("button", html.Attr("id", "next"), html.OnClick(func() { setID("2") })),
		html.El("button", html.Attr("id", "bad"), html.OnClick(func() { setID("bad") })),
		html.El("button", html.Attr("id", "refetch"), html.OnClick(item.Refetch)),
	)
}

func main() {
	ctx := context.Background()
	web.Mount(ctx, "app", App)

	<-ctx.Done()
}
