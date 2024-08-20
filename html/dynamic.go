package html

import (
	"context"
	"syscall/js"

	"github.com/VJftw/vigor"
)

type nodeDynamic struct {
	fnThatReturnsNode func(vigor.Subscriber) Node
}

func Dynamic(fn func(vigor.Subscriber) Node) Node {
	return &nodeDynamic{
		fnThatReturnsNode: fn,
	}
}

func (n *nodeDynamic) DOMObject(ctx context.Context, doc js.Value) js.Value {
	obj := doc.Call("createDocumentFragment")
	fragmentRendered := false

	subscriber := vigor.NewFnSubscriber()
	subscriber.SetFn(func() {
		newObj := n.fnThatReturnsNode(subscriber).DOMObject(ctx, doc)

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
