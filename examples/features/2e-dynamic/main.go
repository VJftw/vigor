//go:build js && wasm

/*
The <Dynamic> tag is useful when you render from data. It lets you pass either a string for a native element or a component function and it will render that with the rest of the provided props.

This is often more compact than writing a number of <Show> or <Switch> components.
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
	selected, setSelected := vigor.CreateSignal("red")

	options := map[string]html.Node{
		"red":   html.El("strong", html.Style("color", "red"), "Red Thing"),
		"green": html.El("strong", html.Style("color", "green"), "Green Thing"),
		"blue":  html.El("strong", html.Style("color", "blue"), "Blue Thing"),
	}
	optionKeys := []any{"red", "green", "blue"} // static order

	return html.El("div",
		html.El("select",
			html.Property("value", selected),
			html.On("input", func(this js.Value, args []js.Value) {
				setSelected(args[0].Get("currentTarget").Get("value").String())
			}),
			html.For(optionKeys, func(i int, v any) html.Node {
				return html.El("option", html.Attr("value", v), html.Text(v))
			}),
		),
		html.Dynamic(func(s vigor.Subscriber) html.Node {
			return options[selected(s).(string)]
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
