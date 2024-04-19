/*
Sometimes you need to deal with conditionals with more than 2 mutual exclusive
outcomes. For this case, we have the <Switch> and <Match> components modeled
roughly after JavaScript's switch/case.

It will try in order to match each condition, stopping to render the first that
evaluates to true. Failing all of them, it will render the fallback.

In the example, we can replace our nested <Show> components with this:

```go
<Switch fallback={<p>{x()} is between 5 and 10</p>}>

	<Match when={x() > 10}>
	  <p>{x()} is greater than 10</p>
	</Match>
	<Match when={5 > x()}>
	  <p>{x()} is less than 5</p>
	</Match>

</Switch>

```
*/
package main

import (
	"context"

	"github.com/VJftw/vigor"
	"github.com/VJftw/vigor/html"
	"github.com/VJftw/vigor/web"
)

func App() html.Node {
	x, _ := vigor.CreateSignal(7)

	return html.Switch(
		html.El("p", x, " is between 5 and 10"),
		html.Case(
			func() bool { return x().(int) > 10 },
			html.El("p", x, " is between 5 and 10"), x,
		),
		html.Case(
			func() bool { return x().(int) < 5 },
			html.El("p", x, " is less than 5"), x,
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
