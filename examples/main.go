package main

import (
	"time"

	"github.com/VJftw/vigor"
	"github.com/VJftw/vigor/html"
	"github.com/VJftw/vigor/web"
)

func main() {
	a, setA := vigor.CreateSignal("a")
	bgClass, setBgClass := vigor.CreateSignal("bg-primary")
	_ = setA
	_ = bgClass

	whenVal, setWhenVal := vigor.CreateSignal(false)
	_ = whenVal

	sliceVal, setSliceVal := vigor.CreateSignal([]any{"Hello", "World"})

	root := html.El("p",
		html.El("div",
			html.Class(bgClass),
			html.El("p",
				html.Text(a),
				html.Portal(html.Text("p a")),
			),
			html.Show(whenVal,
				html.El("div", html.Text("Yee")),
				html.El("div", html.El("p", html.Text("Haw"))),
			),
			html.El("ul",
				html.For(sliceVal, func(v any) html.Node {
					return html.El("li", html.Text(v))
				}),
			),
			html.Switch(
				html.Text("default"),
				html.Case(func() bool {
					return bgClass() == "bg-secondary"
				}, html.Text("1"), bgClass),
				html.Case(func() bool {
					return bgClass() == "bg-primary"
				}, html.Text("5"), bgClass),
			),
			html.Portal(html.Text("another portal")),
		),
	)

	if err := web.RenderToElementID(root, "app"); err != nil {
		panic(err)
	}

	for i := 0; i < 10; i++ {
		time.Sleep(1 * time.Second)
		setA(a().(string) + "b")
		if i%2 == 0 {
			setBgClass("bg-secondary")
			setWhenVal(false)
		} else {
			setBgClass("bg-primary")
			setWhenVal(true)
		}
	}

	setSliceVal([]any{"End", "World", "Again"})
}
