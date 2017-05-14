package frameof_test

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/cstockton/pkg/frameof"
)

func Example() {
	fr := frameof.Call()

	// Print with filepath.Base to strip path for example.
	fmt.Printf("Frame in %v() at %v:%d for package %v.\n",
		fr.Name(), filepath.Base(fr.File), fr.Line, fr.Pkg())

	// Output:
	// Frame in Example() at example_test.go:12 for package frameof_test.
}

func Example_allocs() {
	allocs := testing.AllocsPerRun(1000, func() {
		fr := frameof.Call()
		if exp := 24; exp != fr.Line {
			fmt.Printf("expected line %v; not %v\n", exp, fr.Line)
		}
	})
	fmt.Printf("Acquiring a frame took %v allocations.\n", allocs)

	// Output:
	// Acquiring a frame took 0 allocations.
}

func Example_caller() {
	fn := func() frameof.Frame {
		return frameof.Caller()
	}

	fr := fn()
	fmt.Printf("Frame in %v() at %v:%d for package %v.\n",
		fr.Name(), filepath.Base(fr.File), fr.Line, fr.Pkg())

	// Output:
	// Frame in Example_caller() at example_test.go:40 for package frameof_test.
}

func call(n, depth int, fn func() frameof.Frame) frameof.Frame {
	if depth > 1 {
		return call(n, depth-1, fn)
	}
	return fn()
}

func Example_skip() {
	fr := call(0, 4, func() frameof.Frame {
		return frameof.Skip(5)
	})
	fmt.Printf("Frame in %v() at %v for package %v.\n",
		fr.Name(), filepath.Base(fr.File), fr.Pkg())

	// Output:
	// Frame in Example_skip() at example_test.go for package frameof_test.
}

func calls(n, depth int, fn func() []frameof.Frame) []frameof.Frame {
	if depth > 1 {
		return calls(n, depth-1, fn)
	}
	return fn()
}

func Example_callers() {
	frs := calls(0, 4, func() []frameof.Frame {
		return frameof.Callers()
	})
	for _, fr := range frs[:5] {
		fmt.Printf("Frame in %v() in file %v.\n",
			fr.Name(), filepath.Base(fr.File))
	}

	// Output:
	// Frame in Callers() in file frameof.go.
	// Frame in func1() in file example_test.go.
	// Frame in calls() in file example_test.go.
	// Frame in calls() in file example_test.go.
	// Frame in calls() in file example_test.go.
}
