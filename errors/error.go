package errors

import "fmt"

type (
	// Error defined a addition method for debug trace
	Error interface {
		error
		Debug() string
	}

	err struct {
		stack   []step
		current step
	}
)

func (e *err) Error() string {
	if len(e.stack) > 0 {
		return fmt.Sprintf("%s: %s", e.stack[0].msg, e.current.msg)
	}
	return e.current.msg
}

func (e *err) Debug() string {
	ss := e.Error() + "\n"
	stack := append([]step{e.current}, e.stack...)
	for i, s := range stack {
		ss = fmt.Sprintf("%s%d %s %d: %v\n", ss, i, s.file, s.line, s.msg)
	}

	return ss
}

func (e *err) wrap(msg string) Error {
	var stack []step
	copy(stack, e.stack)
	return &err{
		stack:   append([]step{e.current}, e.stack...),
		current: newStep(msg, 3),
	}
}

// New adapt origin errors.New
func New(msg string) Error {
	return &err{
		current: newStep(msg, 2),
		stack:   []step{},
	}
}

// Wrap create a trace chain for errors
func Wrap(e Error, msg string) Error {
	if e == nil {
		return nil
	}
	return e.(*err).wrap(msg)
}
