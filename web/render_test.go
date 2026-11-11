package web_test

import (
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

	web.RenderToElementID(t.Context(), node, "app")

	assert.Equal(t, `<vigor-info id="vigor-info" hidden=""></vigor-info><vigor-main id="vigor-main"><p>Hello World!</p></vigor-main><vigor-portal-root id="vigor-portal-root"></vigor-portal-root>`, divObj.Get("innerHTML").String())
}
