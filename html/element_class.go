//go:build js && wasm

package html

import (
	"context"
	"strings"
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

func (p *ClassElementPlugin) Render(ctx context.Context, doc, obj js.Value) error {
	if len(p.classList) == 0 {
		return nil
	}
	bind(func(sub vigor.Subscriber) {
		var classes []string
		for _, class := range p.classList {
			if s := strings.TrimSpace(text(value(class, sub))); s != "" {
				classes = append(classes, s)
			}
		}
		// Always set it, so a class that went away is removed. As the
		// attribute: SVG elements' className is read-only.
		className := strings.Join(classes, " ")
		if className == "" && !obj.Call("hasAttribute", "class").Bool() {
			return
		}
		obj.Call("setAttribute", "class", className)
	})

	return nil
}

type elementClass struct {
	v []any
}

func (e *elementClass) values() []any {
	return e.v
}

// Class adds classes: strings, or getters of strings (a getter may return
// several, space-separated, or nothing).
func Class(v ...any) *elementClass {
	return &elementClass{
		v: v,
	}
}

type classCond struct {
	name string
	cond any
}

// ClassIf adds a class while cond (a bool, or a getter of anything truthy)
// holds, e.g. ClassIf("active", isActive).
func ClassIf(name string, cond any) *elementClass {
	return &elementClass{v: []any{classCond{name: name, cond: cond}}}
}
