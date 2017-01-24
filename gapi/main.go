package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"

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

var (
	secretFile string
	scopes     = Scopes{}
)

func init() {
	flag.StringVar(&secretFile, "secret-file", "", "secret file")
	flag.Var(&scopes, "scope", "scope")

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
