package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApp(t *testing.T) {
	fs := http.FileServer(http.Dir("../build"))
	s := httptest.NewServer(fs)
	defer s.Close()

	page := rod.New().MustConnect().
		Timeout(10 * time.Second).
		Logger(rod.DefaultLogger).
		MustPage(s.URL)
	defer page.MustClose()
	page.MustWaitStable()
	page.MustSetWindow(0, 0, 1920, 1080)

	appEl := page.MustElement("#app")

	appHTML, err := appEl.HTML()
	require.NoError(t, err)

	assert.Equal(t,
		`<div id="app"><button>Click Me</button></div>`,
		appHTML,
	)

	logMsgs := make(chan string)
	getLogMsg := func() string {
		select {
		case <-time.After(5 * time.Second):
			assert.Fail(t, "timed out waiting for log message")
			return ""
		case msg := <-logMsgs:
			return strings.Join(strings.Split(msg, " ")[2:], " ")
		}
	}
	go page.EachEvent(func(e *proto.RuntimeConsoleAPICalled) {
		logMsg := e.Args[0].Value.Str()
		logMsgs <- logMsg
	})()

	buttonEl := page.MustElement("button")
	for i := 1; i <= 20; i++ {
		buttonEl.MustClick()
		assert.Equal(t, fmt.Sprintf("The count is now %d", i), getLogMsg())
	}
}
