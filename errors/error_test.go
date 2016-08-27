package errors

import (
	"fmt"
	"testing"
)

func TestWrap(t *testing.T) {

	var e Error

	e = Wrap(e, "ignored")
	if e != nil {
		t.Error(e, "should be nil")
	}

	e = New("a")
	e = Wrap(e, "b")
	e = Wrap(e, "c")

	if e.Error() != "c: b" {
		t.Errorf("expect [c: b], but %v", e)
	}
}

func ExampleError_Debug() {

	var e Error

	e = New("a")
	e = Wrap(e, "b")
	e = Wrap(e, "c")

	fmt.Println(e.Debug())
}
