package html

import (
	"syscall/js"

	"github.com/VJftw/vigor"
)

type nodeFor struct {
	fn func(i int, v any) Node
	v  vigor.GetterFn
}

func For(v vigor.GetterFn, fn func(i int, v any) Node) Node {
	return &nodeFor{
		fn: fn,
		v:  v,
	}
}

func (n *nodeFor) DOMObject(doc js.Value) js.Value {
	obj := doc.Call("createElement", "vigor-for")

	subscriber := vigor.NewFnSubscriber()
	subscriber.SetFn(func() {
		items := n.v(subscriber).([]any)
		itemObjs := make([]any, len(items))
		for i, item := range items {
			itemObjs[i] = n.fn(i, item).DOMObject(doc)
		}

		obj.Call("replaceChildren", itemObjs...)
	}).Run()

	return obj
}
