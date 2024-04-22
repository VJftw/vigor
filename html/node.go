package html

import (
	"syscall/js"
)

type Node interface {
	DOMObject(document js.Value) js.Value
}

type nodeRender struct {
	domObjectFn func(js.Value) js.Value
}

func (n *nodeRender) DOMObject(doc js.Value) js.Value {
	return n.domObjectFn(doc)
}

func NewNode(domObjectFn func(js.Value) js.Value) Node {
	return &nodeRender{
		domObjectFn: domObjectFn,
	}
}
