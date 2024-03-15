package html

type elementClass struct {
	v []any
}

func (e *elementClass) values() []any {
	return e.v
}

func Class(v ...any) *elementClass {
	return &elementClass{
		v: v,
	}
}
