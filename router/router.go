//go:build js && wasm

package router

import (
	"context"
	"net/url"
	"syscall/js"

	"github.com/VJftw/vigor"
	"github.com/VJftw/vigor/html"
)

// The browser has one location, so the router keeps it as package state:
// signals that follow the address bar (and back and forward buttons).
var (
	location    func() Location
	setLocation func(Location)
)

func init() {
	win := js.Global()
	location, setLocation = vigor.NewSignal(current(), vigor.DeepEqual[Location]())
	win.Call("addEventListener", "popstate", js.FuncOf(func(js.Value, []js.Value) any {
		setLocation(current())

		return nil
	}))
}

func current() Location {
	l := js.Global().Get("location")

	return ParseLocation(l.Get("pathname").String() + l.Get("search").String() + l.Get("hash").String())
}

// Path is the current path, e.g. "/p/automations". Reading it in an effect
// or component subscribes to navigation.
func Path() string { return location().Path }

// Query is the current query string's values.
func Query() url.Values { return location().Query }

// QueryParam is a getter of one query parameter ("" when absent), which
// only notifies when that parameter changes.
func QueryParam(name string) func() string {
	return vigor.Memo(func() string { return location().Query.Get(name) })
}

// Hash is the current fragment, without the "#".
func Hash() string { return location().Hash }

// Navigate goes to a URL (a path with optional query), adding a history
// entry, without reloading the page.
func Navigate(to string) { move(to, false) }

// Replace is Navigate without adding a history entry (for settings that
// live in the URL).
func Replace(to string) { move(to, true) }

// Back goes back in history.
func Back() { js.Global().Get("history").Call("back") }

func move(to string, replace bool) {
	l := ParseLocation(to)
	method := "pushState"
	if replace {
		method = "replaceState"
	}
	js.Global().Get("history").Call(method, nil, "", l.String())
	setLocation(l)
}

// SetQuery changes query parameters and keeps the path, replacing the
// current history entry. An empty value removes the parameter.
func SetQuery(params map[string]string) {
	l := vigor.Untrack(location)
	q := url.Values{}
	for k, v := range l.Query {
		q[k] = v
	}
	for k, v := range params {
		if v == "" {
			q.Del(k)
		} else {
			q.Set(k, v)
		}
	}
	l.Query = q
	move(l.String(), true)
}

// Route maps a path pattern to the view shown for it.
type Route struct {
	// Path is a pattern: "/", "/p/:id", "/files/*path".
	Path string
	// View builds the page. param reads a parameter of the matched path
	// reactively, so when only a parameter changes (/p/a to /p/b) the view
	// stays and what reads it updates.
	View func(param func(name string) string) html.Node
}

// Routes shows the view of the first route matching the current path, or
// notFound (if given) when none does. Views are built when navigated to and
// disposed when navigated away from.
func Routes(notFound func() html.Node, routes ...Route) html.Node {
	patterns := make([]Pattern, len(routes))
	for i, r := range routes {
		patterns[i] = ParsePattern(r.Path)
	}

	return html.NewNode(func(ctx context.Context, doc js.Value) js.Value {
		match := vigor.Memo(func() struct {
			index  int
			params Params
		} {
			path := Path()
			for i, p := range patterns {
				if params, ok := p.Match(path); ok {
					return struct {
						index  int
						params Params
					}{i, params}
				}
			}

			return struct {
				index  int
				params Params
			}{-1, nil}
		}, vigor.DeepEqual[struct {
			index  int
			params Params
		}]())
		index := vigor.Memo(func() int { return match().index })

		return html.Reactive(func() html.Node {
			i := index()
			if i < 0 {
				if notFound == nil {
					return nil
				}

				return vigor.Untrack(notFound)
			}
			param := func(name string) string { return match().params[name] }

			return vigor.Untrack(func() html.Node { return routes[i].View(param) })
		}).DOMObject(ctx, doc)
	})
}

// A is a link to a path inside the app: clicking it navigates without
// reloading. Modified clicks (new tab, etc.) and links with a target are
// left to the browser. Give it attributes and children like an element.
func A(href string, children ...any) html.Node {
	opts := append([]any{html.Attr("href", href)}, children...)
	opts = append(opts, html.On("click", func(this js.Value, args []js.Value) {
		if len(args) == 0 {
			return
		}
		ev := args[0]
		if ev.Get("defaultPrevented").Bool() || ev.Get("button").Int() != 0 ||
			ev.Get("metaKey").Bool() || ev.Get("ctrlKey").Bool() ||
			ev.Get("shiftKey").Bool() || ev.Get("altKey").Bool() {
			return
		}
		if t := this.Call("getAttribute", "target"); !t.IsNull() && t.String() != "" && t.String() != "_self" {
			return
		}
		ev.Call("preventDefault")
		Navigate(href)
	}))

	return html.El("a", opts...)
}
