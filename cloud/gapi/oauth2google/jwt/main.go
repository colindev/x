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
	"net/http"
	"os"
	"strings"

	"golang.org/x/oauth2/google"

	"github.com/colindev/x/cloud/gapi/spec"
)

var (
	secretFile string
	scopes     spec.Scopes
	out        string
)

func init() {
	flag.StringVar(&secretFile, "secret-file", "", "secret file path")
	flag.StringVar(&out, "out", "", "output file")
	flag.Var(&scopes, "scope", "specific scope")

	log.SetFlags(log.Lshortfile)
}

// service account
func main() {

	flag.Parse()

	b, err := ioutil.ReadFile(secretFile)
	if err != nil {
		log.Fatal(err)
	}

	projectID := getProjectID(b)

	conf, err := google.JWTConfigFromJSON(b, scopes...)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	src := conf.TokenSource(ctx)

	if out != "" {
		tok, err := src.Token()
		log.Println("Expiry:", tok.Expiry)
		log.Printf("%#v\n%v\n", tok, err)
		ioutil.WriteFile(out, []byte(tok.AccessToken), 0777)
		return
	}

	// just display
	tok, err := src.Token()
	log.Println("Expiry:", tok.Expiry)
	log.Println(tok.AccessToken)
	log.Println("ProjectID:", projectID)

	// get http client
	// use api lib
	//client := oauth2.NewClient(ctx, src)
	//service, err := appengine.New(client)
	//if err != nil {
	//	log.Fatal(err)
	//}

	//oper, err := service.Apps.Create(&appengine.Application{
	//	Id:         "test-20170214-1640",
	//	LocationId: "asia-northeast1",
	//}).Do()
	// oper, err := service.Apps.Get("gcetest-156204").Do()
	//log.Printf("%#v\n%v\n", oper, err)
	//return

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
			fmt.Println("\033[2mstart fetch api\033[m")
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

	_pid, exists := conf["project_id"]
	if !exists {
		return
	}

	pid, _ = _pid.(string)
	return
}
