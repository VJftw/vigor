package html

import (
	"fmt"
	"syscall/js"

	"github.com/VJftw/vigor"
)

type ClassElementPlugin struct {
	classList []any
}

func NewClassElementPlugin() ElementPlugin {
	return &ClassElementPlugin{
		classList: []any{},
	}
}

func (p *ClassElementPlugin) HandleChild(child any) (bool, error) {
	if x, ok := child.(*elementClass); ok {
		p.classList = append(p.classList, x.values()...)

		return true, nil
	}

	return false, nil
}

func (p *ClassElementPlugin) Render(doc, obj js.Value) error {
	subscriber := vigor.NewFnSubscriber()
	subscriber.SetFn(func() {
		className := ""
		for i, class := range p.classList {
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
	return nil
}

type elementClass struct {
	v []any
}

func (e *elementClass) values() []any {
	return e.v
}

func Class(v ...any) *elementClass {
	return &elementClass{
		v: v,
	}
}
