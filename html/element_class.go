package html

import "github.com/VJftw/vigor"

type elementClass struct {
	v any
}

func (e *elementClass) value() any {
	if x, ok := e.v.(vigor.GetterFn); ok {
		return x()
	}

	return e.v
}

func Class(v any) *elementClass {
	return &elementClass{
		v: v,
	}
}
