package html

import (
	"fmt"
	"syscall/js"

	"github.com/VJftw/vigor"
)

type nodeText struct {
	texts []any
}

func Text(texts ...any) Node {
	return &nodeText{
		texts: texts,
	}
}

func (n *nodeText) DOMObject(doc js.Value) js.Value {
	obj := doc.Call("createTextNode", "")
	hasGetterFn := false
	for _, text := range n.texts {
		if _, ok := text.(vigor.GetterFn); ok {
			hasGetterFn = true
			break
		}
		if _, ok := text.(*subscribedText); ok {
			hasGetterFn = true
			break
		}
	}

	if hasGetterFn {
		subscriber := vigor.NewFnSubscriber()
		subscriber.SetFn(func() {
			content := ""
			for _, text := range n.texts {
				if dep, ok := text.(*subscribedText); ok {
					text = dep.fn(subscriber)
				}

				if getter, ok := text.(vigor.GetterFn); ok {
					text = getter(subscriber)
				}

				content += fmt.Sprintf("%v", text)
			}

			obj.Set("textContent", content)
		}).Run()
	} else {
		content := ""
		for _, text := range n.texts {
			content += fmt.Sprintf("%v", text)
		}

		obj.Set("textContent", content)
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
