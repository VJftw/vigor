package html

import (
	"syscall/js"

	"github.com/VJftw/vigor"
)

type ElementPlugin interface {
	HandleChild(child any) (bool, error)
	Render(doc, obj js.Value) error
}

type nodeEl struct {
	namespace string
	name      string

	events  map[string][]JSEventFunc
	plugins map[ElementPlugin]struct{}
}

func NEl(namespace, name string, children ...any) Node {
	n := El(name, children...).(*nodeEl)

	n.namespace = namespace

	return n
}

func El(name string, children ...any) Node {
	n := &nodeEl{
		name:   name,
		events: map[string][]JSEventFunc{},
		plugins: map[ElementPlugin]struct{}{
			NewPropertyElementPlugin():  {},
			NewAttributeElementPlugin(): {},
			NewClassElementPlugin():     {},
			NewChildrenElementPlugin():  {},
			NewStyleElementPlugin():     {},
		},
	}

	for _, o := range children {

		switch x := o.(type) {
		case *elementEvent:
			if _, ok := n.events[x.event]; !ok {
				n.events[x.event] = []JSEventFunc{}
			}
			n.events[x.event] = append(n.events[x.event], x.fn)
			continue
		}

		var handled bool
		var err error
		for plugin := range n.plugins {
			handled, err = plugin.HandleChild(o)
			if err != nil {
				vigor.Log.Fatal("could not handle child", err)
			}
			if handled {
				break
			}
		}

		if !handled {
			vigor.Log.Fatal("unexpected child", o)
		}
	}

	return n
}

func (n *nodeEl) DOMObject(doc js.Value) js.Value {
	var obj js.Value
	if n.namespace == "" {
		obj = doc.Call("createElement", n.name)
	} else {
		obj = doc.Call("createElementNS", n.namespace, n.name)
	}

	for plugin := range n.plugins {
		if err := plugin.Render(doc, obj); err != nil {
			vigor.Log.Fatal("error rendering", err)
		}
	}

	for event, listeners := range n.events {
		for _, listener := range listeners {
			fn := js.FuncOf(func(this js.Value, args []js.Value) any { listener(this, args); return nil })
			obj.Call("addEventListener", event, fn)
		}
	}

	return obj
}
