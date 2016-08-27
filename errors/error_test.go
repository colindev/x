package errors

import "testing"

func TestWrap(t *testing.T) {

	var e Error

	en := Wrap(e, "ignored")
	if en != nil {
		t.Error(en, "should be nil")
	}

	ea := New("a")
	eb := Wrap(ea, "b")
	ec := Wrap(eb, "c")

	t.Log(ec)
}

func TestDebug(t *testing.T) {

	ea := New("a")
	eb := Wrap(ea, "b")
	ec := Wrap(eb, "c")

	t.Log(ec.Debug())
}
