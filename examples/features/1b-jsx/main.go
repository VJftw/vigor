/*
Go's closest alternative to JSX right now is through the use of the
variadic (`...`) argument to functions. [margudk/gomponents](https://github.com/maragudk/gomponents)
is an excellent example of this in practice and an inspiration for the
implementation here. SolidJS provides a [Hyperscript](https://github.com/solidjs/solid/tree/main/packages/solid/h)
which looks a bit like this too.
*/
package main

import (
	"context"

	"github.com/VJftw/vigor/html"
	"github.com/VJftw/vigor/web"
)

func HelloWorld() html.Node {
	name := "Solid"
	svgNs := "http://www.w3.org/2000/svg"
	svg := html.NEl(svgNs, "svg",
		html.Attr("height", 300),
		html.Attr("width", 400),
		html.NEl(svgNs, "defs",
			html.NEl(svgNs, "linearGradient", html.Attr("id", "gr1"), html.Attr("x1", "0%"), html.Attr("y1", "60%"), html.Attr("x2", "100%"), html.Attr("y2", "0%"),
				html.NEl(svgNs, "stop", html.Attr("offset", "5%"), html.Style("stop-color", "rgb(255,255,3)"), html.Style("stop-opacity", "1")),
				html.NEl(svgNs, "stop", html.Attr("offset", "100%"), html.Style("stop-color", "rgb(255,0,0)"), html.Style("stop-opacity", "1")),
			),
		),
		html.NEl(svgNs, "ellipse", html.Attr("cx", "125"), html.Attr("cy", "150"), html.Attr("rx", "100"), html.Attr("ry", "60"), html.Attr("fill", "url(#gr1)")),
		"Sorry but this browser does not support inline SVG.",
	)

	return html.El("div",
		html.El("div", "Hello ", name, "!"),
		svg,
	)
}

func main() {
	if err := web.RenderToElementID(
		context.Background(),
		HelloWorld(),
		"app",
	); err != nil {
		panic(err)
	}
}
