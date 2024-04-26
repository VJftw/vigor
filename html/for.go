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
	isVigorEmpty := false

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

		defer func() {
			currentItemObjs = newItemObjs
		}()

		if len(newItemObjs) < 1 {
			if !isVigorEmpty {
				placeholderObj := doc.Call("createElement", "vigor-empty")
				obj.Call("replaceChildren", placeholderObj)
				obj = placeholderObj
				isVigorEmpty = true
			}

			return
		}

		if len(currentItemObjs) < 1 {
			// empty before, so placeholder or fragment is obj.
			if isVigorEmpty {
				fragmentObj := doc.Call("createDocumentFragment")
				fragmentObj.Call("replaceChildren", newItemObjs...)
				obj.Call("replaceWith", fragmentObj)
				obj = fragmentObj
				isVigorEmpty = false
			} else {
				obj.Call("replaceChildren", newItemObjs...)
			}

			return
		}

		for _, newItemObj := range newItemObjs {
			currentItemObjs[len(currentItemObjs)-1].(js.Value).Call("insertAdjacentElement", "beforebegin", newItemObj)
		}
		for _, currentItemObj := range currentItemObjs {
			currentItemObj.(js.Value).Call("remove")
		}
	}).Run()

	return obj
}
