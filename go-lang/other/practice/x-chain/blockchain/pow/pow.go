package pow

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math"
	"math/big"
)

type pow struct {
	bits    uint8
	maxhash big.Int
}

var POW = NewPOW(24)

func NewPOW(bits uint8) *pow {
	p := new(pow)
	p.SetBits(bits)
	return p
}

func (p *pow) GetBits() uint8 {
	return p.bits
}

func (p *pow) SetBits(bits uint8) {
	p.bits = bits
	p.maxhash.SetInt64(1)
	p.maxhash.Lsh(&p.maxhash, 256-uint(bits))
}

func (p *pow) Hash(b []byte) (nonce uint64, hash []byte) {
	const maxNonce = math.MaxUint64

	var hint big.Int
	for nonce = 0; nonce <= maxNonce; nonce++ {
		h := sha256.New()
		h.Write(b)
		binary.Write(h, binary.BigEndian, p.bits)
		binary.Write(h, binary.BigEndian, nonce)
		hash = h.Sum(nil)
		hint.SetBytes(hash)
		if hint.Cmp(&p.maxhash) < 0 {
			return nonce, hash
		}
	}

	panic(fmt.Sprintf("pow.Hash: %d bits is too big, can not find suitable hash", p.bits))
}

func (p *pow) Validate(b []byte, nonce uint64, hash []byte) bool {
	var hint big.Int
	hint.SetBytes(hash)
	if hint.Cmp(&p.maxhash) >= 0 {
		return false
	}

	h := sha256.New()
	h.Write(b)
	binary.Write(h, binary.BigEndian, p.bits)
	binary.Write(h, binary.BigEndian, nonce)
	newhash := h.Sum(nil)
	if !bytes.Equal(hash, newhash) {
		return false
	}
	return true
}
