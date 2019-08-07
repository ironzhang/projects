package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/bmizerany/pat"
)

func main() {
	addr := ":7000"
	h := newHandler()
	http.Handle("/", h)
	if err := http.ListenAndServe(addr, nil); err != nil {
		fmt.Fprintf(os.Stderr, "http listen and serve on [%s] failed, err[%v]\n", addr, err)
		os.Exit(-1)
	}
}

func newHandler() http.Handler {
	mux := pat.New()
	mux.Post("/ping", http.HandlerFunc(handlePostPing))
	return mux
}

func handlePostPing(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "pong")
}
