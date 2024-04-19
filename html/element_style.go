package html

import (
	"fmt"
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

func (p *StyleElementPlugin) Render(doc, obj js.Value) error {
	if len(p.propertyValues) <= 0 {
		return nil
	}
	subscriber := vigor.NewFnSubscriber()
	subscriber.SetFn(func() {
		styleAttr := ""
		for _, k := range p.propertyOrder {
			v := p.propertyValues[k]
			if x, ok := v.(vigor.GetterFn); ok {
				v = x(subscriber)
			}
			styleAttr += fmt.Sprintf("%s: %s;", k, v)
		}
		obj.Set("style", styleAttr)
	}).Run()

	return nil
}

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
