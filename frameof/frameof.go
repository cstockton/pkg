// Package frameof provides a few simple helpers to access call frames from the
// perspective of this packages caller.
//
// For example Call() returns the current frame of the call site, not the frame
// of the invocation of Call() within this package. This along with the packages
// name should make for intuitive usage, i.e.:
//
//   frameof.Call() // My current frame
//   frameof.Caller() // My caller
//   frameof.Skip(2) // My caller's caller.
//
// All top level functions in this package are safe for concurrent use and
// acquiring a frame will not allocate unless you choose to have it escape.
package frameof

import (
	"fmt"
	"path"
	"runtime"
	"strings"
)

// Frame is simply a runtime.Frame with a few convenience methods.
type Frame runtime.Frame

// At returns this frames file and line.
func (fr Frame) At() (file string, line int) {
	return fr.File, fr.Line
}

// Name returns this frames function name.
func (fr Frame) Name() string {
	base := path.Base(fr.Function)
	if i := strings.LastIndex(base, `.`); i != -1 {
		return base[i+1:]
	}
	return ``
}

// Pkg returns this frames package name.
func (fr Frame) Pkg() string {
	base := path.Base(fr.Function)
	if i := strings.LastIndex(base, `.`); i != -1 {
		return base[:i]
	}
	return base
}

// PkgPath returns this frames import path.
func (fr Frame) PkgPath() string {
	dir, file := path.Split(fr.Function)
	if `` != dir {
		return strings.TrimRight(dir, `/`)
	}
	if i := strings.LastIndex(file, `.`); i != -1 && i > 0 {
		return file[:i]
	}
	return file
}

// String implement fmt.Stringer.
func (fr Frame) String() string {
	return fmt.Sprintf("%v\n\t%v:%v", fr.Function, fr.File, fr.Line)
}

// Call returns the current call frame from the perspective of this functions
// caller.
func Call() Frame {
	var fr Frame
	// Note I purposely did not wrap these all with a single frame loader to
	// prevent needlessly ascending the call stack, each frame is around 100-200ns
	// depending on the current stack layout which accounts for about 1/4th of the
	// duration of this function call.
	fr.PC, fr.File, fr.Line, _ = runtime.Caller(1)
	loadFrame(&fr)
	return fr
}

// Caller returns the current call frame from the perspective of this functions
// caller.
func Caller() Frame {
	var fr Frame
	fr.PC, fr.File, fr.Line, _ = runtime.Caller(2)
	loadFrame(&fr)
	return fr
}

// Skip will ascend n stack frames starting from the current frame. For example
// calls to Call() and Caller() are equivalent to Skip(0) and Skip(1).
func Skip(n int) Frame {
	var fr Frame
	fr.PC, fr.File, fr.Line, _ = runtime.Caller(1 + n)
	loadFrame(&fr)
	return fr
}

// Callers will return all frames in the current call stack.
func Callers() []Frame {
	var pcs [32]uintptr
	n := runtime.Callers(1, pcs[:])
	if n == 0 {
		return nil
	}

	out, frames := make([]Frame, 0, n), runtime.CallersFrames(pcs[:n])
	for {
		fr, ok := frames.Next()
		if !ok {
			return out
		}
		out = append(out, Frame(fr))
	}
}

func loadFrame(fr *Frame) {
	fn := runtime.FuncForPC(fr.PC)
	if fn == nil {
		return
	}
	fr.Func, fr.Function, fr.Entry = fn, fn.Name(), fn.Entry()
	return
}
