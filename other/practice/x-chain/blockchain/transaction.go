package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"math/big"
)

type PublicKey struct {
	X *big.Int
	Y *big.Int
}

type Signature struct {
	R *big.Int
	S *big.Int
}

type TxIn struct {
	TxID      []byte
	OutIndex  int
	PubKey    PublicKey
	Signature Signature
}

func NewTxIn(txid []byte, outIndex int) TxIn {
	return TxIn{
		TxID:     txid,
		OutIndex: outIndex,
	}
}

type TxOut struct {
	Value      int
	PubKeyHash []byte
}

func NewTxOut(value int, address string) TxOut {
	hash, err := base58CheckDecode([]byte(address))
	if err != nil {
		panic(err)
	}
	return TxOut{
		Value:      value,
		PubKeyHash: hash,
	}
}

func (t TxOut) BelongTo(address string) bool {
	hash, err := base58CheckDecode([]byte(address))
	if err != nil {
		panic(err)
	}
	return bytes.Equal(t.PubKeyHash, hash)
}

type UTXO struct {
	TxID       []byte
	Index      int
	Value      int
	PubKeyHash []byte
}

type Transaction struct {
	ID   []byte
	Ins  []TxIn
	Outs []TxOut
}

func NewTransaction(ins []TxIn, outs []TxOut) *Transaction {
	tx := Transaction{
		Ins:  ins,
		Outs: outs,
	}

	var b bytes.Buffer
	gob.NewEncoder(&b).Encode(tx)
	hash := sha256.Sum256(b.Bytes())
	tx.ID = hash[:]
	return &tx
}

func (t *Transaction) Hash() []byte {
	var b bytes.Buffer
	gob.NewEncoder(&b).Encode(tx)
	hash := sha256.Sum256(b.Bytes())
	return hash[:]
}

func (t *Transaction) IsCoinbase() bool {
	return len(t.Ins) == 1 && t.Ins[0].TxID == nil && t.Ins[0].OutIndex == -1
}

func NewCoinbaseTX(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward to %s", to)
	}

	const subsidy = 10
	txin := NewTxIn(nil, -1)
	txout := NewTxOut(subsidy, to)
	return NewTransaction([]TxIn{txin}, []TxOut{txout})
}
