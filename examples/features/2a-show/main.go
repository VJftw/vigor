//go:build js && wasm

/*
JSX allows you to use JavaScript to control the logic flow in the templates.
However, without a Virtual DOM, using things like Array.prototype.map would
wastefully recreate all the DOM nodes on every update. Instead it is common for
Reactive libraries to use template helpers. In Solid we wrap them in components.

The most basic control flow is the conditional. Solid's compiler is smart enough
to optimally handle ternaries (a ? b : c) and boolean expressions (a && b).
However, often it is more readable to use Solid's <Show> component.

In the example, we would like to show only the appropriate button that reflects
the current state (whether the user is logged in). Update the JSX to:

```go

		<Show
	  when={loggedIn()}
	  fallback={<button onClick={toggle}>Log in</button>}

>

	<button onClick={toggle}>Log out</button>

</Show>
```

The fallback prop acts as the else and will show when the condition passed to
when is not truthy.

Now clicking the button will change back and forth like you would expect.
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
	loggedIn, setLoggedIn := vigor.CreateSignal(false)

	toggle := func(this js.Value, args []js.Value) {
		setLoggedIn(!loggedIn().(bool))
	}

	return html.El("div",
		html.Show(loggedIn,
			html.El("button", html.Text("Log out"), html.On("click", toggle)),
			html.El("button", html.Text("Log in"), html.On("click", toggle)),
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
