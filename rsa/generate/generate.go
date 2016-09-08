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
		outKeyFile string
		outPubFile string
		rsaBits    int
	)
	flag.StringVar(&outKeyFile, "key", "key.pem", "outout private key filename")
	flag.StringVar(&outPubFile, "pub", "pub.pem", "outout public key filename")
	flag.IntVar(&rsaBits, "len", 2048, "key size")
	flag.Parse()
	log.SetFlags(log.Lshortfile)

	priv, err := rsa.GenerateKey(rand.Reader, rsaBits)
	checkErr(err)

	keyOut, err := os.Create(outKeyFile)
	checkErr(err)
	pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})
	keyOut.Close()

	pubOut, err := os.Create(outPubFile)
	checkErr(err)
	pub, err := x509.MarshalPKIXPublicKey(&priv.PublicKey)
	checkErr(err)
	pem.Encode(pubOut, &pem.Block{Type: "RSA PUBLIC KEY", Bytes: pub})
	pubOut.Close()
}

func checkErr(err error) {
	if err != nil {
		log.Output(2, err.Error())
		os.Exit(1)
	}
}
