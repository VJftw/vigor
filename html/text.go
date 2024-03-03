package html

import (
	"syscall/js"

	"github.com/VJftw/vigor"
)

type nodeText struct {
	text any
}

func Text(text any) Node {
	return &nodeText{
		text: text,
	}
}

func (n *nodeText) DOMObject(doc js.Value) js.Value {
	obj := doc.Call("createTextNode", "")

	if x, ok := n.text.(vigor.GetterFn); ok {
		subscriber := vigor.NewFnSubscriber()
		subscriber.SetFn(func() {
			obj.Set("textContent", x(subscriber))
		}).Run()
	} else {
		obj.Set("textContent", n.text)
	}

	return obj
}
