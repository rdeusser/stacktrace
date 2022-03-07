package stacktrace

import (
	"fmt"
	"strings"
)

type Error struct {
	wrapped error
	cause   error
	message string
	frames  []frame
}

func (e *Error) Error() string {
	var sb strings.Builder

	sb.WriteString("Caused by: ")
	sb.WriteString(e.message)

	if _, ok := e.cause.(*Error); !ok {
		if e.cause != nil {
			sb.WriteString(" <-- ")
			sb.WriteString(e.cause.Error())
		}
	}

	sb.WriteString("\n")

	for _, frame := range e.frames {
		if frame.file != "" && frame.line != 0 {
			sb.WriteString("    at ")
			sb.WriteString(frame.file)
			sb.WriteString(":")
			sb.WriteString(fmt.Sprint(frame.line))

			if frame.function != "" {
				sb.WriteString(" (")
				sb.WriteString(frame.function)
				sb.WriteString(")")
			}

			sb.WriteString("\n")
		}
	}

	if _, ok := e.cause.(*Error); ok {
		buf := sb.String()
		sb.Reset()
		sb.WriteString(buf)
		sb.WriteString("\n")
		sb.WriteString(e.cause.Error())
	}

	return sb.String()
}

func (e *Error) Wrapped() error {
	if err, ok := e.cause.(*Error); ok {
		e.wrapped = fmt.Errorf("%v: %v", e.message, err.Wrapped())
	} else {
		e.wrapped = fmt.Errorf("%v: %v", e.message, e.cause)
	}

	return e.wrapped
}
