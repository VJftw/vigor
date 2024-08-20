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

func (n *nodePortal) DOMObject(ctx context.Context, doc js.Value) js.Value {
	rootPortal := doc.Call("getElementById", "vigor-portal-root")
	thisPortalObj := doc.Call("createElement", "vigor-portal")
	rootPortal.Call("append", thisPortalObj)

	thisPortalObj.Call("replaceChildren", n.node.DOMObject(ctx, doc))

	return js.Undefined()
}
