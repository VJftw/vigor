package html

import (
	"context"
	"strings"
	"syscall/js"

	"github.com/VJftw/vigor"
)

type StyleElementPlugin struct {
	propertyValues map[string]any
	propertyOrder  []string
}

func NewStyleElementPlugin() ElementPlugin {
	return &StyleElementPlugin{
		propertyValues: map[string]any{},
		propertyOrder:  []string{},
	}
}

func (p *StyleElementPlugin) HandleChild(child any) (bool, error) {
	if x, ok := child.(*elementStyle); ok {
		if _, ok := p.propertyValues[x.property]; !ok {
			p.propertyOrder = append(p.propertyOrder, x.property)
		}
		p.propertyValues[x.property] = x.value

		return true, nil
	}

	return false, nil
}

func (p *StyleElementPlugin) Render(ctx context.Context, doc, obj js.Value) error {
	if len(p.propertyValues) <= 0 {
		return nil
	}
	bind(func(sub vigor.Subscriber) {
		var sb strings.Builder
		for _, k := range p.propertyOrder {
			v := value(p.propertyValues[k], sub)
			if v == nil || v == "" {
				continue
			}
			sb.WriteString(k + ": " + text(v) + ";")
		}
		obj.Set("style", sb.String())
	})

	return nil
}

// Style sets a CSS property, including custom properties
// (Style("--duration", "5s")). value is a value or a getter; nil or empty
// leaves the property out.
//
// https://www.solidjs.com/tutorial/bindings_style
type elementStyle struct {
	property string
	value    any
}

func Style(property string, value any) *elementStyle {
	return &elementStyle{
		property: property,
		value:    value,
	}
}
