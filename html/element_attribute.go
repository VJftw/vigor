package html

import "github.com/VJftw/vigor"

type elementAttribute struct {
	k string
	v any
}

func (e *elementAttribute) value() any {
	if x, ok := e.v.(vigor.GetterFn); ok {
		return x()
	}

	return e.v
}

func Attr(k string, v any) *elementAttribute {
	return &elementAttribute{
		k: k,
		v: v,
	}
}
