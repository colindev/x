package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"log"
	"os"
)

func main() {

	var (
		outFile string
		rsaBits int
	)
	flag.StringVar(&outFile, "o", "key.pem", "outout file name")
	flag.IntVar(&rsaBits, "b", 256, "key size")
	flag.Parse()

	priv, err := rsa.GenerateKey(rand.Reader, rsaBits)
	if err != nil {
		log.Fatal(err)
	}

	keyOut, err := os.Create(outFile)
	if err != nil {
		log.Fatal(err)
	}

	pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})
	keyOut.Close()
}
