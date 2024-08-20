package html

import (
	"context"
	"syscall/js"
)

type nodeStylesheetContent struct {
	content string
}

func StylesheetContent(content string) Node {
	return &nodeStylesheetContent{
		content: content,
	}
}

func (n *nodeStylesheetContent) DOMObject(ctx context.Context, doc js.Value) js.Value {
	head := doc.Get("head")

	style := doc.Call("createElement", "style")
	style.Set("type", "text/css")
	cssText := doc.Call("createTextNode", n.content)
	style.Call("append", cssText)

	head.Call("append", style)

	return js.Undefined()
}
