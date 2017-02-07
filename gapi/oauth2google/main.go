package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/colindev/x/gapi/secret"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
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

func main() {

	flag.Parse()

	ser := new(secret.OAuth2)
	if err := ser.ReadFile(secretFile); err != nil {
		log.Fatal(err)
	}

	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, &http.Client{})

	b, err := ioutil.ReadFile(secretFile)
	conf, err := google.ConfigFromJSON(b, scopes...)

	url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline)
	fmt.Printf("開啟連結並授權存取: %v\n輸入授權碼: ", url)

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatal(err)
	}

	tok, err := conf.Exchange(ctx, code)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Expiry:", tok.Expiry)
	fmt.Println(tok.AccessToken)

	// just sample
	//client := conf.Client(ctx, tok)
	//svc, err := urlshortener.New(client)
	//surl, err := svc.Url.Insert(&urlshortener.Url{
	//	LongUrl: "http://github.com/colindev/x",
	//}).Do()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Println(surl.Id)
}
