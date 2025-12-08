package eval

import (
	"fmt"
	"path/filepath"

	"github.com/helmtk/htkl/parser"
)

// EvalError represents an error that occurred during evaluation
type EvalError struct {
	Message  string
	Filename string
	Line     int
	Col      int
}

func (e *EvalError) Error() string {
	if e.Filename != "" {
		if e.Line > 0 {
			fn := filepath.Base(e.Filename)
			return fmt.Sprintf("[%s %d:%d] %s", fn, e.Line, e.Col, e.Message)
		}
		return fmt.Sprintf("[%s] %s", e.Filename, e.Message)
	}
	return e.Message
}

// errorf creates an error with position information from the node
func errorf(pos parser.Pos, format string, args ...interface{}) error {
	msg := fmt.Sprintf(format, args...)

	if pos.Line > 0 && pos.Filename != "" {
		return &EvalError{
			Message:  msg,
			Filename: pos.Filename,
			Line:     pos.Line,
			Col:      pos.Col,
		}
	}
	if pos.Filename != "" {
		return &EvalError{
			Message:  msg,
			Filename: pos.Filename,
		}
	}
	// No position info available
	return fmt.Errorf("%s", msg)
}
