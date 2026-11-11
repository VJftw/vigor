package vigor

// Map is a reactive map for collections that change one entry at a time,
// like entities arriving over a stream. Reading one key subscribes only to
// that key, so an update re-runs only what shows that entry; Keys, Values
// and Len subscribe to the set of entries changing. Keys keep the order
// they were first added in.
type Map[K comparable, V any] struct {
	items    map[K]V
	order    []K
	keySrc   map[K]*source
	shape    source
	valEqual func(a, b V) bool
}

// NewMap returns an empty Map. Options set how a replaced value is
// compared: equal values don't notify.
func NewMap[K comparable, V any](opts ...Option[V]) *Map[K, V] {
	return &Map[K, V]{
		items:    map[K]V{},
		keySrc:   map[K]*source{},
		valEqual: equalsFn(opts),
	}
}

func (m *Map[K, V]) src(k K) *source {
	s, ok := m.keySrc[k]
	if !ok {
		s = &source{}
		m.keySrc[k] = s
	}

	return s
}

// Get returns an entry, subscribing to it (also when it doesn't exist yet,
// so it re-runs when it appears).
func (m *Map[K, V]) Get(k K) (V, bool) {
	m.src(k).track()
	v, ok := m.items[k]

	return v, ok
}

// Set adds or replaces an entry.
func (m *Map[K, V]) Set(k K, v V) {
	Batch(func() {
		old, existed := m.items[k]
		if existed && safeEq(m.valEqual, old, v) {
			return
		}
		m.items[k] = v
		if !existed {
			m.order = append(m.order, k)
			m.shape.notify()
		}
		m.src(k).notify()
	})
}

// Delete removes an entry.
func (m *Map[K, V]) Delete(k K) {
	Batch(func() {
		if _, ok := m.items[k]; !ok {
			return
		}
		delete(m.items, k)
		for i, o := range m.order {
			if o == k {
				m.order = append(m.order[:i:i], m.order[i+1:]...)

				break
			}
		}
		m.shape.notify()
		m.src(k).notify()
	})
}

// Replace makes the map exactly the given entries, notifying only about
// what actually changed. Entries keep their place; new ones go last, in the
// order of keys (pass an ordered key list to control it).
func (m *Map[K, V]) Replace(entries map[K]V, keys ...K) {
	Batch(func() {
		for _, k := range append([]K(nil), m.order...) {
			if _, keep := entries[k]; !keep {
				m.Delete(k)
			}
		}
		for _, k := range keys {
			if v, ok := entries[k]; ok {
				m.Set(k, v)
			}
		}
		for k, v := range entries {
			m.Set(k, v)
		}
	})
}

// Clear removes every entry.
func (m *Map[K, V]) Clear() { m.Replace(nil) }

// Has reports whether an entry exists, subscribing to that entry.
func (m *Map[K, V]) Has(k K) bool {
	_, ok := m.Get(k)

	return ok
}

// Len is the number of entries, subscribing to entries being added and
// removed.
func (m *Map[K, V]) Len() int {
	m.shape.track()

	return len(m.items)
}

// Keys returns the keys in order, subscribing to entries being added and
// removed (not to their values changing).
func (m *Map[K, V]) Keys() []K {
	m.shape.track()

	return append([]K(nil), m.order...)
}

// Values returns the values in key order, subscribing to entries being
// added or removed, and to every value changing.
func (m *Map[K, V]) Values() []V {
	m.shape.track()
	out := make([]V, 0, len(m.order))
	for _, k := range m.order {
		m.src(k).track()
		out = append(out, m.items[k])
	}

	return out
}
