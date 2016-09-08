package main

import (
	"crypto/rsa"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/colindev/x/rsa/parser"
)

func main() {

	var (
		key  string
		sign string
		size int
		err  error
	)

	flag.StringVar(&key, "key", "key.pem", "RSA Private key filename")
	flag.StringVar(&sign, "sign", "", "signed string")
	flag.IntVar(&size, "sha", 256, "RSA Private key size")
	flag.Parse()
	log.SetFlags(log.Lshortfile)

	b, err := ioutil.ReadFile(key)
	checkErr(err)

	bs, err := ioutil.ReadFile(sign)
	checkErr(err)
	signature, err := parser.Decode(string(bs))
	checkErr(err)

	priv, err := parser.Parse(b)
	checkErr(err)

	hash, hasher := parser.NewHash(size)
	if hasher == nil {
		log.Fatal(hash, "not available")
	}

	in, err := ioutil.ReadAll(os.Stdin)
	checkErr(err)

	hasher.Write(in)

	fmt.Printf("verify: %#v\n", rsa.VerifyPKCS1v15(&priv.PublicKey, hash, hasher.Sum(nil), signature))

}

func checkErr(err error) {
	if err != nil {
		log.Output(2, err.Error())
		os.Exit(1)
	}
}
