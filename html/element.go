package html

import (
	"fmt"
	"syscall/js"

	"github.com/VJftw/vigor"
)

type nodeEl struct {
	name string

	attrs     map[string]any
	attrOrder []string
	classList []any
	events    map[string][]JSEventFunc

	children []Node
}

func El(name string, children ...any) Node {
	n := &nodeEl{
		name:      name,
		attrs:     map[string]any{},
		attrOrder: []string{},
		children:  []Node{},
		events:    map[string][]JSEventFunc{},
	}

	autoTextNodeChildren := []any{}

	for _, o := range children {
		switch x := o.(type) {
		case Node:
			n.children = append(n.children, x)
		case *elementAttribute:
			if _, ok := n.attrs[x.k]; !ok {
				n.attrOrder = append(n.attrOrder, x.k)
			}
			n.attrs[x.k] = x.value()
		case *elementClass:
			n.classList = append(n.classList, x.values()...)
		case *elementEvent:
			if _, ok := n.events[x.event]; !ok {
				n.events[x.event] = []JSEventFunc{}
			}
			n.events[x.event] = append(n.events[x.event], x.fn)
		case bool,
			string,
			int, int8, int16, int32, int64,
			uint, uint8, uint16, uint32, uint64, uintptr,
			float32, float64,
			complex64, complex128:
			// Allow primitive types shorthand to render a Text Node.
			autoTextNodeChildren = append(autoTextNodeChildren, x)
		case vigor.GetterFn:
			// Allow getter to be passed shorthand to render a Text node.
			autoTextNodeChildren = append(autoTextNodeChildren, x)

		default:
			vigor.Log.Fatal("unexpected child type", x)
		}
	}

	if len(autoTextNodeChildren) > 0 {
		n.children = append(n.children, Text(autoTextNodeChildren...))
	}

	return n
}

func (n *nodeEl) DOMObject(doc js.Value) js.Value {
	obj := doc.Call("createElement", n.name)

	subscriber := vigor.NewFnSubscriber()
	subscriber.SetFn(func() {
		for _, k := range n.attrOrder {
			v := n.attrs[k]
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
			classStr := fmt.Sprintf("%s", class)
			if len(classStr) > 0 {
				className += classStr
			}
		}
		if className != "" {
			obj.Set("className", className)
		}
	}).Run()

	for event, listeners := range n.events {
		for _, listener := range listeners {
			fn := js.FuncOf(func(this js.Value, args []js.Value) any { listener(this, args); return nil })
			obj.Call("addEventListener", event, fn)
		}
	}

	for _, c := range n.children {
		if x := c.DOMObject(doc); !x.IsUndefined() {
			obj.Call("append", x)
		}
	}

	return obj
}
