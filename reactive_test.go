package vigor_test

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/VJftw/vigor"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEffectTracksAutomatically(t *testing.T) {
	a, setA := vigor.NewSignal(1)
	b, setB := vigor.NewSignal(10)
	var seen []int
	vigor.CreateRoot(func(dispose func()) any {
		vigor.Effect(func() { seen = append(seen, a()+b()) })

		return nil
	})
	setA(2)
	setB(20)
	setB(20) // unchanged
	assert.Equal(t, []int{11, 12, 22}, seen)
}

func TestEffectDropsStaleDependencies(t *testing.T) {
	useA, setUseA := vigor.NewSignal(true)
	a, setA := vigor.NewSignal(0)
	b, setB := vigor.NewSignal(0)
	runs := 0
	vigor.CreateRoot(func(func()) any {
		vigor.Effect(func() {
			runs++
			if useA() {
				a()
			} else {
				b()
			}
		})

		return nil
	})
	setA(1)
	assert.Equal(t, 2, runs)
	setUseA(false)
	assert.Equal(t, 3, runs)
	setA(2) // no longer read
	assert.Equal(t, 3, runs)
	setB(1)
	assert.Equal(t, 4, runs)
}

func TestBatchRunsEffectsOnce(t *testing.T) {
	a, setA := vigor.NewSignal(1)
	b, setB := vigor.NewSignal(1)
	runs := 0
	vigor.CreateRoot(func(func()) any {
		vigor.Effect(func() { runs++; a(); b() })

		return nil
	})
	vigor.Batch(func() { setA(2); setB(2) })
	assert.Equal(t, 2, runs)
}

func TestMemoIsGlitchFree(t *testing.T) {
	n, setN := vigor.NewSignal(1)
	var log []string
	vigor.CreateRoot(func(func()) any {
		double := vigor.Memo(func() int { return n() * 2 })
		parity := vigor.Memo(func() bool { return double()%2 == 0 })
		vigor.Effect(func() { log = append(log, "n="+string(rune('0'+n()))+" d="+string(rune('0'+double()))) })
		runs := 0
		vigor.Effect(func() { parity(); runs++; log = append(log, "parity") })

		return nil
	})
	setN(2)
	// The memo's result changed, the parity memo's didn't: its effect stays put.
	assert.Equal(t, []string{"n=1 d=2", "parity", "n=2 d=4"}, log)
}

func TestMemoOnlyNotifiesOnChange(t *testing.T) {
	n, setN := vigor.NewSignal(1)
	runs := 0
	vigor.CreateRoot(func(func()) any {
		big := vigor.Memo(func() bool { return n() > 5 })
		vigor.Effect(func() { big(); runs++ })

		return nil
	})
	setN(2)
	setN(3)
	assert.Equal(t, 1, runs)
	setN(9)
	assert.Equal(t, 2, runs)
}

func TestCleanupAndNestedEffectsDisposeOnRerun(t *testing.T) {
	trigger, setTrigger := vigor.NewSignal(0)
	inner, setInner := vigor.NewSignal(0)
	var cleanups, innerRuns int
	vigor.CreateRoot(func(func()) any {
		vigor.Effect(func() {
			trigger()
			vigor.OnCleanup(func() { cleanups++ })
			vigor.Effect(func() { inner(); innerRuns++ })
		})

		return nil
	})
	setInner(1)
	assert.Equal(t, 2, innerRuns)
	setTrigger(1) // re-runs outer: old inner is disposed, a new one created
	assert.Equal(t, 1, cleanups)
	assert.Equal(t, 3, innerRuns)
	setInner(2) // only the new inner runs
	assert.Equal(t, 4, innerRuns)
}

func TestDisposeStopsEffects(t *testing.T) {
	s, set := vigor.NewSignal(0)
	runs, cleaned := 0, 0
	var dispose func()
	vigor.CreateRoot(func(d func()) any {
		dispose = d
		vigor.Effect(func() { s(); runs++ })
		vigor.OnCleanup(func() { cleaned++ })

		return nil
	})
	dispose()
	set(1)
	assert.Equal(t, 1, runs)
	assert.Equal(t, 1, cleaned)
}

func TestUntrack(t *testing.T) {
	a, setA := vigor.NewSignal(0)
	b, setB := vigor.NewSignal(0)
	runs := 0
	vigor.CreateRoot(func(func()) any {
		vigor.Effect(func() {
			a()
			vigor.Untrack(func() int { return b() })
			runs++
		})

		return nil
	})
	setB(1)
	assert.Equal(t, 1, runs)
	setA(1)
	assert.Equal(t, 2, runs)
}

func TestSliceSignalsAlwaysNotifyUnlessDeepEqual(t *testing.T) {
	s, set := vigor.NewSignal([]int{1})
	d, setD := vigor.NewSignal([]int{1}, vigor.DeepEqual[[]int]())
	var sRuns, dRuns int
	vigor.CreateRoot(func(func()) any {
		vigor.Effect(func() { s(); sRuns++ })
		vigor.Effect(func() { d(); dRuns++ })

		return nil
	})
	set([]int{1})
	setD([]int{1})
	assert.Equal(t, 2, sRuns)
	assert.Equal(t, 1, dRuns)
}

func TestInfiniteLoopPanics(t *testing.T) {
	n, setN := vigor.NewSignal(0)
	assert.PanicsWithValue(t, "vigor: effects keep triggering each other (infinite update loop)", func() {
		vigor.CreateRoot(func(func()) any {
			vigor.Effect(func() { setN(n() + 1) })

			return nil
		})
	})
}

func TestContext(t *testing.T) {
	theme := vigor.CreateContext("light")
	var got []string
	vigor.CreateRoot(func(func()) any {
		got = append(got, theme.Use())
		theme.Provide("dark", func() {
			got = append(got, theme.Use())
			vigor.Effect(func() { got = append(got, theme.Use()) })
			theme.Provide("neon", func() { got = append(got, theme.Use()) })
		})

		return nil
	})
	assert.Equal(t, []string{"light", "dark", "dark", "neon"}, got)
}

func TestOwnerContextCancelledOnDispose(t *testing.T) {
	var ctx context.Context
	var dispose func()
	vigor.CreateRoot(func(d func()) any { dispose, ctx = d, vigor.OwnerContext(); return nil })
	require.NoError(t, ctx.Err())
	dispose()
	assert.Error(t, ctx.Err())
}

func TestMap(t *testing.T) {
	m := vigor.NewMap[string, int]()
	var oneRuns, keysRuns int
	var keys []string
	vigor.CreateRoot(func(func()) any {
		vigor.Effect(func() { m.Get("a"); oneRuns++ })
		vigor.Effect(func() { keys = m.Keys(); keysRuns++ })

		return nil
	})
	m.Set("a", 1) // a appears: both run
	m.Set("b", 1) // only keys
	m.Set("b", 2) // neither (b's value isn't read by either)
	m.Set("a", 1) // unchanged
	m.Set("a", 5) // only the "a" reader
	assert.Equal(t, 3, oneRuns)
	assert.Equal(t, 3, keysRuns)
	assert.Equal(t, []string{"a", "b"}, keys)

	m.Replace(map[string]int{"b": 2, "c": 3})
	assert.Equal(t, []string{"b", "c"}, keys)
	_, ok := m.Get("a")
	assert.False(t, ok)
	assert.Equal(t, []int{2, 3}, vigor.Untrack(m.Values))
}

func TestResource(t *testing.T) {
	id, setID := vigor.NewSignal(1)
	var r *vigor.Resource[string]
	var stop func()
	t.Cleanup(func() { vigor.Dispatch(stop) })
	done := make(chan struct{}, 8)
	fetches := atomic.Int32{}
	vigor.Dispatch(func() {
		vigor.CreateRoot(func(d func()) any {
			stop = d
			r = vigor.NewResourceFor(id, func(ctx context.Context, id int) (string, error) {
				fetches.Add(1)
				defer func() { done <- struct{}{} }()
				if id == 3 {
					return "", errors.New("boom")
				}

				return string(rune('a' + id)), nil
			})

			return nil
		})
	})
	assert.True(t, read(r.Loading) || read(r.Loaded))
	<-done
	waitFor(t, func() bool { return !read(r.Loading) })
	assert.Equal(t, "b", read(r.Value))
	assert.True(t, read(r.Loaded))

	vigor.Dispatch(func() { setID(3) })
	<-done
	waitFor(t, func() bool { return read(r.Error) != nil })
	assert.Equal(t, "b", read(r.Value), "keeps the last value on error")

	vigor.Dispatch(r.Refetch)
	<-done
	assert.EqualValues(t, 3, fetches.Load())
}

// read and waitFor touch the graph under Dispatch, as goroutines must.
func read[T any](fn func() T) (v T) {
	vigor.Dispatch(func() { v = fn() })

	return v
}

func waitFor(t *testing.T, cond func() bool) {
	t.Helper()
	require.Eventually(t, cond, 2*time.Second, time.Millisecond)
}

func TestIntervalStopsWithScope(t *testing.T) {
	var n atomic.Int32
	var dispose func()
	vigor.Dispatch(func() {
		vigor.CreateRoot(func(d func()) any {
			dispose = d
			vigor.Interval(5*time.Millisecond, func() { n.Add(1) })

			return nil
		})
	})
	require.Eventually(t, func() bool { return n.Load() >= 2 }, time.Second, time.Millisecond)
	vigor.Dispatch(dispose)
	time.Sleep(20 * time.Millisecond)
	at := n.Load()
	time.Sleep(30 * time.Millisecond)
	assert.Equal(t, at, n.Load())
}

func TestStreamReconnects(t *testing.T) {
	var connects atomic.Int32
	var st *vigor.Stream[int]
	var stop func()
	t.Cleanup(func() { vigor.Dispatch(stop); time.Sleep(20 * time.Millisecond) })
	vigor.Dispatch(func() {
		vigor.CreateRoot(func(d func()) any {
			stop = d
			st = vigor.NewStream(func(ctx context.Context, emit func(int)) error {
				n := connects.Add(1)
				emit(int(n))

				return errors.New("dropped")
			}, vigor.WithReconnectDelay(time.Millisecond, 5*time.Millisecond))

			return nil
		})
	})
	require.Eventually(t, func() bool { return connects.Load() >= 3 }, 2*time.Second, time.Millisecond)
	assert.GreaterOrEqual(t, read(st.Latest), 2)
}
