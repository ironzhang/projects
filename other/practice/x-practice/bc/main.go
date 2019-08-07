package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/ironzhang/practice/x-practice/bc/blockchain"
)

func main() {
	addr := ":8080"

	s := Server{difficulty: 7, chain: blockchain.NewBlockchain()}
	h := pat.New()
	h.Post("/", http.HandlerFunc(s.AddBlock))
	h.Get("/", http.HandlerFunc(s.GetBlockchain))

	log.Printf("listen and serve on %s", addr)
	if err := http.ListenAndServe(addr, h); err != nil {
		log.Fatalf("listen and serve on %s: %v", addr, err)
	}
}

type Server struct {
	difficulty int
	chain      blockchain.Blockchain
}

type Message struct {
	BPM int
}

func (s *Server) AddBlock(w http.ResponseWriter, r *http.Request) {
	var m Message
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	b := s.chain.AddBlock(s.difficulty, m.BPM)
	data, err := json.MarshalIndent(b, "", "\t")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

func (s *Server) GetBlockchain(w http.ResponseWriter, r *http.Request) {
	data, err := json.MarshalIndent(s.chain, "", "\t")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(data)
}
