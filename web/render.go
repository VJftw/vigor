package web

import (
	"context"
	"syscall/js"

	"github.com/VJftw/vigor/html"
)

// RenderToElementID is a helper function to quickly render the given Node to
// a HTML element that has the given ID.
func RenderToElementID(ctx context.Context, n html.Node, id string) {
	document := js.Global().Get("document")
	rootEl := document.Call("getElementById", id)
	RenderTo(ctx, document, n, rootEl)
}

// RenderTo renders the given HTML node to the given JS element.
func RenderTo(ctx context.Context, document js.Value, n html.Node, rootEl js.Value) {
	mainObj := document.Call("createElement", "vigor-main")
	mainObj.Set("id", "vigor-main")

	thisEl := n.DOMObject(ctx, document)
	if !thisEl.IsNull() {
		mainObj.Call("replaceChildren", thisEl)
	}
	rootEl.Call("append", mainObj)

	portalObj := document.Call("createElement", "vigor-portal-root")
	portalObj.Set("id", "vigor-portal-root")
	rootEl.Call("append", portalObj)

	infoObj := document.Call("createElement", "vigor-info")
	infoObj.Set("id", "vigor-info")
	infoObj.Set("hidden", true)
	rootEl.Call("append", infoObj)
}
