package html

import (
	"context"
	"syscall/js"
)

type nodePortal struct {
	node Node
}

func Portal(node Node) Node {
	return &nodePortal{
		node: node,
	}
}

// TODO: Portal should be separate if a separate tree
func (n *nodePortal) DOMObject(ctx context.Context, doc js.Value) js.Value {
	return n.node.DOMObject(ctx, doc)
}
