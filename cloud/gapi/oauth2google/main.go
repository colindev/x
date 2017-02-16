package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
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
	"time"

	"github.com/colindev/x/cloud/gapi/spec"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	secretFile string
	scopes     spec.Scopes
	out        string
)

func init() {
	flag.StringVar(&secretFile, "secret-file", "", "secret file")
	flag.StringVar(&out, "out", "", "output file")
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

	projectID := getProjectID(b)

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

	log.Println("OAuth URL:", url)
	log.Println("Expiry:", tok.Expiry)
	log.Println("Refresh Token:", tok.RefreshToken)
	log.Println(tok.AccessToken)
	log.Println("ProjectID:", projectID)
	log.Println("Scopes:", scopes)

	src := conf.TokenSource(ctx, tok)
	stdin := bufio.NewReader(os.Stdin)
	var (
		method string
		uri    string
		body   io.Reader
	)
	for {
		line, _, err := stdin.ReadLine()

		if err == io.EOF {
			return
		} else if err != nil {
			log.Println(err)
		}

		s := string(line)

		switch strings.ToUpper(s) {
		case "QUIT":
			fmt.Println("bye")
			return
		case "GET", "PUT", "POST", "DELETE":
			fmt.Printf("\033[2mmethod is [\033[m%s\033[2m]\033[m\n", s)
			method = s
		default:
			if strings.HasPrefix(s, "http") {
				uri = strings.Replace(s, "{project}", projectID, -1)
				fmt.Printf("\033[2muri is [\033[m%s\033[2m]\033[m\n", uri)
			} else {
				fmt.Printf("\033[2mbody is [\033[m%s\033[2m]\033[m\n", s)
				body = strings.NewReader(s)
			}
		}

		if method != "" && uri != "" && body != nil {
			fmt.Println("\033[2mstart fetch api\033[m", time.Now())
			req, err := http.NewRequest(strings.ToUpper(method), uri, body)
			if err != nil {
				log.Println(err)
				continue
			}

			tok, err := src.Token()
			if err != nil {
				log.Println(err)
				continue
			}
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tok.AccessToken))

			fmt.Printf("\033[2mrequest is %+v\033[m\n", req)

			res, err := (&http.Client{}).Do(req)
			if err != nil {
				log.Println("\033[31m", err, "\033[m")
			}

			method = ""
			uri = ""
			body = nil

			fmt.Println(readRes(res))
			fmt.Println("\033[32m--------------------\033[m")
		}
	}
}

func readRes(res *http.Response) string {
	if res == nil {
		return ""
	}

	b, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return err.Error()
	}

	buf := bytes.NewBuffer(nil)
	if err := json.Indent(buf, b, "", "  "); err != nil {
		return string(b)
	}

	return buf.String()
}

func getProjectID(b []byte) (pid string) {

	var conf map[string]interface{}

	if err := json.Unmarshal(b, &conf); err != nil {
		return
	}

	_pid, exists := conf["web"].(map[string]interface{})["project_id"]
	if !exists {
		return
	}

	pid, _ = _pid.(string)
	return
}
