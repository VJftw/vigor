package html

import (
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

func (n *nodeOnMount) DOMObject(doc js.Value) js.Value {
	defer n.fn()

	return n.underlyingNode.DOMObject(doc)
}
