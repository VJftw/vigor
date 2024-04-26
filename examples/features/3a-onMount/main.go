//go:build js && wasm

/*
There are only a few Lifecycle functions in Solid as everything lives and dies by the reactive system. The reactive system is created and updates synchronously, so the only scheduling comes down to Effects which are pushed to the end of the update.

We've found that developers doing basic tasks don't often think this way, so to make things a little easier we've made an alias, onMount. onMount is just a createEffect call that is non-tracking, which means it never re-runs. It is just an Effect call but you can use it with confidence that it will run only once for your component, once all initial rendering is done.

Let's use the onMount hook to fetch some photos:

Lifecycles are only run in the browser, so putting code in onMount has the benefit of not running on the server during SSR. Even though we are doing data fetching in this example, usually we use Solid's resources for true server/browser coordination
*/
package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/VJftw/vigor"
	"github.com/VJftw/vigor/html"
	"github.com/VJftw/vigor/web"
)

func App() html.Node {
	photos, setPhotos := vigor.CreateSignal([]any{})

	type photo struct {
		ThumbnailUrl string `json:"thumbnailUrl"`
		Title        string `json:"title"`
	}

	n := html.El("div",
		html.El("h1", "Photo album"),
		html.El("div", html.Class("photos"),
			html.For(photos, func(i int, v any) html.Node {
				p := v.(photo)

				return html.El("figure",
					html.El("img", html.Property("src", p.ThumbnailUrl), html.Attr("alt", p.Title)),
					html.El("figcaption", p.Title),
				)
			}),
		),
	)

	return html.OnMount(n, func() {
		resp, err := http.Get("https://jsonplaceholder.typicode.com/photos?_limit=20")
		if err != nil {
			log.Printf("ERROR %s", err)
		}
		fetchedPhotos := []photo{}
		if err := json.NewDecoder(resp.Body).Decode(&fetchedPhotos); err != nil {
			log.Printf("ERROR %s", err)
		}

		setPhotos(fetchedPhotos)
	})
}

func main() {
	if err := web.RenderToElementID(
		context.Background(),
		App(), "app",
	); err != nil {
		panic(err)
	}
}
