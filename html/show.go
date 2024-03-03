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
	obj := doc.Call("createElement", "vigor-show")

	subscriber := vigor.NewFnSubscriber()
	subscriber.SetFn(func() {
		if n.when(subscriber) == true {
			obj.Call("replaceChildren", n.truthy.DOMObject(doc))
		} else {
			obj.Call("replaceChildren", n.falsey.DOMObject(doc))
		}
	}).Run()

	return obj
}
