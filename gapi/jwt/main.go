package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/colindev/x/gapi/secret"

	"golang.org/x/oauth2/jwt"
)

type Scopes []string

func (s *Scopes) Set(v string) error {
	*s = append(*s, v)
	return nil
}

func (s *Scopes) String() string {
	return fmt.Sprintf("%v", ([]string)(*s))
}

type Exp time.Duration

func (exp *Exp) Set(s string) (err error) {
	d, err := time.ParseDuration(s)
	*exp = Exp(d)
	return
}

func (exp *Exp) String() string {
	return time.Duration(*exp).String()
}

var (
	secretFile string
	scopes     = Scopes{}
	exp        Exp
)

func init() {
	flag.StringVar(&secretFile, "secret-file", "", "secret file")
	flag.Var(&scopes, "scope", "scope")
	flag.Var(&exp, "exp", "expires")

	log.SetFlags(log.Lshortfile)
}

func main() {

	flag.Parse()

	ser := new(secret.Service)
	if err := ser.ReadFile(secretFile); err != nil {
		log.Fatal(err)
	}

	conf := &jwt.Config{
		Email:      ser.ClientEmail,
		PrivateKey: []byte(ser.PrivateKey),
		TokenURL:   ser.TokenURI,
		Scopes:     []string(scopes),
		// TODO 目前測試 超過 1h 會被拒絕, 5m 會被忽略一樣是 1h
		Expires: time.Duration(exp),
	}

	tok, err := conf.TokenSource(context.Background()).Token()
	if err != nil {
		log.Fatal(err)
	}

	if !tok.Valid() {
		log.Fatal("valid fail")
	}

	log.Println("Expiry at:", tok.Expiry)
	fmt.Println(tok.AccessToken)

	// c := conf.Client(context.Background())
	// aes, err := appengine.New(c)
}
