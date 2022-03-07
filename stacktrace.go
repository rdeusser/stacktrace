package stacktrace

import (
	"fmt"
	"os"
	"runtime"
)

const maxDepth = 5

// Propagate propagates errors up through the call stack.
func Propagate(err error, format string, args ...interface{}) *Error {
	pc := make([]uintptr, maxDepth)
	_ = runtime.Callers(2, pc[:])
	frames := runtime.CallersFrames(pc)

	serr := &Error{
		wrapped: err,
		cause:   err,
		message: fmt.Sprintf(format, args...),
		frames:  make([]frame, 0),
	}

	for {
		f, more := frames.Next()
		if !more {
			break
		}

		serr.frames = append(serr.frames, frame{
			file:     f.File,
			function: f.Function,
			line:     f.Line,
		})
	}

	return serr
}

// Throw panics, recovers, and prints the stacktrace to standard out.
func Throw(err error) {
	if r := recover(); r == nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	panic(err)
}

// Unwrap returns the wrapped error much like `github.com/pkg/errors`
// does (i.e. no stacktrace).
func Unwrap(err error) error {
	if e, ok := err.(*Error); ok {
		return e.Wrapped()
	}
	return err
}
