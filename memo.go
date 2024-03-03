package vigor

// Memo creates a new Signal that depends on the given deps. This can be used
// for combining other Signals or caching costly calculations.
func CreateMemo(fn func() any, deps ...GetterFn) GetterFn {
	getMemo, setMemo := CreateSignal(fn())

	CreateEffect(func() {
		setMemo(fn())
	}, deps...)

	return getMemo
}
