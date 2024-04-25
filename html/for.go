package html

import (
	"syscall/js"

	"github.com/VJftw/vigor"
)

type nodeFor struct {
	fn func(i int, v any) Node
	v  any
}

func For(v any, fn func(i int, v any) Node) Node {
	return &nodeFor{
		fn: fn,
		v:  v,
	}
}

func (n *nodeFor) DOMObject(doc js.Value) js.Value {
	obj := doc.Call("createDocumentFragment")

	currentItemObjs := []any{}

	subscriber := vigor.NewFnSubscriber()
	subscriber.SetFn(func() {
		items := []any{}
		if getterFn, ok := n.v.(vigor.GetterFn); ok {
			items = getterFn(subscriber).([]any)
		} else {
			items = n.v.([]any)
		}

		newItemObjs := make([]any, len(items))
		for i, item := range items {
			newItemObjs[i] = n.fn(i, item).DOMObject(doc)
		}

		if len(currentItemObjs) < 1 {
			obj.Call("replaceChildren", newItemObjs...)
		} else {
			for _, newItemObj := range newItemObjs {
				currentItemObjs[len(currentItemObjs)-1].(js.Value).Call("insertBefore", newItemObj)
			}
			for _, currentItemObj := range currentItemObjs {
				currentItemObj.(js.Value).Call("remove")
			}
		}

		currentItemObjs = newItemObjs
	}).Run()

	return obj
}
