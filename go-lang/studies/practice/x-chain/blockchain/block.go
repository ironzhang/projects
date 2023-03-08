package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"io"
	"time"

	"github.com/ironzhang/practice/x-chain/blockchain/pow"
	"github.com/ironzhang/x-pearls/log"
)

type Block struct {
	Time         time.Time      // 出块时间
	Prev         []byte         // 前一个区块的hash
	Bits         uint8          // 出块难度
	Nonce        uint64         // nonce
	Hash         []byte         // 当前区块的hash
	Transactions []*Transaction // 事物列表
}

func (b Block) Validate() bool {
	buf := blockPrefixBytes(b.Time, b.Prev, b.Transactions)
	return pow.NewPOW(b.Bits).Validate(buf, b.Nonce, b.Hash)
}

func (b Block) Write(w io.Writer) {
	fmt.Fprintf(w, "Time: %s\n", b.Time)
	fmt.Fprintf(w, "Prev: %x\n", b.Prev)
	fmt.Fprintf(w, "Diff: %d\n", b.Bits)
	fmt.Fprintf(w, "Nonce: %d\n", b.Nonce)
	fmt.Fprintf(w, "Hash: %x\n", b.Hash)
	for i, tx := range b.Transactions {
		fmt.Fprintf(w, "Transactions[%d]:\n", i)
		for idx, in := range tx.Ins {
			fmt.Fprintf(w, "\tIns[%d]: TxID=%x, OutIndex=%d\n", idx, in.TxID, in.OutIndex)
		}
		for idx, out := range tx.Outs {
			fmt.Fprintf(w, "\tOuts[%d]: Value=%d, PubKeyHash=%s\n", idx, out.Value, hexString(out.PubKeyHash))
		}
	}
}

func (b Block) Marshal() ([]byte, error) {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(b); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (b *Block) Unmarshal(data []byte) error {
	return gob.NewDecoder(bytes.NewReader(data)).Decode(b)
}

func newBlock(prev []byte, txs []*Transaction) Block {
	log.Tracew("new block", "prev", hexString(prev))

	now := time.Now().UTC()
	buf := blockPrefixBytes(now, prev, txs)
	bits := pow.POW.GetBits()
	nonce, hash := pow.POW.Hash(buf)
	return Block{
		Time:         now,
		Prev:         prev,
		Bits:         bits,
		Nonce:        nonce,
		Hash:         hash,
		Transactions: txs,
	}
}

var genesisPrevHash = bytes.Repeat([]byte{0}, 32)

func newGenesisBlock(coinbase *Transaction) Block {
	return newBlock(genesisPrevHash, []*Transaction{coinbase})
}

func blockPrefixBytes(t time.Time, prev []byte, txs []*Transaction) []byte {
	var b bytes.Buffer
	buf, _ := t.MarshalBinary()
	b.Write(buf)
	b.Write(prev)
	b.Write(hashTransactions(txs))
	return b.Bytes()
}

func hashTransactions(txs []*Transaction) []byte {
	var hashes [][]byte
	for _, tx := range txs {
		hashes = append(hashes, tx.ID)
	}
	hash := sha256.Sum256(bytes.Join(hashes, []byte{}))
	return hash[:]
}

func hexString(b []byte) string {
	return fmt.Sprintf("0x%x", b)
}
