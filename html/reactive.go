//go:build js && wasm

package html

import (
	"context"
	"fmt"
	"reflect"
	"syscall/js"

	"github.com/VJftw/vigor"
)

// accessorType reports whether v is a zero-argument function with one
// result, i.e. a getter from vigor.NewSignal or vigor.Memo, or any func()
// T. Reading it inside an effect subscribes to what it reads.
func accessorType(v any) (reflect.Value, bool) {
	if v == nil {
		return reflect.Value{}, false
	}
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Func || rv.IsNil() {
		return reflect.Value{}, false
	}
	t := rv.Type()
	if t.NumIn() != 0 || t.NumOut() != 1 || t.IsVariadic() {
		return reflect.Value{}, false
	}

	return rv, true
}

// isReactive reports whether v changes over time: a getter of either kind.
func isReactive(v any) bool {
	switch v.(type) {
	case vigor.GetterFn, *subscribedText:
		return true
	}
	_, ok := accessorType(v)

	return ok
}

// value reads v: getters (legacy ones subscribing sub, new ones being
// tracked by the running effect) are called, anything else is itself.
func value(v any, sub vigor.Subscriber) any {
	switch x := v.(type) {
	case vigor.GetterFn:
		return x(sub)
	case *subscribedText:
		return value(x.fn(sub), sub)
	case classCond:
		if truthy(value(x.cond, sub)) {
			return x.name
		}

		return ""
	}
	if rv, ok := accessorType(v); ok {
		return rv.Call(nil)[0].Interface()
	}

	return v
}

func truthy(v any) bool {
	switch x := v.(type) {
	case nil:
		return false
	case bool:
		return x
	case string:
		return x != ""
	}
	rv := reflect.ValueOf(v)
	if rv.CanInt() {
		return rv.Int() != 0
	}
	if rv.CanFloat() {
		return rv.Float() != 0
	}
	switch rv.Kind() {
	case reflect.Slice, reflect.Map:
		return rv.Len() > 0
	case reflect.Pointer, reflect.Interface:
		return !rv.IsNil()
	}

	return true
}

// bind runs apply now and again whenever a getter it reads changes, whether
// it's a legacy GetterFn (subscribed through sub) or a vigor.NewSignal
// accessor (tracked by the effect).
func bind(apply func(sub vigor.Subscriber)) {
	sub := vigor.NewFnSubscriber()
	sub.SetFn(func() { apply(sub) })
	vigor.Effect(func() { apply(sub) })
}

// text formats a value as text: nothing for nil, otherwise %v.
func text(v any) string {
	if v == nil {
		return ""
	}

	return fmt.Sprintf("%v", v)
}

// childNode makes a Node of anything allowed as an element's child.
func childNode(c any) Node {
	switch x := c.(type) {
	case nil:
		return nil
	case Node:
		return x
	}

	return Text(c)
}

// nodesOf returns the nodes a rendered value consists of: a fragment's
// children (a fragment is emptied when inserted, so they have to be noted
// first), or the value itself.
func nodesOf(v js.Value) []js.Value {
	if v.IsUndefined() || v.IsNull() {
		return nil
	}
	if v.Get("nodeType").Int() == 11 {
		kids := v.Get("childNodes")
		out := make([]js.Value, kids.Length())
		for i := range out {
			out[i] = kids.Index(i)
		}

		return out
	}

	return []js.Value{v}
}

func clearBetween(start, end js.Value) {
	for {
		n := start.Get("nextSibling")
		if n.IsNull() || n.Equal(end) {
			return
		}
		n.Call("remove")
	}
}

// Fragment groups nodes without a wrapping element. Strings, numbers and
// getters are allowed, as in El.
func Fragment(children ...any) Node {
	return NewNode(func(ctx context.Context, doc js.Value) js.Value {
		frag := doc.Call("createDocumentFragment")
		for _, c := range children {
			if n := childNode(c); n != nil {
				if o := n.DOMObject(ctx, doc); !o.IsUndefined() {
					frag.Call("appendChild", o)
				}
			}
		}

		return frag
	})
}

// Reactive shows whatever fn returns, calling it again whenever a signal it
// reads changes. What the previous call created is disposed first (effects
// stopped, OnCleanup callbacks run). fn may return nil for nothing.
func Reactive(fn func() Node) Node {
	return NewNode(func(ctx context.Context, doc js.Value) js.Value {
		frag := doc.Call("createDocumentFragment")
		start := doc.Call("createComment", "vigor")
		end := doc.Call("createComment", "/vigor")
		frag.Call("appendChild", start)
		frag.Call("appendChild", end)

		vigor.Effect(func() {
			var content js.Value
			if n := fn(); n != nil {
				// Building the content shouldn't subscribe this effect to
				// what the content reads: it has effects of its own.
				content = vigor.Untrack(func() js.Value { return n.DOMObject(ctx, doc) })
			}
			clearBetween(start, end)
			if !content.IsUndefined() && !content.IsNull() {
				end.Get("parentNode").Call("insertBefore", content, end)
			}
		})

		return frag
	})
}

// When shows then while cond is true, and otherwise (if given) while it's
// false. The branches are built when they're shown, and only rebuilt when
// cond flips, not whenever what cond reads changes.
func When(cond func() bool, then func() Node, otherwise ...func() Node) Node {
	return NewNode(func(ctx context.Context, doc js.Value) js.Value {
		on := vigor.Memo(cond)

		return Reactive(func() Node {
			if on() {
				return vigor.Untrack(then)
			}
			if len(otherwise) > 0 && otherwise[0] != nil {
				return vigor.Untrack(otherwise[0])
			}

			return nil
		}).DOMObject(ctx, doc)
	})
}

// Await shows a Resource's three states: pending until it first loads (or
// fails), failed with the error if it hasn't loaded and its last fetch
// failed, and ready once it has a value; ready reads the value reactively,
// so a refetch updates it without rebuilding. A nil pending or failed shows
// nothing.
func Await[T any](
	r *vigor.Resource[T],
	pending func() Node,
	failed func(err error) Node,
	ready func(value func() T) Node,
) Node {
	return NewNode(func(ctx context.Context, doc js.Value) js.Value {
		const (
			isPending = iota
			isFailed
			isReady
		)
		state := vigor.Memo(func() int {
			switch {
			case r.Loaded():
				return isReady
			case r.Error() != nil:
				return isFailed
			default:
				return isPending
			}
		})

		return Reactive(func() Node {
			s := state()

			return vigor.Untrack(func() Node {
				switch s {
				case isReady:
					return ready(r.Value)
				case isFailed:
					if failed != nil {
						return failed(r.Error())
					}
				default:
					if pending != nil {
						return pending()
					}
				}

				return nil
			})
		}).DOMObject(ctx, doc)
	})
}

// Raw inserts HTML (or SVG, with its <svg> element) as it is: html is a
// string, or a getter of one. Only use it with markup you trust.
func Raw(html any) Node {
	build := func(doc js.Value, s string) js.Value {
		tpl := doc.Call("createElement", "template")
		tpl.Set("innerHTML", s)

		return tpl.Get("content")
	}

	return NewNode(func(ctx context.Context, doc js.Value) js.Value {
		if !isReactive(html) {
			return build(doc, text(html))
		}

		return Reactive(func() Node {
			s := text(value(html, nil))

			return NewNode(func(_ context.Context, doc js.Value) js.Value { return build(doc, s) })
		}).DOMObject(ctx, doc)
	})
}

// Each renders a list, keyed: when items change, the rows of keys still
// present are kept (and moved if they were reordered), new keys are
// rendered, and gone keys are disposed. A row's item and index are
// getters, so a row whose item changed updates in place.
//
// items is tracked; render is not (read signals inside the row's own
// effects). If items is empty, fallback (if given) is shown.
func Each[T any, K comparable](
	items func() []T,
	key func(T) K,
	render func(item func() T, index func() int) Node,
	fallback ...func() Node,
) Node {
	type row struct {
		scope *vigor.Owner
		nodes []js.Value
		item  func() T
		set   func(T)
		idx   func() int
		setI  func(int)
	}

	return NewNode(func(ctx context.Context, doc js.Value) js.Value {
		frag := doc.Call("createDocumentFragment")
		start := doc.Call("createComment", "vigor_each")
		end := doc.Call("createComment", "/vigor_each")
		frag.Call("appendChild", start)
		frag.Call("appendChild", end)

		parent := vigor.CurrentOwner()
		rows := map[K]*row{}
		var empty *row

		build := func(draw func() Node) *row {
			r := &row{}
			vigor.RunWithOwner(parent, func() {
				r.scope = vigor.CreateScope(func() {
					if n := draw(); n != nil {
						r.nodes = nodesOf(n.DOMObject(ctx, doc))
					}
				})
			})

			return r
		}
		drop := func(r *row) {
			r.scope.Dispose()
			for _, n := range r.nodes {
				n.Call("remove")
			}
		}
		vigor.OnCleanup(func() {
			for _, r := range rows {
				r.scope.Dispose()
			}
			if empty != nil {
				empty.scope.Dispose()
			}
		})

		vigor.Effect(func() {
			list := items()
			if len(list) == 0 && len(fallback) > 0 && fallback[0] != nil {
				for k, r := range rows {
					drop(r)
					delete(rows, k)
				}
				if empty == nil {
					empty = build(fallback[0])
					for _, n := range empty.nodes {
						end.Get("parentNode").Call("insertBefore", n, end)
					}
				}

				return
			}
			if empty != nil {
				drop(empty)
				empty = nil
			}

			next := make(map[K]*row, len(list))
			order := make([]*row, len(list))
			fresh := map[*row]bool{}
			for i, it := range list {
				k := key(it)
				r, ok := rows[k]
				if next[k] != nil {
					panic(fmt.Sprintf("vigor: duplicate key %v in Each", k))
				}
				if !ok {
					item, setItem := vigor.NewSignal(it, vigor.AlwaysNotify[T]())
					idx, setIdx := vigor.NewSignal(i)
					r = build(func() Node { return render(item, idx) })
					r.item, r.set, r.idx, r.setI = item, setItem, idx, setIdx
					fresh[r] = true
				} else {
					vigor.Batch(func() {
						r.set(it)
						r.setI(i)
					})
				}
				next[k] = r
				order[i] = r
			}
			for k, r := range rows {
				if next[k] != r {
					drop(r)
				}
			}
			rows = next

			// Place rows back to front, moving only the ones out of place.
			parentNode := end.Get("parentNode")
			ref := end
			for i := len(order) - 1; i >= 0; i-- {
				r := order[i]
				if len(r.nodes) == 0 {
					continue
				}
				last := r.nodes[len(r.nodes)-1]
				if fresh[r] || !last.Get("nextSibling").Equal(ref) {
					for _, n := range r.nodes {
						parentNode.Call("insertBefore", n, ref)
					}
				}
				ref = r.nodes[0]
			}
		})

		return frag
	})
}
