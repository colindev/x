package errors

import "testing"

func TestNewStep(t *testing.T) {
	s := newStep("test msg", 1)

	t.Logf("%#v", s)
}
