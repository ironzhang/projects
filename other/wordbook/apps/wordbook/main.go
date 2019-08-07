package main

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

	"github.com/bmizerany/pat"
	"github.com/ironzhang/wordbook/cores/bookserver"
	"github.com/ironzhang/wordbook/cores/storage"
	"github.com/ironzhang/x-pearls/log"
	"github.com/ironzhang/x/dbutil"
)

func GetOptionalInt(v url.Values, name string, def int) (int, error) {
	if s := v.Get(name); s != "" {
		return strconv.Atoi(s)
	}
	return def, nil
}

type handler struct {
	svr *bookserver.Server
}

func (h *handler) ListWords(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	offset, err := GetOptionalInt(params, "offset", 0)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if offset < 0 {
		offset = 0
	}
	limit, err := GetOptionalInt(params, "limit", 50)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if limit <= 0 {
		limit = 50
	}
	order, err := GetOptionalInt(params, "order", 1)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	words, err := h.svr.ListWords(offset, limit, order != 0)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err = json.NewEncoder(w).Encode(words); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *handler) LookupWord(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	word := params.Get(":word")
	wd, ok, err := h.svr.LookupWord(word)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if !ok {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	if err = json.NewEncoder(w).Encode(wd); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *handler) AdjustWordPriority(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	word := params.Get(":word")
	n, err := GetOptionalInt(params, "n", 0)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err = h.svr.AdjustWordPriority(word, n); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	db, err := dbutil.OpenDatabase("mysql", storage.DefaultOptions)
	if err != nil {
		log.Errorf("open database: %v", err)
		return
	}
	svr := bookserver.NewServer(db)
	hdr := handler{svr: svr}

	mux := pat.New()
	mux.Get("/api/words", http.HandlerFunc(hdr.ListWords))
	mux.Get("/api/words/:word", http.HandlerFunc(hdr.LookupWord))
	mux.Put("/api/words/:word/p", http.HandlerFunc(hdr.AdjustWordPriority))
	http.ListenAndServe(":8080", mux)
}
