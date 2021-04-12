package blockchain

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func TestCalculateHash(t *testing.T) {
	ts := time.Now()
	var blocks = []Block{
		{},
		{Index: 1, Time: ts, PrevHash: "", BPM: 2},
		{Index: 1, Time: ts, PrevHash: "0", BPM: 2},
	}
	for _, b := range blocks {
		//fmt.Fprintf(os.Stdout, "[%s,%d,%d,%v,%s,%s]\n", b.Time.UTC().Format(time.RFC3339Nano), b.Index, b.Difficulty, b.BPM, b.Nonce, b.PrevHash)
		if got, want := calculateHash(b), calculateHash(b); got != want {
			t.Fatalf("hash: %v != %v", got, want)
		} else {
			t.Logf("hash: %v == %v", got, want)
		}
	}
}

func TestNewBlock(t *testing.T) {
	tests := []struct {
		difficulty int
		index      int
		BPM        int
		prevHash   string
	}{
		{difficulty: 0, index: 0, BPM: 0, prevHash: ""},
		{difficulty: 1, index: 0, BPM: 0, prevHash: ""},
		{difficulty: 4, index: 0, BPM: 0, prevHash: ""},
	}
	for _, tt := range tests {
		b := NewBlock(tt.difficulty, tt.index, tt.BPM, tt.prevHash)
		if isValidHash(b.Hash, tt.difficulty) {
			t.Logf("valid hash: hash=%q, difficulty=%d", b.Hash, tt.difficulty)
		} else {
			t.Errorf("invalid hash: hash=%q, difficulty=%d", b.Hash, tt.difficulty)
		}
	}
}

func TestBlockchain(t *testing.T) {
	chain := NewBlockchain()
	chain.AddBlock(1, 1)
	chain.AddBlock(1, 2)
	chain.AddBlock(2, 3)
	chain.AddBlock(3, 4)
	data, _ := json.MarshalIndent(chain, "", "\t")
	fmt.Printf("%s\n", data)
}
