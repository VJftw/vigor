package vigor

// CreateEffect calls the given `fn` when any of the given `deps` change.
func CreateEffect(fn func(), deps ...GetterFn) {
	if len(deps) < 1 {
		Log.Fatal("no deps means the effect will never run")
	}

	subscriber := NewFnSubscriber().SetFn(fn)
	for _, dep := range deps {
		dep(subscriber)
	}
}
