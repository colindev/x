package parser

import "testing"

func TestNewHash(t *testing.T) {
	h, hr := NewHash(666)

	t.Log(h, hr)
}
