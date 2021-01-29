package stacktrace

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

const maxDepth = 10

type frame struct {
	file     string
	line     int
	function string
}

type stacktrace struct {
	wrapped error
	cause   error
	message string
	frames  []frame
}

func (s *stacktrace) Error() string {
	var b strings.Builder

	b.WriteString("Caused by: ")
	b.WriteString(s.message)

	if _, ok := s.cause.(*stacktrace); !ok {
		if s.cause != nil {
			b.WriteString(" <-- ")
			b.WriteString(s.cause.Error())
		}
	}

	b.WriteString("\n")

	for _, frame := range s.frames {
		if frame.file != "" && frame.line != 0 {
			b.WriteString("    at ")
			b.WriteString(frame.file)
			b.WriteString(":")
			b.WriteString(fmt.Sprint(frame.line))

			if frame.function != "" {
				b.WriteString(" (")
				b.WriteString(frame.function)
				b.WriteString(")")
			}

			b.WriteString("\n")
		}
	}

	if _, ok := s.cause.(*stacktrace); ok {
		buf := b.String()
		b.Reset()
		b.WriteString(buf)
		b.WriteString("\n")
		b.WriteString(s.cause.Error())
	}

	return b.String()
}

func (s *stacktrace) Wrapped() error {
	if e, ok := s.cause.(*stacktrace); ok {
		s.wrapped = fmt.Errorf("%v: %v", s.message, e.Wrapped())
	} else {
		s.wrapped = fmt.Errorf("%v: %v", s.message, s.cause)
	}

	return s.wrapped
}

// Propagate propagates errors up through the call stack.
func Propagate(err error, format string, args ...interface{}) error {
	pc := make([]uintptr, maxDepth)
	_ = runtime.Callers(2, pc[:])
	frames := runtime.CallersFrames(pc)

	st := &stacktrace{
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

		st.frames = append(st.frames, frame{
			file:     f.File,
			line:     f.Line,
			function: f.Function,
		})
	}

	return st
}

// Throw panics, recovers, and prints the stacktrace to standard out.
func Throw(err error) {
	if r := recover(); r == nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	panic(err)
}

// Error returns the wrapped error much like `github.com/pkg/errors`
// does (i.e. no stacktrace).
func Error(err error) error {
	if e, ok := err.(*stacktrace); ok {
		return e.Wrapped()
	}
	return err
}
