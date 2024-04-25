package main_test

import (
	"fmt"
	"testing"

	"github.com/VJftw/vigor/examples/features"
	"github.com/chromedp/chromedp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test1gMemos(t *testing.T) {
	page := features.BuildServeAndGetVigorPage(t)

	var appElHtml string
	require.NoError(t, chromedp.Run(page,
		chromedp.OuterHTML("#app", &appElHtml),
	))

	getExpectedHTML := func(count int) string {
		s := `<div id="app"><div>`
		s += `<button>Count: ` + fmt.Sprintf("%d", count) + `</button>`
		fib := fibonacci(count)
		for i := 1; i <= 10; i++ {
			s += fmt.Sprintf(`<div>%d. %d %d %d %d %d</div>`, i, fib, fib, fib, fib, fib)
		}
		s += `</div></div>`

		return s
	}

	assert.Equal(t,
		getExpectedHTML(10),
		appElHtml,
	)

	for i := 10; i < 20; i++ {
		require.NoError(t, chromedp.Run(page,
			chromedp.Click("button"),
			chromedp.OuterHTML("#app", &appElHtml),
		))

		assert.Equal(t, getExpectedHTML(i+1), appElHtml)
	}
}

func fibonacci(num int) int {
	if num <= 1 {
		return 1
	}

	return fibonacci(num-1) + fibonacci(num-2)
}
