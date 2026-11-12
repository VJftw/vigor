package html_test

import (
	"syscall/js"
	"testing"

	"github.com/VJftw/vigor"
	"github.com/VJftw/vigor/html"
	"github.com/stretchr/testify/assert"
)

func TestEl(t *testing.T) {
	aSig, _ := vigor.CreateSignal("signal!")

	tests := []struct {
		description string
		inEl        string
		inChildren  []any
		outContent  string
	}{
		{"div no children", "div", []any{}, "<div></div>"},
		{"div with children", "div", []any{html.El("p"), html.El("p")}, `<div><p></p><p></p></div>`},
		{"div with attribute", "div", []any{html.Attr("data-foo", "bar")}, `<div data-foo="bar"></div>`},
		{"div with signal attribute", "div", []any{html.Attr("data-foo", aSig)}, `<div data-foo="signal!"></div>`},
		{"div with multiple attributes", "div", []any{html.Attr("data-foo", "bar"), html.Attr("data-bar", "foo")}, `<div data-foo="bar" data-bar="foo"></div>`},
		{"div with a class", "div", []any{html.Class("header")}, `<div class="header"></div>`},
		{"div with a signal class", "div", []any{html.Class(aSig)}, `<div class="signal!"></div>`},
		{"div with multiple classes", "div", []any{html.Class("header"), html.Class("strong")}, `<div class="header strong"></div>`},
		{"div with multiple classes alternate", "div", []any{html.Class("header", "strong")}, `<div class="header strong"></div>`},
		{"div with text", "div", []any{html.Text("Hello World!")}, `<div>Hello World!</div>`},
		{"div with text alternate", "div", []any{"Hello World!"}, `<div>Hello World!</div>`},
		{"div with signal text alternate", "div", []any{"Hello ", aSig}, `<div>Hello signal!</div>`},
		{"div with event", "div", []any{html.On("mouseover", func(this js.Value, args []js.Value) {})}, `<div></div>`},
		// TODO: 100% coverage:
		// - dispatch event
		// - invalid child
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			elNode := html.El(tt.inEl, tt.inChildren...)

			document := js.Global().Get("document")
			textContent := elNode.DOMObject(t.Context(), document).Get("outerHTML").String()

			assert.Equal(t, tt.outContent, textContent)
		})
	}
}

const svgNS = "http://www.w3.org/2000/svg"

func TestSVGClass(t *testing.T) {
	document := js.Global().Get("document")

	el := html.NEl(svgNS, "g", html.Class("node", "big")).DOMObject(t.Context(), document)
	assert.Equal(t, `<g class="node big"></g>`, el.Get("outerHTML").String())
	assert.Equal(t, svgNS, el.Get("namespaceURI").String())

	// A class that follows a signal comes and goes.
	selected, setSelected := vigor.NewSignal(false)
	el = html.NEl(svgNS, "circle", html.Class("node"), html.ClassIf("selected", selected)).DOMObject(t.Context(), document)
	assert.Equal(t, "node", el.Call("getAttribute", "class").String())
	setSelected(true)
	assert.Equal(t, "node selected", el.Call("getAttribute", "class").String())
	assert.True(t, el.Get("classList").Call("contains", "selected").Bool())
	setSelected(false)
	assert.Equal(t, "node", el.Call("getAttribute", "class").String())
}
