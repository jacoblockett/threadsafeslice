package threadsafeslice

import (
	"reflect"
	"slices"
	"sync"
	"testing"
)

func TestInitialize(t *testing.T) {
	result1 := Initialize([]string{"test"})
	result2 := Initialize([]string{"test"})

	if result1 == result2 {
		t.Error("expected result1 and result2 to be different entities")
	}
}

func TestLength(t *testing.T) {
	tss1 := Initialize([]int{1})
	tss2 := Initialize([]int{1, 2, 3})
	if tss1.Length() != 1 {
		t.Error("expected tss1 to have a length of 1")
	}
	if tss2.Length() != 3 {
		t.Error("expected tss2 to have a length of 3")
	}
}

func TestShift(t *testing.T) {
	tss := Initialize([]int{1, 2})
	v, isEmpty := tss.Shift()
	if v != 1 {
		t.Errorf("expected v (%d) to equal 1", v)
	}
	if isEmpty {
		t.Error("expected isEmpty to be false")
	}

	v, isEmpty = tss.Shift()
	if v != 2 {
		t.Errorf("expected v (%d) to equal 2", v)
	}
	if isEmpty {
		t.Error("expected isEmpty to be false")
	}

	_, isEmpty = tss.Shift()
	if !isEmpty {
		t.Error("expected isEmpty to be true")
	}
}

func TestPop(t *testing.T) {
	tss := Initialize([]int{1, 2})
	v, isEmpty := tss.Pop()
	if v != 2 {
		t.Errorf("expected v (%d) to equal 2", v)
	}
	if isEmpty {
		t.Error("expected isEmpty to be false")
	}

	v, isEmpty = tss.Pop()
	if v != 1 {
		t.Errorf("expected v (%d) to equal 1", v)
	}
	if isEmpty {
		t.Error("expected isEmpty to be false")
	}

	_, isEmpty = tss.Pop()
	if !isEmpty {
		t.Error("expected isEmpty to be true")
	}
	if tss.Length() != 0 {
		t.Error("expected length of slice to be 0")
	}
}

func TestAt(t *testing.T) {
	tss := Initialize([]int{1, 2})
	f := tss.At(0)
	l := tss.At(1)
	if f != 1 {
		t.Error("expected f to equal 1")
	}
	if l != 2 {
		t.Error("expected l to equal 2")
	}

	f = tss.At(-2)
	l = tss.At(-1)
	if f != 1 {
		t.Error("expected f (neg index) to equal 1")
	}
	if l != 2 {
		t.Error("expected l (neg index) to equal 2")
	}

	outFront := tss.At(-3)
	outEnd := tss.At(2)
	if outFront != 0 || outEnd != 0 {
		t.Error("expected outFront and outEnd to be zero values")
	}
}

func TestUnshift(t *testing.T) {
	tss := Initialize([]int{1, 2})
	tss.Unshift(0)
	if !slices.Equal(tss.Get(), []int{0, 1, 2}) {
		t.Error("tss.Unshift(0) failed")
	}

	tss.Unshift(1).Unshift(2)
	if !slices.Equal(tss.Get(), []int{2, 1, 0, 1, 2}) {
		t.Error("tss.Unshift(1).Unshift(2) failed")
	}

	tss.Unshift(7, 8, 9)
	if !slices.Equal(tss.Get(), []int{7, 8, 9, 2, 1, 0, 1, 2}) {
		t.Error("tss.Unshift(7, 8, 9) failed")
	}
}

func TestPush(t *testing.T) {
	tss := Initialize([]int{1, 2})
	tss.Push(0)
	if !slices.Equal(tss.Get(), []int{1, 2, 0}) {
		t.Error("tss.Push(0) failed")
	}

	tss.Push(1).Push(2)
	if !slices.Equal(tss.Get(), []int{1, 2, 0, 1, 2}) {
		t.Error("tss.Push(1).Push(2) failed")
	}

	tss.Push(7, 8, 9)
	if !slices.Equal(tss.Get(), []int{1, 2, 0, 1, 2, 7, 8, 9}) {
		t.Error("tss.Push(1).Push(2) failed")
	}
}

func TestClear(t *testing.T) {
	tss := Initialize([]int{1, 2, 3})
	tss.Clear()
	if tss.Length() != 0 {
		t.Error("expected the slice to be empty")
	}
}

func TestSet(t *testing.T) {
	tss := Initialize([]int{1, 2, 3})
	tss.Set([]int{3, 2, 1})
	if tss.At(0) != 3 {
		t.Error("expected the first element to be 3")
	}
}

func TestGet(t *testing.T) {
	tss := Initialize([]int{1, 2, 3})
	slice := tss.Get()
	if reflect.TypeOf(slice).Kind() != reflect.Slice {
		t.Error("expected slice to be a slice type")
	}
	if slice[0] != 1 {
		t.Error("expected the first element to be 1")
	}
}

func TestMap(t *testing.T) {
	tss := Initialize([]int{1, 2, 3})
	m := tss.Map(func(v, i int, s []int) int {
		return v + 1
	})
	if tss.At(-1) != 4 {
		t.Error("expected the last element to be 4")
	}
	if tss != m {
		t.Error("expected tss and m to be the same entity")
	}
}

func TestMapCopy(t *testing.T) {
	tss := Initialize([]int{1, 2, 3})
	m := tss.MapCopy(func(v, i int, s []int) int {
		return v + 1
	})
	if tss.At(-1) != 4 {
		t.Error("expected the last element to be 4")
	}
	if tss == m {
		t.Error("expected tss and m to be distinct entities")
	}
}

func TestThreadSafe(t *testing.T) {
	// not sure how naiive this test is, but it passes, so...... :)
	tss := Initialize([]int{})
	gc := 1000
	var wg sync.WaitGroup
	wg.Add(gc)

	for i := 0; i < gc; i++ {
		go func(val int) {
			defer wg.Done()
			tss.Push(val)
		}(i)
	}

	wg.Wait()

	finalLength := tss.Length()
	if finalLength != gc {
		t.Errorf("Expected length %d, got %d", gc, finalLength)
	}
}
