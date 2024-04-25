package tests

import (
	"syscall/js"
	"testing"

	"github.com/VJftw/vigor"
	"github.com/VJftw/vigor/html"
	"github.com/stretchr/testify/assert"
)

func TestText(t *testing.T) {
	aSig, _ := vigor.CreateSignal("signal!")
	tests := []struct {
		description string
		inTexts     []any
		outContent  string
	}{
		{"string", []any{"hello world!"}, "hello world!"},
		{"integer", []any{1}, "1"},
		{"bool", []any{true}, "true"},
		{"empty", []any{}, ""},
		{"multiple texts", []any{"hello", 1, true}, "hello1true"},
		{"signal", []any{"hello ", aSig}, "hello signal!"},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			textNode := html.Text(tt.inTexts...)

			document := js.Global().Get("document")

			textContent := textNode.DOMObject(document).Get("textContent").String()

			assert.Equal(t, tt.outContent, textContent)
		})
	}
}
