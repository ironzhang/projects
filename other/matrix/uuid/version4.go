package uuid

import (
	crand "crypto/rand"
	mrand "math/rand"
	"time"
)

// NewRandom returns a Random (Version 4) UUID
func NewRandom() UUID {
	uuid := make([]byte, 16)
	randBytes(uuid)
	uuid[6] = (uuid[6] & 0x0f) | 0x40 // Version 4
	uuid[8] = (uuid[8] & 0x3f) | 0x80 // Variant is 10
	return uuid
}

var seeded = false

func randBytes(b []byte) {
	if n, err := crand.Read(b); err != nil || n != len(b) {
		if !seeded {
			mrand.Seed(time.Now().UnixNano())
			seeded = true
		}

		for i := 0; i < len(b); i++ {
			b[i] = byte(mrand.Int31n(256))
		}
	}
}
