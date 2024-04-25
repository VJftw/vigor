//go:build js && wasm

/*
The <For> component is the best way to loop over an array of objects. As the
array changes, <For> updates or moves items in the DOM rather than recreating
them. Let's look at an example.

```go
<For each={cats()}>{(cat, i) =>

	  <li>
	    <a target="_blank" href={`https://www.youtube.com/watch?v=${cat.id}`}>
	      {i() + 1}: {cat.name}
	    </a>
	  </li>
	}</For>

```

There is one prop on the <For> component: each, where you pass the array to loop
over.

Then, instead of writing nodes directly between <For> and </For>, you pass a
callback. This is a function similar to JavaScript's map callback. For each
element in the array, the callback is called with the element as the first
argument and the index as the second. (cat and i in this example.) You can then
make use of those in the callback, which should return a node to be rendered.

Note that the index is a signal, not a constant number. This is because <For> is
"keyed by reference": each node that it renders is coupled to an element in the
array. In other words, if an element changes placement in the array, rather than
being destroyed and recreated, the corresponding node will move too and its
index will change.

The each prop expects an array, but you can turn other iterable objects into
arrays with utilities like Array.from, Object.keys, or spread syntax.
*/
package main

import (
	"context"

	"github.com/VJftw/vigor"
	"github.com/VJftw/vigor/html"
	"github.com/VJftw/vigor/web"
)

func App() html.Node {
	type idToCat struct {
		id   string
		name string
	}
	cats, _ := vigor.CreateSignal([]any{
		idToCat{id: "J---aiyznGQ", name: "Keyboard Cat"},
		idToCat{id: "z_AbfPXTKms", name: "Maru"},
		idToCat{id: "OUtn3pvWmpg", name: "Henri The Existential Cat"},
	})

	return html.El("ul",
		html.For(cats, func(i int, v any) html.Node {
			cat := v.(idToCat)
			return html.El("li",
				html.El("a",
					html.Attr("target", "_blank"),
					html.Attr("href", "https://www.youtube.com/watch?v="+cat.id),
					html.Text(i+1, ": ", cat.name),
				),
			)
		}),
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
