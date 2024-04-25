package main_test

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/VJftw/vigor/examples/features"
	"github.com/chromedp/chromedp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test1eEffects(t *testing.T) {
	logsCh := make(chan string, 10)
	page := features.BuildServeAndGetVigorPage(t,
		chromedp.WithDebugf(func(s string, i ...interface{}) {
			event := map[string]any{}
			if err := json.Unmarshal(i[0].([]byte), &event); err != nil {
				return
			}
			if method, ok := event["method"]; ok && method == "Runtime.consoleAPICalled" {
				line := event["params"].(map[string]any)["args"].([]any)[0].(map[string]any)["value"].(string)
				logsCh <- line
			}
		}),
	)

	var appElHtml string
	require.NoError(t, chromedp.Run(page,
		chromedp.OuterHTML("#app", &appElHtml),
	))

	assert.Equal(t,
		`<div id="app"><button>Click Me</button></div>`,
		appElHtml,
	)

	getLogMsg := func() string {
		select {
		case <-time.After(3 * time.Second):
			require.Fail(t, "timed out waiting for log message")
			return ""
		case msg := <-logsCh:
			return strings.Join(strings.Split(msg, " ")[2:], " ")
		}
	}

	for i := 1; i <= 20; i++ {
		require.NoError(t, chromedp.Run(page,
			chromedp.Click("button"),
		))
		assert.Equal(t, fmt.Sprintf("The count is now %d", i), getLogMsg())
	}
}
