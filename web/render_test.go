package web_test

import (
	"context"
	"syscall/js"
	"testing"

	"github.com/VJftw/vigor/html"
	"github.com/VJftw/vigor/web"
	"github.com/stretchr/testify/assert"
)

func TestRenderToElementID(t *testing.T) {
	node := html.El("p", "Hello World!")

	document := js.Global().Get("document")
	divObj := document.Call("createElement", "div")
	divObj.Set("id", "app")
	document.Get("body").Call("append", divObj)

	web.RenderToElementID(context.TODO(), node, "app")

	assert.Equal(t, `<vigor-main id="vigor-main"><p>Hello World!</p></vigor-main><vigor-portal-root id="vigor-portal-root"></vigor-portal-root><vigor-info id="vigor-info" hidden=""></vigor-info>`, divObj.Get("innerHTML").String())
}
