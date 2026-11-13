//go:build js && wasm

package html

import (
	"context"
	"syscall/js"
)

type DOMObjectFn func(context.Context, js.Value) js.Value

type Node interface {
	// DOMObject returns the DOM Object for the Node, with the given context and
	// global 'document' value.
	DOMObject(ctx context.Context, document js.Value) js.Value
}

type nodeRender struct {
	domObjectFn DOMObjectFn
}

func (n *nodeRender) DOMObject(ctx context.Context, doc js.Value) js.Value {
	return n.domObjectFn(ctx, doc)
}

func NewNode(domObjectFn DOMObjectFn) Node {
	return &nodeRender{
		domObjectFn: domObjectFn,
	}
}
