package html

import (
	"context"
	"syscall/js"
)

type nodeErrBoundary struct {
	fallbackFn func(err any) Node
	children   []Node
}

func ErrorBoundary(fallbackFn func(err any) Node, children ...Node) Node {
	return &nodeErrBoundary{
		fallbackFn: fallbackFn,
		children:   children,
	}
}

func (n *nodeErrBoundary) DOMObject(ctx context.Context, doc js.Value) (obj js.Value) {
	obj = doc.Call("createDocumentFragment")

	defer func() {
		if r := recover(); r != nil {
			fallbackObj := n.fallbackFn(r).DOMObject(ctx, doc)
			obj.Call("replaceChildren", fallbackObj)
		}
	}()

	newItemObjs := make([]any, len(n.children))
	for i, item := range n.children {
		newItemObjs[i] = item.DOMObject(ctx, doc)
	}

	obj.Call("replaceChildren", newItemObjs...)

	return obj
}
