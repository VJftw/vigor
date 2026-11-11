//go:build js && wasm

/*
The router package maps the address bar to views: routes with parameters
(/items/:id), links that navigate without reloading, the back button, and the
query string as signals. A view is built when navigated to and disposed when
navigated away from; when only a parameter changes, the view stays and what
reads the parameter updates.
*/
package main

import (
	"context"

	"github.com/VJftw/vigor/html"
	"github.com/VJftw/vigor/router"
	"github.com/VJftw/vigor/web"
)

func App() html.Node {
	builds := 0
	sort := router.QueryParam("sort")

	return html.El("div",
		html.El("nav",
			router.A("/", html.Attr("id", "to-home"), "home"),
			router.A("/items/1", html.Attr("id", "to-1"), "one"),
			router.A("/items/2", html.Attr("id", "to-2"), "two"),
		),
		html.El("p", html.Attr("id", "path"), router.Path),
		html.El("p", html.Attr("id", "sort"), sort),
		html.El("button", html.Attr("id", "sort-asc"), "sort", html.OnClick(func() { router.SetQuery(map[string]string{"sort": "asc"}) })),
		html.El("main", router.Routes(
			func() html.Node { return html.El("p", html.Attr("id", "view"), "not found") },
			router.Route{Path: "/", View: func(func(string) string) html.Node { return html.El("p", html.Attr("id", "view"), "home") }},
			router.Route{Path: "/items/:id", View: func(param func(string) string) html.Node {
				builds++
				n := builds

				return html.El("p", html.Attr("id", "view"), "item ", func() string { return param("id") }, " (built ", n, ")")
			}},
		)),
	)
}

func main() {
	ctx := context.Background()
	web.Mount(ctx, "app", App)

	<-ctx.Done()
}
