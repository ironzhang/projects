package blockchain

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"errors"

	"github.com/ironzhang/practice/golang/base58"
	"golang.org/x/crypto/ripemd160"
)

const version = byte(0x00)

type Wallet struct {
	PrivateKey *ecdsa.PrivateKey
	PublicKey  []byte
}

func NewWallet() *Wallet {
	private, public := newKeyPair()
	return &Wallet{
		PrivateKey: private,
		PublicKey:  public,
	}
}

func (w *Wallet) GetAddress() []byte {
	return pubKey2Address(w.PublicKey)
}

func newKeyPair() (*ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()
	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		panic(err)
	}
	public := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)
	return private, public
}

func pubKey2Address(pubKey []byte) []byte {
	hash := hashPubKey(pubKey)
	return base58CheckEncode(hash)
}

func hashPubKey(pubKey []byte) []byte {
	pubSHA256 := sha256.Sum256(pubKey)
	pubRIPEMD160 := ripemd160Sum(pubSHA256[:])
	return pubRIPEMD160
}

func ripemd160Sum(b []byte) []byte {
	h := ripemd160.New()
	h.Write(b)
	return h.Sum(nil)
}

func base58CheckEncode(payload []byte) []byte {
	data := append([]byte{version}, payload...)
	sum := checksum(data)
	return base58.StdEncoding.Encode(append(data, sum...))
}

func base58CheckDecode(data []byte) (payload []byte, err error) {
	dec := base58.StdEncoding.Decode(data)
	if dec[0] != version {
		return nil, errors.New("base58 check decode: version is unmatch")
	}
	if sum := checksum(dec[:len(dec)-4]); !bytes.Equal(sum, dec[len(dec)-4:]) {
		return nil, errors.New("base58 check decode: checksum is wrong")
	}
	return dec[1 : len(dec)-4], nil
}

func checksum(b []byte) []byte {
	firstSHA256 := sha256.Sum256(b)
	secondSHA256 := sha256.Sum256(firstSHA256[:])
	return secondSHA256[:4]
}
