package blockchain

import (
	"crypto/ecdsa"
	"crypto/rand"
	"errors"
	"fmt"

	"github.com/boltdb/bolt"
)

var (
	ErrBlocksBucketNotExist = errors.New("blocks bucket not exist")
	ErrLastBlockKeyNotExist = errors.New("last block key not exist")
)

var (
	blocksBucket = []byte("blocks")
	lastBlockKey = []byte("last")
)

type Chain struct {
	db *bolt.DB
}

func NewChain(path, address string) (*Chain, error) {
	db, err := bolt.Open(path, 0666, nil)
	if err != nil {
		return nil, err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(blocksBucket)
		if b == nil {
			coinbase := NewCoinbaseTX(address, "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks")
			genesis := newGenesisBlock(coinbase)
			blockBytes, err := genesis.Marshal()
			if err != nil {
				return err
			}
			if b, err = tx.CreateBucket(blocksBucket); err != nil {
				return err
			}
			if err = b.Put(genesis.Hash, blockBytes); err != nil {
				return err
			}
			if err = b.Put(lastBlockKey, genesis.Hash); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &Chain{db: db}, nil
}

func (p *Chain) Close() error {
	return p.db.Close()
}

func (p *Chain) AddBlock(txs []*Transaction) error {
	err := p.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(blocksBucket)
		if b == nil {
			return ErrBlocksBucketNotExist
		}
		lastHash := b.Get(lastBlockKey)
		if lastHash == nil {
			return ErrLastBlockKeyNotExist
		}

		block := newBlock(lastHash, txs)
		blockBytes, err := block.Marshal()
		if err != nil {
			return err
		}
		if err = b.Put(block.Hash, blockBytes); err != nil {
			return err
		}
		if err = b.Put(lastBlockKey, block.Hash); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (p *Chain) LastBlock() (Block, error) {
	var block Block
	err := p.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(blocksBucket)
		if b == nil {
			return ErrBlocksBucketNotExist
		}
		lastHash := b.Get(lastBlockKey)
		if lastHash == nil {
			return ErrLastBlockKeyNotExist
		}
		blockBytes := b.Get(lastHash)
		if blockBytes == nil {
			return fmt.Errorf("not found %x key's value", lastHash)
		}
		if err := block.Unmarshal(blockBytes); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return Block{}, err
	}
	return block, nil
}

func (p *Chain) PrevBlock(b Block) (Block, bool, error) {
	return p.GetBlock(b.Prev)
}

func (p *Chain) GetBlock(hash []byte) (block Block, exist bool, err error) {
	err = p.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(blocksBucket)
		if b == nil {
			return ErrBlocksBucketNotExist
		}
		blockBytes := b.Get(hash)
		if blockBytes == nil {
			return nil
		}
		if err := block.Unmarshal(blockBytes); err != nil {
			return err
		}
		exist = true
		return nil
	})
	if err != nil {
		return Block{}, false, err
	}
	return block, exist, nil
}

func (p *Chain) FindUTXO(address string) ([]UTXO, error) {
	var utxos []UTXO
	var stxos TxOutSet

	var ok bool
	block, err := p.LastBlock()
	if err != nil {
		return nil, err
	}
	for {
		for _, tx := range block.Transactions {
			// 查找未花费的TxOut
			for idx, out := range tx.Outs {
				// 是否已花费
				if stxos.HasTxOut(tx.ID, idx) {
					continue
				}

				// 是否属于address
				if !out.BelongTo(address) {
					continue
				}

				utxo := UTXO{
					TxID:       tx.ID,
					Index:      idx,
					Value:      out.Value,
					PubKeyHash: out.PubKeyHash,
				}

				utxos = append(utxos, utxo)
			}

			// 将以花费的TxOut加入已花费集合
			if !tx.IsCoinbase() {
				for _, in := range tx.Ins {
					stxos.AddTxOut(in.TxID, in.OutIndex)
				}
			}
		}

		block, ok, err = p.PrevBlock(block)
		if err != nil {
			return nil, err
		}
		if !ok {
			break
		}
	}
	return utxos, nil
}

func (p *Chain) FindSpendableUTXO(address string, amount int) (int, []UTXO, error) {
	utxos, err := p.FindUTXO(address)
	if err != nil {
		return 0, nil, err
	}
	accumulated := 0
	for i, utxo := range utxos {
		accumulated += utxo.Value
		if accumulated >= amount {
			return accumulated, utxos[:i+1], nil
		}
	}
	return 0, nil, fmt.Errorf("balance(%d) is not enough", accumulated)
}

func (p *Chain) NewTransaction(from, to string, amount int) (*Transaction, error) {
	accumulated, utxos, err := p.FindSpendableUTXO(from, amount)
	if err != nil {
		return nil, err
	}
	var ins []TxIn
	for _, utxo := range utxos {
		ins = append(ins, NewTxIn(utxo.TxID, utxo.Index))
	}
	var outs []TxOut
	outs = append(outs, NewTxOut(amount, to))
	if accumulated > amount {
		outs = append(outs, NewTxOut(accumulated-amount, from))
	}
	return NewTransaction(ins, outs), nil
}

func (p *Chain) SignTransaction(tx *Transaction, pk *ecdsa.PrivateKey) error {
	for _, in := range tx.Ins {
		in.PubKey.X = pk.PublicKey.X
		in.PubKey.Y = pk.PublicKey.Y
	}

	hash := tx.Hash()
	for _, in := range tx.Ins {
		r, s, err := ecdsa.Sign(rand.Reader, pk, hash)
		if err != nil {
			return err
		}
	}

	return nil
}

type TxOutSet struct {
	m map[string][]int
}

func (p *TxOutSet) AddTxOut(txid []byte, outidx int) {
	if p.m == nil {
		p.m = make(map[string][]int)
	}

	key := hexString(txid)
	p.m[key] = append(p.m[key], outidx)
}

func (p *TxOutSet) HasTxOut(txid []byte, outidx int) bool {
	key := hexString(txid)
	if indexs, ok := p.m[key]; ok {
		for _, idx := range indexs {
			if idx == outidx {
				return true
			}
		}
	}
	return false
}
