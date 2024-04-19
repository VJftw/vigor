package html

import (
	"reflect"
	"syscall/js"

	"github.com/VJftw/vigor"
)

type nodeIndex struct {
	fn func(i int, v any) Node
	v  vigor.GetterFn
}

func Index(v vigor.GetterFn, fn func(i int, v any) Node) Node {
	return &nodeIndex{
		fn: fn,
		v:  v,
	}
}

func (n *nodeIndex) DOMObject(doc js.Value) js.Value {
	obj := doc.Call("createElement", "vigor-index")

	currentItems := []any{}
	currentItemObjs := []js.Value{}

	subscriber := vigor.NewFnSubscriber()
	subscriber.SetFn(func() {
		newItems := n.v(subscriber).([]any)

		for i := 0; i < len(newItems) && i < len(currentItems); i++ {
			if !reflect.DeepEqual(newItems[i], currentItems[i]) {
				currentItems[i] = newItems[i]

				newItemObj := n.fn(i, newItems[i]).DOMObject(doc)
				currentItemObjs[i].Call("replaceWith", newItemObj)
			}
		}

		if len(newItems) > len(currentItems) {
			for i := len(currentItems); i < len(newItems); i++ {
				currentItems = append(currentItems, newItems[i])

				newItemObj := n.fn(i, newItems[i]).DOMObject(doc)
				currentItemObjs = append(currentItemObjs, newItemObj)
				obj.Call("append", newItemObj)
			}
		}

		if len(currentItems) > len(newItems) {
			for i := len(newItems); i < len(currentItems); i++ {
				itemObjToRemove := currentItemObjs[i]
				itemObjToRemove.Call("remove")
			}

			currentItems = currentItems[0:len(newItems):len(newItems)]
			currentItemObjs = currentItemObjs[0:len(newItems):len(newItems)]
		}
	}).Run()

	return obj
}
