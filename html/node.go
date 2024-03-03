package html

import (
	"syscall/js"
)

type Node interface {
	DOMObject(document js.Value) js.Value
}
