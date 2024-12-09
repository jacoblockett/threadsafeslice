package threadsafeslice

import (
	"sort"
	"sync"
)

type ThreadSafeSlice[T any] struct {
	mu    sync.Mutex
	slice []T
}

// Mapping callback that passes in a (v)alue, (i)ndex, and (s)lice.
type MapCallback[T any] func(v T, i int, s []T) T

type SortComparatee[T any] struct {
	Index int
	Value T
}

// Sorting callback that passes in comparison for side a and side b as integers.
type SortCallback[T any] func(a, b SortComparatee[T]) bool

// A boolean representing if a slice is empty or not.
type IsEmpty bool

// Removes and returns the first value of the slice. IsEmpty is true if
// no more values can be shifted from the slice.
func (t *ThreadSafeSlice[T]) Shift() (T, IsEmpty) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if len(t.slice) == 0 {
		var zero T
		return zero, true
	}

	v := t.slice[0]
	t.slice = t.slice[1:]

	return v, false
}

// Inserts the given value(s) at the beginning of the slice. Returns the
// slice for chaining.
func (t *ThreadSafeSlice[T]) Unshift(v ...T) *ThreadSafeSlice[T] {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.slice = append(v, t.slice...)

	return t
}

// Removes and returns the last value of the slice. IsEmpty is true if
// no more values can be popped from the slice.
func (t *ThreadSafeSlice[T]) Pop() (T, IsEmpty) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if len(t.slice) == 0 {
		var zero T
		return zero, true
	}

	v := t.slice[len(t.slice)-1]
	t.slice = t.slice[:len(t.slice)-1]

	return v, false
}

// Inserts the given value(s) at the end of the slice. Returns the slice for chaining.
func (t *ThreadSafeSlice[T]) Push(v ...T) *ThreadSafeSlice[T] {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.slice = append(t.slice, v...)

	return t
}

// Clears the content of the slice. Returns the slice for chaining.
func (t *ThreadSafeSlice[T]) Clear() *ThreadSafeSlice[T] {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.slice = []T{}

	return t
}

// Sets the underlying slice to the given slice. Returns the slice for chaining.
func (t *ThreadSafeSlice[T]) Set(s []T) *ThreadSafeSlice[T] {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.slice = s

	return t
}

// Returns a copied snapshot of the underlying slice.
func (t *ThreadSafeSlice[T]) Get() []T {
	t.mu.Lock()
	defer t.mu.Unlock()

	snap := make([]T, len(t.slice))
	copy(snap, t.slice)

	return snap
}

// Returns the value at the given index. Negative indices will map from the
// end of the slice (i.e. -1 is the last element, -2 the second to last, and so on).
func (t *ThreadSafeSlice[T]) At(i int) T {
	t.mu.Lock()
	defer t.mu.Unlock()

	if i < 0 {
		i = len(t.slice) + i
	}
	if i < 0 || i > len(t.slice)-1 {
		var zero T
		return zero
	}

	return t.slice[i]
}

// Returns the length of the slice.
func (t *ThreadSafeSlice[T]) Length() int {
	t.mu.Lock()
	defer t.mu.Unlock()

	return len(t.slice)
}

// Maps over the slice, replacing each value with the result
// of the given callback. Returns the slice for chaining.
func (t *ThreadSafeSlice[T]) Map(callback MapCallback[T]) *ThreadSafeSlice[T] {
	t.mu.Lock()
	defer t.mu.Unlock()

	for i, v := range t.slice {
		snap := make([]T, len(t.slice))
		copy(snap, t.slice)

		t.slice[i] = callback(v, i, snap)
	}

	return t
}

// Maps over the slice, replacing each value with the result
// of the given callback. Does not affect the original slice.
// Returns a new *ThreadSafeSlice[T], distinct from the original.
func (t *ThreadSafeSlice[T]) MapCopy(callback MapCallback[T]) *ThreadSafeSlice[T] {
	t.mu.Lock()
	defer t.mu.Unlock()

	tss := Initialize(t.slice)

	return tss.Map(callback)
}

// Sorts the slice using `sort.Slice`. Returns the slice for chaining.
func (t *ThreadSafeSlice[T]) Sort(callback SortCallback[T]) *ThreadSafeSlice[T] {
	t.mu.Lock()
	defer t.mu.Unlock()

	sort.Slice(t.slice, func(i, j int) bool {
		a := SortComparatee[T]{Index: i, Value: t.slice[i]}
		b := SortComparatee[T]{Index: j, Value: t.slice[j]}

		return callback(a, b)
	})

	return t
}

// Sorts the slice using `sort.Slice`. Does not affect the original slice.
// Returns a new *ThreadSafeSlice[T], distinct from the original.
func (t *ThreadSafeSlice[T]) SortCopy(callback SortCallback[T]) *ThreadSafeSlice[T] {
	t.mu.Lock()
	defer t.mu.Unlock()

	tss := Initialize(t.slice)

	return tss.Sort(callback)
}

// Initializes a new *ThreadSafeSlice[T].
func Initialize[T any](s []T) *ThreadSafeSlice[T] {
	tss := &ThreadSafeSlice[T]{
		slice: make([]T, len(s)),
	}
	copy(tss.slice, s)

	return tss
}
