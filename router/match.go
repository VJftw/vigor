// Package router is client-side routing for vigor apps: routes with
// parameters (/p/:id), the query string as signals (?area=kitchen), links
// that navigate without reloading, and the browser's back button.
package router

import (
	"net/url"
	"strings"
)

// Params are what a route's pattern captured, by name.
type Params map[string]string

type segment struct {
	literal string
	param   string // ":name"
	rest    bool   // "*name": the rest of the path
}

// Pattern is a compiled route path: literal segments, ":name" segments that
// match one segment, and a final "*name" that matches the rest (possibly
// nothing).
type Pattern struct {
	raw  string
	segs []segment
}

// ParsePattern compiles a route path such as "/", "/p/:id" or "/files/*path".
func ParsePattern(pattern string) Pattern {
	p := Pattern{raw: pattern}
	for _, s := range split(pattern) {
		switch {
		case strings.HasPrefix(s, ":"):
			p.segs = append(p.segs, segment{param: s[1:]})
		case strings.HasPrefix(s, "*"):
			p.segs = append(p.segs, segment{param: s[1:], rest: true})
		default:
			p.segs = append(p.segs, segment{literal: s})
		}
	}

	return p
}

func split(path string) []string {
	path = strings.Trim(path, "/")
	if path == "" {
		return nil
	}

	return strings.Split(path, "/")
}

// Match reports whether path fits the pattern, and what it captured.
// Captured values are URL-decoded.
func (p Pattern) Match(path string) (Params, bool) {
	parts := split(path)
	params := Params{}
	for i, seg := range p.segs {
		if seg.rest {
			if i > len(parts) {
				return nil, false
			}
			params[seg.param] = unescape(strings.Join(parts[i:], "/"))

			return params, true
		}
		if i >= len(parts) {
			return nil, false
		}
		switch {
		case seg.param != "":
			params[seg.param] = unescape(parts[i])
		case seg.literal != parts[i]:
			return nil, false
		}
	}

	if len(parts) != len(p.segs) {
		return nil, false
	}

	return params, true
}

func unescape(s string) string {
	if u, err := url.PathUnescape(s); err == nil {
		return u
	}

	return s
}

// Build fills the pattern's parameters in, giving a path
// (Build("/p/:id", {"id": "a b"}) is "/p/a%20b").
func (p Pattern) Build(params Params) string {
	var b strings.Builder
	for _, seg := range p.segs {
		b.WriteByte('/')
		switch {
		case seg.rest:
			b.WriteString(params[seg.param])
		case seg.param != "":
			b.WriteString(url.PathEscape(params[seg.param]))
		default:
			b.WriteString(seg.literal)
		}
	}
	if b.Len() == 0 {
		return "/"
	}

	return b.String()
}

// Location is a parsed URL: where the app is.
type Location struct {
	Path  string
	Query url.Values
	Hash  string
}

// ParseLocation splits "/p/x?a=1#top" into its parts. The path always
// starts with a slash.
func ParseLocation(href string) Location {
	loc := Location{Query: url.Values{}}
	rest := href
	if i := strings.IndexByte(rest, '#'); i >= 0 {
		loc.Hash = rest[i+1:]
		rest = rest[:i]
	}
	if i := strings.IndexByte(rest, '?'); i >= 0 {
		if q, err := url.ParseQuery(rest[i+1:]); err == nil {
			loc.Query = q
		}
		rest = rest[:i]
	}
	if !strings.HasPrefix(rest, "/") {
		rest = "/" + rest
	}
	loc.Path = rest

	return loc
}

// String is the location as a relative URL.
func (l Location) String() string {
	s := l.Path
	if len(l.Query) > 0 {
		s += "?" + l.Query.Encode()
	}
	if l.Hash != "" {
		s += "#" + l.Hash
	}

	return s
}
