package errors

import "runtime"

type step struct {
	msg  string
	file string
	line int
}

func newStep(msg string, skip int) step {

	_, file, line, _ := runtime.Caller(skip)

	return step{
		msg:  msg,
		line: line,
		file: file,
	}
}
