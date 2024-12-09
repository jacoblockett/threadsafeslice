# ThreadSafeSlice

Some utility functions on a wrapper for generic slices with thread safety. I may implement all of the built-in functions from the `slices` package and some seen in other languages at some point, but this is it for now until they become relevant to my workflow.

### Installation

```bash
go get github.com/jacoblockett/threadsafeslice
```

You can read the godoc [here](https://pkg.go.dev/github.com/jacoblockett/threadsafeslice) for detailed documentation.

### Quickstart

```go
package main

import tss "github.com/jacoblockett/threadsafeslice"

func main() {
    s := tss.Initialize([]int{1, 2, 3})

	firstElem, _ := s.Shift() // 1
	lastElem, _ := s.Pop()    // 3

	fmt.Println(firstElem, lastElem, s.Length()) // 1 3 1 (length is 1 because Shift and Pop are destructive)

	s.Unshift(1).Push(3, 4, 5)

	fmt.Println(s.Get())                        // [1 2 3 4 5]
	fmt.Println(s.Set([]int{10, 11, 12}).Get()) // [10 11 12]
	fmt.Println(s.Map(func(v int, i int, s []int) int {
		return v + 1
	}).Get()) // [11 12 13]
}
```
