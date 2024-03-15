/*
Signals are the cornerstone of reactivity in Vigor. They contain values that
change over time; when you change a signal's value, it automatically updates
anything that uses it.

To create a signal, let's use `vigor.CreateSignal` and call it from our Counter
component like this:

```go
count, setCount := vigor.CreateSignal(0)
```

The argument passed to the create call is the initial value, and the return
value is an array with two functions, a getter and a setter. By destructuring,
we can name these functions whatever we like. In this case, we name the getter
count and the setter setCount.

It is important to notice that the first returned value is a getter (a function
returning the current value) and not the value itself. This is because the
framework needs to keep track of where that signal is read so it can update
things accordingly.
*/
package main

import (
	"context"
	"time"

	"github.com/VJftw/vigor"
	"github.com/VJftw/vigor/html"
	"github.com/VJftw/vigor/web"
)

func Counter() html.Node {
	count, setCount := vigor.CreateSignal(0)

	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for {

			<-ticker.C
			setCount(count().(int) + 1)
		}
	}()

	return html.El("div",
		html.Text("Count: "), html.Text(count),
	)
}

func main() {
	if err := web.RenderToElementID(
		context.Background(),
		Counter(), "app",
	); err != nil {
		panic(err)
	}
}
