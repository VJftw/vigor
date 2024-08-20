package html

import (
	"context"
	"syscall/js"
)

type nodeOnMount struct {
	underlyingNode Node
	fn             func()
}

func OnMount(n Node, fn func()) Node {
	return &nodeOnMount{
		underlyingNode: n,
		fn:             fn,
	}
}

func (n *nodeOnMount) DOMObject(ctx context.Context, doc js.Value) js.Value {
	defer func() { go n.fn() }()

	return n.underlyingNode.DOMObject(ctx, doc)
}
