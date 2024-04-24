//go:build js && wasm

/*
Signals are trackable values, but they are only one half of the equation. To
complement those are observers that can be updated by those trackable values. An
effect is one such observer; it runs a side effect that depends on signals.

An effect can be created by using `vigor.CreateEffect` and providing it a
function as well as the signals which that functions depends on. The effect
is subscribed to all of the dependencies and reruns when any of them change.

So let's create an Effect that reruns whenever count changes:
Signals are the cornerstone of reactivity in Vigor. They contain values that
change over time; when you change a signal's value, it automatically updates
anything that uses it.

So let's create an Effect that reruns whenever `count` changes.

```go

	vigor.CreateEffect(func() {
		log.Printf("The count is now: %s", count())
	}, count)

```

To update our `count` Signal, we'll attach a click handler on our button:

```go

	html.El("button",
		html.On("click", func(this js.Value, args []js.Value) any {
			setCount(count().(int) + 1)
			return nil
		}),
		html.Text("Click Me"),
	)

```

Now clicking the button writes to the console. This is a relatively simple
example, but to understand how Vigor works, you should imagine that every
expression in JSX is its own effect that re-executes whenever its dependent
signals change. This is how all rendering works in Vigor: from Vigor's
perspective, all rendering is just a side effect of the reactive system.
*/
package main

import (
	"context"
	"log"
	"syscall/js"

	"github.com/VJftw/vigor"
	"github.com/VJftw/vigor/html"
	"github.com/VJftw/vigor/web"
)

func Counter() html.Node {
	count, setCount := vigor.CreateSignal(0)

	vigor.CreateEffect(func() {
		log.Printf("The count is now %d", count())
	}, count)

	return html.El("button",
		html.On("click", func(this js.Value, args []js.Value) {
			setCount(count().(int) + 1)
		}),
		html.Text("Click Me"),
	)
}

func main() {
	if err := web.RenderToElementID(
		context.Background(),
		Counter(), "app",
	); err != nil {
		panic(err)
	}
}
