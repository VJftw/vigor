package html

import (
	"context"
	"syscall/js"

	"github.com/VJftw/vigor"
)

type elementAttribute struct {
	k string
	v any
}

type AttributeElementPlugin struct {
	attrs     map[string]any
	attrOrder []string
}

func NewAttributeElementPlugin() ElementPlugin {
	return &AttributeElementPlugin{
		attrs:     map[string]any{},
		attrOrder: []string{},
	}
}

func (p *AttributeElementPlugin) HandleChild(child any) (bool, error) {
	if x, ok := child.(*elementAttribute); ok {
		if _, ok := p.attrs[x.k]; !ok {
			p.attrOrder = append(p.attrOrder, x.k)
		}
		p.attrs[x.k] = x.v

		return true, nil
	}

	return false, nil
}

func (p *AttributeElementPlugin) Render(ctx context.Context, doc, obj js.Value) error {
	if len(p.attrOrder) == 0 {
		return nil
	}
	// Each attribute has its own effect, so changing one doesn't rewrite
	// the others.
	for _, k := range p.attrOrder {
		v := p.attrs[k]
		bind(func(sub vigor.Subscriber) { setAttribute(obj, k, value(v, sub)) })
	}

	return nil
}

// setAttribute sets an attribute the way HTML means it: a boolean attribute
// (disabled, checked, hidden...) is there or not, and nil removes it.
func setAttribute(obj js.Value, k string, v any) {
	switch x := v.(type) {
	case nil:
		obj.Call("removeAttribute", k)
	case bool:
		if x {
			obj.Call("setAttribute", k, "")
		} else {
			obj.Call("removeAttribute", k)
		}
	default:
		obj.Call("setAttribute", k, v)
	}
}

// Attr sets an attribute. v is a string, number, or a getter of one. A bool
// sets or removes the attribute (Attr("disabled", isOff)), and nil removes
// it.
func Attr(k string, v any) *elementAttribute {
	return &elementAttribute{
		k: k,
		v: v,
	}
}
