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
	obj := doc.Call("createDocumentFragment")
	fragmentRendered := false

	subscriber := vigor.NewFnSubscriber()
	subscriber.SetFn(func() {
		var newObj js.Value

		caseMatched := false
		for _, c := range n.cases {
			if c.when() {
				newObj = c.node.DOMObject(doc)
				caseMatched = true
				break
			}
		}

		if !caseMatched {
			newObj = n.defaultNode.DOMObject(doc)
		}

		if !fragmentRendered {
			obj.Call("replaceChildren", newObj)
			fragmentRendered = true
		} else {
			obj.Call("replaceWith", newObj)
		}

		obj = newObj
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
