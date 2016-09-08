package main

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/colindev/x/rsa/parser"
)

func main() {
	var (
		key  string
		size int
		err  error
	)
	flag.StringVar(&key, "key", "key.pem", "rsa private key filename")
	flag.IntVar(&size, "sha", 256, "rsa bits size")
	flag.Parse()
	log.SetFlags(log.Lshortfile | log.Ldate)

	b, err := ioutil.ReadFile(key)
	checkErr(err)
	priv, err := parser.Parse(b)

	in, err := ioutil.ReadAll(os.Stdin)
	checkErr(err)

	hash, hasher := parser.NewHash(size)
	if hasher == nil {
		log.Fatal(hash, "not available")
	}
	hasher.Write(in)

	signed, err := rsa.SignPKCS1v15(rand.Reader, priv, hash, hasher.Sum(nil))
	if err != nil {
		log.Fatalf("rsa.SignPKCS1v15 failed: %v\n%s", err, in)
	}

	fmt.Println(strings.TrimRight(base64.URLEncoding.EncodeToString(signed), "="))
}

func checkErr(err error) {
	if err != nil {
		log.Output(2, err.Error())
	}
}
