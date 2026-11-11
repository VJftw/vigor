package html

import (
	"context"
	"strings"
	"syscall/js"

	"github.com/VJftw/vigor"
)

type nodeText struct {
	texts []any
}

// Text is a text node of its parts joined: strings, numbers and so on, and
// getters, which keep it up to date.
func Text(texts ...any) Node {
	return &nodeText{
		texts: texts,
	}
}

func (n *nodeText) DOMObject(ctx context.Context, doc js.Value) js.Value {
	obj := doc.Call("createTextNode", "")

	render := func(sub vigor.Subscriber) {
		var sb strings.Builder
		for _, t := range n.texts {
			sb.WriteString(text(value(t, sub)))
		}
		obj.Set("textContent", sb.String())
	}

	reactive := false
	for _, t := range n.texts {
		if isReactive(t) {
			reactive = true

			break
		}
	}
	if reactive {
		bind(render)
	} else {
		render(nil)
	}

	return obj
}

type SubscribedTextFunc func(vigor.Subscriber) any

type subscribedText struct {
	fn SubscribedTextFunc
}

func SubscribedText(fn SubscribedTextFunc) *subscribedText {
	return &subscribedText{
		fn: fn,
	}
}
