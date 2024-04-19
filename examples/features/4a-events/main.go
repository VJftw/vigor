/*
Events in Solid are attributes prefixed with on. They are treated specially in a
few ways. First, they do not follow the normal heuristics for wrapping. In many
cases, it is difficult to determine the difference between a Signal and an event
handler. And so, since events are called and don't require reactivity to update,
they are only bound initially. You can always just have your handler run
different code based on the current state of your app.

Common UI events (that bubble and are composed) are automatically delegated to
the document. To improve delegated performance, Solid supports an array syntax
to call the handler with data (as the first argument) without creating
additional closures:

```go

	const handler = (data, event) => ...

	<button onClick={[handler, data]}>Click Me</button>

```

In the example, let's attach the handler to the mousemove event:

```go

<div onMouseMove={handleMouseMove}>

	The mouse position is {pos().x} x {pos().y}

</div>

```

All on bindings are case insensitive which means that event names need to be in
lowercase. For example, onMouseMove monitors the event name mousemove. If you
need to support other casings or not use event delegation, you can use on:
namespace to match event handlers that follows the colon:

```go
<button on:DOMContentLoaded={() => ...} >Click Me</button>
```
*/
package main

import (
	"context"
	"syscall/js"

	"github.com/VJftw/vigor"
	"github.com/VJftw/vigor/html"
	"github.com/VJftw/vigor/web"
)

func App() html.Node {
	// pos, setPos := vigor.CreateSignal([]any{0, 0})
	type coords struct {
		x int
		y int
	}
	pos, setPos := vigor.CreateSignal(coords{x: 0, y: 0})

	handleMouseEvent := func(this js.Value, args []js.Value) {
		event := args[0]
		setPos(coords{x: event.Get("clientX").Int(), y: event.Get("clientY").Int()})
	}

	return html.El("div",
		html.On("mousemove", handleMouseEvent),
		html.Text("The mouse position is: "), html.Text(
			html.SubscribedText(func(s vigor.Subscriber) any { return pos(s).(coords).x }),
			" x ",
			html.SubscribedText(func(s vigor.Subscriber) any { return pos(s).(coords).y }),
		),
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
