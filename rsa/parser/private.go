package parser

import (
	"crypto"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"hash"
)

func Parse(key []byte) (*rsa.PrivateKey, error) {
	var err error

	// Parse PEM block
	var block *pem.Block
	if block, _ = pem.Decode(key); block == nil {
		return nil, errors.New("key must be private")
	}

	var parsedKey interface{}
	if parsedKey, err = x509.ParsePKCS1PrivateKey(block.Bytes); err != nil {
		if parsedKey, err = x509.ParsePKCS8PrivateKey(block.Bytes); err != nil {
			return nil, err
		}
	}

	var pkey *rsa.PrivateKey
	var ok bool
	if pkey, ok = parsedKey.(*rsa.PrivateKey); !ok {
		return nil, errors.New("not RSA private key")
	}

	return pkey, nil
}

func NewHash(size int) (crypto.Hash, hash.Hash) {
	var (
		h  crypto.Hash
		hr hash.Hash
	)
	switch size {
	case 256:
		h = crypto.SHA256
	case 384:
		h = crypto.SHA384
	case 512:
		h = crypto.SHA512
	}

	if h.Available() {
		hr = h.New()
	}

	return h, hr
}
