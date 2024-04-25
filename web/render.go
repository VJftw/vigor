package web

import (
	"context"
	"syscall/js"

	"github.com/VJftw/vigor/html"
)

// RenderToElementID is a helper function to quickly render the given Node to
// a HTML element that has the given ID.
func RenderToElementID(ctx context.Context, n html.Node, id string) error {
	document := js.Global().Get("document")
	rootEl := document.Call("getElementById", id)
	return RenderTo(ctx, document, n, rootEl)
}

// RenderTo renders the given HTML node to the given JS element.
func RenderTo(ctx context.Context, document js.Value, n html.Node, rootEl js.Value) error {
	portalObj := document.Call("createElement", "vigor-portal-root")
	portalObj.Set("id", "vigor-portal-root")
	document.Get("body").Call("append", portalObj)

	thisEl := n.DOMObject(document)
	if !thisEl.IsNull() {
		rootEl.Call("replaceChildren", thisEl)
	}

	infoObj := document.Call("createElement", "vigor-info")
	infoObj.Set("id", "vigor-info")
	infoObj.Set("hidden", true)
	document.Get("body").Call("append", infoObj)

	<-ctx.Done()

	return ctx.Err()
}
