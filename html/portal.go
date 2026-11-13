//go:build js && wasm

package html

import (
	"context"
	"syscall/js"

	"github.com/VJftw/vigor"
)

type nodePortal struct {
	node Node
}

// Portal renders its node in the page's portal root (outside the component
// it's written in), e.g. for a popup that must escape a clipped container.
// It is removed when the surrounding scope is disposed. Without a portal root
// (the page wasn't rendered with web.RenderTo), it renders in place.
func Portal(node Node) Node {
	return &nodePortal{
		node: node,
	}
}

func (n *nodePortal) DOMObject(ctx context.Context, doc js.Value) js.Value {
	root := doc.Call("getElementById", "vigor-portal-root")
	if root.IsNull() {
		return n.node.DOMObject(ctx, doc)
	}

	holder := doc.Call("createElement", "vigor-portal")
	root.Call("append", holder)
	if content := n.node.DOMObject(ctx, doc); !content.IsUndefined() && !content.IsNull() {
		holder.Call("replaceChildren", content)
	}
	vigor.OnCleanup(func() { holder.Call("remove") })

	return js.Undefined()
}
