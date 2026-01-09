package html_test

import (
	"context"
	"syscall/js"
	"testing"
	"time"

	"github.com/VJftw/vigor"
	"github.com/VJftw/vigor/html"
	"github.com/VJftw/vigor/web"
	"github.com/stretchr/testify/assert"
)

func TestFor_Static(t *testing.T) {
	tests := []struct {
		description string
		inSlice     any
		inForFn     html.ForFn
		opts        []html.ForOpt
		output      string
	}{
		{
			"slice primitive",
			[]int{1, 2, 3, 4, 5},
			func(i int, v any) html.Node {
				return html.Text(i, ": ", v, ";")
			}, nil,
			`<!--vigor_for-start-->0: 1;1: 2;2: 3;3: 4;4: 5;<!--vigor_for-end-->`,
		},
		{
			"slice element",
			[]int{1, 2, 3, 4, 5},
			func(i int, v any) html.Node {
				return html.El("li", i, ": ", v)
			}, nil,
			`<!--vigor_for-start--><li>0: 1</li><li>1: 2</li><li>2: 3</li><li>3: 4</li><li>4: 5</li><!--vigor_for-end-->`,
		},
		{
			"empty slice",
			[]int{},
			func(i int, v any) html.Node {
				return html.El("li", i, ": ", v)
			}, nil,
			`<!--vigor_for-start--><!--vigor_for-end-->`,
		},
		{
			"empty slice with fallback",
			[]int{},
			func(i int, v any) html.Node {
				return html.El("li", i, ": ", v)
			},
			[]html.ForOpt{html.WithForFallback(html.Text("Hello World"))},
			`<!--vigor_for-start-->Hello World<!--vigor_for-end-->`,
		},
	}

	document := js.Global().Get("document")
	appObj := document.Call("createElement", "div")
	appObj.Set("id", "app")
	document.Get("body").Call("append", appObj)

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			t.Cleanup(func() { appObj.Call("replaceChildren") })
			forNode := html.For(tt.inSlice, tt.inForFn, tt.opts...)
			ctx, cancel := context.WithCancel(context.TODO())
			t.Cleanup(cancel)
			web.RenderTo(ctx, document, forNode, appObj)

			assert.Eventually(t, func() bool {
				innerHTML := document.Call("getElementById", "vigor-main").Get("innerHTML").String()

				return tt.output == innerHTML
			}, time.Second, time.Millisecond)
		})
	}
}

func TestFor_Dynamic(t *testing.T) {
	tests := []struct {
		description string
		inSlices    []any
		inForFn     html.ForFn
		opts        []html.ForOpt
		outputs     []string
	}{
		{
			"slice primitive",
			[]any{
				[]int{1, 2, 3, 4, 5},
				[]int{5, 4, 3, 2, 1},
				[]int{5},
				[]int{},
				[]int{1},
			},
			func(i int, v any) html.Node {
				return html.Text(i, ": ", v, ";")
			}, nil,
			[]string{
				`<!--vigor_for-start-->0: 1;1: 2;2: 3;3: 4;4: 5;<!--vigor_for-end-->`,
				`<!--vigor_for-start-->0: 5;1: 4;2: 3;3: 2;4: 1;<!--vigor_for-end-->`,
				`<!--vigor_for-start-->0: 5;<!--vigor_for-end-->`,
				`<!--vigor_for-start--><!--vigor_for-end-->`,
				`<!--vigor_for-start-->0: 1;<!--vigor_for-end-->`,
			},
		},
		{
			"slice element",
			[]any{
				[]int{1, 2, 3, 4, 5},
				[]int{5},
				[]int{},
				[]int{5, 4, 3, 2, 1},
			},
			func(i int, v any) html.Node {
				return html.El("li", i, ": ", v)
			}, nil,
			[]string{
				`<!--vigor_for-start--><li>0: 1</li><li>1: 2</li><li>2: 3</li><li>3: 4</li><li>4: 5</li><!--vigor_for-end-->`,
				`<!--vigor_for-start--><li>0: 5</li><!--vigor_for-end-->`,
				`<!--vigor_for-start--><!--vigor_for-end-->`,
				`<!--vigor_for-start--><li>0: 5</li><li>1: 4</li><li>2: 3</li><li>3: 2</li><li>4: 1</li><!--vigor_for-end-->`,
			},
		},
		{
			"slice element with fallback",
			[]any{
				[]int{1, 2, 3, 4, 5},
				[]int{5},
				[]int{},
				[]int{5, 4, 3, 2, 1},
			},
			func(i int, v any) html.Node {
				return html.El("li", i, ": ", v)
			},
			[]html.ForOpt{html.WithForFallback(html.Text("Hello World"))},
			[]string{
				`<!--vigor_for-start--><li>0: 1</li><li>1: 2</li><li>2: 3</li><li>3: 4</li><li>4: 5</li><!--vigor_for-end-->`,
				`<!--vigor_for-start--><li>0: 5</li><!--vigor_for-end-->`,
				`<!--vigor_for-start-->Hello World<!--vigor_for-end-->`,
				`<!--vigor_for-start--><li>0: 5</li><li>1: 4</li><li>2: 3</li><li>3: 2</li><li>4: 1</li><!--vigor_for-end-->`,
			},
		},
	}

	document := js.Global().Get("document")
	appObj := document.Call("createElement", "div")
	appObj.Set("id", "app")
	document.Get("body").Call("append", appObj)

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			t.Cleanup(func() { appObj.Call("replaceChildren") })

			sliceSignal, setSliceSignal := vigor.CreateSignal([]any{})

			forNode := html.For(sliceSignal, tt.inForFn, tt.opts...)
			ctx, cancel := context.WithCancel(context.TODO())
			t.Cleanup(cancel)
			web.RenderTo(ctx, document, forNode, appObj)

			for i, output := range tt.outputs {
				setSliceSignal(tt.inSlices[i])

				assert.Eventuallyf(t, func() bool {
					innerHTML := document.Call("getElementById", "vigor-main").Get("innerHTML").String()

					return output == innerHTML
				}, time.Second, time.Millisecond, "last innerHTML: %s", document.Call("getElementById", "vigor-main").Get("innerHTML").String())
			}
		})
	}
}
