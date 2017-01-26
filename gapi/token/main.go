package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"time"

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

// https://cloud.google.com/appengine/docs/admin-api/creating-an-application
func main() {

	flag.Parse()

	b, err := ioutil.ReadFile(secretFile)
	if err != nil {
		log.Fatal(err)
	}

	var ser Secret
	if err := json.Unmarshal(b, &ser); err != nil {
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

type Secret struct {
	Type         string `json:"type"`
	ProjectID    string `json:"project_id"`
	PrivateKeyID string `json:"project_key_id"`
	PrivateKey   string `json:"private_key"`
	ClientEmail  string `json:"client_email"`
	ClientID     string `json:"client_id"`
	AuthURI      string `json:"auth_uri"`
	TokenURI     string `json:"token_uri"`
}
