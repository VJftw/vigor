package html

import (
	"syscall/js"

	"github.com/VJftw/vigor"
)

type nodeEl struct {
	name string

	attrs     map[string]any
	classList []any

	children []Node
}

func El(name string, children ...any) Node {
	n := &nodeEl{
		name:     name,
		attrs:    map[string]any{},
		children: []Node{},
	}

	for _, o := range children {
		switch x := o.(type) {
		case Node:
			n.children = append(n.children, x)
		case *elementAttribute:
			n.attrs[x.k] = x.value()
		case *elementClass:
			n.classList = append(n.classList, x.value())
		default:
			vigor.Log.Fatal("unexpected child type")
		}
	}

	return n
}

func (n *nodeEl) DOMObject(doc js.Value) js.Value {
	obj := doc.Call("createElement", n.name)

	subscriber := vigor.NewFnSubscriber()
	subscriber.SetFn(func() {
		for k, v := range n.attrs {
			if x, ok := v.(vigor.GetterFn); ok {
				v = x(subscriber)
			}
			obj.Call("setAttribute", k, v)
		}

		className := ""
		for i, class := range n.classList {
			if i > 0 {
				className += " "
			}
			if x, ok := class.(vigor.GetterFn); ok {
				class = x(subscriber)
			}
			className += class.(string)
		}
		if className != "" {
			obj.Set("className", className)
		}
	}).Run()

	for _, c := range n.children {
		if x := c.DOMObject(doc); !x.IsUndefined() {
			obj.Call("append", x)
		}
	}

	return obj
}
