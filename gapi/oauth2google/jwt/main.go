package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"golang.org/x/oauth2/google"

	"github.com/colindev/x/gapi/spec"
)

var (
	secretFile string
	scopes     spec.Scopes
)

func init() {
	flag.StringVar(&secretFile, "secret-file", "", "secret file path")
	flag.Var(&scopes, "scope", "specific scope")
}

// service account
func main() {

	flag.Parse()

	b, err := ioutil.ReadFile(secretFile)
	if err != nil {
		log.Fatal(err)
	}

	conf, err := google.JWTConfigFromJSON(b, scopes...)
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	src := conf.TokenSource(ctx)

	stdin := bufio.NewReader(os.Stdin)
	for {
		tok, err := src.Token()
		fmt.Println("Expiry:", tok.Expiry)
		fmt.Printf("%#v\n%v\n", tok, err)
		line, _, err := stdin.ReadLine()
		if err == io.EOF {
			return
		} else if err != nil {
			log.Println(err)
		}
		if strings.ToUpper(string(line)) == "QUIT" {
			return
		}
	}

}
