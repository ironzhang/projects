package pow

import "testing"

func TestPow(t *testing.T) {
	tests := []struct {
		bits uint8
		data []byte
	}{
		{bits: 1, data: []byte("1")},
		{bits: 2, data: []byte("1")},
		{bits: 3, data: []byte("1")},
		{bits: 4, data: []byte("1")},
		{bits: 5, data: []byte("1")},
		{bits: 6, data: []byte("1")},
		{bits: 7, data: []byte("1")},
		{bits: 8, data: []byte("1")},
		{bits: 16, data: []byte("1")},
		//{bits: 24, data: []byte("1")},
	}
	for _, tt := range tests {
		p := NewPOW(tt.bits)
		nonce, hash := p.Hash(tt.data)
		t.Logf("bits: %d, nonce=%d, hash=%x", p.GetBits(), nonce, hash)
	}
}
