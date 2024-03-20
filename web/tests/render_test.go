package tests

import (
	"context"
	"syscall/js"
	"testing"
	"time"

	"github.com/VJftw/vigor/html"
	"github.com/VJftw/vigor/web"
)

func TestRenderToElementID(t *testing.T) {
	node := html.El("p", "Hello World!")

	document := js.Global().Get("document")
	divObj := document.Call("createElement", "div")
	divObj.Set("id", "app")
	document.Get("body").Call("append", divObj)

	ctx, cancel := context.WithCancel(context.TODO())
	go func() {
		<-time.After(1 * time.Second)
		cancel()
	}()
	web.RenderToElementID(ctx, node, "app")
}
