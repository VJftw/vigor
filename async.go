package vigor

import (
	"context"
	"time"
)

// Resource is a value loaded asynchronously, like Solid's createResource.
// Its fields are signals: read them in effects and in rendering.
type Resource[T any] struct {
	// Value is the latest loaded value (the zero value until the first
	// load; it stays while reloading).
	Value func() T
	// Loading is true while a fetch is in flight.
	Loading func() bool
	// Error is the last fetch's error, nil if it succeeded.
	Error func() error
	// Loaded is true once a fetch has succeeded.
	Loaded func() bool

	refetch func()
	mutate  func(T)
}

// Refetch fetches again now.
func (r *Resource[T]) Refetch() { r.refetch() }

// Mutate sets the value without fetching (e.g. optimistic updates, or after
// a write returned the new value).
func (r *Resource[T]) Mutate(v T) { r.mutate(v) }

// ResourceOption configures a Resource.
type ResourceOption func(*resourceConfig)

type resourceConfig struct {
	poll time.Duration
}

// WithPoll fetches again every d after the last fetch finished.
func WithPoll(d time.Duration) ResourceOption {
	return func(c *resourceConfig) { c.poll = d }
}

// NewResource fetches once, then on Refetch (and on WithPoll's schedule).
// The context given to fetch is cancelled when a newer fetch starts or the
// current scope is disposed.
func NewResource[T any](fetch func(ctx context.Context) (T, error), opts ...ResourceOption) *Resource[T] {
	return NewResourceFor(func() struct{} { return struct{}{} },
		func(ctx context.Context, _ struct{}) (T, error) { return fetch(ctx) }, opts...)
}

// NewResourceFor fetches whenever source changes (source is tracked), so a
// resource can follow a route parameter or a selection. A source that
// returns an error-free zero value still fetches.
func NewResourceFor[S, T any](source func() S, fetch func(ctx context.Context, s S) (T, error), opts ...ResourceOption) *Resource[T] {
	cfg := resourceConfig{}
	for _, o := range opts {
		o(&cfg)
	}

	value, setValue := NewSignal[T](*new(T), AlwaysNotify[T]())
	loading, setLoading := NewSignal(true)
	failure, setFailure := NewSignal[error](nil)
	loaded, setLoaded := NewSignal(false)
	trigger, bump := NewSignal(0)

	owner := currentOwner
	ownerCtx := owner.Context()

	Effect(func() {
		s := source()
		trigger() // Refetch re-runs this effect.

		ctx, cancel := context.WithCancel(ownerCtx)
		OnCleanup(cancel) // a newer fetch, or disposal, cancels this one

		Untrack(func() struct{} { setLoading(true); return struct{}{} })
		go func() {
			var timer *time.Timer
			for {
				v, err := fetch(ctx, s)
				if ctx.Err() != nil {
					return
				}
				Dispatch(func() {
					if ctx.Err() != nil { // superseded while waiting
						return
					}
					if err != nil {
						setFailure(err)
					} else {
						setFailure(nil)
						setValue(v)
						setLoaded(true)
					}
					setLoading(false)
				})
				if cfg.poll <= 0 {
					return
				}
				timer = time.NewTimer(cfg.poll)
				select {
				case <-ctx.Done():
					timer.Stop()

					return
				case <-timer.C:
				}
				Dispatch(func() { setLoading(true) })
			}
		}()
	})

	r := &Resource[T]{Value: value, Loading: loading, Error: failure, Loaded: loaded}
	r.refetch = func() { bump(Untrack(trigger) + 1) }
	r.mutate = func(v T) { setValue(v) }

	return r
}

// Interval calls fn every d, in the current scope, until it's disposed. fn
// runs as if inside Dispatch, so it may set signals.
func Interval(d time.Duration, fn func()) {
	ctx := OwnerContext()
	go func() {
		t := time.NewTicker(d)
		defer t.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				Dispatch(func() {
					if ctx.Err() == nil {
						fn()
					}
				})
			}
		}
	}()
}

// Timeout calls fn once after d, unless the current scope is disposed
// first. The returned function cancels it.
func Timeout(d time.Duration, fn func()) (cancel func()) {
	ctx, stop := context.WithCancel(OwnerContext())
	go func() {
		t := time.NewTimer(d)
		defer t.Stop()
		select {
		case <-ctx.Done():
		case <-t.C:
			Dispatch(func() {
				if ctx.Err() == nil {
					fn()
				}
			})
		}
	}()

	return stop
}

// Clock is a signal of the current time, updated on every multiple of
// every (e.g. time.Second ticks on the second), for clocks and for "x
// minutes ago" text.
func Clock(every time.Duration) func() time.Time {
	now, setNow := NewSignal(time.Now(), AlwaysNotify[time.Time]())
	ctx := OwnerContext()
	go func() {
		for {
			n := time.Now()
			wait := n.Truncate(every).Add(every).Sub(n)
			select {
			case <-ctx.Done():
				return
			case <-time.After(wait):
			}
			t := time.Now()
			Dispatch(func() {
				if ctx.Err() == nil {
					setNow(t)
				}
			})
		}
	}()

	return now
}

// Stream is a value fed by a long-lived connection (server-sent events, a
// watch stream), reconnecting when it drops.
type Stream[T any] struct {
	// Latest is the last value emitted.
	Latest func() T
	// Connected is true while run is running (a connection is up or being
	// made).
	Connected func() bool
	// Error is why the last connection ended, nil while connected.
	Error func() error
}

// StreamOption configures a Stream.
type StreamOption func(*streamConfig)

type streamConfig struct {
	delay     time.Duration
	maxDelay  time.Duration
	onConnect func()
	stopOn    func(error) bool
}

// WithReconnectDelay sets the wait before reconnecting (it doubles after
// each failure, up to max; the default is 1s up to 30s).
func WithReconnectDelay(initial, max time.Duration) StreamOption {
	return func(c *streamConfig) { c.delay, c.maxDelay = initial, max }
}

// OnConnect calls fn, as in Dispatch, each time a connection starts (e.g.
// to clear state a snapshot is about to refill).
func OnConnect(fn func()) StreamOption {
	return func(c *streamConfig) { c.onConnect = fn }
}

// StopOn stops reconnecting when the connection ends with an error for
// which fn is true (e.g. "unauthorised").
func StopOn(fn func(error) bool) StreamOption {
	return func(c *streamConfig) { c.stopOn = fn }
}

// NewStream calls run, in a goroutine, until the current scope is disposed,
// calling it again after it returns (with backoff). run reads its stream
// and calls emit for each value; emit is safe to call from run.
func NewStream[T any](run func(ctx context.Context, emit func(T)) error, opts ...StreamOption) *Stream[T] {
	cfg := streamConfig{delay: time.Second, maxDelay: 30 * time.Second}
	for _, o := range opts {
		o(&cfg)
	}

	latest, setLatest := NewSignal[T](*new(T), AlwaysNotify[T]())
	connected, setConnected := NewSignal(false)
	failure, setFailure := NewSignal[error](nil)
	ctx := OwnerContext()

	go func() {
		delay := cfg.delay
		for ctx.Err() == nil {
			Dispatch(func() {
				setConnected(true)
				setFailure(nil)
				if cfg.onConnect != nil {
					cfg.onConnect()
				}
			})
			start := time.Now()
			err := run(ctx, func(v T) {
				Dispatch(func() {
					if ctx.Err() == nil {
						setLatest(v)
					}
				})
			})
			if ctx.Err() != nil {
				return
			}
			Dispatch(func() {
				setConnected(false)
				setFailure(err)
			})
			if err != nil && cfg.stopOn != nil && cfg.stopOn(err) {
				return
			}
			if time.Since(start) > cfg.maxDelay { // it had been healthy
				delay = cfg.delay
			}
			select {
			case <-ctx.Done():
				return
			case <-time.After(delay):
			}
			delay = min(delay*2, cfg.maxDelay)
		}
	}()

	return &Stream[T]{Latest: latest, Connected: connected, Error: failure}
}
