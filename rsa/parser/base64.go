package parser

import (
	"encoding/base64"
	"strings"
)

func Encode(b []byte) string {
	return strings.TrimRight(base64.URLEncoding.EncodeToString(b), "=")
}

func Decode(s string) ([]byte, error) {
	s = strings.TrimSpace(s)
	if l := len(s) % 4; l != 0 {
		s += strings.Repeat("=", 4-l)
	}

	return base64.URLEncoding.DecodeString(s)
}
