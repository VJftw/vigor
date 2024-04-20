package html

import (
	"syscall/js"

	"github.com/VJftw/vigor"
)

type elementProperty struct {
	k string
	v any
}

type PropertyElementPlugin struct {
	props map[string]any
}

func NewPropertyElementPlugin() ElementPlugin {
	return &PropertyElementPlugin{
		props: map[string]any{},
	}
}

func (p *PropertyElementPlugin) HandleChild(child any) (bool, error) {
	if x, ok := child.(*elementProperty); ok {
		p.props[x.k] = x.v

		return true, nil
	}

	return false, nil
}

func (p *PropertyElementPlugin) Render(doc, obj js.Value) error {
	subscriber := vigor.NewFnSubscriber()
	subscriber.SetFn(func() {
		for k, v := range p.props {
			if x, ok := v.(vigor.GetterFn); ok {
				v = x(subscriber)
			}
			obj.Set(k, v)
		}
	}).Run()

	return nil
}

func Property(k string, v any) *elementProperty {
	return &elementProperty{
		k: k,
		v: v,
	}
}
