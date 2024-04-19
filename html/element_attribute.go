package html

import (
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

func (p *AttributeElementPlugin) Render(doc, obj js.Value) error {
	subscriber := vigor.NewFnSubscriber()
	subscriber.SetFn(func() {
		for _, k := range p.attrOrder {
			v := p.attrs[k]
			if x, ok := v.(vigor.GetterFn); ok {
				v = x(subscriber)
			}
			obj.Call("setAttribute", k, v)
		}
	}).Run()

	return nil
}

func Attr(k string, v any) *elementAttribute {
	return &elementAttribute{
		k: k,
		v: v,
	}
}
