/*
Go's closest alternative to JSX right now is through the use of the
variadic (`...`) argument to functions. [margudk/gomponents](https://github.com/maragudk/gomponents)
is an excellent example of this in practice and an inspiration for the
implementation here. SolidJS provides a [Hyperscript](https://github.com/solidjs/solid/tree/main/packages/solid/h)
which looks a bit like this too.

TODO: add support for SVG, and maybe plain/raw HTML with templating (templating
in a separate package though).
*/
package main

import (
	"context"

	"github.com/VJftw/vigor/html"
	"github.com/VJftw/vigor/web"
)

func main() {
	if err := web.RenderToElementID(
		context.Background(),
		// TODO: Add SVG support.
		html.El("div"),
		"app",
	); err != nil {
		panic(err)
	}
}
