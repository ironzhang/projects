package main

import (
	"os"

	"github.com/ironzhang/practice/x-chain/blockchain"
	"github.com/ironzhang/practice/x-chain/blockchain/pow"
	"github.com/ironzhang/practice/x-chain/command"
	"github.com/ironzhang/x-pearls/log"
)

func init() {
	//log.SetLevel("trace")
	pow.POW.SetBits(16)
}

func main() {
	chain, err := blockchain.NewChain("blockchain.db", "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa")
	if err != nil {
		log.Fatalw("new chain", "error", err)
	}
	defer chain.Close()
	command.NewExecuter(chain).Execute(os.Args)
}
