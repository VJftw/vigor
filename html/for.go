package html

import (
	"context"
	"syscall/js"

	"github.com/VJftw/vigor"
)

type ForFn func(i int, v any) Node

type nodeFor struct {
	fn ForFn
	v  any

	fallback Node
}

type ForOpt func(*nodeFor)

func WithForFallback(n Node) ForOpt {
	return func(nf *nodeFor) {
		nf.fallback = n
	}
}

func For(v any, fn ForFn, opts ...ForOpt) Node {
	nf := &nodeFor{
		fn: fn,
		v:  vigor.Normalise(v),
	}

	for _, opt := range opts {
		opt(nf)
	}

	return nf
}

func (n *nodeFor) DOMObject(ctx context.Context, doc js.Value) js.Value {
	endObj := doc.Call("createComment", "vigor_for-end")
	startObj := doc.Call("createComment", "vigor_for-start")

	itemsCh := make(chan []any, 1)

	currentItemObjs := []any{}

	go func() {
		parentNode, err := getParentNode(ctx, endObj)
		if err != nil {
			return
		}

		parentNode.Call("insertBefore", startObj, endObj)

		for {
			select {
			case <-ctx.Done():
				return
			case newItemObjs := <-itemsCh:
				if len(newItemObjs) < 1 && n.fallback != nil {
					newItemObjs = append(newItemObjs, n.fallback.DOMObject(ctx, doc))
				}

				for _, newItemObj := range newItemObjs {
					parentNode.Call("insertBefore", newItemObj, endObj)
				}

				for _, currentItemObj := range currentItemObjs {
					currentItemObj.(js.Value).Call("remove")
				}

				currentItemObjs = newItemObjs
			}
		}
	}()

	subscriber := vigor.NewFnSubscriber()
	subscriber.SetFn(func() {
		items := []any{}
		if getterFn, ok := n.v.(vigor.GetterFn); ok {
			items = getterFn(subscriber).([]any)
		} else {
			items = n.v.([]any)
		}

		newItemObjs := make([]any, len(items))
		for i, item := range items {
			newItemObjs[i] = n.fn(i, item).DOMObject(ctx, doc)
		}

		itemsCh <- newItemObjs
	}).Run()

	return endObj
}

func getParentNode(ctx context.Context, node js.Value) (js.Value, error) {
	for {
		select {
		case <-ctx.Done():
			return js.Null(), ctx.Err()
		default:
			parentNode := node.Get("parentNode")
			if !parentNode.IsNull() {
				return parentNode, nil
			}
		}
	}
}
