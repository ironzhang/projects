package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
	"testing"
)

func TestGenerateKey(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 256)
	if err != nil {
		t.Fatal(err)
	}

	pri := x509.MarshalPKCS1PrivateKey(key)
	pub, err := x509.MarshalPKIXPublicKey(key.Public())
	if err != nil {
		t.Fatal(err)
	}

	priBlock := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: pri}
	pubBlock := &pem.Block{Type: "PUBLIC KEY", Bytes: pub}
	pem.Encode(os.Stdout, priBlock)
	pem.Encode(os.Stdout, pubBlock)
}
