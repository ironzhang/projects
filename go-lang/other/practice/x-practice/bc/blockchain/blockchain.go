package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
)

type Block struct {
	Time       time.Time
	Difficulty int
	Index      int
	BPM        int
	Nonce      string
	Hash       string
	PrevHash   string
}

func calculateHash(b Block) string {
	h := sha256.New()
	fmt.Fprintf(h, "[%s,%d,%d,%v,%s,%s]", b.Time.UTC().Format(time.RFC3339Nano), b.Difficulty, b.Index, b.BPM, b.Nonce, b.PrevHash)
	hash := h.Sum(nil)
	return hex.EncodeToString(hash)
}

func isValidHash(hash string, difficulty int) bool {
	prefix := strings.Repeat("0", difficulty)
	return strings.HasPrefix(hash, prefix)
}

func NewBlock(difficulty, index, BPM int, prevHash string) Block {
	b := Block{
		Time:       time.Now(),
		Difficulty: difficulty,
		Index:      index,
		BPM:        BPM,
		PrevHash:   prevHash,
	}
	for i := 0; ; i++ {
		b.Nonce = fmt.Sprint(i)
		hash := calculateHash(b)
		if !isValidHash(hash, difficulty) {
			continue
		}
		b.Hash = hash
		return b
	}
}

type Blockchain []Block

func NewBlockchain() Blockchain {
	return Blockchain{NewBlock(0, 0, 0, "")}
}

func (p *Blockchain) AddBlock(difficulty int, BPM int) Block {
	i := len(*p) - 1
	prev := (*p)[i]
	block := NewBlock(difficulty, i+1, BPM, prev.Hash)
	*p = append(*p, block)
	return block
}
