package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-rod/rod"
	"github.com/stretchr/testify/assert"
)

func TestApp(t *testing.T) {
	fs := http.FileServer(http.Dir("../build"))
	s := httptest.NewServer(fs)
	defer s.Close()

	page := rod.New().MustConnect().Timeout(10 * time.Second).Logger(rod.DefaultLogger).MustPage(s.URL)
	defer page.MustClose()
	page.Race().Element("#vigor-info").MustDo()

	page.MustSetWindow(0, 0, 640, 480)

	appEl := page.MustElement("#app")

	appHTML := appEl.MustHTML()
	buttonEl := page.MustElement("button")

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
		appHTML,
	)

	for i := 10; i < 20; i++ {
		buttonEl.MustClick()
		appHTML = appEl.MustHTML()
		assert.Equal(t, getExpectedHTML(i+1), appHTML)
	}
}

func fibonacci(num int) int {
	if num <= 1 {
		return 1
	}

	return fibonacci(num-1) + fibonacci(num-2)
}
