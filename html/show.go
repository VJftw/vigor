package html

import (
	"syscall/js"

	"github.com/VJftw/vigor"
)

type nodeShow struct {
	when   vigor.GetterFn
	truthy Node
	falsey Node
}

func Show(
	when vigor.GetterFn,
	truthy Node,
	falsey Node,
) Node {
	return &nodeShow{
		when:   when,
		truthy: truthy,
		falsey: falsey,
	}
}

func (n *nodeShow) DOMObject(doc js.Value) js.Value {
	obj := doc.Call("createDocumentFragment")
	fragmentRendered := false

	subscriber := vigor.NewFnSubscriber()
	subscriber.SetFn(func() {
		var newObj js.Value
		if n.when(subscriber).(bool) {
			newObj = n.truthy.DOMObject(doc)
		} else {
			newObj = n.falsey.DOMObject(doc)
		}

		if !fragmentRendered {
			obj.Call("replaceChildren", newObj)
			fragmentRendered = true
		} else {
			obj.Call("replaceWith", newObj)
		}

		obj = newObj
	}).Run()

	return obj
}
