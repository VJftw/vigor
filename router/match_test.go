package router_test

import (
	"testing"

	"github.com/VJftw/vigor/router"
	"github.com/stretchr/testify/assert"
)

func TestPatternMatch(t *testing.T) {
	tests := []struct {
		pattern, path string
		params        router.Params
		ok            bool
	}{
		{"/", "/", router.Params{}, true},
		{"/", "", router.Params{}, true},
		{"/", "/x", nil, false},
		{"/display", "/display/", router.Params{}, true},
		{"/p/:id", "/p/automations", router.Params{"id": "automations"}, true},
		{"/p/:id", "/p/a%20b", router.Params{"id": "a b"}, true},
		{"/p/:id", "/p", nil, false},
		{"/p/:id", "/p/a/b", nil, false},
		{"/a/:x/b/:y", "/a/1/b/2", router.Params{"x": "1", "y": "2"}, true},
		{"/files/*path", "/files/a/b/c", router.Params{"path": "a/b/c"}, true},
		{"/files/*path", "/files", router.Params{"path": ""}, true},
		{"/files/*path", "/other/a", nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.pattern+" "+tt.path, func(t *testing.T) {
			params, ok := router.ParsePattern(tt.pattern).Match(tt.path)
			assert.Equal(t, tt.ok, ok)
			assert.Equal(t, tt.params, params)
		})
	}
}

func TestPatternBuild(t *testing.T) {
	assert.Equal(t, "/", router.ParsePattern("/").Build(nil))
	assert.Equal(t, "/p/a%20b", router.ParsePattern("/p/:id").Build(router.Params{"id": "a b"}))
	assert.Equal(t, "/files/x/y", router.ParsePattern("/files/*p").Build(router.Params{"p": "x/y"}))
}

func TestParseLocation(t *testing.T) {
	l := router.ParseLocation("/display?area=kitchen&idle=0#top")
	assert.Equal(t, "/display", l.Path)
	assert.Equal(t, "kitchen", l.Query.Get("area"))
	assert.Equal(t, "top", l.Hash)
	assert.Equal(t, "/display?area=kitchen&idle=0#top", l.String())
	assert.Equal(t, "/x", router.ParseLocation("x").Path)
}
