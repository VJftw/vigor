package html

import (
	"syscall/js"

	"github.com/VJftw/vigor"
)

type nodeSwitch struct {
	defaultNode Node
	cases       []*switchCase
}

func Switch(defaultNode Node, cases ...*switchCase) Node {
	return &nodeSwitch{
		defaultNode: defaultNode,
		cases:       cases,
	}
}

func (n *nodeSwitch) DOMObject(doc js.Value) js.Value {
	obj := doc.Call("createElement", "vigor-switch")

	subscriber := vigor.NewFnSubscriber()
	subscriber.SetFn(func() {
		for _, c := range n.cases {
			if c.when() {
				obj.Call("replaceChildren", c.node.DOMObject(doc))
				return
			}
		}

		obj.Call("replaceChildren", n.defaultNode.DOMObject(doc))
	}).Run()

	for _, c := range n.cases {
		for _, dep := range c.deps {
			dep(subscriber)
		}
	}

	return obj
}

type switchCase struct {
	when func() bool
	node Node
	deps []vigor.GetterFn
}

func Case(when func() bool, node Node, deps ...vigor.GetterFn) *switchCase {
	return &switchCase{
		when: when,
		node: node,
		deps: deps,
	}
}
