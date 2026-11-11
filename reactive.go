package vigor

import (
	"context"
	"reflect"
	"sync"
)

// This file is vigor's auto-tracking reactive core, modelled on SolidJS:
//
//   - Signals hold values. Reading one inside an Effect or Memo subscribes
//     that computation to it; no dependency lists.
//   - Effects and Memos re-run when what they read changes, in creation
//     order, once per Batch.
//   - Everything created while a computation runs is owned by it, and is
//     disposed (OnCleanup callbacks run, nested effects stop) before it
//     re-runs and when its owner is disposed.
//
// The graph is not goroutine-safe: use it from one goroutine at a time. In
// WebAssembly that is naturally so; goroutines that report back (timers,
// fetches, streams) must do so through Dispatch.

var (
	currentOwner    *Owner
	currentObserver *computation
	batchDepth      int
	pending         = map[*computation]struct{}{}
	nextID          uint64
	dispatchMu      sync.Mutex
)

// maxFlushPasses guards against effects that keep triggering each other.
const maxFlushPasses = 100_000

// Owner is a scope that owns effects, cleanups and child scopes. It is
// disposed as a unit.
type Owner struct {
	parent   *Owner
	children []*Owner
	cleanups []func()
	contexts map[any]any
	disposed bool

	ctx    context.Context
	cancel context.CancelFunc
}

func newOwner(parent *Owner) *Owner {
	o := &Owner{parent: parent}
	if parent != nil {
		parent.children = append(parent.children, o)
	}

	return o
}

// CurrentOwner returns the scope code is currently running in (nil outside
// any).
func CurrentOwner() *Owner { return currentOwner }

// Disposed reports whether the scope has been disposed.
func (o *Owner) Disposed() bool { return o == nil || o.disposed }

// Dispose runs the scope's cleanups (children first, then its own, in
// reverse order of registration) and stops everything it owns.
func (o *Owner) Dispose() {
	if o == nil || o.disposed {
		return
	}
	o.disposed = true

	for i := len(o.children) - 1; i >= 0; i-- {
		o.children[i].Dispose()
	}
	o.children = nil

	if o.cancel != nil {
		o.cancel()
	}

	// Run each cleanup untracked and unowned, so cleanups can read signals
	// without subscribing to them.
	prevObs, prevOwn := currentObserver, currentOwner
	currentObserver, currentOwner = nil, nil
	defer func() { currentObserver, currentOwner = prevObs, prevOwn }()

	cleanups := o.cleanups
	o.cleanups = nil
	for i := len(cleanups) - 1; i >= 0; i-- {
		cleanups[i]()
	}

	if o.parent != nil {
		p := o.parent
		for i, c := range p.children {
			if c == o {
				p.children = append(p.children[:i], p.children[i+1:]...)

				break
			}
		}
	}
}

// Context returns a context.Context that is cancelled when the scope is
// disposed; hand it to goroutines and requests the scope starts.
func (o *Owner) Context() context.Context {
	if o == nil {
		return context.Background()
	}
	if o.ctx == nil {
		o.ctx, o.cancel = context.WithCancel(context.Background())
		if o.disposed {
			o.cancel()
		}
	}

	return o.ctx
}

// OwnerContext is CurrentOwner().Context().
func OwnerContext() context.Context { return currentOwner.Context() }

// CreateRoot runs fn in a new top-level scope that isn't disposed with
// whatever scope it's created in. fn gets the function that disposes it.
func CreateRoot[T any](fn func(dispose func()) T) T {
	o := newOwner(nil)
	prevObs, prevOwn := currentObserver, currentOwner
	currentObserver, currentOwner = nil, o
	defer func() { currentObserver, currentOwner = prevObs, prevOwn }()

	return fn(o.Dispose)
}

// CreateScope runs fn in a new scope owned by the current one, returning
// the scope so it can be disposed early (it is disposed with its parent
// anyway). Unlike CreateRoot it keeps the current context values.
func CreateScope(fn func()) *Owner {
	o := newOwner(currentOwner)
	RunWithOwner(o, fn)

	return o
}

// RunWithOwner runs fn as part of the given scope (e.g. later, from a
// callback, to create effects that belong to a component), untracked.
func RunWithOwner(o *Owner, fn func()) {
	prevObs, prevOwn := currentObserver, currentOwner
	currentObserver, currentOwner = nil, o
	defer func() { currentObserver, currentOwner = prevObs, prevOwn }()
	fn()
}

// OnCleanup registers fn to run when the current scope is disposed, or
// before the current effect re-runs.
func OnCleanup(fn func()) {
	if currentOwner == nil {
		return
	}
	if currentOwner.disposed {
		fn()

		return
	}
	currentOwner.cleanups = append(currentOwner.cleanups, fn)
}

// Untrack runs fn without subscribing the current computation to anything
// fn reads.
func Untrack[T any](fn func() T) T {
	prev := currentObserver
	currentObserver = nil
	defer func() { currentObserver = prev }()

	return fn()
}

// Batch groups updates: effects run once, when the outermost Batch ends.
func Batch(fn func()) {
	batchDepth++
	func() {
		defer func() {
			batchDepth--
			if batchDepth == 0 {
				flush()
			}
		}()
		fn()
	}()
}

// Dispatch runs fn from a goroutine other than the one running the
// interface (a timer, a fetch), serialised with other Dispatches, and
// batched. Don't call it from reactive code.
func Dispatch(fn func()) {
	dispatchMu.Lock()
	defer dispatchMu.Unlock()
	Batch(fn)
}

// computation is an effect or a memo.
type computation struct {
	id      uint64
	fn      func()
	sources map[*source]struct{}
	owner   *Owner // the scope it creates its own children in
	parent  *Owner // the scope it belongs to
	dead    bool
}

type source struct {
	observers map[*computation]struct{}
}

func (s *source) track() {
	c := currentObserver
	if c == nil || c.dead {
		return
	}
	if s.observers == nil {
		s.observers = map[*computation]struct{}{}
	}
	s.observers[c] = struct{}{}
	c.sources[s] = struct{}{}
}

func (s *source) notify() {
	for c := range s.observers {
		pending[c] = struct{}{}
	}
	if batchDepth == 0 {
		flush()
	}
}

func (c *computation) clear() {
	for s := range c.sources {
		delete(s.observers, c)
	}
	c.sources = map[*source]struct{}{}
	if c.owner != nil {
		c.owner.Dispose()
		c.owner = nil
	}
}

func (c *computation) run() {
	if c.dead {
		return
	}
	c.clear()
	c.owner = newOwner(c.parent)

	prevObs, prevOwn := currentObserver, currentOwner
	currentObserver, currentOwner = c, c.owner
	defer func() { currentObserver, currentOwner = prevObs, prevOwn }()
	c.fn()
}

func (c *computation) dispose() {
	c.dead = true
	c.clear()
	delete(pending, c)
}

func newComputation(fn func()) *computation {
	nextID++

	return &computation{id: nextID, fn: fn, sources: map[*source]struct{}{}, parent: currentOwner}
}

var flushing bool

func flush() {
	if flushing {
		return
	}
	flushing = true
	defer func() { flushing = false }()

	// Run the earliest-created pending computation first, so a memo is
	// brought up to date before the effects that read it, and each
	// computation runs once however many of its sources changed.
	for runs := 0; len(pending) > 0; runs++ {
		if runs > maxFlushPasses {
			clear(pending)
			panic("vigor: effects keep triggering each other (infinite update loop)")
		}
		var next *computation
		for c := range pending {
			if next == nil || c.id < next.id {
				next = c
			}
		}
		delete(pending, next)
		next.run()
	}
}

// Effect runs fn now, and again whenever a signal it read changes. What fn
// creates, and what it registers with OnCleanup, is disposed before each
// re-run. It stops when the scope it was created in is disposed; the
// returned function stops it earlier.
func Effect(fn func()) (stop func()) {
	c := newComputation(fn)
	OnCleanup(c.dispose)
	// Run untracked by any outer computation; the effect tracks for itself.
	Untrack(func() struct{} { c.run(); return struct{}{} })

	return c.dispose
}

// Option configures a Signal or Memo.
type Option[T any] func(*signalConfig[T])

type signalConfig[T any] struct {
	equals func(a, b T) bool
}

// WithEquals sets how a new value is compared with the old; equal values
// don't notify. The default is == for comparable types, and "always
// different" for slices, maps and the like (which can't be compared).
func WithEquals[T any](eq func(a, b T) bool) Option[T] {
	return func(c *signalConfig[T]) { c.equals = eq }
}

// AlwaysNotify makes every Set notify, even with an equal value.
func AlwaysNotify[T any]() Option[T] {
	return WithEquals(func(T, T) bool { return false })
}

// DeepEqual compares values with reflect.DeepEqual, so a slice or struct set
// to an equal copy doesn't notify.
func DeepEqual[T any]() Option[T] {
	return WithEquals(func(a, b T) bool { return reflect.DeepEqual(a, b) })
}

func defaultEquals[T any]() func(a, b T) bool {
	if reflect.TypeFor[T]().Comparable() {
		return func(a, b T) bool {
			// An interface-typed T can still hold an uncomparable value.
			defer func() { _ = recover() }()

			return any(a) == any(b)
		}
	}

	return func(T, T) bool { return false }
}

func equalsFn[T any](opts []Option[T]) func(a, b T) bool {
	cfg := signalConfig[T]{}
	for _, o := range opts {
		o(&cfg)
	}
	if cfg.equals == nil {
		return defaultEquals[T]()
	}

	return cfg.equals
}

// NewSignal returns a getter and a setter for a value. Reading it inside an
// Effect or Memo subscribes that computation; setting it (to a different
// value) re-runs them.
func NewSignal[T any](initial T, opts ...Option[T]) (get func() T, set func(T)) {
	value := initial
	src := &source{}
	eq := equalsFn(opts)

	get = func() T {
		src.track()

		return value
	}
	set = func(next T) {
		if eq != nil && safeEq(eq, value, next) {
			return
		}
		value = next
		src.notify()
	}

	return get, set
}

func safeEq[T any](eq func(a, b T) bool, a, b T) (equal bool) {
	defer func() {
		if recover() != nil {
			equal = false
		}
	}()

	return eq(a, b)
}

// Update sets a signal from its current value, without subscribing to it.
func Update[T any](get func() T, set func(T), fn func(T) T) {
	set(fn(Untrack(get)))
}

// Memo is a signal computed from other signals: fn re-runs when they
// change, and what reads the memo is notified only if the result differs.
func Memo[T any](fn func() T, opts ...Option[T]) func() T {
	var value T
	src := &source{}
	eq := equalsFn(opts)
	first := true

	c := newComputation(nil)
	c.fn = func() {
		next := fn()
		if first {
			first = false
			value = next

			return
		}
		if safeEq(eq, value, next) {
			return
		}
		value = next
		src.notify()
	}
	OnCleanup(c.dispose)
	Untrack(func() struct{} { c.run(); return struct{}{} })

	return func() T {
		src.track()

		return value
	}
}

// On makes an effect run only when deps change: deps is tracked, fn is
// not. Like Solid's on().
func On(deps func(), fn func()) func() {
	return Effect(func() {
		deps()
		Untrack(func() struct{} { fn(); return struct{}{} })
	})
}

// --- context ---------------------------------------------------------

// Context passes values down the scope tree without threading them through
// every component: Provide it above, Use it below.
type Context[T any] struct {
	key *byte
	def T
}

// CreateContext makes a context whose value is def where nothing provides
// one.
func CreateContext[T any](def T) *Context[T] {
	return &Context[T]{key: new(byte), def: def}
}

// Provide runs fn in a child scope where Use returns value.
func (c *Context[T]) Provide(value T, fn func()) {
	o := newOwner(currentOwner)
	o.contexts = map[any]any{c.key: value}
	RunWithOwner(o, fn)
}

// Use returns the nearest provided value, or the default.
func (c *Context[T]) Use() T {
	for o := currentOwner; o != nil; o = o.parent {
		if v, ok := o.contexts[c.key]; ok {
			return v.(T)
		}
	}

	return c.def
}
