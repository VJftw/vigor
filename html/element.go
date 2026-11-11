package html

import (
	"context"
	"syscall/js"

	"github.com/VJftw/vigor"
)

type ElementPlugin interface {
	HandleChild(child any) (bool, error)
	Render(ctx context.Context, doc, obj js.Value) error
}

type nodeEl struct {
	namespace string
	name      string

	events  []*elementEvent
	refs    []func(js.Value)
	plugins []ElementPlugin
}

// elementGroup is several options and children, e.g. what BindValue returns.
type elementGroup struct{ items []any }

// Group bundles children and options (Attr, Class, On...) into one value,
// so a helper can return them together.
func Group(items ...any) *elementGroup { return &elementGroup{items: items} }

type elementRef struct{ fn func(js.Value) }

// Ref calls fn with the element once it is built (its attributes and
// children are set, but it isn't in the page yet), e.g. to focus it later or
// measure it.
func Ref(fn func(el js.Value)) *elementRef { return &elementRef{fn: fn} }

func NEl(namespace, name string, children ...any) Node {
	n := El(name, children...).(*nodeEl)

	n.namespace = namespace

	return n
}

func El(name string, children ...any) Node {
	n := &nodeEl{
		name: name,
		// In this order: children before properties, because setting a
		// select's value needs its options to exist.
		plugins: []ElementPlugin{
			NewAttributeElementPlugin(),
			NewClassElementPlugin(),
			NewStyleElementPlugin(),
			NewChildrenElementPlugin(),
			NewPropertyElementPlugin(),
		},
	}

	n.add(children)

	return n
}

func (n *nodeEl) add(children []any) {
	for _, o := range children {
		switch x := o.(type) {
		case *elementGroup:
			n.add(x.items)

			continue
		case *elementEvent:
			n.events = append(n.events, x)

			continue
		case *elementRef:
			n.refs = append(n.refs, x.fn)

			continue
		}

		var handled bool
		var err error
		for _, plugin := range n.plugins {
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
}

func (n *nodeEl) DOMObject(ctx context.Context, doc js.Value) js.Value {
	var obj js.Value
	if n.namespace == "" {
		obj = doc.Call("createElement", n.name)
	} else {
		obj = doc.Call("createElementNS", n.namespace, n.name)
	}

	for _, plugin := range n.plugins {
		if err := plugin.Render(ctx, doc, obj); err != nil {
			vigor.Log.Fatal("error rendering", err)
		}
	}

	for _, ev := range n.events {
		listen(obj, ev.event, ev.opts, ev.fn)
	}

	for _, ref := range n.refs {
		ref(obj)
	}

	return obj
}
