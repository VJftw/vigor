package html

import (
	"context"
	"syscall/js"

	"github.com/VJftw/vigor"
)

type elementProperty struct {
	k string
	v any
}

type PropertyElementPlugin struct {
	props map[string]any
	order []string
}

func NewPropertyElementPlugin() ElementPlugin {
	return &PropertyElementPlugin{
		props: map[string]any{},
	}
}

func (p *PropertyElementPlugin) HandleChild(child any) (bool, error) {
	if x, ok := child.(*elementProperty); ok {
		if _, ok := p.props[x.k]; !ok {
			p.order = append(p.order, x.k)
		}
		p.props[x.k] = x.v

		return true, nil
	}

	return false, nil
}

func (p *PropertyElementPlugin) Render(ctx context.Context, doc, obj js.Value) error {
	for _, k := range p.order {
		v := p.props[k]
		bind(func(sub vigor.Subscriber) {
			next := js.ValueOf(value(v, sub))
			// Leave it alone when it's already right: setting a text
			// input's value moves the caret.
			if cur := obj.Get(k); cur.Equal(next) {
				return
			}
			obj.Set(k, next)
		})
	}

	return nil
}

// Property sets a DOM property (value, checked, scrollTop...) rather than
// an attribute. v is a value or a getter.
func Property(k string, v any) *elementProperty {
	return &elementProperty{
		k: k,
		v: v,
	}
}
