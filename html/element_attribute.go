package html

type elementAttribute struct {
	k string
	v any
}

func (e *elementAttribute) value() any {
	return e.v
}

func Attr(k string, v any) *elementAttribute {
	return &elementAttribute{
		k: k,
		v: v,
	}
}
