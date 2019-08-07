package main

import (
	"log"

	"github.com/btcsuite/btcd/rpcclient"
)

func main() {
	cfg := rpcclient.ConnConfig{
		Host: "localhost:18333",
		//Endpoint:     "ws",
		User:         "kek",
		Pass:         "kek",
		DisableTLS:   true,
		HTTPPostMode: true,
	}
	c, err := rpcclient.New(&cfg, nil)
	if err != nil {
		log.Fatalf("rpcclient.New: %v", err)
	}
	unspent, err := c.ListUnspent()
	if err != nil {
		log.Fatalf("c.ListUnspent: %v", err)
	}
	for _, res := range unspent {
		log.Printf("utxo: %v\n", res)
	}
}
