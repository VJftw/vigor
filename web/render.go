package web

import (
	"context"
	"syscall/js"

	"github.com/VJftw/vigor"
	"github.com/VJftw/vigor/html"
)

// RenderToElementID is a helper function to quickly render the given Node to
// a HTML element that has the given ID. It returns the function that
// unmounts it (also called when ctx is done).
func RenderToElementID(ctx context.Context, n html.Node, id string) (unmount func()) {
	document := js.Global().Get("document")
	rootEl := document.Call("getElementById", id)

	return RenderTo(ctx, document, n, rootEl)
}

// RenderTo renders the given HTML node to the given JS element, as the root
// of a reactive scope: everything it creates is disposed when the returned
// function is called, or when ctx is done.
func RenderTo(ctx context.Context, document js.Value, n html.Node, rootEl js.Value) (unmount func()) {
	infoObj := document.Call("createElement", "vigor-info")
	infoObj.Set("id", "vigor-info")
	infoObj.Set("hidden", true)
	rootEl.Call("append", infoObj)

	mainObj := document.Call("createElement", "vigor-main")
	mainObj.Set("id", "vigor-main")

	portalObj := document.Call("createElement", "vigor-portal-root")
	portalObj.Set("id", "vigor-portal-root")

	// Both are in the page before rendering, so Portals can find theirs.
	rootEl.Call("append", mainObj)
	rootEl.Call("append", portalObj)

	dispose := vigor.CreateRoot(func(dispose func()) func() {
		if thisEl := n.DOMObject(ctx, document); !thisEl.IsNull() && !thisEl.IsUndefined() {
			mainObj.Call("replaceChildren", thisEl)
		}

		return dispose
	})

	unmount = func() {
		dispose()
		mainObj.Call("replaceChildren")
		portalObj.Call("replaceChildren")
	}
	go func() {
		<-ctx.Done()
		unmount()
	}()

	return unmount
}

// Mount builds the app inside its own reactive scope and renders it to the
// element with the given ID: signals, effects and streams that app creates
// are owned by it, and stopped when the returned function is called (or ctx
// is done).
func Mount(ctx context.Context, id string, app func() html.Node) (unmount func()) {
	return RenderToElementID(ctx, html.NewNode(func(ctx context.Context, doc js.Value) js.Value {
		return app().DOMObject(ctx, doc)
	}), id)
}
