package html

import (
	"syscall/js"

	"github.com/VJftw/vigor"
)

type ChildrenElementPlugin struct {
	children []any
}

func NewChildrenElementPlugin() ElementPlugin {
	return &ChildrenElementPlugin{
		children: []any{},
	}
}

func (p *ChildrenElementPlugin) HandleChild(child any) (bool, error) {
	switch x := child.(type) {
	case Node:
		p.children = append(p.children, x)
	case bool,
		string,
		int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64, uintptr,
		float32, float64,
		complex64, complex128:
		p.children = append(p.children, x)

		// Allow primitive types shorthand to render a Text Node.
		// autoTextNodeChildren = append(autoTextNodeChildren, x)
	case vigor.GetterFn:
		p.children = append(p.children, x)

		// Allow getter to be passed shorthand to render a Text node.
		// autoTextNodeChildren = append(autoTextNodeChildren, x)

	default:
		return false, nil
	}

	return true, nil
}

func (p *ChildrenElementPlugin) Render(doc, obj js.Value) error {
	for i := 0; i < len(p.children); i++ {
		child := p.children[i]
		if childNode, ok := child.(Node); ok {
			if x := childNode.DOMObject(doc); !x.IsUndefined() {
				obj.Call("append", x)
			}
		} else {
			// build new text node and append it
			textNodeArgs := []any{}
			for i < len(p.children) {
				child := p.children[i]
				if _, ok := child.(Node); !ok {
					textNodeArgs = append(textNodeArgs, child)
				} else {
					break
				}
				i++
			}
			childNode := Text(textNodeArgs...)
			if x := childNode.DOMObject(doc); !x.IsUndefined() {
				obj.Call("append", x)
			}
		}
	}
	return nil
}
