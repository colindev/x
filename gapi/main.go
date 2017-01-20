package main

import (
	"flag"
	"fmt"
	"log"

	"golang.org/x/oauth2/jwt"
)

type Scopes []string

func (s Scopes) Set(v string) error {
	s = append(s, v)
	return nil
}

func (s Scopes) String() string {
	return fmt.Sprintf("%v", ([]string)(s))
}

var (
	secretFile string
	scopes     = Scopes{}
)

func init() {
	flag.StringVar(&secretFile, "secret-file", "", "secret file")
	flag.Var(&scopes, "scope", "scope")

	log.SetFlags(log.Lshortfile)
}

// https://develoopers.google.com/identity/protocols/OAuth2ServiceAccount#expiration
// https://sourcegraph.com/github.com/golang/oauth2@HEAD/-/blob/jwt/jwt.go#L72-80
func main() {

	flag.Parse()

	conf := &jwt.Config{
		Scopes: ([]string)(scopes),
	}
}
