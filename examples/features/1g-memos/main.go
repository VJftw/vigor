//go:build js && wasm

/*
Most of the time, composing derived signals is sufficient. However, it is
sometimes beneficial to cache values in order to reduce duplicated work. We can
use memos to evaluate a function and store the result until its dependencies
change. This is great for caching calculations for effects that have other
dependencies and mitigating the work required for expensive operations like DOM
node creation.

Memos are both an observer, like an effect, and a read-only signal. Since they
are aware of both their dependencies and their observers, they can ensure that
they run only once for any change. This makes them preferable to registering
effects that write to signals. Generally, what can be derived, should be
derived.

Creating a Memo is as simple as passing a function that returns a value and its,
dependencies to `vigor.CreateMemo`. In the example, recalculating the value gets
increasingly more expensive with each click. If we wrap it in `vigor.CreateMemo`,
it recalculates only once per click:

```go

	fib := vigor.CreateMemo(func() any {
		return fibonacci(count().(int))
	}, count)

```

Place a `log.Printf` inside the `fib` function to confirm how often it runs:
To update our `count` Signal, we'll attach a click handler on our button:

```go

	fib := vigor.CreateMemo(func() any {
		log.Printf("Calculating Fibonacci")
		return fibonacci(count().(int))
	}, count)

```
*/
package main

import (
	"context"
	"log"
	"syscall/js"

	"github.com/VJftw/vigor"
	"github.com/VJftw/vigor/html"
	"github.com/VJftw/vigor/web"
)

func fibonacci(num int) int {
	if num <= 1 {
		return 1
	}

	return fibonacci(num-1) + fibonacci(num-2)
}

func Counter() html.Node {
	count, setCount := vigor.CreateSignal(10)

	fib := vigor.CreateMemo(func() any {
		log.Printf("Calculating Fibonacci")
		return fibonacci(count().(int))
	}, count)

	return html.El("div",
		html.El("button",
			html.On("click", func(this js.Value, args []js.Value) {
				setCount(count().(int) + 1)
			}),
			"Count: ", count,
		),
		html.El("div", "1. ", fib, " ", fib, " ", fib, " ", fib, " ", fib),
		html.El("div", "2. ", fib, " ", fib, " ", fib, " ", fib, " ", fib),
		html.El("div", "3. ", fib, " ", fib, " ", fib, " ", fib, " ", fib),
		html.El("div", "4. ", fib, " ", fib, " ", fib, " ", fib, " ", fib),
		html.El("div", "5. ", fib, " ", fib, " ", fib, " ", fib, " ", fib),
		html.El("div", "6. ", fib, " ", fib, " ", fib, " ", fib, " ", fib),
		html.El("div", "7. ", fib, " ", fib, " ", fib, " ", fib, " ", fib),
		html.El("div", "8. ", fib, " ", fib, " ", fib, " ", fib, " ", fib),
		html.El("div", "9. ", fib, " ", fib, " ", fib, " ", fib, " ", fib),
		html.El("div", "10. ", fib, " ", fib, " ", fib, " ", fib, " ", fib),
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
