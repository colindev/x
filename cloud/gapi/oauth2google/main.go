package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strings"

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

// server side web app
func main() {

	flag.Parse()

	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, &http.Client{})

	b, err := ioutil.ReadFile(secretFile)
	if err != nil {
		log.Fatal(err)
	}

	conf, err := google.ConfigFromJSON(b, scopes...)
	if err != nil {
		log.Fatal(err)
	}

	const callbackAddr = ":8000"
	code := make(chan string, 1)
	listener, err := net.Listen("tcp", callbackAddr)
	if err != nil {
		log.Fatal(err)
	}
	go http.Serve(listener, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		code <- r.FormValue("code")
		w.Write([]byte("thx and please close window"))
		listener.Close()

	}))

	url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline, oauth2.SetAuthURLParam("prompt", "consent"))
	if err := exec.Command("google-chrome", url).Run(); err != nil {
		log.Fatal(err)
	}
	log.Println("open browser...")

	c := <-code
	log.Println("code is:", c)
	tok, err := conf.Exchange(ctx, c)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Expiry:", tok.Expiry)
	log.Println("Refresh Token:", tok.RefreshToken)
	fmt.Println(tok.AccessToken)

	fmt.Println("按下 Enter 重新調用 TokenSource")
	var src oauth2.TokenSource
	stdin := bufio.NewReader(os.Stdin)
	for {
		src = conf.TokenSource(ctx, tok)
		tok, err = src.Token()
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
