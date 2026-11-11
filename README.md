# vigor

This package is a **highly experimental** interpretation and implementation of [SolidJS](https://www.solidjs.com/) tenets in Golang for use in building reactive web interfaces powered by Web Assembly (WASM).

## TODO

- Restructure to new Docs https://docs.solidjs.com/

## Building real apps

Beyond static reactive nodes, a live app needs state that follows a stream,
pages that come and go without leaking, routing, async data and forms. What
vigor offers for each, and the example that shows it:

| Need | vigor | Example |
| ---- | ----- | ------- |
| State without dependency lists | `vigor.NewSignal`, `Memo`, `Effect` track what they read; `Batch`, `Untrack`, `On` | `reactive_test.go` |
| Pages that come and go without leaking | Scopes: `OnCleanup`, `CreateRoot`, `OwnerContext` (a `context.Context` cancelled on disposal); effects and nested effects are disposed on re-run | `5g-context-cleanup` |
| Passing values down (clients, themes) | `vigor.CreateContext[T]` → `Provide` / `Use` | `5g-context-cleanup` |
| Collections that change one entry at a time | `vigor.Map[K,V]`: reading a key subscribes to that key only; `Keys()` to membership | `5b-map` |
| Lists that keep their rows | `html.Each(items, key, render, fallback)`: keyed, rows keep their DOM and state, reorder in place | `5a-each` |
| A live connection that reconnects | `vigor.NewStream` (backoff, `OnConnect`, `StopOn`) | `5e-stream` |
| Fetching data, polling | `vigor.NewResource` / `NewResourceFor(source, fetch)`, `WithPoll`; `html.Await` for pending / failed / ready | `5d-resource` |
| Clocks and timers | `vigor.Clock(time.Second)`, `Interval`, `Timeout`, all stopped with their scope | `5g-context-cleanup` |
| Routes, links, query parameters | package `router`: `Routes`, `A`, `Path`, `QueryParam`, `SetQuery`, `Navigate`, back button | `5c-router` |
| Remembered settings | `web.LocalStorage(key, default)` | `5h-local-storage` |
| Conditional UI | `html.When`, `html.Reactive` (build only what's shown, dispose the rest) | `5g-context-cleanup` |
| Forms and attributes | `html.BindValue`, `BindChecked`, `OnInput`, `OnChange`; boolean attributes (`Attr("disabled", cond)`), `ClassIf`, CSS custom properties via `Style("--x", v)` | `5f-bindings` |
| Touch and keyboard handling | `html.OnWith` (passive or not), `OnWindow`, `OnDocument` (removed with the component), `PreventDefault`, `Ref` | `5i-ref` |
| SVG, icons, trusted markup | `html.Raw`, and `html.NEl` for SVG elements | |
| Calling an API | `net/http` already works under WebAssembly (fetch); results come back through `vigor.Dispatch` (Resource and Stream do it for you) | |

Everything takes plain Go functions as getters: anywhere an element accepts a
value (attributes, classes, styles, text, children) it also accepts a
`func() T`, and re-renders just that part when the signals inside change.

The older callback-style API (`CreateSignal`, `CreateEffect`, `Show`, `For`...)
still works alongside.

Signals are not goroutine-safe; in WebAssembly there is one thread, and
goroutines that report back (timers, fetches, streams) must use
`vigor.Dispatch`.
