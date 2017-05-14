package frameof

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"
)

func TestCreate(t *testing.T) {
	tests := []struct {
		fr Frame
		ln int
		fn string
	}{
		{Call(), 16, `github.com/cstockton/pkg/frameof.TestCreate`},
		{testCaller(), 17, `github.com/cstockton/pkg/frameof.TestCreate`},
		{testSkip(1, 1), 18, `github.com/cstockton/pkg/frameof.TestCreate`},
		{testSkip(2, 2), 19, `github.com/cstockton/pkg/frameof.TestCreate`},
		{testSkip(3, 3), 20, `github.com/cstockton/pkg/frameof.TestCreate`},
		{testCallers(t), 21, `github.com/cstockton/pkg/frameof.TestCreate`},
	}
	for _, test := range tests {
		fr := test.fr

		if exp, got := test.ln, fr.Line; exp != got {
			t.Fatalf(`exp line %v; got %v`, exp, got)
		}
		if exp, got := test.fn, fr.Function; exp != got {
			t.Fatalf(`exp Function field %v; got %v`, exp, got)
		}
		if exp, got := `frameof_test.go`, filepath.Base(fr.File); exp != got {
			t.Fatalf(`exp File field %v; got %v`, exp, got)
		}
		if 0 == fr.Entry {
			t.Fatal(`exp non-zero Entry field`)
		}
		if got := fr.String(); !strings.Contains(got, fr.File) {
			t.Fatalf(`exp String() return value %q to contain %v`, got, fr.File)
		}
		exp := fmt.Sprint(`:`, fr.Line)
		if got := fr.String(); !strings.Contains(got, exp) {
			t.Fatalf(`exp String() return value %q to contain %v`, got, exp)
		}
	}
}

func TestAllocations(t *testing.T) {
	const (
		expFn = `github.com/cstockton/pkg/frameof.TestAllocations.func1`
		expLn = 54
	)
	allocs := testing.AllocsPerRun(1000, func() {
		fr := Call()
		if exp, got := expLn, fr.Line; exp != got {
			t.Fatalf(`exp line %v; got %v`, exp, got)
		}
		if exp, got := expFn, fr.Function; exp != got {
			t.Fatalf(`exp Function field %v; got %v`, exp, got)
		}
		fr = testCaller()
		if exp, got := expLn+7, fr.Line; exp != got {
			t.Fatalf(`exp line %v; got %v`, exp, got)
		}
		if exp, got := expFn, fr.Function; exp != got {
			t.Fatalf(`exp Function field %v; got %v`, exp, got)
		}
	})
	if allocs != 0 {
		t.Fatalf(`exp zero allocs; got %v`, allocs)
	}
}

func TestNegative(t *testing.T) {
	if exp, got := uintptr(0), Skip(10000).PC; exp != got {
		t.Errorf(`exp fromCaller(10000) return value %v; got %v`, exp, got)
	}
	var fr Frame
	loadFrame(&fr)
	if exp, got := uintptr(0), fr.PC; exp != got {
		t.Errorf(`exp loadFrame(0).PC return value %v; got %v`, exp, got)
	}
}

func TestFrame(t *testing.T) {
	tests := []struct {
		given, path, pkg, name string
	}{
		{`example.org/path/of/pkg.Func`,
			`example.org/path/of`, `pkg`, `Func`},
		{`example.org/path/of/pkg.F`,
			`example.org/path/of`, `pkg`, `F`},
		{`example.org/path/of/pkg.`,
			`example.org/path/of`, `pkg`, ``},
		{`example.org/path/of/pkg`,
			`example.org/path/of`, `pkg`, ``},
		{`/path/of/pkg.Func`, `/path/of`, `pkg`, `Func`},
		{`/path/of/pkg.F`, `/path/of`, `pkg`, `F`},
		{`/path/of/pkg.`, `/path/of`, `pkg`, ``},
		{`/path/of/pkg`, `/path/of`, `pkg`, ``},
		{`../of/pkg.Func`, `../of`, `pkg`, `Func`},
		{`../of/pkg.F`, `../of`, `pkg`, `F`},
		{`../of/pkg.`, `../of`, `pkg`, ``},
		{`../of/pkg`, `../of`, `pkg`, ``},
		{`of/pkg.Func`, `of`, `pkg`, `Func`},
		{`of/pkg.F`, `of`, `pkg`, `F`},
		{`of/pkg.`, `of`, `pkg`, ``},
		{`of/pkg`, `of`, `pkg`, ``},
		{`pkg.Func`, `pkg`, `pkg`, `Func`},
		{`pkg.F`, `pkg`, `pkg`, `F`},
		{`pkg.`, `pkg`, `pkg`, ``},
		{`pkg`, `pkg`, `pkg`, ``},
		{`p.`, `p`, `p`, ``},
		{`p`, `p`, `p`, ``},
		{`.`, `.`, ``, ``},
	}
	t.Run(`Name`, func(t *testing.T) {
		for idx, test := range tests {
			t.Logf(`test #%v - given %q ...`, idx, test.given)
			fr := Frame{Function: test.given}
			if exp, got := test.name, fr.Name(); exp != got {
				t.Errorf(`exp Name() return value %q; got %q`, exp, got)
			}
		}
	})
	t.Run(`PkgName`, func(t *testing.T) {
		for idx, test := range tests {
			t.Logf(`test #%v - given %q ...`, idx, test.given)
			fr := Frame{Function: test.given}
			if exp, got := test.pkg, fr.Pkg(); exp != got {
				t.Errorf(`exp PkgName() return value %q; got %q`, exp, got)
			}
		}
	})
	t.Run(`PkgPath`, func(t *testing.T) {
		for idx, test := range tests {
			t.Logf(`test #%v - given %q ...`, idx, test.given)
			fr := Frame{Function: test.given}
			if exp, got := test.path, fr.PkgPath(); exp != got {
				t.Errorf(`exp PkgPath() return value %q; got %q`, exp, got)
			}
		}
	})
}

func testCallers(t *testing.T) Frame {
	frames := Callers()
	if exp, got := 3, len(frames); exp > got {
		t.Fatalf(`exp %v frames; got %v`, exp, got)
	}
	return frames[2]
}

func testCaller() Frame {
	return Caller()
}

func testSkip(n, depth int) Frame {
	if depth > 1 {
		return testSkip(n, depth-1)
	}
	return Skip(n)
}
