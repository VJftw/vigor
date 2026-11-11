//go:build js && wasm

/*
vigor.NewStream feeds a signal from a long-lived connection and reconnects,
with backoff, when it drops. Latest is the last value, Connected says whether
the run function is running, and Error is why the last one ended. OnConnect
runs at each (re)connection, e.g. to clear state a fresh snapshot refills.
*/
package main

import (
	"context"
	"errors"
	"time"

	"github.com/VJftw/vigor"
	"github.com/VJftw/vigor/html"
	"github.com/VJftw/vigor/web"
)

func App() html.Node {
	connects, setConnects := vigor.NewSignal(0)
	attempt := 0

	stream := vigor.NewStream(func(ctx context.Context, emit func(int)) error {
		attempt++
		emit(attempt * 10)
		if attempt < 3 {
			return errors.New("dropped") // reconnects
		}
		<-ctx.Done() // the third stays up

		return ctx.Err()
	},
		vigor.WithReconnectDelay(10*time.Millisecond, 50*time.Millisecond),
		vigor.OnConnect(func() { setConnects(vigor.Untrack(connects) + 1) }),
	)

	return html.El("div",
		html.El("p", html.Attr("id", "latest"), stream.Latest),
		html.El("p", html.Attr("id", "connected"), func() bool { return stream.Connected() }),
		html.El("p", html.Attr("id", "error"), func() string {
			if err := stream.Error(); err != nil {
				return err.Error()
			}

			return "none"
		}),
		html.El("p", html.Attr("id", "connects"), connects),
	)
}

func main() {
	ctx := context.Background()
	web.Mount(ctx, "app", App)

	<-ctx.Done()
}
