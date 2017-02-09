package spec

import "fmt"

type Scopes []string

func (s *Scopes) Set(v string) error {
	*s = append(*s, v)
	return nil
}

func (s *Scopes) String() string {
	return fmt.Sprintf("%v", ([]string)(*s))
}
